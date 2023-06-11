package api

import (
	db "github.com/Richd0tcom/SafeX-Pay/db/sqlc"
	"github.com/gin-gonic/gin"
)

//serves all http request in the banking application
type Server struct {
	store db.Store
	serverRouter *gin.Engine
}

//Creates a new server and sets up routing to handle request
func NewServer(store db.Store) *Server {
	server:= &Server{store: store}
	router := gin.Default() //router

	router.POST("/account/create", server.createAccount)
	router.GET("/account/list", server.listAccounts)
	router.GET("/account/:id", server.getAccount)


	server.serverRouter = router
	return server
}

//Starts the created sever
func (server *Server) Start (address string) error {
	return server.serverRouter.Run(address)
}

func buildErrorResponse(err error) gin.H {
	return gin.H{"error": err.Error()}
}