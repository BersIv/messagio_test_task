package router

import (
	"messagio_test_task/internal/message"
	"net/http"

	_ "messagio_test_task/docs"

	"github.com/gorilla/mux"
	httpSwagger "github.com/swaggo/http-swagger"
)

type Option func(router *mux.Router)

// InitRouter инициализирует маршрутизатор с предоставленными опциями
func InitRouter(options ...Option) *mux.Router {
	r := mux.NewRouter()

	for _, option := range options {
		option(r)
	}
	r.PathPrefix("/swagger/").Handler(httpSwagger.WrapHandler)

	return r
}

func Start(addr string, r *mux.Router) error {
	return http.ListenAndServe(addr, r)
}

// MessageRouter добавляет маршруты для Message обработчиков
func MessageRouter(messageHandler *message.MessageHandler) Option {
	return func(r *mux.Router) {
		r.HandleFunc("/messages/create", messageHandler.CreateMessage).Methods("POST")
	}
}
