package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"time"
)

func main() {
	start := time.Now()
	f, err := os.Open("file1.csv")
	if err != nil {
		fmt.Println("Error", err)
	}
	defer func() {
		if err = f.Close(); err != nil {
			log.Fatal(err)
		}
	}()
	r := bufio.NewScanner(f)
	for r.Scan() {
		fmt.Println(r.Text())
	}
	fmt.Println("here")
	elapsed := time.Since(start)
	fmt.Printf("Binomial took %s", elapsed)

}
