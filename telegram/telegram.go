package telegram

import (
	"fmt"
	"strings"
	"time"

	"github.com/kaoriEl/go-tdlib/client"
	"github.com/kaoriEl/go-tdlib/tdlib"
)

func CollectInfoFromTelegram() error {
	var err error
	var rec Record
	client.SetLogVerbosityLevel(1)
	client.SetFilePath("./errors.txt")
	cli := client.NewClient(client.Config{
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

	currentState, err := cli.Authorize()
	if err != nil {
		return err
	}
	for ; currentState.GetAuthorizationStateEnum() != tdlib.AuthorizationStateReadyType; currentState, _ = cli.Authorize() {
		time.Sleep(300 * time.Millisecond)
	}

	last, err := cli.GetChatHistory(-1001678455451, 0, 0, 1, false)

	msgs, err := cli.GetChatHistory(-1001678455451, last.Messages[0].ID, 0, 10, false)
	if err != nil {
		return err
	}

	for i := 0; i < int(msgs.TotalCount); i++ {
		historyMsg := (msgs.Messages[i].Content).(*tdlib.MessagePhoto)
		err := fillTheRecord(&rec, historyMsg, cli)
		if err != nil {
			return err
		}
	}
	eventFilter := func(msg *tdlib.TdMessage) bool {
		updateMsg := (*msg).(*tdlib.UpdateNewMessage)
		if updateMsg.Message.IsChannelPost == true {
			result := updateMsg.Message.ChatID == -1001678455451
			return result
		}
		return false
	}

	receiver := cli.AddEventReceiver(&tdlib.UpdateNewMessage{}, eventFilter, 5)
	for newMsg := range receiver.Chan {
		updateMsg := (newMsg).(*tdlib.UpdateNewMessage)
		msg := updateMsg.Message.Content.(*tdlib.MessagePhoto)
		err := fillTheRecord(&rec, msg, cli)
		if err != nil {
			return err
		}
	}
	return nil
}

func fillTheRecord(rec *Record, msg *tdlib.MessagePhoto, cli *client.Client) error {
	rec.Name = strings.Split(msg.Caption.Text, "\n")[0]
	rec.Date = strings.Split(msg.Caption.Text, "\n")[1]
	rec.Size = strings.Split(msg.Caption.Text, "\n")[2]
	rec.Price = strings.Split(msg.Caption.Text, "\n")[3]
	entity := fmt.Sprintf("%s", msg.Caption.Entities[1].Type)
	rec.Buy = entity[27 : len(entity)-1]
	link, err := cli.CreateChatInviteLink(-1001678455451, 0, 0)
	if err != nil {
		return err
	}
	rec.Source = link.InviteLink
	_, err = insert(rec.Name, rec.Size, rec.Date, rec.Price, rec.Buy, rec.Source)
	if err != nil {
		return err
	}
	return nil
}
