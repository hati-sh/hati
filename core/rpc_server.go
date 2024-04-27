package core

import (
	"context"
	"io"
	"net/http"
	"net/rpc"
)

type Hello struct{}
type Args struct{}

func (t *Hello) Hi(args *Args, reply *string) error {
	// Fill reply pointer to send the data back
	*reply = "hi rpc"

	return nil
}

type RpcServerConfig struct {
	Host    string
	Port    string
	Enabled bool
}

type RpcServer struct {
	ctx            context.Context
	config         *RpcServerConfig
	payloadHandler PayloadHandler
}

func NewRpcServer(ctx context.Context, config *RpcServerConfig, payloadHandler PayloadHandler) *RpcServer {

	return &RpcServer{
		ctx:            ctx,
		config:         config,
		payloadHandler: payloadHandler,
	}
}

func (s *RpcServer) Start() error {
	go func(s *RpcServer) {
		hello := new(Hello)
		rpc.Register(hello)

		http.HandleFunc("/", func(w http.ResponseWriter, req *http.Request) {
			defer req.Body.Close()
			w.Header().Set("Content-Type", "application/json")
			res := newRPCRequest(req.Body).Call()
			io.Copy(w, res)
		})

		if err := http.ListenAndServe(s.config.Host+":"+s.config.Port, nil); err != nil {
			panic(err)
		}
	}(s)

	return nil
}
