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

const maxRoutines = 500
const maxWaitCount = 12

var currentRoutines uint32
var waitCount int

type SMS struct {
	queue *queue.Config

}

func isSMSReady(q *queue.Config) bool {
	_, err := q.SMS.Receive()
	if err != nil {
		return false
	}
	return true
}

func initSMS(q *queue.Config) {
	fmt.Println("here")
	//var smsWaitGroup sync.WaitGroup
	sub := q.SMS.Subscribe()
	//smsWaitGroup.Add(1)
	go readSMSQueue(sub)
	//smsWaitGroup.Wait()
}

func readSMSQueue(subscribedChannel *redis.PubSub) {
	var wg sync.WaitGroup
	ch := subscribedChannel.Channel()
	for waitCount < maxWaitCount {
		fmt.Println("for")
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
				fmt.Println("default")
				waitCount++
				time.Sleep(time.Second * 2)
			}
		} else {
			fmt.Println("else")
			waitCount++
			time.Sleep(time.Second * 2)
		}
	}
	fmt.Println("Done with sending sms")
	wg.Wait()
	//smsWaitGroup.Done()
}

func sendSMS(phone, message string, wg *sync.WaitGroup) {
	rand.Seed(time.Now().UnixNano())
	//time.Sleep(time.Millisecond * time.Duration(rand.Intn(2000)))
	atomic.AddUint32(&currentRoutines, ^uint32(0))
	fmt.Println(message, " sent to ", phone, atomic.LoadUint32(&currentRoutines))
	wg.Done()
}
