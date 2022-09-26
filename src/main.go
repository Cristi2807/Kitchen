package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"
)

const TIMEUNIT int = 200

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
	fmt.Printf("%v\n\n", order)
	order.CookingTime = time.Now().UnixMilli()
	ordersMap.Store(order.Id, order)

	for i := 0; i < len(order.Items); i++ {
		jobsRank[menu[order.Items[i]-1].Complexity-1] <- OrderFood{OrderId: order.Id, FoodId: order.Items[i]}
	}

}

func main() {
	http.HandleFunc("/order", getOrder)

	ParseMenu()
	ParseCooks()

	InitCookRankChs()

	go HandleDoneJobs()

	for i := 0; i < len(cooks); i++ {
		for j := 0; j < cooks[i].Proficiency; j++ {
			go HandleCook(i, cooks[i].Rank)
		}
	}

	fmt.Printf("Server Kitchen started on PORT 8010\n")
	if err := http.ListenAndServe(":8010", nil); err != nil {
		log.Fatal(err)
	}
}
