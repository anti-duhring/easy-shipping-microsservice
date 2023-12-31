package main

import (
	"database/sql"
	"encoding/json"
	"fmt"

	_ "github.com/go-sql-driver/mysql"

	"github.com/anti-duhring/easy-shipping-microsservice/internal/entity"
	"github.com/anti-duhring/easy-shipping-microsservice/internal/infra/repository"
	"github.com/anti-duhring/easy-shipping-microsservice/internal/usecase"
	"github.com/anti-duhring/easy-shipping-microsservice/pkg/kafka"
	ckafka "github.com/confluentinc/confluent-kafka-go/kafka"
)

// func main() {
// 	myRoute := &entity.Route{}

// 	fmt.Printf("My rout before: %+v\n", myRoute)

// 	err := gofactory.Instantiate(myRoute, &struct{ Name string }{Name: "John Doe"})
// 	if err != nil {
// 		fmt.Println("Error:", err)
// 	}
// 	fmt.Printf("My route after: %+v\n", myRoute)
// }

func main() {
	db, err := sql.Open("mysql", "root:root@tcp(host.docker.internal:3306)/routes?parseTime=true")

	if err != nil {
		panic(err)
	}
	defer db.Close()

	channel := make(chan *ckafka.Message)
	topics := []string{"routes"}
	servers := "host.docker.internal:9094"

	go kafka.Consume(topics, servers, channel)

	repository := repository.NewRouteRepositoryMysql(db)
	var price float64 = 10
	freight := entity.NewFreight(price)
	createRouteUseCase := usecase.NewCreateRouteUseCase(repository, freight)
	changeRouteStatusUseCase := usecase.NewChangeRouteStatusUseCase(repository)

	for msg := range channel {
		input := usecase.CreateRouteInput{}
		json.Unmarshal(msg.Value, &input)

		switch input.Event {
		case "RouteCreated":
			output, err := createRouteUseCase.Execute(input)
			if err != nil {
				fmt.Println(err)
			}
			fmt.Println(output)

		case "RouteStarted", "RouteFinished":
			input := usecase.ChangeRouteStatusInput{}
			json.Unmarshal(msg.Value, &input)
			output, err := changeRouteStatusUseCase.Execute(input)

			if err != nil {
				fmt.Println(err)
			}
			fmt.Println(output)
		}

	}
}
