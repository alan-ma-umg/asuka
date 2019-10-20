package helper

type EnvConfig struct {
	BloomFilterPath           string
	WEBPassword               string
	WEBListen                 string
	LocalTransport            bool
	Redis                     Redis
	MysqlDSN                  string
	WechatSendMessagePassword string
	SingleProject             string
}

type Redis struct {
	Network     string
	Addr        string
	Password    string
	DB          int
	URLQueueKey string
}
