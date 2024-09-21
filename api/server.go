package api

import (
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
	db "himavisoft.simple_bank/db/sqlc"
	"himavisoft.simple_bank/token"
	"himavisoft.simple_bank/util"
)

type Server struct {
	config     util.Config
	tokenMaker token.Maker
	store      db.Store
	router     *gin.Engine
}

func NewServer(config util.Config, store db.Store) (server *Server, err error) {

	tokenMaker, err := token.NewPaseto(config.TokenSecret)
	if err != nil {
		return
	}

	server = &Server{store: store, tokenMaker: tokenMaker, config: config}
	server.setupRoutes()

	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		v.RegisterValidation("currency", validCurrency)
	}

	return
}

func (server *Server) setupRoutes() {

	router := gin.Default()

	// routes for users (they don't use auth)
	router.POST("/users", server.CreateUser)
	router.POST("/users/login", server.CreateUser)

	authRoutes := router.Group("/").Use(authMiddleware(server.tokenMaker))

	// routes for User
	authRoutes.GET("/users", server.GetUser)

	// routes for accounts
	authRoutes.POST("/accounts", server.CreateAccount)
	authRoutes.GET("/accounts/:id", server.GetAccount)
	authRoutes.GET("/accounts", server.ListAccounts)

	// routes for transfers
	authRoutes.POST("/transfers", server.TransferAmount)

	server.router = router

}

func (s *Server) Start(address string) error {
	return s.router.Run(address)
}

func errorResponse(err error) gin.H {
	return gin.H{
		"error": err.Error(),
	}

}
