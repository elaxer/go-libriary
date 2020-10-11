package main

import (
	"fmt"
	"time"
)

func parser(site string, c chan string) {
	fmt.Println("Parsing", site)
	c <- site
}

func logger(c chan string) {
	for {
		site := <-c
		fmt.Println("Parsed", site)
	}
}

func main() {
	c := make(chan string)

	sites := []string{"google.com", "pornhub.com", "vk.com"}

	go logger(c)
	for _, s := range sites {
		go parser(s, c)
	}

	<-time.After(time.Second * 1)
}
