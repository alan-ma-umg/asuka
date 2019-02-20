package helper

type EnvConfig struct {
	BloomFilterPath string
	WEBPassword     string
	LocalTransport  LocalTransport //using http.DefaultTransport
	Redis           Redis
	MysqlDSN        string
}

type LocalTransport struct {
	Enable bool
}

type Redis struct {
	Network     string
	Addr        string
	Password    string
	DB          int
	URLQueueKey string
}
