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

	"github.com/go-playground/assert"
	"github.com/nats-io/stan.go"
)

func TestLoad(t *testing.T) {

	fakedata := fakedata.Generate(1000)

	sc, err := stan.Connect("hello", "load-test-client")
	if err != nil {
		log.Fatal("failed to conn:", err)
	}

	t.Log("[+] Start publish fake data")
	subj := "orders"
	for i := 0; i < len(fakedata); i++ {
		encoded, err := json.Marshal(fakedata[i])
		if err != nil {
			panic(err)
		}
		sc.Publish(subj, encoded)
	}

	time.Sleep(time.Second*10)

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
			assert.Equal(t, http.StatusOK, resp.StatusCode)
			var receivedOrder entity.Order

			json.NewDecoder(resp.Body).Decode(&receivedOrder)

			assert.Equal(t, fakedata[i].OrderUID, receivedOrder.OrderUID)

		}()
	}

	wg.Wait()

	elapsed := time.Since(start)
	log.Printf("All requests completed in %s\n", elapsed)
}


func TestPubGet(t *testing.T){
	fakedata := fakedata.Generate(10)
	sc, err := stan.Connect("hello", "load-test-client2")
	if err != nil {
		log.Fatal("failed to conn:", err)
	}


	subj := "orders"

	clientHTTP := &http.Client{}

	for i := 0; i < len(fakedata); i++ {
		encoded, err := json.Marshal(fakedata[i])
		if err != nil {
			panic(err)
		}
		sc.Publish(subj, encoded)

		time.Sleep(time.Millisecond*50)
		url := fmt.Sprintf("http://127.0.0.1:8888/order?id=%s", fakedata[i].OrderUID)
		resp, err := clientHTTP.Get(url)
		if err != nil {
			log.Printf("Error fetching order %s", err)
			return
		}
		defer resp.Body.Close()
		assert.Equal(t, http.StatusOK, resp.StatusCode)
		var receivedOrder entity.Order

		json.NewDecoder(resp.Body).Decode(&receivedOrder)

		assert.Equal(t, fakedata[i].OrderUID, receivedOrder.OrderUID)

	}


}