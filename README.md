## Описание

```
    Простой http server для отправки транзакций в сетях TRON и TON
    взаимодействует с любыми валютами используя смарт контракт отправляемого токена 

    (tron/transaction? currency={currency}&fromAddress={fromAddress}&toAddress={toAddress}&amount={amount}),
    в response – {hash: hash}
    
    (ton/transaction? currency={currency}&fromAddress={fromAddress}&toAddress={toAddress}&amount={amount}),
    в response – {hash: hash}

    не обязательные параметры - были установлены по ТЗ
    currency
    fromAddress
```

## Запуск 

```
    Инициализировать переменные окружения внутри ОС
    TON_SEED - мнемоническая фраза кошелька
    TON_ADDRESS - адресс кошелька
    TRON_SECRET - приватный ключ кошелька
    TRON_ADDRESS - адресс кошелька

    git clone https://github.com/HeczZots/ton-trx_transcactions.git

    cd ton-trx_transcactions/cmd

    Освободить порт 8083 если он занят

    go run .
``