package main

import (
	"fmt"
	"log"
	"time"

	"golang.org/x/net/context"
)

func main() {
	ctx, cancel := context.WithTimeout(context.Background(), 5500*time.Millisecond)
	defer func() {
		cancel()
		time.Sleep(1 * time.Second)
	}()
	countDone := make(chan bool, 1)

	ctx = context.WithValue(ctx, "value", map[string]interface{}{"n": 5, "done": countDone})
	ctx2 := context.WithValue(ctx, "value", map[string]interface{}{"n": 3, "done": countDone})

	go counterCtx(ctx)
	go counterCtx(ctx2)

	select {
	case <-ctx.Done():
		log.Fatal(ctx.Err())
	case <-countDone:
	}
}

func counterCtx(ctx context.Context) {
	value := ctx.Value("value").(map[string]interface{})
	n := value["n"].(int)
	done := value["done"].(chan bool)
	for i := 1; i <= n; i++ {
		select {
		case <-time.After(1 * time.Second):
			fmt.Println(i)
		case <-ctx.Done():
			fmt.Println(ctx.Err())
			return
		}
	}
	done <- true
}
