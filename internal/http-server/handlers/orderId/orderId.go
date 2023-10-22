package orderId

import (
	"encoding/json"
	"github.com/go-chi/chi/v5"
	"github.com/sema0205/WB_L0/internal/models"
	"net/http"
)

type OrdersGetter interface {
	GetOrders() ([]models.OrderInfo, error)
}

func New(cache map[string]models.OrderInfo) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		uuid := chi.URLParam(r, "uuid")
		model, ok := cache[uuid]
		if !ok {
			w.Write([]byte("no such uuid published"))
			return
		} else {
			jsonResp, _ := json.Marshal(model)
			w.Write(jsonResp)
			return
		}

	}
}
