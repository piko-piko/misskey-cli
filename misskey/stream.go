package misskey

import (
	"errors"
	"fmt"
	"log"
	"net/url"
	"os"
	"os/signal"

	"github.com/buger/jsonparser"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"github.com/sacOO7/gowebsocket"
)

// gettimeline.goに依存した実装

func (c *Client) GetStream(plainPrint bool, mode string) error {

	fmt.Println("Stream: " + mode + "  @" + c.InstanceInfo.UserName + " (" + c.InstanceInfo.Host + ")")
	if !plainPrint {
		printLine()
	}

	parsedUrl, err := url.Parse(c.InstanceInfo.Host)
	if err != nil {
		log.Fatal(err)
	}

	wsUrl := "wss://" + parsedUrl.Host + "/streaming?i=" + c.InstanceInfo.Token

	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)

	socket := gowebsocket.New(wsUrl)

	socket.OnConnected = func(socket gowebsocket.Socket) {
		log.Println("Connected to server")
	}

	socket.OnConnectError = func(err error, socket gowebsocket.Socket) {
		log.Println("Received connect error ", err)
	}

	socket.OnTextMessage = func(message string, socket gowebsocket.Socket) {
		printNote(message)
	}

	socket.OnDisconnected = func(err error, socket gowebsocket.Socket) {
		log.Println("Disconnected from server.")
		socket.Connect()
		initialConnect(socket, mode)
	}

	socket.OnPingReceived = func(data string, socket gowebsocket.Socket) {
		socket.SendBinary([]byte{websocket.PongMessage})
	}

	socket.Connect()
	initialConnect(socket, mode)

	for {
		select {
		case <-interrupt:
			log.Println("interrupt")
			socket.Close()
			return nil
		}

	}

}

func initialConnect(socket gowebsocket.Socket, mode string) error {
	uu, err := uuid.NewRandom()
	if err != nil {
		return err
	}
	mainChId := uu.String()

	uu, err = uuid.NewRandom()
	if err != nil {
		return err
	}
	tlChId := uu.String()

	socket.SendText("{\"type\":\"connect\",\"body\":{\"channel\":\"main\",\"id\":\"" + mainChId + "\"}}")

	var channelText string

	if mode == "local" || mode == "global" || mode == "home" {
		channelText = "{\"type\":\"connect\",\"body\":{\"channel\":\"" + mode + "Timeline\",\"id\":\"" + tlChId + "\"}}"
	} else {
		return errors.New("Please select mode in local/home/global")
	}
	socket.SendText(channelText)
	return nil
}

func printNote(message string) {
	var err error

	messageBody, _, _, err := jsonparser.Get([]byte(message), "body", "body")
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return
	}

	// とりあえずTextを持ってきてみる
	_, err = jsonparser.GetString(messageBody, "renoteId")

	var note *Note

	if err != nil {
		note, err = NewNote(messageBody)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			return
		}

		_, err = jsonparser.GetString(messageBody, "replyId")
		if err == nil {
			replyParentValue, _, _, _ := jsonparser.Get(messageBody, "reply")
			replyParent, err := NewNote(replyParentValue)
			if err != nil {
				fmt.Fprintln(os.Stderr, err)
				return
			}
			fmt.Fprintln(os.Stdout, replyParent)
			note.Offset = "    "
		}

	} else { // renoteだったら

		renoteValue, _, _, _ := jsonparser.Get(messageBody, "renote")

		note, err = NewNote(renoteValue)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			return
		}

		note.User.Name = "[RN]" + note.User.Name

	}

	str := fmt.Sprint(note)

	fmt.Fprintln(os.Stdout, str)

	return
}
