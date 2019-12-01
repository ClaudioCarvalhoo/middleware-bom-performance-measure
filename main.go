package main

import (
	"bufio"
	"middleware-bom/model"
	"middleware-bom/publisher"
	"middleware-bom/subscriber"
	"os"
	"strconv"
	"strings"
	"time"
)

func main()  {
	println("1 for measurer, 2 for measuree")
	reader := bufio.NewReader(os.Stdin)
	choice, _ := reader.ReadString('\n')
	if strings.HasPrefix(choice, "1") {
		println("how many?")
		num, _ := reader.ReadString('\n')
		num = strings.ReplaceAll(num, "\n", "")
		num = strings.ReplaceAll(num, "\r", "")
		n, _ := strconv.ParseInt(num, 10, 64)
		p, _ := publisher.NewPublisher("a", "localhost:7447")
		s, _ := subscriber.NewSubscriber("b", "localhost:7447")
		c := s.Subscribe()
		rec := make(chan interface{}, 5 * n)
		go func() {
			for{
				msg, more := <- c
				if more {
					rec <- msg
				} else {
					break
				}
			}
		}()
		start := time.Now()
		go func() {
			for i:= int64(0); i<n; i++ {
				p.Publish(model.Content{Content:"oi"})
			}
		}()
		for j:= int64(0); j<n; j++ {
			_, more := <- rec
			if more {
				continue
			}
		}
		elapsed := time.Since(start)
		print("Elapsed time: ")
		print(elapsed.Nanoseconds() / 1000000)
		print("ms")
	}else {
		s, _ := subscriber.NewSubscriber("a", "localhost:7447")
		p, _ := publisher.NewPublisher("b", "localhost:7447")
		c := s.Subscribe()
		for {
			_, more := <- c
			if more {
				go p.Publish(model.Content{Content:"alo"})
			} else {
				break
			}
		}
	}
}
