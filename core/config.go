package core

// Config
type Config struct {
	ServerTcp *ServerTcpConfig
}

type ServerTcpConfig struct {
	Host       string
	Port       string
	TlsEnabled bool
}
