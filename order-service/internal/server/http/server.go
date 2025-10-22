package httpserver

import (
	"context"
	"fmt"
	"net/http"

	// _ "github.com/ZaiiiRan/backend_labs/order-service/docs"
	"github.com/ZaiiiRan/backend_labs/order-service/internal/server/http/controllers"
	httpSwagger "github.com/swaggo/http-swagger"
)

type Server struct {
	srv             *http.Server
	orderController *controllers.OrderController
}

func NewServer(port int, orderController *controllers.OrderController) *Server {
	s := &Server{orderController: orderController}

	s.srv = &http.Server{
		Addr:    ":" + fmt.Sprint(port),
		Handler: registerRoutes(orderController),
	}

	return s
}

func (s *Server) Start() error {
	return s.srv.ListenAndServe()
}

func (s *Server) Stop(ctx context.Context) error {
	return s.srv.Shutdown(ctx)
}

func (s *Server) Addr() string {
	return s.srv.Addr
}

func registerRoutes(orderController *controllers.OrderController) http.Handler {
	mux := http.NewServeMux()

	mux.Handle("/swagger/", httpSwagger.WrapHandler)

	mux.HandleFunc("/api/v1/order/batch-create", orderController.BatchCreate)
	mux.HandleFunc("/api/v1/order/query", orderController.QueryOrders)
	mux.HandleFunc("/api/v1/audit-log/order/batch-create", orderController.AuditLogOrderBatchCreate)

	return mux
}
