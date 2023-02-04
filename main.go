package main

import (
	"bufio"
	"context"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/gojek/courier-go"
)

type ChatMessage struct {
	From      string    `json:"from"`
	To        string    `json:"to"`
	CreatedAt time.Time `json:"created_at"`
	Data      string    `json:"data"`
}

func recvMsg(chatClient *courier.Client, sender, receiver string) {

	cb := func(ctx context.Context, ps courier.PubSub, m *courier.Message) {
		msg := new(ChatMessage)
		if err := m.DecodePayload(msg); err != nil {
			log.Println(err)
		}
		log.Println(msg.Data)

	}

	err := chatClient.SubscribeMultiple(context.Background(), map[string]courier.QOSLevel{
		fmt.Sprintf("/chats/%s/with/%s", sender, receiver): courier.QOSTwo,
		fmt.Sprintf("/chats/%s/with/%s", receiver, sender): courier.QOSTwo,
	}, cb)
	if err != nil {
		log.Println(err)
	}
	time.Sleep(1000 * time.Hour)
}

func sendMsg(chatClient *courier.Client, topic string, loginUser, secondUser string) {

	scanner := bufio.NewScanner(os.Stdin)

	for scanner.Scan() {
		fmt.Print(" > ")
		msg := scanner.Text()

		err := chatClient.Publish(context.Background(), topic, ChatMessage{
			From: loginUser,
			To:   secondUser,
			Data: msg,
		}, courier.QOSTwo)
		if err != nil {
			log.Println(err)
		}

	}

}

func main() {

	// connecting to broker

	credential := strings.Split(os.Args[1], ":")
	loginUser := credential[0]
	loginPassword := credential[1]
	secondUser := os.Args[2]
	clientId := credential[2]

	chatClient, err := courier.NewClient(
		courier.WithAddress("127.0.0.1", 1883),
		courier.WithClientID(clientId),
		courier.WithUsername(loginUser),
		courier.WithPassword(loginPassword),
	)

	if err != nil {
		panic(err)
	}

	if err := chatClient.Start(); err != nil {
		panic(err)
	}

	fmt.Println("connected", chatClient.IsConnected())

	if os.Args[3] == "send" {
		topic := fmt.Sprintf("/chats/%s/with/%s", loginUser, secondUser)
		fmt.Println(topic)

		sendMsg(chatClient, topic, loginUser, secondUser)
	} else {
		topic := fmt.Sprintf("/chats/%s/%s", secondUser, loginUser)
		fmt.Println(topic)
		recvMsg(chatClient, loginUser, secondUser)
	}

}
