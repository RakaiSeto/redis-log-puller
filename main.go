package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	duitrapi "github.com/rakaiseto/redis-log-puller/consumer/DuitRapi"
	ocr_marketplace "github.com/rakaiseto/redis-log-puller/consumer/OCR_Marketplace"
	utils "github.com/rakaiseto/redis-log-puller/utils"
)

func main() {
	// 1. Initialize Redis
	rdb, err := utils.NewRedisClient()
	if err != nil {
		fmt.Printf("Error connecting to Redis: %v\n", err)
		return
	}
	fmt.Println("Connected to Redis")
	defer rdb.Close()

	// 2. Initialize Router and Consumers
	ocrConsumer, err := ocr_marketplace.NewOCRMarketplaceConsumer("ocr_marketplace")
	if err != nil {
		fmt.Printf("Error initializing OCR Marketplace consumer: %v\n", err)
		return
	}

	duitRapiConsumer, err := duitrapi.NewDuitRapiConsumer("duitrapi")
	if err != nil {
		fmt.Printf("Error initializing DuitRapi consumer: %v\n", err)
		return
	}

	routing := map[string]utils.Consumer{
		"queue:ocr_marketplace": ocrConsumer,
		"queue:duitrapi":        duitRapiConsumer,
	}

	router := utils.NewRouter()
	for queue, consumer := range routing {
		router.Register(queue, consumer)
	}

	// 3. Get all queues to listen to
	queues := router.GetQueues()
	if len(queues) == 0 {
		fmt.Println("No queues registered. Exiting.")
		return
	}
	fmt.Printf("Listening on queues: %v\n", queues)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Handle graceful shutdown
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		<-sigChan
		fmt.Println("\nReceived shutdown signal...")
		cancel()
	}()

	// 4. Main BLPop Loop
	for {
		select {
		case <-ctx.Done():
			fmt.Println("Shutting down...")
			return
		default:
			// BLPop blocks. We use a timeout (e.g., 1 second) so we can check ctx.Done() periodically
			// or we can use 0 for infinite block if we handle cancellation differently.
			// With go-redis, passing a context that gets cancelled will unblock BLPop.
			// So we can use 0 timeout.
			result, err := rdb.BLPop(ctx, 0, queues...).Result()
			if err != nil {
				if err != context.Canceled {
					fmt.Printf("Redis error: %v\n", err)
					time.Sleep(time.Second) // Backoff on error
				}
				continue
			}

			// result[0] is the queue name, result[1] is the value
			if len(result) == 2 {
				queue := result[0]
				data := result[1]
				if err := router.Route(ctx, queue, data); err != nil {
					fmt.Printf("Error routing message from %s: %v\n", queue, err)
				}
			}
		}
	}
}
