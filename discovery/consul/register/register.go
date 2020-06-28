package register

import (
	"fmt"
	"github.com/cxqlkk/library/discovery/consul/config"
	"github.com/cxqlkk/library/discovery/consul/env"
	"github.com/fvbock/endless"
	"github.com/gin-gonic/gin"
	consulapi "github.com/hashicorp/consul/api"
	"strconv"
	"time"
)

type register struct {
	registration      *consulapi.AgentServiceRegistration
	agentServiceCheck *consulapi.AgentServiceCheck
	consuleClient     *consulapi.Client
}

func Register(engine *gin.Engine) {
	register := &register{
		registration:      new(consulapi.AgentServiceRegistration),
		agentServiceCheck: new(consulapi.AgentServiceCheck),
		consuleClient:     config.GetConsulClient(),
	}
	register.register()
	register.check()
	register.consuleClient.Agent().ServiceRegister(register.registration)
	engine.GET("/check", func(ctx *gin.Context) {
		ctx.JSON(200, map[string]string{"status": "ok"})
	})
	s := endless.NewServer(":"+strconv.Itoa(env.ServerPort), engine)
	err := s.ListenAndServe()
	if err != nil {
		fmt.Println(err)
	}
}

func (r *register) register() {
	r.registration.ID = env.ServerName + time.Now().Format("150405")
	r.registration.Name = env.ServerName
	r.registration.Port = env.ServerPort
	//registration.Tags = []string{"tags"}
	r.registration.Address = env.LocalAddr
}

func (r *register) check() {
	r.agentServiceCheck.HTTP = fmt.Sprintf("http://%s:%d%s", r.registration.Address, r.registration.Port, "/check")
	//设置超时 5s。
	r.agentServiceCheck.Timeout = "5s"
	//设置间隔 5s。
	r.agentServiceCheck.Interval = "5s"
	//注册check服务。
	r.registration.Check = r.agentServiceCheck

}
