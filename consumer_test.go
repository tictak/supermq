package supermq

import (
	"fmt"
	"testing"
)

type MyHandler struct {
	num int
}

func (h *MyHandler) HandleMessage(message *Message) error {
	fmt.Printf("[M]%s: %s\n", message.NSQDAddress, string(message.Body))
	return nil
}

func TestTconsumer(t *testing.T) {
	h := MyHandler{}
	r, _ := NewTconsumer("dlresult", "ch1", h.HandleMessage)
	fmt.Println(r.ConnectToNSQLookupds(
		[]string{"192.168.29.78:4161", "192.168.29.162:4161"}))
	fmt.Println(r.IsStarved())
	stop := make(chan int)
	<-stop
}
