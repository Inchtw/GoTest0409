// client.go
package main

import (
	"log"
	"os"
	"os/signal"

	"github.com/go-redis/redis"
	"github.com/gorilla/websocket"
)

var done chan interface{}
var interrupt chan os.Signal

func receiveHandler(connection *websocket.Conn) {
	// redisClient := NewRedisClient()
	defer close(done)
	for {
		_, msg, err := connection.ReadMessage()
		if err != nil {
			log.Println("Error in receive:", err)
			return
		}
		log.Printf("receive: %s\n", msg)
	}
}

func main() {

	done = make(chan interface{})    // Channel to indicate that the receiverHandler is done
	interrupt = make(chan os.Signal) // Channel to listen for interrupt signal to terminate gracefully

	signal.Notify(interrupt, os.Interrupt) // Notify the interrupt channel for SIGINT

	socketUrl := "wss://stream.yshyqxx.com/stream" + "?streams=btcusdt@aggTrade"
	conn, _, err := websocket.DefaultDialer.Dial(socketUrl, nil)
	if err != nil {
		log.Fatal("Error connecting to Websocket Server:", err)
	}
	defer conn.Close()
	go receiveHandler(conn)

	// Our main loop for the client
	// We send our relevant packets here
	for {
		<-interrupt
		// We received a SIGINT (Ctrl + C). Terminate gracefully...
		log.Println("Received SIGINT interrupt signal. Closing all pending connections")

		// Close our websocket connection
		err := conn.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
		if err != nil {
			log.Println("Error during closing websocket:", err)
			return
		}

		<-done
		log.Println("Receiver Channel Closed! Exiting....")
		return
	}
}

func NewRedisClient() *redis.Client { // 實體化redis.Client 並返回實體的位址
	client := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "", // no password set
		DB:       0,  // use default DB
	})
	return client
}
