package inventory

import (
	"artOfDevPractise/internal/apperror"
	"artOfDevPractise/internal/handlers"
	"artOfDevPractise/pkg/logging"
	"github.com/julienschmidt/httprouter"
	"net/http"
)

var _ handlers.Handler = &handler{}

const inventoryURL = "/inventory"

type handler struct {
	inv    *Inventory
	logger *logging.Logger
}

func NewHandler(inv *Inventory, logger *logging.Logger) handlers.Handler {
	return &handler{inv, logger}
}

func (h *handler) Register(router *httprouter.Router) {
	router.HandlerFunc(http.MethodGet, inventoryURL, apperror.Middleware(h.GetInventory))
}

func (h *handler) GetInventory(w http.ResponseWriter, r *http.Request) error {

	StructureInventory := h.inv.StructureInventory(true)
	w.Write([]byte(StructureInventory))

	h.logger.Debug("Get structure inventory")

	return nil
}
