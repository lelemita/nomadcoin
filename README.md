# Nomadcoin

Making a Cryptocurrency using the Go programming language.  
https://nomadcoders.co/nomadcoin/lobby

### Features

* Mining
* Transactions
* Database Backend
* Wallets
* REST API
* HTML Explorer
* P2P (Websockets)
* Unit Testing

### Run example
```bash
go run main.go -mode=rest -port=4000
```

### Test
```bash
go test ./... -v
go test -v -coverprofile cover.out  ./...
# 테스트 하는 비율 % 표시 됨
go tool cover -html='cover.out'
# 녹색은 테스트 된거, 적색은 안된거
```