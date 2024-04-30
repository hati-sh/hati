package core

// Config structure for hati
type Config struct {
	// DataDir absolute path to hati data directory
	DataDir   string
	ServerTcp *TcpServerConfig
	ServerRpc *RpcServerConfig
}
