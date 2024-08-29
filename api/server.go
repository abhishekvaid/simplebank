package api

import (
	"github.com/gin-gonic/gin"
	db "himavisoft.simple_bank/db/sqlc"
)

type Server struct {
	store  *db.Store
	router *gin.Engine
}

func NewServer(store *db.Store) (server *Server) {
	server = &Server{store: store}
	router := gin.Default()

	router.POST("/accounts", server.CreateAccount)
	router.GET("/accounts/:id", server.GeAccount)
	router.GET("/accounts", server.GeAccounts)

	server.router = router

	return
}

func (s *Server) Start(address string) error {
	return s.router.Run(address)
}

func errorResponse(err error) gin.H {
	return gin.H{
		"error": err.Error(),
	}

}
