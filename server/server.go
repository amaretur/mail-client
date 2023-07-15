package server 

import (
	"log"
	"time"
	"net/http"
	"context"
)

// Кастомная структура сервера
type HttpServer struct {
	httpServer *http.Server
}

// Запуск сервера
func (s *HttpServer) Run(port string, handler http.Handler) {
	s.httpServer = &http.Server { 
		Addr:			":" + port,			// Порт
		Handler: 		handler,			// Маршрутизация
		MaxHeaderBytes: 1 << 20, 			// 1 МБ
		ReadTimeout: 	10 * time.Second,	// Таймаут для чтения
		WriteTimeout: 	10 * time.Second,	// Таймаут для записи
	}

	log.Printf(
		"Server is running on %s port... (quit the server with <Ctrl-C>)\n", 
		port,
	)

	err := s.httpServer.ListenAndServe()
	if err != nil && err != http.ErrServerClosed {
		log.Printf("Server terminated with an error: %v\n", err)
	}
}

func (s *HttpServer) Shutdown(ctx context.Context) error {
	return s.httpServer.Shutdown(ctx)
}

