package sender

import "goDocker/queue"

func Init(q *queue.Config) bool {
	return initSMS(q) && initWhatsapp(q)

}
