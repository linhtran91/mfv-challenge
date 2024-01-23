test:
	go test ./...

mock:
	mockgen -source=internals/handlers/loan.go -destination=mocks/handlers/loan.go
	mockgen -source=internals/handlers/repayment.go -destination=mocks/handlers/repayment.go
	mockgen -source=internals/handlers/customer.go -destination=mocks/handlers/customer.go

dep:
	go mod tidy
	go mod vendor

init-db:
	migrate create -ext sql -dir migrations -seq init_tables