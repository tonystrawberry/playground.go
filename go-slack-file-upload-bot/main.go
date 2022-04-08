package main

import (
	"fmt"
	"os"

	"github.com/slack-go/slack"
)

func main() {
  os.Setenv("SLACK_BOT_TOKEN", "xoxb-1029611081895-3370050029985-hZkTU8aoAkG2ovSIUwRdP8RD")
  os.Setenv("CHANNEL_ID", "C03B16MRE8Z")

  api := slack.New(os.Getenv("SLACK_BOT_TOKEN"))
  channelArr := []string{os.Getenv("CHANNEL_ID")}
  fileArr := []string{"sample.pdf"}

  for i := 0; i < len(fileArr); i++ {
    params := slack.FileUploadParameters{
      Channels: channelArr,
      File: fileArr[i],
    }

    file, err := api.UploadFile(params)
    
    if err != nil {
      fmt.Printf("%s\n", err)
    }

    fmt.Printf("Name: %s, URL: %s\n", file.Name, file.URL)
  }
}