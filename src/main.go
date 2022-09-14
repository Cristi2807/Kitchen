package main

import (
	"bytes"
	"container/list"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"sync"
	"time"
)

const TIMEUNIT int = 300

type MenuItem struct {
	Id       int    `json:"id"`
	Name     string `json:"name"`
	PrepTime int    `json:"preparation-time"`
}

type Order struct {
	Id      int   `json:"order_id"`
	Items   []int `json:"items"`
	MaxWait int   `json:"max_wait"`
}

var Menu = []MenuItem{
	{Id: 1, Name: "Pizza", PrepTime: 20 * TIMEUNIT},
	{Id: 2, Name: "Salad", PrepTime: 10 * TIMEUNIT},
	{Id: 3, Name: "Zeama", PrepTime: 7 * TIMEUNIT},
	{Id: 4, Name: "Scallop Sashimi with Meyer Lemon Confit", PrepTime: 32 * TIMEUNIT},
	{Id: 5, Name: "Island Duck with Mulberry Mustard", PrepTime: 35 * TIMEUNIT},
	{Id: 6, Name: "Waffles", PrepTime: 10 * TIMEUNIT},
	{Id: 7, Name: "Aubergine", PrepTime: 20 * TIMEUNIT},
	{Id: 8, Name: "Lasagna", PrepTime: 30 * TIMEUNIT},
	{Id: 9, Name: "Burger", PrepTime: 15 * TIMEUNIT},
	{Id: 10, Name: "Gyros", PrepTime: 15 * TIMEUNIT},
	{Id: 11, Name: "Kebab", PrepTime: 15 * TIMEUNIT},
	{Id: 12, Name: "Unagi Maki", PrepTime: 20 * TIMEUNIT},
	{Id: 13, Name: "Tobacco Chicken", PrepTime: 30 * TIMEUNIT},
}

var OrdersList = list.New()

func getOrder(w http.ResponseWriter, r *http.Request) {

	if r.URL.Path != "/order" {
		http.Error(w, "404 not found.", http.StatusNotFound)
		return
	}

	if r.Method != "POST" {
		http.Error(w, "Method is not supported.", http.StatusNotFound)
		return
	}

	var order Order
	json.NewDecoder(r.Body).Decode(&order)
	fmt.Printf("%d %v %d \n\n", order.Id, order.Items, order.MaxWait)
	OrdersList.PushBack(order)

}

var m sync.Mutex

func tryGetOrder() (Order, bool) {
	m.Lock()
	defer m.Unlock()
	if OrdersList.Front() != nil {
		return OrdersList.Remove(OrdersList.Front()).(Order), false
	}
	return Order{}, true
}

func cookOrder() {
	for {
		a, err := tryGetOrder()

		if err == false {
			time.Sleep(time.Duration(a.MaxWait) * time.Millisecond)

			orderMarshalled, _ := json.Marshal(a)
			responseBody := bytes.NewBuffer(orderMarshalled)

			http.Post("http://dinninghall:8020/distribution", "application/json", responseBody)
		}
	}
}

func main() {
	http.HandleFunc("/order", getOrder)

	go cookOrder()
	go cookOrder()
	go cookOrder()

	fmt.Printf("Server Kitchen started on PORT 8010\n")
	if err := http.ListenAndServe(":8010", nil); err != nil {
		log.Fatal(err)
	}
}
