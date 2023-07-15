package app

import (
	"os"
	"os/signal"
	"log"
	"context"

	"github.com/amaretur/mail-client/server"
)

type App struct {
	repo		*Repository
	service		*Service
	usecase		*Usecase
	transport	*Transport
}

func New() *App {
	return &App{}
}

func (a *App) Init() {

	db, err := NewPostgresDB(PgConfig{
		Host:		"127.0.0.1",
		Port:		"5432",
		DBname:		"mailclient",
		Username:	"mailadmin",
		Password:	"admin",
		SSLmode:	"disable",
	})
	if err != nil {
		log.Fatalf("error initializing database: %s", err.Error())
	}

	a.repo = NewReposotory(db)

	a.service = NewService(&ServiceCongif{
		Repo: a.repo,
		SecretKey: "",
		AccessExpires: 60*24*7,		// 15 min
		RefreshExpires: 60*24*3,	// 3 days
	})

	a.usecase = NewUsecase(a.service)

	a.transport = NewTransport(a.usecase, a.service.Auth)
}

func (a* App) Run() {

	s := new(server.HttpServer)

	go s.Run("8080", a.transport.Http.Router())

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	<-quit

	if err := s.Shutdown(context.Background()); err != nil {
		log.Printf("HTTP server Shutdown: %v", err)
	}

	log.Printf("Server has successfully shut down!")
}
