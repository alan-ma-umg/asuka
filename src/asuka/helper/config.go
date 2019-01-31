package helper

type EnvConfig struct {
	TemplatePath    string
	BloomFilterPath string
	WEBPassword     string
	LocalTransport  LocalTransport //using http.DefaultTransport
	SsServers       []*SsServer
	Redis           Redis
	MysqlDSN        string
}

type LocalTransport struct {
	Enable   bool
	Interval float64
	Name     string
	Group    string
}

type SsServer struct {
	Enable     bool
	EnablePing bool
	Interval   float64
	Name       string
	Group      string
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
	Server      string
	DB          int
	URLQueueKey string
}
