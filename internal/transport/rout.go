package transport

import (
	"net"

	"github.com/DblMOKRQ/auth-service/internal/service"
	auth "github.com/DblMOKRQ/auth-service/pkg/api"
	"google.golang.org/grpc"
)

type Router struct {
	server  *grpc.Server
	service *service.Service
}

func NewRouter(server *grpc.Server, service *service.Service) *Router {
	return &Router{server: server, service: service}
}
func (r *Router) Run(addr string) error {
	listener, err := net.Listen("tcp", addr)
	if err != nil {
		return err
	}
	auth.RegisterAuthServer(r.server, r.service)
	return r.server.Serve(listener)
}
