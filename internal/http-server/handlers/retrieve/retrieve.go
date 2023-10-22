package retrieve

import (
	"encoding/json"
	"github.com/sema0205/WB_L0/internal/models"
	"net/http"
)

type OrdersGetter interface {
	GetOrders() ([]models.OrderInfo, error)
}

func New(ordersGetter OrdersGetter) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		var models []models.OrderInfo

		models, err := ordersGetter.GetOrders()
		if err != nil {
			return
		}

		if len(models) == 0 {
			w.Write([]byte("nothing was published yet"))
		}

		for _, model := range models {
			jsonResp, _ := json.Marshal(model)
			w.Write(jsonResp)
			w.Write([]byte{})
		}
	}
}
