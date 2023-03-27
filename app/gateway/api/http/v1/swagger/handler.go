package swagger

import (
	"github.com/go-chi/chi"
	swdoc "github.com/leandroag/desafio/docs/swagger"
	swhttp "github.com/swaggo/http-swagger"
)

type SwaggerHandler struct{}

func NewSwaggerHandler() *SwaggerHandler {
	return &SwaggerHandler{}
}

func (h SwaggerHandler) RegisterRoutes(router *chi.Mux) {
	router.Get("/docs/swagger/*", swhttp.Handler(
		swhttp.InstanceName(swdoc.SwaggerInfov1.InfoInstanceName),
	))
}
