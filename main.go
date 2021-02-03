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

// Build Version
const version float32 = 1.18

// Files Paths
const (
	sitesSourceFile = "sites.txt"
	logsTargetFile  = "logs.txt"
)

// Commands
const (
	monitoringCMD = 1
	showLogsCMD   = 2
	exitCMD       = 9
)

func main() {
	showIntroduction()

	for {
		showMenu()

		command := readCommand()

		handleCommand(command)
	}
}

// Flow Controller

func handleCommand(command int) {
	switch command {

	case monitoringCMD:
		initMonitoring()
		break

	case showLogsCMD:
		showLogs()
		break

	case exitCMD:
		exit()
		break

	default:
		errAndDie("Command not found...")
		break

	}
}

// Exit Process Handlers

func exit() {
	showExitMessage("Exiting Program!!!")
	os.Exit(0)
}

func errAndDie(msg string) {
	showExitMessage(msg)
	os.Exit(-1)
}

// User Iteraction

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

func showRequestStatus(pos int, size int, site string, code int, status string) {
	fmt.Println("- - - - - - - - - - - - - - - - - - - - - -")
	fmt.Println("Position [", pos, "/", size, "]")
	fmt.Println("Site     [", site, "]")
	fmt.Println("Code     [", code, "]")
	fmt.Println("Status   [", status, "]")
}

func showExitMessage(msg string) {
	fmt.Println(msg)
	fmt.Println("Thanks for your time!!!")
}

func showLogs() {
	for _, val := range readListFromFile(logsTargetFile) {
		fmt.Println(val)
	}
}

func readCommand() int {
	fmt.Print(">>> ")
	var command int
	fmt.Scan(&command)
	return command
}

// Http

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

// File IO

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
