package controllers

import (
	"bufio"
	"bytes"
	"container/list"
	"fmt"
	"io"
	"net"
	"os/exec"
	"strings"
	"sync"
	"time"

	"github.com/revel/revel"
)

var mutex = &sync.Mutex{}

type ServiceSocket struct {
	*revel.Controller
}

var processSocketData = make(chan []byte, 512)
var processSocketActive = false
var processSocketChecked = false

var processSocketLocalLineBuffer = list.New()
var activeConnectionsList = list.New()

func init() {
	go initProcessSocket()
}

func sendToAllConnections(msg []byte) {
	i := 0
	for e := activeConnectionsList.Front(); e != nil; e = e.Next() {
		i++
		fmt.Println("Sending Data to ", i)
		el := e.Value.(chan []byte)
		el <- msg
	}
}
func (c ServiceSocket) HandleWSConnection(ws revel.ServerWebSocket) revel.Result {

	if ws == nil {
		fmt.Println("Error 1 socket not valid")
		return nil
	}

	fmt.Println("Sending old console output to client")
	for e := processSocketLocalLineBuffer.Front(); e != nil; e = e.Next() {
		str := e.Value.([]byte)
		fmt.Printf("%s\n", str)
		if ws.MessageSendJSON(string(str)) != nil {
			fmt.Println("Error 2 sending to client")
			break
		}
	}
	messagesToSocket := make(chan []byte)
	connectionListEntry := activeConnectionsList.PushBack(messagesToSocket)
	//recieve messages into this channel
	messagesFromSocket := make(chan []byte)
	go func() {
		var msg []byte
		for {
			err := ws.MessageReceiveJSON(&msg)
			if err != nil {
				close(messagesFromSocket)
				return
			}
			if strings.Compare(string(msg), "StartService") == 0 {
				launchProcessSocket()
				continue
			}
			if strings.Compare(string(msg), "DumpSocketBuffer") == 0 {
				fmt.Println("DumpSocketBuffer")
				for e := processSocketLocalLineBuffer.Front(); e != nil; e = e.Next() {
					str := e.Value.([]byte)
					fmt.Printf("%s\n", str)
				}
				continue
			}
			fmt.Println(msg)
			messagesFromSocket <- []byte(msg)
		}
	}()

	// OneSecChannel := time.NewTicker(time.Second * 1)
	// tickCount := 0
	// maxTicks := 10

	//do something with new messages
SocketLoop:
	for {
		select {
		// case tickTime := <-OneSecChannel.C:
		// 	tickCount++
		// 	if tickCount > maxTicks {
		// 		OneSecChannel.Stop()
		// 	}
		// 	if ws.MessageSendJSON(fmt.Sprintf("%3d", tickCount)+" "+tickTime.String()) != nil {
		// 		// socket disconnected
		// 		fmt.Println("Error 3 sending to client")
		// 		break
		// 	}
		case msg := <-messagesToSocket:
			ws.MessageSendJSON(string(msg))
		case msg, notClosed := <-messagesFromSocket:
			// If the channel is closed, they disconnected.
			if !notClosed {
				fmt.Println("WS Closed 2!")
				break SocketLoop
			}

			fmt.Printf("WS Recieved Message\n%s\n", msg)
			processSocketData <- msg
		}
	}

	fmt.Println("websocket closed!")
	activeConnectionsList.Remove(connectionListEntry)
	// if ws.MessageSendJSON("Done") != nil {
	// 	// socket disconnected
	// 	fmt.Println("Error 2 sending to client")
	// }

	// data := make(map[string]interface{})
	// data["success"] = true
	return c.RenderJSON(nil)
}

func launchProcessSocket() {
	if processSocketActive == true {
		return
	}
	fmt.Println("Launching Process!")
	cmd := exec.Command("/home/miner/gocode/src/github.com/ryandrew/process-output-socket/process-output-socket")
	err := cmd.Start()
	if err != nil {
		fmt.Println("Error launching Process Socket")
	}
	fmt.Println("Launched!")
	processSocketActive = true
}

func initProcessSocket() {
	fmt.Println("\nInitProcessSocket\n")
	if processSocketChecked == true && processSocketActive == false {
		//launchProcessSocket()
		time.Sleep(200 * time.Millisecond)
	}
	processSocketChecked = true

	conn, err := net.Dial("tcp4", "127.0.0.1:25560")
	if err != nil {
		fmt.Println("\nSocket Connect Error! Trying again in 1 sec\n", err)
		time.Sleep(time.Second)
		initProcessSocket()
		return
	}
	//timeoutTimer := time.NewTimer(time.Second * 4)
	if !socketPizzaHandshake(conn) {
		fmt.Println("Handshake Failed. Closing Connection")
		conn.Close()
		return
	}
	fmt.Println("\nSocket Connected!\n", err)
	processSocketClosed := make(chan struct{})
	processSocketActive = true

	go func() {
		fmt.Println("Process Connection Beginning to stream!")
		scanner := bufio.NewScanner(conn)
		for scanner.Scan() {

			socketData := scanner.Bytes()
			fmt.Printf("Read %d bytes from socket %s\n", len(socketData), conn.RemoteAddr())

			fmt.Printf("%s\n", socketData)

			processSocketData <- socketData
		}
		err = scanner.Err()
		if err != nil {
			fmt.Println("Socket Scanner Error: ", err)
		}
		processSocketClosed <- struct{}{}
	}()

	for {
		select {
		case <-processSocketClosed:
			processSocketActive = false
			break
		case msg := <-processSocketData:
			//withNewLine := append(msg, '\n')

			mutex.Lock()
			processSocketLocalLineBuffer.PushBack(msg)
			mutex.Unlock()
			sendToAllConnections(msg)
			//conn.Write(withNewLine)
		}
	}

}

func socketPizzaHandshake(conn net.Conn) bool {
	start := time.Now()

	byteReader := bufio.NewReaderSize(conn, 2)
	challengeBytes := make([]byte, 2)
	challengeBytesExpected := []byte{'P', 'I'} //[]byte("PI")

	_, err := io.ReadFull(byteReader, challengeBytes)
	if err != nil {
		fmt.Println("Pizza Handshake error 1.1: Socket Read Error!\n", err)
		return false
	}
	if bytes.Compare(challengeBytesExpected, challengeBytes) != 0 {
		fmt.Println("Pizza Handshake error 2.1: Invalid Challenge Expecting \"PI\"")
		return false
	}

	conn.Write([]byte("ZZ\r"))

	scanner := bufio.NewScanner(conn)
	scanner.Scan()
	if err := scanner.Err(); err != nil {
		fmt.Println("Pizza Handshake error 3: Error Waiting for challenge success\n", err)
		return false
	}

	if strings.Compare(scanner.Text(), "PIZZA!") != 0 {
		fmt.Println("Pizza Handshake error 4: Invalid challenge success. Expecting \"PIZZA!\"")
		return false
	}

	elapsed := time.Since(start)
	fmt.Printf("Pizza Handshake Complete. Elapsed: %s\n", elapsed)

	return true
}
