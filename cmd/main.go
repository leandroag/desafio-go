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
	"github.com/leandroag/desafio/app/domain/usescases/login"
	"github.com/leandroag/desafio/app/domain/usescases/transfer"

	handlerAccount "github.com/leandroag/desafio/app/gateway/api/http/v1/account"
	handlerLogin "github.com/leandroag/desafio/app/gateway/api/http/v1/login"
	handlerTransfer "github.com/leandroag/desafio/app/gateway/api/http/v1/transfer"

	"github.com/leandroag/desafio/app/gateway/bcrypt"
	"github.com/leandroag/desafio/app/gateway/db/postgres"
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

	// inicializando crypt service
	bcryptService := bcrypt.NewCryptService([]byte(cfg.JwtKey))

	// inicializa repositórios
	accountRepository := postgres.NewAccountRepository(conn)
	transferRepository := postgres.NewTransferRepository(conn)

	// inicializando casos de uso
	accountUseCase := account.NewAccountService(accountRepository, bcryptService)
	transferUseCase := transfer.NewTransferService(accountRepository, transferRepository, bcryptService)
	loginUseCase := login.NewLoginService(accountRepository, bcryptService)

	// inicializando os handlers
	accountHandler := handlerAccount.NewAccountHandler(accountUseCase)
	transferHandler := handlerTransfer.NewTransferHandler(transferUseCase, bcryptService)
	loginHandler := handlerLogin.NewLoginHandler(loginUseCase)

	// cria o roteador
	router := mux.NewRouter()

	// registra as rotas do handler de contas
	accountHandler.RegisterRoutes(router)
	transferHandler.RegisterRoutes(router)
	loginHandler.RegisterRoutes(router)

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
