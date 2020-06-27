package env

import "flag"

/**
@Desc:包含consul的相关配置
@Param: ConsulAddr:consul 注册地址
*/

var (
	//config
	ConsulAddr string
	KVPrefix   string

	//register
	LocalAddr  string
	ServerPort int
	ServerName string

	//discovery


)

func init() {

	//config
	flag.StringVar(&ConsulAddr, "consuladdr", "dev.codenai.com:8500", "consul 注册地址")
	flag.StringVar(&KVPrefix, "kvprefix", "foo", "配置首项")

	//register
	flag.StringVar(&LocalAddr, "localaddr", "127.0.0.1", "本地服务地址")
	flag.IntVar(&ServerPort, "serverport", 9050, "本地服务地址")
	flag.StringVar(&ServerName, "servername", "hello", "服务名称")

	//discovery


}
