package v1

import (
	"github.com/go-chi/chi"
	swdoc "github.com/leandroag/desafio/docs"
	swhttp "github.com/swaggo/http-swagger"
)

type handler interface {
	RegisterRoutes(router *chi.Mux)
}

// API
// @Title Desafio
// @Description Desafio REST API.
// @Version 0.0.1
// @Schemes https
func Setup(router *chi.Mux, handlers ...handler) {
	for _, handler := range handlers {
		handler.RegisterRoutes(router)
	}

	router.Get("/docs/v1/swagger/*", swhttp.Handler(
		swhttp.InstanceName(swdoc.SwaggerInfov1.InfoInstanceName),
	))
}
