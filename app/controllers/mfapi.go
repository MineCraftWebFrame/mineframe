package controllers

import (
	"fmt"
	"log"
	"os/exec"

	"github.com/revel/revel"
)

type MfApi struct {
	*revel.Controller
}

func (c MfApi) ServicetStatus() revel.Result {
	data := make(map[string]interface{})
	data["ServerStatus"] = "Stopped"
	data["success"] = true

	return c.RenderJSON(data)
}

func (c MfApi) ServiceStart() revel.Result {
	//https://golang.org/pkg/os/exec/

	cmdPrep := "java -jar spigot-1.12.2.jar"
	cmdOutput := exec.Command("bash", "-c", cmdPrep)
	cmdOutput.Dir = "/home/miner/server"
	stdoutStderr, err := cmdOutput.CombinedOutput()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("%s\n", stdoutStderr)

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
