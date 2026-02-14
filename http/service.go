package http

import (
	"database/sql"
	"log"
	"net/http"
	"test-lbc/http/handlers"
	"time"

	"github.com/gin-gonic/gin"
)

type Server struct {
	db       *sql.DB
	bindAddr string

	router *gin.Engine
}

func New(db *sql.DB, bindAddr string) *Server {
	return &Server{
		db:       db,
		bindAddr: bindAddr,
	}
}

func (s *Server) Start() error {
	log.Printf("start http server on port %s", s.bindAddr)
	gin.SetMode(gin.ReleaseMode)
	s.router = gin.New()
	s.loadRoutes()
	server := &http.Server{
		Addr:           ":8080",
		Handler:        s.router,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}
	server.ListenAndServe()
	return nil
}

func (s *Server) loadRoutes() {
	// load fizzBuzz routes
	fbGroup := s.router.Group("/fizzbuzz")
	fbGroup.Handle("POST", "/run", func(ctx *gin.Context) {
		handlers.FizzBuzzRun(ctx, s.db)
	})
	fbStatsGroup := fbGroup.Group("/stats")
	fbStatsGroup.Handle("GET", "/most-requested", func(ctx *gin.Context) {
		handlers.FizzBuzzStats(ctx, s.db)
	})
}
