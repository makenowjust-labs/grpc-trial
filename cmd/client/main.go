package main

import (
	"flag"
	"fmt"
	"github.com/MakeNowJust-Labo/grpc-trial/echo"
	"github.com/pkg/errors"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"log"
	"sync"
	"time"
)

func main() {
	var addr string
	flag.StringVar(&addr, "addr", ":4567", "server address")
	flag.Parse()

	conn, err := grpc.Dial(addr, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("failed to dial: %v", err)
	}
	defer conn.Close()

	client := echo.NewEchoClient(conn)
	errCh := make(chan error, 10)
	var wg sync.WaitGroup

	ctx, cancel := context.WithCancel(context.Background())
	connect := func(id int) {
		wg.Add(1)
		go func() {
			err := callConnect(client, ctx, id)
			if err != nil {
				errCh <- err
			}
			wg.Done()
		}()
	}

	for i := 0; i < 10; i++ {
		connect(i)
	}

	go func() {
		time.Sleep(1000 * time.Millisecond)
		log.Print("cancel")
		cancel()
	}()

	wg.Wait()
	close(errCh)

	for err := range errCh {
		log.Print(err)
	}
}

func callConnect(client echo.EchoClient, ctx context.Context, id int) error {
	stream, err := client.Connect(ctx)
	if err != nil {
		return errors.Wrapf(err, "failed to call Connect")
	}

	tick := time.NewTicker(500 * time.Millisecond)

	for {
		select {
		case <-tick.C:
			err := stream.Send(&echo.Message{Body: fmt.Sprintf("%2d: Hello, World!", id)})
			if err != nil {
				return errors.Wrapf(err, "failed to send")
			}

			msg, err := stream.Recv()
			if err != nil {
				return errors.Wrapf(err, "failed to receive")
			}

			log.Printf("recv %2d: %#v", id, msg)

		case <-ctx.Done():
			return nil
		}
	}
}
