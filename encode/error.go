package encode

var (
	OK = add(0) // 正确

	AuthorizeFailed = add(1001) //无权限
	AuthorizeExpire = add(1002) //授权过期
	InvalidParam    = add(1003) //参数校验失败
	ServerErr     = add(1004) //系统错误

	//代理

	ServiceNotFound = add(1005) // 代理服务未找到
	InterfaceError  = add(1006) //第三方接口调用失败
)
