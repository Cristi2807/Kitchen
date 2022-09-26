package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"sync"
	"time"
)

type PreparedFood struct {
	FoodId int `json:"food_id"`
	CookId int `json:"cook_id"`
}

type OrderFood struct {
	OrderId int
	FoodId  int
	CookId  int
}

type Order struct {
	Id             int            `json:"order_id"`
	Items          []int          `json:"items"`
	MaxWait        int            `json:"max_wait"`
	TableId        int            `json:"table_id"`
	WaiterId       int            `json:"waiter_id"`
	Priority       int            `json:"priority"`
	PickUpTime     int64          `json:"pick_up_time"`
	CookingDetails []PreparedFood `json:"cooking_details"`
	CookingTime    int64          `json:"cooking_time"`
}

var ordersMap sync.Map

func HandleDoneJobs() {
	for jobDone := range jobsDone {

		value, _ := ordersMap.Load(jobDone.OrderId)
		order := value.(Order)

		order.CookingDetails = append(order.CookingDetails, PreparedFood{FoodId: jobDone.FoodId, CookId: jobDone.CookId})

		if len(order.CookingDetails) == len(order.Items) {
			ordersMap.Delete(order.Id)

			order.CookingTime = time.Now().UnixMilli() - order.CookingTime

			orderMarshalled, _ := json.Marshal(order)
			responseBody := bytes.NewBuffer(orderMarshalled)

			http.Post("http://dinninghall:8020/distribution", "application/json", responseBody)

		}

		ordersMap.Store(jobDone.OrderId, order)

	}
}
