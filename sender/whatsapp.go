package sender

import (
	"fmt"
	"goDocker/queue"
)

func initWhatsapp(q *queue.Config) bool {
	sub := q.Whatsapp.Subscribe()

	_, err := sub.Receive()
	if err != nil {
		return false
	}
	ch := sub.Channel()

	go func() {
		for {
			fmt.Println("sddata", <-ch)
		}
	}()
	return true

}
