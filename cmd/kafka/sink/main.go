package main

import (
	"bufio"
	"context"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/confluentinc/confluent-kafka-go/kafka"
	"github.com/lastexile/kepler"
	kkafka "github.com/lastexile/kepler/kafka"
)

func main() {
	log.SetFlags(log.LstdFlags | log.Lmicroseconds)

	log.Println("starting...")

	cfg := &kafka.ConfigMap{
		"metadata.broker.list": "localhost:9092",
		"group.id":             "odd",
		"default.topic.config": kafka.ConfigMap{"auto.offset.reset": "earliest"},
	}

	s, err := kkafka.NewSink("test", cfg, func(m kepler.Message) ([]byte, error) {
		return []byte(m.Value().(string)), nil
	})

	if err != nil {
		log.Fatalf("Unable to create kafkaspring: %v\n", err)
	}

	spring := kepler.NewSpring("odd", func(ctx context.Context, ch chan<- kepler.Message) {

		i := 1
		for {
			ch <- kepler.NewMessage("odd", strconv.Itoa(i))
			i++
			time.Sleep(1 * time.Second)
			log.Println(i)
		}
	})

	spring.LinkTo(s, kepler.Allways)

	reader := bufio.NewReader(os.Stdin)
	log.Print("Enter text: ")
	text, _ := reader.ReadString('\n')
	log.Println(text)
}
