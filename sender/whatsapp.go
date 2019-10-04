package sender

import (
	"fmt"
	"goDocker/queue"
)
func isWhatsappReady(q *queue.Config) bool {
	_, err := q.SMS.Receive()
	if err != nil {
		return false
	}
	return true
}

func initWhatsapp(q *queue.Config)  {
	sub := q.Whatsapp.Subscribe()
	ch := sub.Channel()
	go func() {
		for {
			fmt.Println("sddata", <-ch)
		}
	}()
}
