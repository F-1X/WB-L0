package tests

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"sync"
	"testing"
	"time"
	"wb/backend/internal/fakedata"
	"wb/backend/internal/fakedata/fakemodel"

	"github.com/nats-io/stan.go"
)

type OrderResponse struct {
	Message string `json:"message"`
	Error   string `json:"error"`
}

func TestLoad(t *testing.T) {

	fakedata := fakedata.Generate(1000)

	// nats.PingInterval(20 * time.Second), nats.MaxPingsOutstanding(5)

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
	time.Sleep(time.Millisecond * 50)


	clientHTTP := &http.Client{}
	var wg sync.WaitGroup
	start := time.Now()
	for i := 0; i < 1000; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()

			url := fmt.Sprintf("http://127.0.0.1:8888/order?id=%s", fakedata[i].OrderUID)
			resp, err := clientHTTP.Post(url, "Content-type: application/json", nil)
			if err != nil {
				log.Printf("Error fetching order %s", err)
				return
			}
			defer resp.Body.Close()

			var receivedOrder fakemodel.OrderFake
			if err := json.NewDecoder(resp.Body).Decode(&receivedOrder); err != nil {
				log.Printf("Error decoding response body for order %s: %v", fakedata[i].OrderUID, err)
				return
			}

			if receivedOrder.OrderUID != fakedata[i].OrderUID {
				log.Printf("Order UID mismatch for order %s: expected %s, got %s", fakedata[i].OrderUID, fakedata[i].OrderUID, receivedOrder.OrderUID)
			}

		}()

	}

	wg.Wait()

	elapsed := time.Since(start)
	log.Printf("All requests completed in %s\n", elapsed)
}
