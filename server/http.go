package server

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"

	chain "umbrellaX/chains"

	"github.com/gorilla/mux"
)

type Server struct {
	chains map[string]chain.Chain
}

func New(chains ...chain.Chain) *Server {
	s := &Server{
		chains: make(map[string]chain.Chain),
	}
	for _, chain := range chains {
		s.chains[chain.Name()] = chain
	}

	return s
}

func (s *Server) Start(port string) error {
	r := mux.NewRouter()
	r.HandleFunc("/ton/transaction", s.handleTonTransaction).Methods("GET")
	r.HandleFunc("/tron/transaction", s.handleTronTransaction).Methods("GET")
	return http.ListenAndServe(":"+port, r)
}

type txRequest struct {
	Currency string `json:"currency"`
	From     string `json:"fromAddress"`
	To       string `json:"toAddress"`
	Amount   string `json:"amount"`
	FeeLimit string `json:"feeLimit"`
}

func (s *Server) handleTonTransaction(w http.ResponseWriter, r *http.Request) {
	txReq := getQueryTxParams(r.URL.Query())

	if txReq.To == "" || txReq.Amount == "" {
		http.Error(w, "missing required parameters", http.StatusBadRequest)
		return
	}

	var amountFloat float64
	_, err := fmt.Sscan(txReq.Amount, &amountFloat)
	if err != nil {
		http.Error(w, "invalid amount", http.StatusBadRequest)
		return
	}

	hash, err := s.chains["ton"].SendTx(txReq.Currency, txReq.To, amountFloat, 0)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	response := map[string]string{"hash": hash}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func (s *Server) handleTronTransaction(w http.ResponseWriter, r *http.Request) {
	txReq := getQueryTxParams(r.URL.Query())

	if txReq.To == "" || txReq.Amount == "" {
		http.Error(w, "missing required parameters", http.StatusBadRequest)
		return
	}

	var amountFloat float64
	_, err := fmt.Sscan(txReq.Amount, &amountFloat)
	if err != nil {
		http.Error(w, "invalid amount", http.StatusBadRequest)
		return
	}

	hash, err := s.chains["tron"].SendTx(txReq.Currency, txReq.To, amountFloat, 10)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"hash": hash})
}

func getQueryTxParams(urlPath url.Values) txRequest {
	return txRequest{
		Currency: urlPath.Get("currency"),
		From:     urlPath.Get("fromAddress"),
		To:       urlPath.Get("toAddress"),
		Amount:   urlPath.Get("amount"),
		FeeLimit: urlPath.Get("feeLimit"),
	}
}
