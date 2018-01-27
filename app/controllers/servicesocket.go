package controllers

import (
	"fmt"
	"time"

	"github.com/revel/revel"
)

type ServiceSocket struct {
	*revel.Controller
}

func init() {
	fmt.Println("\n\n\nInitProcessSocket\n\n\n")
}

func (c ServiceSocket) HandleWSConnection(ws revel.ServerWebSocket) revel.Result {

	if ws == nil {
		fmt.Println("invalid socket")
		return nil
	}

	OneSecChannel := time.NewTicker(time.Second * 1)
	tickCount := 0
	maxTicks := 10

	//recieve messages into this channel
	messagesFromSocket := make(chan string)
	go func() {
		var msg string
		for {
			err := ws.MessageReceiveJSON(&msg)
			if err != nil {
				close(messagesFromSocket)
				return
			}
			fmt.Println(msg)
			messagesFromSocket <- msg
		}
	}()

	//do something with new messages
	for {
		select {
		case tickTime := <-OneSecChannel.C:
			tickCount++
			if tickCount > maxTicks {
				OneSecChannel.Stop()
			}
			if ws.MessageSendJSON(fmt.Sprintf("%3d", tickCount)+" "+tickTime.String()) != nil {
				// socket disconnected
				fmt.Println("Error 1 sending to client")
				break
			}
		case msg, ok := <-messagesFromSocket:
			// If the channel is closed, they disconnected.
			if !ok {
				break
			}

			fmt.Printf("Recieved Message\n%s\n", msg)
			ws.MessageSendJSON(msg)
		}
	}

	// if ws.MessageSendJSON("Done") != nil {
	// 	// socket disconnected
	// 	fmt.Println("Error 2 sending to client")
	// }

	// data := make(map[string]interface{})
	// data["success"] = true
	// return c.RenderJSON(data)
}
