package controllers

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"time"

	"github.com/ryandrew/cmd"

	"github.com/revel/revel"
)

type MfApi struct {
	*revel.Controller
}

var serverHomeDir = "/home/miner/server"
var serverConfigFile = serverHomeDir + "/server.properties"

// var serverCmd Cmd
// var serverCmdStdin ReadCloser
// var serverCmdStdout ReadCloser

func (c MfApi) ServicetStatus() revel.Result {
	data := make(map[string]interface{})
	data["ServerStatus"] = "Stopped"
	data["success"] = true

	return c.RenderJSON(data)
}

func (c MfApi) ServiceStart() revel.Result {
	//https://golang.org/pkg/os/exec/

	// serverCmdText := "java -jar spigot-1.12.2.jar >> log.txt"
	// serverCmd := exec.Command("bash", "-c", serverCmdText)
	// serverCmd.Dir = "/home/miner/server"
	// serverCmdStdin, err := serverCmd.StdinPipe()
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// serverCmdStdout, err := serverCmd.StdoutPipe()
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// serverCmdStderr, err := serverCmd.StderrPipe()
	// if err != nil {
	// 	log.Fatal(err)
	// }

	// err := serverCmd.Start()
	// if err != nil {
	// 	log.Fatal(err)
	// }

	//Launch minecraft in a goroutine
	go func() {
		pidFound := false

		serverCmd := cmd.NewCmd("java", "-jar", "spigot-1.12.2.jar")
		//statusChan :=
		serverCmd.Dir = serverHomeDir
		serverCmd.Start()

		ticker := time.NewTicker(time.Second * 1)

		//for tick := range ticker.C {
		//fmt.Println("tick: ", tick)

		for range ticker.C {
			status := serverCmd.Status()

			if status.PID == 0 {
				return
			}
			if pidFound == false {
				fmt.Printf("pid: %d\n", status.PID)
				pidFound = true
			}

			//fmt.Println("Error: " + status.Error.Error())
			//fmt.Println("Cmd: " + status.Cmd)
			//fmt.Printf("Complete: %t\n", status.Complete)

			// out := len(status.Stdout)
			// //fmt.Printf("Stdout len: %d\n", n)
			// if out > 0 {
			for _, line := range status.Stdout {
				fmt.Println(line)
			}
			// }
			//serverCmd.stdout.buf.Reset()

			// err := len(status.Stderr)
			// //fmt.Printf("Stdout len: %d\n", n)
			// if err > 0 {
			for _, line := range status.Stderr {
				fmt.Println(line)
			}

			if status.Complete == true {
				fmt.Printf("process complete\nexit code: %d\n", status.Exit)
				ticker.Stop()
			}
			// }
			//serverCmd.stderr.buf.Reset()
		}
		// // Check if command is done
		// switch {
		// case finalStatus := <-statusChan:
		// // yes!
		// default:
		// // no, still running
		// }
	}()

	outputData := make(map[string]interface{})
	outputData["ServerStatus"] = "Running"
	outputData["success"] = true

	return c.RenderJSON(outputData)
}

func (c MfApi) ServiceStop() revel.Result {
	data := make(map[string]interface{})
	data["ServerStatus"] = "Running"
	data["success"] = true

	return c.RenderJSON(data)
}

func (c MfApi) ServiceRestart() revel.Result {
	data := make(map[string]interface{})
	data["ServerStatus"] = "Running"
	data["success"] = true

	return c.RenderJSON(data)
}

func (c MfApi) ServerConfigRead() revel.Result {
	configFileContents, err := ioutil.ReadFile(serverConfigFile)
	if err != nil {
		panic(err)
	}

	data := make(map[string]interface{})
	data["configFile"] = serverConfigFile
	data["config"] = fmt.Sprintf("%s", configFileContents)
	data["success"] = true

	return c.RenderJSON(data)
}

func (c MfApi) ServerConfigUpdate() revel.Result {
	var jsonData map[string]interface{}
	c.Params.BindJSON(&jsonData)

	// c.Validation.Required(config)
	// c.Validation.MaxSize(config, 5000)
	// c.Validation.MinSize(config, 100)
	config := jsonData["config"]

	fmt.Println(config)

	//https://www.devdungeon.com/content/working-files-go
	file, err := os.OpenFile(
		serverConfigFile,
		os.O_WRONLY|os.O_TRUNC|os.O_CREATE,
		0666,
	)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	// Write bytes to file
	_, err = file.WriteString(config.(string))
	if err != nil {
		log.Fatal(err)
	}

	return c.ServerConfigRead()
	// data := make(map[string]interface{})
	// data["success"] = true

	// return c.RenderJSON(data)
}
