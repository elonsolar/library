package discovery

import (
	"fmt"
	"github.com/fvbock/endless"
	"github.com/gin-gonic/gin"
	"io/ioutil"
	"testing"
)

func TestSerivce(t *testing.T){



	r := gin.New()
	r.POST("/xx", func(c *gin.Context) {
		bts, _ := ioutil.ReadAll(c.Request.Body)
		fmt.Println(string(bts))
		c.String(200, "x")
	})
	s := endless.NewServer(":8128", r)
	 s.ListenAndServe()
}
