package main

import (
	"bufio"
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
)

func handleInstructions() {
	fmt.Println("1 - 👀 Start monitoring")
	fmt.Println("2 - ❓ Display logs")
	fmt.Println("0 - 🏃 Exit")
	fmt.Println("")

}

func handleCommand() int {
	var command int
	fmt.Scan(&command)
	return command
}

func handleErrorsProgram(code int) {
	os.Exit(code)
}

func handleFetchUrl(url string) {
	response, error := http.Get(url)

	if error == nil {
		if response.StatusCode == 200 || response.StatusCode == 204 {
			fmt.Println("✅ url:", url, "status:", response.StatusCode)
			handleRegisterLogs(url, true)
		} else {
			fmt.Println("❌ url:", url, "error:", response.StatusCode)
			handleRegisterLogs(url, false)
		}
	}
}

func handleReadFile(fileName string) []string {
	var urls []string
	file, err := os.Open(fileName)
	reader := bufio.NewReader(file)

	if err != nil {
		fmt.Println("❌ something is wrong:", err)
	}

	for {
		liner, err := reader.ReadString('\n')
		urls = append(urls, strings.TrimSpace(liner))
		if err == io.EOF {
			break
		}
	}

	file.Close()

	return urls
}

func handleRegisterLogs(url string, status bool) {
	CURRENT_DATE_HOUR := time.Now().Format("2006-01-02 15:04:05")
	file, err := os.OpenFile("logs.txt", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)

	if err != nil {
		fmt.Println("❌ something is wrong:", err)
	}

	file.WriteString(CURRENT_DATE_HOUR + " - " + "url: " + url + " - " + strconv.FormatBool(status) + "\n")
	file.Close()
}

func handleLogs() {
	fmt.Println("❓ logs!")
	logs := handleReadFile("logs.txt")
	for _, log := range logs {
		fmt.Println(log)
	}
}

func handleMonitoring() {
	fmt.Println("👀 monitoring...")
	urls := handleReadFile("urls.txt")
	const DELAY = 1 * time.Second
	const MONITORING = 10

	for i := 0; i < MONITORING; i++ {
		for _, url := range urls {
			handleFetchUrl(url)
		}
		time.Sleep(DELAY)
	}

	fmt.Println("")

}

func main() {

	for {
		handleInstructions()
		command := handleCommand()

		switch command {
		case 0:
			handleErrorsProgram(0)
		case 1:
			handleMonitoring()
		case 2:
			handleLogs()
		default:
			fmt.Println("Desconhecido")
			handleErrorsProgram(-1)
		}
	}

}
