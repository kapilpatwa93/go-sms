package sms_reader

import (
	"bufio"
	"fmt"
	"github.com/spf13/viper"
	"goDocker/queue"
	"log"
	"os"
	"path/filepath"
)

func Init(config *viper.Viper, q *queue.Config) {
	path, _ := filepath.Abs("files/")
	f, err := os.Open(path + "/" + config.GetString("file.sms"))
	if err != nil {
		panic(fmt.Errorf("Error while reading config.json : %s \n", err))
	}
	defer func() {
		fmt.Println("Completed reading file")
		if err = f.Close(); err != nil {
			log.Fatal(err)
		}
	}()
	r := bufio.NewScanner(f)
	for r.Scan() {
		err := q.SMS.Publish(r.Text())
		if err != nil {
			fmt.Printf("Error while pushing : %s", r.Text())
		}
	}
	err = q.SMS.Close()
	fmt.Println("closed err :", err)

}
