package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"mfv-challenge/config"
	"mfv-challenge/internal/must"
	"mfv-challenge/internal/repositories"
	"mfv-challenge/internal/services"
	"mfv-challenge/internal/usecases"

	"github.com/gorilla/mux"
)

func main() {
	cfgPath, err := config.ParseFlags()
	if err != nil {
		log.Fatal(err)
	}
	cfg, err := config.NewConfig(cfgPath)
	if err != nil {
		log.Fatal(err)
	}

	server := NewServer(cfg)
	run(server)
}

func NewServer(cfg *config.Config) *http.Server {
	db, err := must.Connect("postgres", cfg.BuildDSN())
	if err != nil {
		panic(err)
	}

	userRepository := repositories.NewUser(db)
	transactionRepository := repositories.NewTransaction(db)
	accountRepository := repositories.NewAccount(db)

	userUC := usecases.NewUser(userRepository)
	accountUC := usecases.NewAccount(accountRepository)
	transactionUC := usecases.NewTransaction(transactionRepository, accountRepository)

	// tokenBuilder := token.NewjwtHMACBuilder(cfg.JWT.Secret, cfg.JWT.Duration)
	userService := services.NewUser(userUC)
	accountService := services.NewAccount(accountUC)
	transactionService := services.NewTransaction(transactionUC)

	router := mux.NewRouter()
	router.HandleFunc("/api/users/{user_id}/transactions", transactionService.List).Methods("GET")
	router.HandleFunc("/api/users/{user_id}/transactions", transactionService.Create).Methods("POST")
	router.HandleFunc("/api/users/{user_id}", userService.Get).Methods("GET")
	router.HandleFunc("/api/users/{user_id}/accounts", userService.ListAccounts).Methods("GET")
	router.HandleFunc("/api/accounts/{account_id}", accountService.Get).Methods("GET")
	// router.HandleFunc("/api/login", authenHandler.Login).Methods("POST")

	return &http.Server{
		Addr:    cfg.Server.Host + ":" + cfg.Server.Port,
		Handler: router,
	}
}

func run(server *http.Server) {
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, os.Interrupt, syscall.SIGTERM)
	log.Printf("Server is starting on %s\n", server.Addr)

	go func() {
		if err := server.ListenAndServe(); err != nil {
			log.Fatalf("Server failed to start due to err: %v", err)
		}
	}()

	interrupt := <-ch
	log.Printf("Server is shutting down due to %+v\n", interrupt)
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()
	if err := server.Shutdown(ctx); err != nil {
		log.Fatalf("Server was unable to gracefully shutdown due to err: %+v", err)
	}
}
