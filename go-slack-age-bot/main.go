package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/shomali11/slacker"
)

func printCommandEvents(analyticsChannel <-chan *slacker.CommandEvent){
	for event := range analyticsChannel{
		fmt.Println("Command Events")
		fmt.Println(event.Timestamp)
		fmt.Println(event.Command)
		fmt.Println(event.Parameters)
		fmt.Println(event.Event)
	}
}

func main() {
	os.Setenv("SLACK_BOT_TOKEN", "xoxb-1029611081895-3370050029985-X3OqbpM41A9B7ogHhopIIlqu")
	os.Setenv("SLACK_APP_TOKEN", "xapp-1-A03AANE2ZSS-3370031835985-97eb2b2cb077b201b94470e2be38250f9a10e1dc93f46099dae3c99c49add67d")

	bot := slacker.NewClient(os.Getenv("SLACK_BOT_TOKEN"), os.Getenv("SLACK_APP_TOKEN"))

	go printCommandEvents(bot.CommandEvents())

	bot.Command("My YOB is <year>", &slacker.CommandDefinition {
		Description: "YOB calculator",
		Example: "My YOB is 2020",
		Handler: func(botCtx slacker.BotContext, request slacker.Request, response slacker.ResponseWriter){
			year := request.Param("year")
			yob, err := strconv.Atoi(year)

			if err != nil {
				fmt.Println("Error", err)
			}

			age := 2022 - yob
			r := fmt.Sprintf("Your age is %d", age)
			response.Reply(r)

		},
	})


	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	err := bot.Listen(ctx)
	if err != nil {
		log.Fatal(err)
	}

}