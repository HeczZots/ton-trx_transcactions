package server

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"net/http"

	"umbrellaX/network/ton"
	"umbrellaX/network/tron"

	"github.com/fbsobreira/gotron-sdk/pkg/proto/api"
	"github.com/gorilla/mux"
	"google.golang.org/protobuf/proto"
)

type Server struct {
	ton  *ton.Client
	tron *tron.Client
}

func New(tonClient *ton.Client, tronClient *tron.Client) *Server {
	return &Server{
		ton:  tonClient,
		tron: tronClient,
	}
}

func (s *Server) Start(port string) error {
	r := mux.NewRouter()
	r.HandleFunc("/ton/transaction", s.handleTonTransaction).Methods("GET")
	r.HandleFunc("/tron/transaction", s.handleTronTransaction).Methods("GET")
	return http.ListenAndServe(":"+port, r)
}

func (s *Server) handleTonTransaction(w http.ResponseWriter, r *http.Request) {
	currency := r.URL.Query().Get("currency")
	fromAddress := r.URL.Query().Get("fromAddress")
	toAddress := r.URL.Query().Get("toAddress")
	amount := r.URL.Query().Get("amount")

	if fromAddress == "" || toAddress == "" || amount == "" {
		http.Error(w, "missing required parameters", http.StatusBadRequest)
		return
	}

	var amountFloat float64
	_, err := fmt.Sscan(amount, &amountFloat)
	if err != nil {
		http.Error(w, "invalid amount", http.StatusBadRequest)
		return
	}

	var txHash string
	if currency == "" {
		txHash, err = s.ton.SendTon(fromAddress, toAddress, amountFloat, 0)
	} else {
		txHash, err = s.ton.SendJetton(fromAddress, toAddress, currency, amountFloat, 0)
	}

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	response := map[string]string{"hash": txHash}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func (s *Server) handleTronTransaction(w http.ResponseWriter, r *http.Request) {
	currency := r.URL.Query().Get("currency")
	fromAddress := r.URL.Query().Get("fromAddress")
	toAddress := r.URL.Query().Get("toAddress")
	amount := r.URL.Query().Get("amount")

	if fromAddress == "" || toAddress == "" || amount == "" {
		http.Error(w, "missing required parameters", http.StatusBadRequest)
		return
	}

	var amountFloat float64
	_, err := fmt.Sscan(amount, &amountFloat)
	if err != nil {
		http.Error(w, "invalid amount", http.StatusBadRequest)
		return
	}

	var tx *api.TransactionExtention
	if currency == "" || currency == "TRX" {
		tx, err = s.tron.CreateTxTRX(fromAddress, toAddress, 100, amountFloat)
	} else {
		tx, err = s.tron.CreateTxTRC20(currency, fromAddress, toAddress, 100, amountFloat)
	}
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	rawData, err := proto.Marshal(tx.Transaction.GetRawData())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	
	hash := sha256.Sum256(rawData)	

	response := map[string]string{"hash": hex.EncodeToString(hash[:])}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}
