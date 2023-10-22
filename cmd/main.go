package main

import (
	"encoding/json"
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/nats-io/stan.go"
	"github.com/sema0205/WB_L0/internal/config"
	"github.com/sema0205/WB_L0/internal/http-server/handlers/orderId"
	"github.com/sema0205/WB_L0/internal/http-server/handlers/publish"
	"github.com/sema0205/WB_L0/internal/http-server/handlers/retrieve"
	"github.com/sema0205/WB_L0/internal/models"
	"github.com/sema0205/WB_L0/internal/storage/postgres"
	"log"
	"net/http"
)

func main() {

	cfg := config.MustLoad()

	storage, err := postgres.NewStorage(fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		cfg.PostgresAuth.Host, cfg.PostgresAuth.Port, cfg.PostgresAuth.User, cfg.PostgresAuth.Password, cfg.PostgresAuth.DbName))

	cache := storage.RestoreCache()

	if err != nil {
		log.Fatal(err)
	}

	sc, _ := stan.Connect(cfg.NatsAuth.StanClusterId, cfg.NatsAuth.ClientId)

	sc.Subscribe(cfg.NatsAuth.ChannelName, func(msg *stan.Msg) {
		order := models.OrderInfo{}
		err = json.Unmarshal(msg.Data, &order)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println(order)
		err = storage.SaveOrder(order)
		cache[order.OrderUID] = order
		if err != nil {
			log.Fatal(err)
		}
	})

	router := chi.NewRouter()
	router.Use(middleware.Logger)
	router.Get("/", retrieve.New(storage))
	router.Get("/id/{uuid}", orderId.New(cache))
	router.Get("/publish", publish.New(sc, cfg.NatsAuth.ChannelName))

	http.ListenAndServe(":3333", router)
}
