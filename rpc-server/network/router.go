package network

import (
	"github.com/gin-gonic/gin"
	"rpc-server/config"
	"rpc-server/gRPC/client"
	"rpc-server/service"
)

type Network struct {
	cfg        *config.Config
	service    *service.Service
	gRPCClient *client.GRPCClient
	engin      *gin.Engine
}

func NewNetwork(cfg *config.Config, service *service.Service, gRPCClient *client.GRPCClient) (*Network, error) {
	n := &Network{
		cfg:        cfg,
		service:    service,
		gRPCClient: gRPCClient,
		engin:      gin.New(),
	}

	n.engin.POST("/login", n.login)
	n.engin.GET("/verify", n.verifyLogin(), n.verify)

	return n, nil
}

func (n *Network) StartServer() {
	n.engin.Run(":8080")
}
