package server

import (
	"context"
	"errors"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"log"
	"net"
	"rpc-server/config"
	"rpc-server/gRPC/paseto"
	auth "rpc-server/gRPC/proto"
	"time"
)

type GRPCServer struct {
	auth.UnimplementedAuthServiceServer
	pasetoMaker    *paseto.PasetoMaker
	tokenVerifyMap map[string]*auth.AuthData
}

func NewGRPCServer(cfg *config.Config) error {
	if lis, err := net.Listen("tcp", cfg.GRPC.URL); err != nil {
		return err
	} else {
		server := grpc.NewServer([]grpc.ServerOption{}...)

		auth.RegisterAuthServiceServer(server, &GRPCServer{
			pasetoMaker:    paseto.NewPasetoMaker(cfg),
			tokenVerifyMap: make(map[string]*auth.AuthData),
		})

		reflection.Register(server)

		go func() {
			log.Println("Start GRPC Server")
			if err = server.Serve(lis); err != nil {
				panic(err)
			}
		}()

		return nil
	}
}

func (g *GRPCServer) CreateAuth(_ context.Context, req *auth.CreateTokenReq) (*auth.CreateTokenRes, error) {
	data := req.Auth
	token := data.Token

	g.tokenVerifyMap[token] = data

	return &auth.CreateTokenRes{Auth: data}, nil
}

func (g *GRPCServer) VerifyAuth(_ context.Context, req *auth.VerifyTokenReq) (*auth.VerifyTokenRes, error) {
	token := req.Token

	res := &auth.VerifyTokenRes{
		V: &auth.Verify{
			Auth: nil,
		}}

	if authData, ok := g.tokenVerifyMap[token]; !ok {
		res.V.Status = auth.ResponseType_FAILED
		return res, errors.New("Not Existed At TokenVerifyMap")
	} else if err := g.pasetoMaker.VerifyToken(token); err != nil {
		return nil, errors.New("Failed Verify Token")
	} else if authData.ExpireDate < time.Now().Unix() {
		delete(g.tokenVerifyMap, token)
		res.V.Status = auth.ResponseType_EXPIRED_DATE
		return res, errors.New("Expired Time Over")
	} else {
		res.V.Status = auth.ResponseType_SUCCESS
		return res, nil
	}
}
