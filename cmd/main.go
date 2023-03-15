package main

import (
	// "log"
	// "net/http"

	// "github.com/gorilla/mux"

	"context"
	"errors"
	"fmt"
	"net/http"
	"os/signal"
	"syscall"

	"github.com/gorilla/mux"
	"github.com/leandroag/desafio/app"
	"golang.org/x/sync/errgroup"

	"github.com/leandroag/desafio/app/domain/usescases/account"
	handlerAccount "github.com/leandroag/desafio/app/gateway/api/http/v1/account"
	"github.com/leandroag/desafio/app/gateway/bcrypt"
	"github.com/leandroag/desafio/app/gateway/db/postgres"
	// "github.com/leandroag/desafio/app/gateway/api/account"
)

func main() {
	// carrega as configurações
	cfg := app.ReadConfig(".env")

	// conexão com o postgres
	conn, err := postgres.New(cfg.Postgres.URL(), cfg.Postgres.PoolMinSize, cfg.Postgres.PoolMaxSize)
	if err != nil {
		fmt.Println("error", err)
	}
	defer conn.Close()

	accountRepository := postgres.NewAccountRepository(conn)

	bcryptService := bcrypt.NewCryptService([]byte(cfg.JwtKey))

	// inicializa o caso de uso de contas
	accountUseCase := account.NewAccountService(accountRepository, bcryptService)

	// inicializa o handler de contas
	accountHandler := handlerAccount.NewAccountHandler(accountUseCase)

	// cria o roteador
	router := mux.NewRouter()

	// registra as rotas do handler de contas
	accountHandler.RegisterRoutes(*router)

	// Server
	fmt.Printf(cfg.Http.Address)
	srv := &http.Server{
		Addr:              cfg.Http.Address,
		Handler:           router,
		WriteTimeout:      cfg.Http.WriteTimeout,
		ReadTimeout:       cfg.Http.ReadTimeout,
		ReadHeaderTimeout: cfg.Http.ReadTimeout,
	}

	signalCtx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer cancel()

	g, gCtx := errgroup.WithContext(signalCtx)
	g.Go(func() error {
		fmt.Printf("server listening: %s", cfg.Http.Address)

		return srv.ListenAndServe()
	})
	g.Go(func() error {
		<-gCtx.Done()
		fmt.Printf("interrupt signal received; shutting down server...")

		timeoutCtx, cancel := context.WithTimeout(context.Background(), srv.WriteTimeout)
		defer cancel()

		return srv.Shutdown(timeoutCtx)
	})

	if err := g.Wait(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		fmt.Printf("failed to listen for connections: %s", err)
	}
}
