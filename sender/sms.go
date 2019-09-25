package sender

import (
	"fmt"
	"github.com/go-redis/redis"
	"goDocker/queue"
	"math/rand"
	"strings"
	"sync"
	"sync/atomic"
	"time"
)

const maxRoutines = 100
const maxWaitCount = 12

var currentRoutines uint32
var waitCount int

type SMS struct {
	queue *queue.Config
}

func initSMS(q *queue.Config) bool {
	var smsWaitGroup sync.WaitGroup
	sub := q.SMS.Subscribe()
	_, err := sub.Receive()
	if err != nil {
		return false
	}
	smsWaitGroup.Add(1)
	go readSMSQueue(sub, &smsWaitGroup)
	smsWaitGroup.Wait()
	return true
}

func readSMSQueue(subscribedChannel *redis.PubSub, smsWaitGroup *sync.WaitGroup) {
	var wg sync.WaitGroup
	ch := subscribedChannel.Channel()
	for waitCount < maxWaitCount {
		if atomic.LoadUint32(&currentRoutines) < maxRoutines {
			select {
			case payload := <-ch:
				atomic.AddUint32(&currentRoutines, 1)
				wg.Add(1)
				payloadArr := strings.Split(payload.Payload, ",")
				phone := payloadArr[0]
				message := payloadArr[1]
				go sendSMS(phone, message, &wg)
				waitCount = 0
			default:
				waitCount++
				time.Sleep(time.Second * 2)
			}
		}
	}
	wg.Wait()
	smsWaitGroup.Done()
}

func sendSMS(phone, message string, wg *sync.WaitGroup) {
	rand.Seed(time.Now().UnixNano())
	atomic.AddUint32(&currentRoutines, 1)
	time.Sleep(time.Millisecond * time.Duration(rand.Intn(5000)))
	atomic.AddUint32(&currentRoutines, ^uint32(0))
	fmt.Println(message, " sent to ", phone, atomic.LoadUint32(&currentRoutines))
	wg.Done()
}
