package user

import (
	"artOfDevPractise/internal/apperror"
	"artOfDevPractise/internal/handlers"
	"artOfDevPractise/pkg/logging"
	"github.com/julienschmidt/httprouter"
	"net/http"
)

var _ handlers.Handler = &handler{}

const userURL = "/user"

type handler struct {
	logger *logging.Logger
}

func NewHandler(logger *logging.Logger) handlers.Handler {
	return &handler{
		logger,
	}
}

func (h *handler) Register(router *httprouter.Router) {
	router.HandlerFunc(http.MethodGet, userURL, apperror.Middleware(h.GetList))
}

func (h *handler) GetList(w http.ResponseWriter, r *http.Request) error {
	return apperror.ErrorNotFound
}
