package Mqtt

import (
	"fmt"
	"github.com/eclipse/paho.mqtt.golang"
	"log"
)

var opts *mqtt.ClientOptions
var client mqtt.Client
var token mqtt.Token

type Callback func(string, []byte)

func Connect(server_url string) {
	opts = mqtt.NewClientOptions().AddBroker(server_url)

	client = mqtt.NewClient(opts)
	if token = client.Connect(); token.Wait() && token.Error() != nil {
		panic(token.Error())
//		return token.Error()
	} else {
		log.Println("Mqtt connection successful")
//		return nil
	}
}

func messageReceivedCallback(client_ mqtt.Client, message_ mqtt.Message)  {
	fmt.Printf("mqtt. Topic: %s. Message: %s\r\n", message_.Topic(), message_.Payload())
}

func Subscribe(topic string) {
	if token = client.Subscribe(topic, 0, messageReceivedCallback); token.Wait() && token.Error() != nil {
		panic(token.Error())
//		return token.Error()
	} else {
//		return nil
	}
}

func SubscribeClbk(topic string, clbk Callback) {
	if token = client.Subscribe(topic, 0, func(client mqtt.Client, msg mqtt.Message) { clbk(msg.Topic(), msg.Payload()) }); token.Wait() && token.Error() != nil {
		panic(token.Error())
//		return token.Error()
	} else {
//		return nil
	}
}

func Unsubscribe(topic string) error {
	if token := client.Unsubscribe(topic); token.Wait() && token.Error() != nil {
		return token.Error()
	} else {
		return nil
	}
}
