package context

import (
	"context"
	"fmt"
	"time"
)

func HandleRequest(ctx context.Context) {
	go writeRedis(ctx)
	go writeDB(ctx)
	for {
		select {
		case <-ctx.Done():
			fmt.Println("Handle request done...")
			return
		default:
			fmt.Println("Handle request running...")
			time.Sleep(2 * time.Second)
		}
	}
}

func writeRedis(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			fmt.Println("write redis done...")
			return
		default:
			fmt.Println("write redis running...")
			time.Sleep(2 * time.Second)
		}
	}
}

func writeDB(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			fmt.Println("write db done...")
			return
		default:
			fmt.Println("write db running...")
			time.Sleep(2 * time.Second)
		}
	}
}

func main(){
	ctx,cancel := context.WithCancel(context.Background())
	go HandleRequest(ctx)

	time.Sleep(5*time.Second)
	fmt.Println("time to stop handle")
	cancel()
	time.Sleep(5*time.Second)
}
