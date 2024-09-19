package api

import (
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
	db "himavisoft.simple_bank/db/sqlc"
)

type Server struct {
	store  db.Store
	router *gin.Engine
}

func NewServer(store db.Store) (server *Server) {

	server = &Server{store: store}
	router := gin.Default()

	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		v.RegisterValidation("currency", validCurrency)
	}

	// routes for accounts
	router.POST("/accounts", server.CreateAccount)
	router.GET("/accounts/:id", server.GeAccount)
	router.GET("/accounts", server.GeAccounts)

	// routes for transfers
	router.POST("/transfers", server.TransferAmount)

	// routes for users
	router.POST("/users", server.CreateUser)
	router.GET("/users/:username", server.GetUser)

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
