package main

import (
	"encoding/json"
	"fmt"
	"log"
	"time"
	"wb/backend/internal/fakedata"

	"github.com/nats-io/stan.go"
)

func main() {

	fakedata := fakedata.Generate(1000)
	var fakedataMarshal [][]byte

	for i := 0; i < len(fakedata); i++ {
		encoded, err := json.Marshal(fakedata[i])
		if err != nil {
			panic(err)
		}
		fakedataMarshal = append(fakedataMarshal, encoded)
	}

	// asyncPublish(fakedataMarshal) // 17.761125ms  1000 seed
	syncPublish(fakedataMarshal) // 722.940512ms 1000 seed
}

func asyncPublish(data [][]byte) {

	sc, err := stan.Connect("hello", "client-flooder")
	if err != nil {
		panic(err)
	}


	ackHandler := func(ackedNuid string, err error) {
		if err != nil {
			log.Printf("Warning: error publishing msg id %s: %v\n", ackedNuid, err.Error())
		} else {
			log.Printf("Received ack for msg id %s\n", ackedNuid)
		}
	}

	iterator := 0
	timestart := time.Now()
	for i := 0; i < 1000; i++ {
		sc.PublishAsync("orders", data[i], ackHandler)
		log.Println("published:", iterator)
		iterator++
		// time.Sleep(time.Millisecond*1000)
	}

	fmt.Println(time.Since(timestart))
}

func syncPublish(data [][]byte) {

	sc, err := stan.Connect("hello", "client-flooder")
	if err != nil {
		panic(err)
	}
	log.Println("Start publish fake data")

	subj := "orders"
	iterator := 0

	timestart := time.Now()
	for i := 0; i < len(data); i++ {
		sc.Publish(subj, data[i])
		log.Println("published:", iterator)
		iterator++
		time.Sleep(time.Millisecond*1000)
	}
	fmt.Println(time.Since(timestart))
}
