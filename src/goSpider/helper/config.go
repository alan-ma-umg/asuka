package helper

type EnvConfig struct {
	TemplatePath    string
	BloomFilterFile string
	LocalTransport  LocalTransport //using http.DefaultTransport
	SsServers       []SsServer
	Redis           Redis
}

type LocalTransport struct {
	Enable bool
	Name   string
}

type SsServer struct {
	Enable   bool
	Name     string
	Addr     string
	Password string
	Cipher   string
}

type Redis struct {
	Addr        string
	DB          int
	URLQueueKey string
}
