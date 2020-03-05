package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"

	"github.com/go-redis/redis"
	"github.com/gorilla/mux"
	"github.com/streadway/amqp"
)

var conn *amqp.Connection
var sendQueue amqp.Queue
var ch *amqp.Channel

func main() {
	conn, _ = amqp.Dial(os.Getenv("QUEUE_STRING"))
	ch, _ = conn.Channel()
	defer ch.Close()

	recvQueue, _ := ch.QueueDeclare(
		"queue2", // name
		false,    // durable
		false,    // delete when unused
		false,    // exclusive
		false,    // no-wait
		nil,      // arguments
	)

	sendQueue, _ = ch.QueueDeclare(
		"queue1", // name
		false,    // durable
		false,    // delete when unused
		false,    // exclusive
		false,    // no-wait
		nil,      // arguments
	)

	msgs, _ := ch.Consume(
		recvQueue.Name, // queue
		"",             // consumer
		true,           // auto-ack
		false,          // exclusive
		false,          // no-local
		false,          // no-wait
		nil,            // args
	)

	go func() {
		for d := range msgs {
			client := redis.NewClient(&redis.Options{
				Addr:     os.Getenv("REDIS_STRING"),
				Password: "", // no password set
				DB:       0,  // use default DB
			})
			client.HSet("posts", d.Headers["id"], d.Body)
			fmt.Println(d.Headers["id"])
		}
	}()

	router := mux.NewRouter()

	router.HandleFunc("/posts/{id}", addItem).Methods("POST")

	router.HandleFunc("/posts", getAll).Methods("GET")

	router.HandleFunc("/posts/{id}", get).Methods("GET")

	http.ListenAndServe(":8081", router)
}

func get(w http.ResponseWriter, r *http.Request) {
	var idParam string = mux.Vars(r)["id"]
	id, err := strconv.Atoi(idParam)
	if err != nil {
		// there was an error
		w.WriteHeader(400)
		w.Write([]byte("ID could not be converted to integer"))
		return
	}

	client := redis.NewClient(&redis.Options{
		Addr:     os.Getenv("REDIS_STRING"),
		Password: "", // no password set
		DB:       0,  // use default DB
	})
	val, err := client.HGet("posts", strconv.Itoa(id)).Result()

	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(val))
}

func getAll(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	client := redis.NewClient(&redis.Options{
		Addr:     os.Getenv("REDIS_STRING"),
		Password: "", // no password set
		DB:       0,  // use default DB
	})
	val, _ := client.HGetAll("posts").Result()
	val1, _ := json.Marshal(val)
	w.Write(val1)
}

func addItem(w http.ResponseWriter, r *http.Request) {
	var idParam string = mux.Vars(r)["id"]
	id, err := strconv.Atoi(idParam)
	if err != nil {
		// there was an error
		w.WriteHeader(400)
		w.Write([]byte("ID could not be converted to integer"))
		return
	}
	body, _ := ioutil.ReadAll(r.Body)
	client := redis.NewClient(&redis.Options{
		Addr:     os.Getenv("REDIS_STRING"),
		Password: "", // no password set
		DB:       0,  // use default DB
	})

	client.HSet("posts", strconv.Itoa(id), body)
	w.Header().Set("Content-Type", "application/json")

	args := make(amqp.Table)
	args["id"] = id
	_ = ch.Publish(
		"",             // exchange
		sendQueue.Name, // routing key
		false,          // mandatory
		false,          // immediate
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        []byte(body),
			Headers:     args,
		})
}
