package config

type Config struct {
	HttpPort   string `json:"httpPort"`
	GinMode    string `json:"ginMode"`
	RPCName    string `json:"rpcName"`
	RPCAddress string `json:"rpcAddress"`
	RPCUsage   string `json:"rpcUsage"`
	AppID      string `json:"app_id"`
	AppSecret  string `json:"app_secret"`
}
