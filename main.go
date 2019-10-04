package main

import (
	"encoding/json"
	"fmt"
	"github.com/go-redis/redis"
	"github.com/spf13/viper"
	"goDocker/queue"
	"goDocker/sender"
	sms_reader "goDocker/sms-reader"
	"math/rand"
	"sync"
	"time"
)

func getConfig() *viper.Viper {
	var vi = viper.New()
	vi.SetConfigName("config")
	vi.AddConfigPath(".")
	err := vi.ReadInConfig()
	if err != nil {
		panic(fmt.Errorf("Error while reading config.json : %s \n", err))
	}
	return vi
}

func getRedisClient(config *viper.Viper) *redis.Client {
	client := redis.NewClient(&redis.Options{
		Addr:     config.GetString("redis.host") + ":" + config.GetString("redis.port"),
		Password: "", // no password set
		DB:       0,  // use default DB
	})
	return client
}
func main() {

	// fetching config
	config := getConfig()
	// initialising redis client
	redisClient := getRedisClient(config)
	// initialising queue
	q := queue.Init(redisClient)
	ok := sender.IsReady(q)
	if ok {
		sender.Init(q)
		sms_reader.Init(config, q)
	}
	ch := make(chan int)
	<-ch

}

type messageData struct {
	Name string `json:"name"`
	Age  int    `json:"age"`
}

func publish(wg *sync.WaitGroup, client *redis.Client) {
	fmt.Println("here")

	for i := 0; i < 20; i++ {
		ma := messageData{Name: "Kapil", Age: i}
		jsonString, _ := json.Marshal(ma)
		wg.Add(1)
		err := client.Publish("Students", jsonString).Err()
		if err != nil {
			fmt.Println(err)
		}
	}
}
func mainOld() {
	var wg sync.WaitGroup
	var client = ExampleNewClient()
	pong, err := client.Ping().Result()
	fmt.Println(pong, err)

	go publish(&wg, client)

	go func() {
		sub := client.Subscribe("Students")
		ch := sub.Channel()

		var mutex sync.Mutex
		var count int
		var limit = 5
		//for {
		for {
			//fmt.Println("inside consume")

			if count < limit {

				var data = <-ch
				//wg.Add(1)
				mutex.Lock()
				count++
				mutex.Unlock()
				go consumer(data, &count, &mutex, &wg)
			}

		}

	}()
	wg.Wait()

	var chan1 = make(chan int)
	<-chan1
}
func consumer(message *redis.Message, count *int, mutex *sync.Mutex, wg *sync.WaitGroup) {
	rand.Seed(time.Now().UnixNano())
	var msg messageData
	err := json.Unmarshal([]byte(message.Payload), &msg)
	time.Sleep(time.Millisecond * time.Duration(rand.Intn(5000)))
	fmt.Println(msg.Name, err, message.Payload)
	mutex.Lock()
	*count--
	mutex.Unlock()
	wg.Done()
}
func printMessage(msg *redis.Message) {
	fmt.Println(msg.Channel, msg.Payload)
}

func ExampleNewClient() *redis.Client {
	client := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "", // no password set
		DB:       0,  // use default DB
	})
	return client

	// Output: PONG <nil>
}
