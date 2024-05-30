package tests

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"sync"
	"testing"
	"time"
	"wb/backend/internal/domain/entity"
	"wb/backend/internal/fakedata"

	"github.com/nats-io/stan.go"
)

func TestLoad(t *testing.T) {

	fakedata := fakedata.Generate(1000)

	sc, err := stan.Connect("hello", "load-test-client")
	if err != nil {
		log.Fatal("failed to conn:", err)
	}

	log.Println("[+] Start publish fake data")
	subj := "orders"
	for i := 0; i < len(fakedata); i++ {
		encoded, err := json.Marshal(fakedata[i])
		if err != nil {
			panic(err)
		}
		sc.Publish(subj, encoded)
		log.Println("published:", fakedata[i].OrderUID)
	}

	// задержка для записи в базу последнего заказа
	time.Sleep(time.Millisecond * 500)

	clientHTTP := &http.Client{}
	var wg sync.WaitGroup
	start := time.Now()
	for i := 0; i < len(fakedata); i++ {

		wg.Add(1)
		go func() {
			defer wg.Done()

			url := fmt.Sprintf("http://127.0.0.1:8888/order?id=%s", fakedata[i].OrderUID)
			resp, err := clientHTTP.Get(url)
			if err != nil {
				log.Printf("Error fetching order %s", err)
				return
			}
			defer resp.Body.Close()

			var receivedOrder entity.Order

			json.NewDecoder(resp.Body).Decode(&receivedOrder)
			if receivedOrder.OrderUID != fakedata[i].OrderUID {
				log.Printf("Order UID mismatch for order %s: expected %s, got %s\n", fakedata[i].OrderUID, fakedata[i].OrderUID, receivedOrder.OrderUID)
			}
		}()
		// time.Sleep(time.Second * 3)
	}

	wg.Wait()

	elapsed := time.Since(start)
	log.Printf("All requests completed in %s\n", elapsed)
}
