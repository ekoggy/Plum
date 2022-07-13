package main

import (
	"strings"
	"time"

	"github.com/kaoriEl/go-tdlib/client"
	"github.com/kaoriEl/go-tdlib/tdlib"
	//	"time"
	_ "github.com/lib/pq"
)

func CollectInfoFromTelegram() error {
	var err error
	client.SetLogVerbosityLevel(1)
	client.SetFilePath("./errors.txt")
	client := client.NewClient(client.Config{
		APIID:               "187786",
		APIHash:             "e782045df67ba48e441ccb105da8fc85",
		SystemLanguageCode:  "en",
		DeviceModel:         "Server",
		SystemVersion:       "1.0.0",
		ApplicationVersion:  "1.0.0",
		UseMessageDatabase:  true,
		UseFileDatabase:     true,
		UseChatInfoDatabase: true,
		UseTestDataCenter:   false,
		DatabaseDirectory:   "./tdlib-db",
		FileDirectory:       "./tdlib-files",
		IgnoreFileNames:     false,
	})

	currentState, err := client.Authorize()
	if err != nil {
		return err
	}
	for ; currentState.GetAuthorizationStateEnum() != tdlib.AuthorizationStateReadyType; currentState, _ = client.Authorize() {
		time.Sleep(300 * time.Millisecond)
	}
	//var msgs []tdlib.Message
	//for i := 0; i < 50; i++ {
	//	msg, err := client.GetChatHistory(-1001678455451, 0, int32(-i), 50, false)
	//	if err != nil {
	//		break
	//	}
	//	msgs = append(msgs, msg.Messages[0])
	//}
	//mmm := (*msgs).(*tdlib.UpdateNewMessage)

	eventFilter := func(msg *tdlib.TdMessage) bool {
		updateMsg := (*msg).(*tdlib.UpdateNewMessage)
		if updateMsg.Message.IsChannelPost == true {
			result := updateMsg.Message.ChatID == -1001678455451
			return result
		}
		return false
	}

	link, err := client.CreateChatInviteLink(-1001678455451, 0, 0)
	if err != nil {
		return err
	}
	var rec Record
	rec.Source = link.InviteLink
	receiver := client.AddEventReceiver(&tdlib.UpdateNewMessage{}, eventFilter, 5)
	for newMsg := range receiver.Chan {
		updateMsg := (newMsg).(*tdlib.UpdateNewMessage)
		msg := updateMsg.Message.Content.(*tdlib.MessagePhoto)
		rec.Name = strings.Split(msg.Caption.Text, "\n")[0]
		rec.Date = strings.Split(msg.Caption.Text, "\n")[1]
		rec.Size = strings.Split(msg.Caption.Text, "\n")[2]
		rec.Price = strings.Split(msg.Caption.Text, "\n")[3]
		rec.Buy = "127.0.0.1:3000"
		_, err := insert(rec.Name, rec.Size, rec.Date, rec.Price, rec.Buy, rec.Source)
		if err != nil {
			return err
		}
	}
	return nil
}
