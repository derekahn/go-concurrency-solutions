//////////////////////////////////////////////////////////////////////
//
// Given is a producer-consumer szenario, where a producer reads in
// tweets from a mockstream and a consumer is processing the
// data. Your task is to change the code so that the producer as well
// as the consumer can run concurrently
//

package main

import (
	"fmt"
	"sync"
	"time"
)

func producer(out chan<- Tweet, stream Stream) {
	for {
		tweet, err := stream.Next()
		if err == ErrEOF {
			close(out)
			return
		}
		out <- *tweet
	}
}

func consumer(tweets <-chan Tweet, wg *sync.WaitGroup) {
	for t := range tweets {
		if t.IsTalkingAboutGo() {
			fmt.Println(t.Username, "\ttweets about golang")
		} else {
			fmt.Println(t.Username, "\tdoes not tweet about golang")
		}
		wg.Done()
	}
}

func main() {
	start := time.Now()
	stream := GetMockStream()
	amount := len(stream.tweets)

	var wg sync.WaitGroup
	wg.Add(amount)
	ch := make(chan Tweet, amount)

	go producer(ch, stream)
	go consumer(ch, &wg)
	wg.Wait()

	fmt.Printf("Process took %s\n", time.Since(start))
}
