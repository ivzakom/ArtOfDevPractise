package inventory

import (
	"github.com/julienschmidt/httprouter"
	"golang_lessons/internal/handlers"
	"golang_lessons/pkg/logging"
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
	router.GET(inventoryURL, h.GetInventory)
}

func (h *handler) GetInventory(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {

	StructureInventory := h.inv.StructureInventory(true)
	w.Write([]byte(StructureInventory))

	h.logger.Debug("Get structure inventory")
}
