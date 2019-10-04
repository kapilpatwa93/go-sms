package sender

import "goDocker/queue"

func IsReady(q *queue.Config) bool {
	return isSMSReady(q) && isWhatsappReady(q)
}

func Init(q *queue.Config) {
	initSMS(q)
	initWhatsapp(q)
}
