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
	Enable     bool
	Name       string
	Server     string
	ServerPort string
	Password   string
	Method     string

	//ssr only
	Obfs          string
	ObfsParam     string
	ProtocolParam string
	Protocol      string
}

type Redis struct {
	Server        string
	DB          int
	URLQueueKey string
}