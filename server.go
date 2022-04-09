package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/go-redis/redis"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{} // use default options

func socketHandler(w http.ResponseWriter, r *http.Request) {

	redisClient := NewRedisClient()

	// Upgrade our raw HTTP connection to a websocket based one
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Print("Error during connection upgradation:", err)
		return
	}
	defer conn.Close()

	// The event loop
	for {
		val, err := redisClient.Get("streams=btcusdt@aggTrade").Result() // => GET key
		if err != nil {
			panic(err)
		}
		err = conn.WriteMessage(websocket.TextMessage, []byte(val))
		if err != nil {
			log.Println(err)
			return
		}
	}
}

func home(w http.ResponseWriter, r *http.Request) {
	redisClient := NewRedisClient()
	val, err := redisClient.Get("streams=btcusdt@aggTrade").Result() // => GET key
	if err != nil {
		panic(err)
	}
	fmt.Fprintf(w, val)
}

func main() {
	http.HandleFunc("/socket", socketHandler)
	http.HandleFunc("/", home)
	log.Fatal(http.ListenAndServe("localhost:8080", nil))
}
func NewRedisClient() *redis.Client { // 實體化redis.Client 並返回實體的位址
	client := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "", // no password set
		DB:       0,  // use default DB
	})
	return client
}
