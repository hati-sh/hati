package core

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net"
	"net/http"
	"net/rpc"
	"time"

	"github.com/hati-sh/hati/common/logger"
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
	storageManager *storage.Manager
}

func NewRpcServer(ctx context.Context, storageManager *storage.Manager, config *RpcServerConfig) *RpcServer {

	return &RpcServer{
		ctx:            ctx,
		config:         config,
		storageManager: storageManager,
	}
}

func (s *RpcServer) Start() error {
	go func(sr *RpcServer) {
		type Storage struct {
			RpcStorageService
		}
		rpcStorage := new(Storage)
		rpcStorage.RpcStorageService.storageManager = sr.storageManager

		rpc.Register(rpcStorage)

		http.HandleFunc("/", func(w http.ResponseWriter, req *http.Request) {
			defer req.Body.Close()
			w.Header().Set("Content-Type", "application/json")
			w.Header().Set("Access-Control-Allow-Origin", "*")
			res := newRPCRequest(req.Body).Call()
			io.Copy(w, res)
		})

		logger.Debug("rpc server listening at: " + sr.config.Host + ":" + sr.config.Port)

		listener, err := net.Listen("tcp", sr.config.Host+":"+sr.config.Port)
		if err != nil {
			panic(err)
		}

		srv := http.Server{
			ReadHeaderTimeout: time.Second * 5,
			ReadTimeout:       time.Second * 10,
			Handler:           nil,
		}

		go func(l net.Listener) {
			if err := srv.Serve(l); err != nil {
				if !errors.Is(err, http.ErrServerClosed) {
					panic(err)
				} else {
					l.Close()
					logger.Debug("rpc server stopped")
				}
			}
		}(listener)

		select {
		case <-s.ctx.Done():
			{
				logger.Debug("rpc server shutting down")

				if err := srv.Shutdown(s.ctx); err != nil {
					logger.Error(fmt.Sprintf("rpc server graceful shutdown failed: %v\n", err))
				}
			}
		}
	}(s)

	return nil
}
