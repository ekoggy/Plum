package telegram

import (
	"fmt"
	"github.com/ekoggy/Plum/postgre"
	"github.com/kaoriEl/go-tdlib/client"
	"github.com/kaoriEl/go-tdlib/tdlib"
	"strings"
)

func CollectInfoFromTelegram() error {
	var err error
	var rec postgre.Record
	client.SetLogVerbosityLevel(1)
	client.SetFilePath("./errors.txt")
	cli := client.NewClient(client.Config{
		APIID:               "15728153",
		APIHash:             "7a959e9a1e68300cd6a1bbfcea3b7a96",
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

	for {
		currentState, err := cli.Authorize()
		if err != nil {
			return err
		}
		fmt.Println(currentState)
		if currentState.GetAuthorizationStateEnum() == tdlib.AuthorizationStateWaitPhoneNumberType {
			fmt.Print("Enter phone: ")
			var number string
			fmt.Scanln(&number)
			_, err := cli.SendPhoneNumber(number)
			if err != nil {
				fmt.Printf("Error sending phone number: %v", err)
			}
		} else if currentState.GetAuthorizationStateEnum() == tdlib.AuthorizationStateWaitCodeType {
			fmt.Print("Enter code: ")
			var code string
			fmt.Scanln(&code)
			_, err := cli.SendAuthCode(code)
			if err != nil {
				fmt.Printf("Error sending auth code : %v", err)
			}
		} else if currentState.GetAuthorizationStateEnum() == tdlib.AuthorizationStateWaitPasswordType {
			fmt.Print("Enter Password: ")
			var password string
			fmt.Scanln(&password)
			_, err := cli.SendAuthPassword(password)
			if err != nil {
				fmt.Printf("Error sending auth password: %v", err)
			}
		} else if currentState.GetAuthorizationStateEnum() == tdlib.AuthorizationStateReadyType {
			fmt.Println("Authorization Ready! Let's rock")
			break
		}
	}

	chat, err := cli.SearchPublicChat("TestDataBaseLeaks")
	if err != nil {
		return err
	}

	last, err := cli.GetChatHistory(chat.ID, 0, 0, 1, false)
	if err != nil {
		return err
	}

	msgs, err := cli.GetChatHistory(chat.ID, last.Messages[0].ID, 0, 10, false)
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

func fillTheRecord(rec *postgre.Record, msg *tdlib.MessagePhoto, cli *client.Client) error {
	rec.Name = strings.Split(msg.Caption.Text, "\n")[0]
	rec.Date = strings.Split(msg.Caption.Text, "\n")[1]
	rec.Size = strings.Split(msg.Caption.Text, "\n")[2]
	rec.Price = strings.Split(msg.Caption.Text, "\n")[3]
	entity := fmt.Sprintf("%s", msg.Caption.Entities[1].Type)
	rec.Buy = entity[27 : len(entity)-1]
	rec.Source = "t.me/TestDataBaseLeaks"
	_, err := postgre.Insert(rec.Name, rec.Size, rec.Date, rec.Price, rec.Buy, rec.Source)
	if err != nil {
		return err
	}
	return nil
}
