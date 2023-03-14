package main

// import (
// 	"log"
// 	"net/http"

// 	"github.com/gorilla/mux"

// 	"github.com/leandroag/desafio/app/gateway/db"
// 	"github.com/leandroag/desafio/app/domain/usescases"
// 	"github.com/leandroag/desafio/app/gateway/api/account"
// )

// func main() {
// 	// cria o caso de uso de contas
// 	accountUseCase := usescases.NewAccountUseCase(accountRepository)

// 	// cria o handler de contas
// 	accountHandler := handler.NewAccountHandler(accountUseCase)

// 	// cria o roteador
// 	router := mux.NewRouter()

// 	// registra as rotas do handler de contas
// 	accountHandler.RegisterRoutes(router)

// 	// inicia o servidor
// 	log.Fatal(http.ListenAndServe(":8080", router))
// }
