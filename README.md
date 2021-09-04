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
go tool cover -html='cover.out'
# 녹색은 테스트한 거, 적색은 안 된거
```

### Make Documents
```bash
# go get golang.org/x/tools/cmd/godoc
# godoc -http=:6060
```

### TO-DO
* db.getDbName(): 실행문에 -port=4000 없으면 에러나는 부분
