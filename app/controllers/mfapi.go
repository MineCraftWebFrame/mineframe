package controllers

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"github.com/revel/revel"
	"github.com/ryandrew/cmd"
)

type MfApi struct {
	*revel.Controller
}

var (
	serverDir        = "mineframe"
	spigotVersionURL = "https://raw.githubusercontent.com/MineCraftWebFrame/mineframe/master/latestSpigot.txt"
)

func getMinecraftConfigFile() string {
	return filepath.FromSlash(getMinecraftDir() + "/server.properties")
}
func getMinecraftDir() string {
	return filepath.FromSlash(getServerDir() + "/minecraft")
}
func getServerConfigFile() string {
	return filepath.FromSlash(getServerDir() + "/mineframe.json")
}
func checkIfFileExists(fileName string) bool {

	_, err := os.Stat(fileName)
	if err != nil {
		if os.IsNotExist(err) {
			return false
		}
	}
	return true
}
func checkifMinecraftServerExists(serverPath string, config map[string]interface{}) {

	if !checkIfFileExists(serverPath) {

		err := downloadFile(serverPath, config["url"].(string))
		if err != nil {
			fmt.Println("Error downloading  server jar file")
			panic(err)
		}
	}
}
func checkJavaVersion() {
	fmt.Println("Checking if Java is installed")

	cmd := exec.Command("java", "-version")

	stdoutStderr, err := cmd.CombinedOutput()
	fmt.Printf("%s\n", stdoutStderr)
	if err != nil {
		fmt.Println("Error! Java Not Installed!")
		panic(err)
	}

}
func checkifMinecraftDirExists() {
	serverDir := getMinecraftDir()

	if !checkIfFileExists(serverDir) {
		err := makeServerDir(serverDir)
		if err != nil {
			errStr := "Error making server dir"
			fmt.Println(errStr)
			fmt.Println(serverDir)
			fmt.Println(err)
			panic(errStr)
		}
	}
}
func makeServerDir(dir string) error {
	err := os.MkdirAll(dir, 0777)
	return err
}
func getLatestSpigotVersion() (url string, name string) {

	resp, err := http.Get(spigotVersionURL)
	if err != nil {
		panic("Fetch Spigot Error 1 bad HTTP request! " + err.Error())
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		panic("Fetch Spigot Error 2 bad HTTP response! " + err.Error())
	}

	responseBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		panic("Fetch Spigot Error 3 cannot read response! " + err.Error())
	}
	response := string(responseBytes)
	response = strings.TrimSpace(response)

	lastSlash := strings.LastIndex(response, "/")
	if lastSlash == -1 {
		panic("last slash not found")
	}
	lastSlash++

	fmt.Println(response)

	name = response[lastSlash:]

	fmt.Println(name)

	url = response
	return url, name
}
func downloadFile(filepath string, url string) (err error) {

	fmt.Println("Downloading File: ")
	fmt.Print("From: ")
	fmt.Println(url)
	fmt.Print("To: ")
	fmt.Println(filepath)

	// Create the file
	out, err := os.Create(filepath)
	if err != nil {
		return err
	}
	defer out.Close()

	// Get the data
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// Writer the body to file
	_, err = io.Copy(out, resp.Body)
	if err != nil {
		return err
	}

	return nil
}

func (c MfApi) ServicetStatus() revel.Result {
	data := make(map[string]interface{})
	data["ServerStatus"] = "Stopped"
	data["success"] = true

	return c.RenderJSON(data)
}

func ServiceStartGoRoutine() {

	//Launch minecraft in a goroutine
	go func() {
		pidFound := false

		minecraftDir := getMinecraftDir()
		config := readConfig()

		minecraftServerJar := config["minecraftServerJar"].(string)

		checkifMinecraftServerExists(minecraftDir+"/"+minecraftServerJar, config)

		checkJavaVersion()

		serverCmd := cmd.NewCmd("java", "-jar", minecraftServerJar)
		//statusChan :=
		serverCmd.Dir = minecraftDir
		serverCmd.Start()
		//err := serverCmd.Start()
		//if err != nil {
		//	fmt.Println("Error starting  server jar file. Java probably not installed")
		//	panic(err)
		//}

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
}

func (c MfApi) ServiceStart() revel.Result {

	checkifMinecraftDirExists()

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

	ServiceStartGoRoutine()

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
func writeConfig(config map[string]interface{}) {
	var err error

	configFile, err := os.Create(getServerConfigFile())
	if err != nil {
		panic("error writing config file 1! " + err.Error())
	}
	defer configFile.Close()

	configFileWriter := io.Writer(configFile)

	enc := json.NewEncoder(configFileWriter)
	err = enc.Encode(config)
	if err != nil {
		panic("error writing config file 2! " + err.Error())
	}
}

func readConfig() (config map[string]interface{}) {
	config = make(map[string]interface{})

	checkifMinecraftDirExists()
	if !checkIfFileExists(getServerConfigFile()) {
		fmt.Println("Config file doesn't exist")
		url, name := getLatestSpigotVersion()

		fmt.Print("name: ")
		fmt.Println(name)
		fmt.Print("url: ")
		fmt.Println(url)

		config["minecraftServerJar"] = name
		config["url"] = url

		writeConfig(config)

		return config
	}

	fmt.Println("reading Config file")

	//filename is the path to the json config file
	file, err := os.Open(getServerConfigFile())
	if err != nil {
		panic("Could not read server config file!")
	}
	decoder := json.NewDecoder(file)
	err = decoder.Decode(&config)
	if err != nil {
		panic("Could not decode server config file!")
	}
	fmt.Print("config minecraftServerJar: ")
	fmt.Println(config["minecraftServerJar"])
	fmt.Print("config url: ")
	fmt.Println(config["url"])

	return config
}

func (c MfApi) MinecraftConfigRead() revel.Result {
	checkifMinecraftDirExists()

	var serverConfigFile = getMinecraftConfigFile()

	data := make(map[string]interface{})

	configFileContents, err := ioutil.ReadFile(serverConfigFile)
	if err == nil {
		data["config"] = fmt.Sprintf("%s", configFileContents)
	} else {
		data["config"] = "Error Reading Config file!"
	}

	data["configFile"] = serverConfigFile
	data["success"] = true
	return c.RenderJSON(data)
}

func (c MfApi) MinecraftConfigUpdate() revel.Result {
	var jsonData map[string]interface{}
	c.Params.BindJSON(&jsonData)

	config := jsonData["config"]

	c.Validation.Required(config)
	c.Validation.MaxSize(config, 5000)
	c.Validation.MinSize(config, 100)

	fmt.Print("config: ")
	fmt.Println(config)

	if c.Validation.HasErrors() {
		return c.outputJsonError("config invalid")
	}

	//https://www.devdungeon.com/content/working-files-go
	file, err := os.OpenFile(
		getServerConfigFile(),
		os.O_WRONLY|os.O_TRUNC|os.O_CREATE,
		0666,
	)
	if err != nil {
		fmt.Println(err.Error())
		return c.outputJsonError("Error 1 while saving config file")
	}
	defer file.Close()

	// Write bytes to file
	_, err = file.WriteString(config.(string))
	if err != nil {
		fmt.Println(err.Error())
		return c.outputJsonError("Error 2 while saving config file")
	}

	return c.MinecraftConfigRead()
	// data := make(map[string]interface{})
	// data["success"] = true

	// return c.RenderJSON(data)
}

func (c MfApi) outputJsonError(errorString string) revel.Result {
	data := make(map[string]interface{})
	data["success"] = false
	data["error"] = errorString
	return c.RenderJSON(data)
}
