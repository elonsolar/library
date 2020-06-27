package config

import (
	"flag"
	"fmt"
	"github.com/fvbock/endless"
	"github.com/gin-gonic/gin"
	"io/ioutil"
	"testing"
)

type Db struct {
	User string `dc:"user"`
}

type Conf struct {
	Cnd *string `dc:"cnd"`
	Dns string  `dc:"dns"`
	Url string  `dc:"url"`
	Db  *Db     `dc:"db"`
}

func TestConfig(t *testing.T) {

	flag.Parse()
	var c = Conf{}
	InitConfig(&c)
	fmt.Println(c)

	select {}

}

func TestGin(t *testing.T) {
	flag.Parse()
	var c = Conf{}
	InitConfig(&c)
	fmt.Println("开始：：：", c.Db.User)
	r := gin.New()
	r.POST("/xx", func(c *gin.Context) {
		bts, _ := ioutil.ReadAll(c.Request.Body)
		fmt.Println(string(bts))
		c.String(200, "x")
	})
	s := endless.NewServer(":8128", r)
	err := s.ListenAndServe()
	if err != nil {
		fmt.Println(err)
	}
}
