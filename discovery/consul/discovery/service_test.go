package discovery

import (
	"fmt"
	"testing"
)

func TestSerivce(t *testing.T){
	fmt.Println(NewService().GetService("hello"))

	select {

	}
}
