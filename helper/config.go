package helper

type EnvConfig struct {
	BloomFilterPath           string
	WEBPassword               string
	WEBListen                 string
	Redis                     Redis
	MysqlDSN                  string
	WechatSendMessagePassword string
	BloomFilterClient         string
	BloomFilterServer         string
}

type Redis struct {
	Network     string
	Addr        string
	Password    string
	DB          int
	URLQueueKey string
}
