package helper

type EnvConfig struct {
	WorkspacePath        string
	BloomFilterFile      string
	LocalTransportEnable bool //using http.DefaultTransport
	LocalTransportWeight int
	SsServers            []SsServer
	Redis                Redis
}

type SsServer struct {
	Addr     string
	Password string
	Cipher   string
	Weight   int
}

type Redis struct {
	Addr        string
	DB          int
	URLQueueKey string
}
