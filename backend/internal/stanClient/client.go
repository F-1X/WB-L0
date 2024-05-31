package stanClient

import (
	"context"
	"log"

	"github.com/nats-io/nats.go"
	"github.com/nats-io/stan.go"
)

type StanClient struct {
	sc stan.Conn
}

func New(ctx context.Context, url string, natsOpts []nats.Option, stanClusterID string, clientID string, stanOpts []stan.Option) *StanClient {

	nc, err := nats.Connect(url, natsOpts...)
	if err != nil {
		log.Fatal(err)
	}

	stanOpts = append(stanOpts, stan.NatsConn(nc), stan.SetConnectionLostHandler(func(_ stan.Conn, reason error) {
		log.Fatalf("Connection lost, reason: %v", reason)
	}))

	sc, err := stan.Connect(stanClusterID, clientID, stanOpts...)

	if err != nil {
		log.Fatalf("Can't connect: %v.\nMake sure a NATS Streaming Server is running at: %s", err, url)
	}

	log.Println("[+] Succussfully connected to stan server")

	go func(){
		<-ctx.Done()
		log.Println("[!] stan gracefully shutdown")
		sc.Close()
	}()
	return &StanClient{sc: sc}
}

func (sc *StanClient) Publish(subj string, msg []byte) error {
	err := sc.sc.Publish(subj, msg)
	if err != nil {
		log.Fatalf("Error during publish: %v\n", err)
		return err
	}
	return nil
}

func (sc *StanClient) Close() error {
	return sc.sc.Close()
}

func (sc *StanClient) NatsConn() *nats.Conn {
	return sc.sc.NatsConn()
}

func (sc *StanClient) PublishAsync(subject string, data []byte, ah stan.AckHandler) (string, error) {
	return sc.sc.PublishAsync(subject, data, ah)
}

func (sc *StanClient) QueueSubscribe(subject string, qgroup string, cb stan.MsgHandler, opts ...stan.SubscriptionOption) (stan.Subscription, error) {
	return sc.sc.QueueSubscribe(subject, qgroup, cb, opts...)
}

func (sc *StanClient) Subscribe(subject string, cb stan.MsgHandler, opts ...stan.SubscriptionOption) (stan.Subscription, error) {
	return sc.sc.Subscribe(subject, cb, opts...)
}
