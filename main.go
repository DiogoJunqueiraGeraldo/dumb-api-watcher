package main

import (
	"bufio"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"time"
)

const version float32 = 1.18
const sitesSourceFile = "sites.txt"
const logsTargetFile = "logs.txt"

func main() {
	showIntroduction()

	for {
		showMenu()

		command := readCommand()

		executeBusinessRules(command)
	}
}

func showIntroduction() {
	fmt.Println("\n\n\nWelcome to the simplest program I have ever wrote!")
	fmt.Printf("This program is running on go version %.2f\n\n", version)
}

func showMenu() {
	fmt.Println("=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=")
	fmt.Println("Choose an option:")
	fmt.Println("1 - Start monitoring")
	fmt.Println("2 - Show logs")
	fmt.Println("9 - Exit Program")
	fmt.Println("=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=")
}

func readCommand() int {
	fmt.Print(">>> ")
	var command int
	fmt.Scan(&command)
	return command
}

func executeBusinessRules(command int) {
	switch command {
	case 1:
		initMonitoring()
	case 2:
		showLogs()
	case 9:
		exit()
	default:
		errAndDie("Command not found...")
	}
}

func showLogs() {
	for _, val := range readListFromFile(logsTargetFile) {
		fmt.Println(val)
	}
}

func initMonitoring() {
	sites := readListFromFile(sitesSourceFile)

	for i, site := range sites {
		resp, err := http.Get(site)

		if err != nil {
			errAndDie(fmt.Sprintf("%s", err))
		}

		pos := i + 1
		size := len(sites)
		code := resp.StatusCode
		status := resp.Status
		showRequestStatus(
			pos,
			size,
			site,
			code,
			status,
		)
		registerLog(site, code)
	}
}

func appendLineOnFile(target string, line string) {
	file, err := os.OpenFile(target, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0666)

	if err != nil {
		errAndDie(fmt.Sprintf("%s", err))
	}

	file.WriteString(line)
	file.Close()
}

func registerLog(site string, code int) {
	now := time.Now().Format("02/01/2006 15:04:05")
	log := fmt.Sprintf("[%s] %d - %s\n", now, code, site)
	appendLineOnFile(logsTargetFile, log)
}

func showRequestStatus(pos int, size int, site string, code int, status string) {
	fmt.Println("- - - - - - - - - - - - - - - - - - - - - -")
	fmt.Println("Position [", pos, "/", size, "]")
	fmt.Println("Site     [", site, "]")
	fmt.Println("Code     [", code, "]")
	fmt.Println("Status   [", status, "]")
}

func readListFromFile(target string) []string {
	var sites []string

	file, err := os.OpenFile(target, os.O_RDWR|os.O_CREATE, 0666)
	if err != nil {
		errAndDie(fmt.Sprintf("%s", err))
	}

	reader := bufio.NewReader(file)

	for {
		line, err := reader.ReadString('\n')
		if err == io.EOF {
			break
		}

		if err != nil {
			errAndDie(fmt.Sprintf("%s", err))
		}

		site := strings.TrimSpace(line)
		sites = append(sites, site)

	}

	file.Close()
	return sites
}

func exit() {
	fmt.Println("Exiting program...")
	os.Exit(0)
}

func errAndDie(msg string) {
	fmt.Println(msg)
	os.Exit(-1)
}
