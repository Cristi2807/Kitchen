package main

import (
	"encoding/json"
	"fmt"
	"os"
	"time"
)

var jobsRank [3]chan OrderFood
var jobsDone = make(chan OrderFood, 500)

type Cooks []struct {
	Rank        int `json:"rank"`
	Proficiency int `json:"proficiency"`
}

var cooks Cooks

func InitCookRankChs() {
	for i := 0; i < 3; i++ {
		jobsRank[i] = make(chan OrderFood, 100)
	}
}

func ParseCooks() {
	menuFile, _ := os.Open("cooks.json")
	jsonParser := json.NewDecoder(menuFile)
	jsonParser.Decode(&cooks)
}

func PrepareFood(foodRank OrderFood, CookId int) {
	time.Sleep(time.Duration(menu[foodRank.FoodId-1].PreparationTime*TIMEUNIT) * time.Millisecond)

	fmt.Printf("Cook %d DONE food %d from Order Nr. %d\n", CookId, foodRank.FoodId, foodRank.OrderId)
	foodRank.CookId = CookId
	jobsDone <- foodRank
}

func HandleCook(CookId int, Rank int) {

	if Rank == 1 {
		for foodRank := range jobsRank[0] {
			PrepareFood(foodRank, CookId)
		}
	}

	if Rank == 2 {
		for {
			select {
			case foodRank := <-jobsRank[1]:
				{
					PrepareFood(foodRank, CookId)
				}
			default:
				select {
				case foodRank := <-jobsRank[1]:
					{
						PrepareFood(foodRank, CookId)
					}
				case foodRank := <-jobsRank[0]:
					{
						PrepareFood(foodRank, CookId)
					}
				}
			}
		}

	}

	if Rank == 3 {
		for {
			select {
			case foodRank := <-jobsRank[2]:
				{
					PrepareFood(foodRank, CookId)
				}
			default:
				select {
				case foodRank := <-jobsRank[2]:
					{
						PrepareFood(foodRank, CookId)
					}
				case foodRank := <-jobsRank[1]:
					{
						PrepareFood(foodRank, CookId)
					}
				default:

					select {
					case foodRank := <-jobsRank[2]:
						{
							PrepareFood(foodRank, CookId)
						}
					case foodRank := <-jobsRank[1]:
						{
							PrepareFood(foodRank, CookId)
						}
					case foodRank := <-jobsRank[0]:
						{
							PrepareFood(foodRank, CookId)
						}
					}

				}

			}
		}

	}
}
