package core

import (
	"context"
	"io"
	"net/http"
	"net/rpc"

	"github.com/hati-sh/hati/storage"
)

type RpcServerConfig struct {
	Host    string
	Port    string
	Enabled bool
}

type RpcServer struct {
	ctx            context.Context
	config         *RpcServerConfig
	storageManager *storage.StorageManager
}

func NewRpcServer(ctx context.Context, storageManager *storage.StorageManager, config *RpcServerConfig) *RpcServer {

	return &RpcServer{
		ctx:            ctx,
		config:         config,
		storageManager: storageManager,
	}
}

func (s *RpcServer) Start() error {
	go func(s *RpcServer) {

		type Storage struct {
			RpcStorageService
		}

		rpcStorage := new(Storage)
		rpcStorage.RpcStorageService.storageManager = s.storageManager

		rpc.Register(rpcStorage)

		http.HandleFunc("/", func(w http.ResponseWriter, req *http.Request) {
			defer req.Body.Close()
			w.Header().Set("Content-Type", "application/json")
			w.Header().Set("Access-Control-Allow-Origin", "*")
			res := newRPCRequest(req.Body).Call()
			io.Copy(w, res)
		})

		if err := http.ListenAndServe(s.config.Host+":"+s.config.Port, nil); err != nil {
			panic(err)
		}
	}(s)

	return nil
}
