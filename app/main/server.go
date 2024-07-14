package main

import (
	"flag"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"

	"github.com/dmast3r/go-forge/app/constants"
	"github.com/dmast3r/go-forge/app/handlers"
)

func main() {
	parseCommandLineArguments()

	l, err := net.Listen("tcp", ":"+constants.Port)
	if err != nil {
		log.Fatalf("Failed to bind to port %s: %v", constants.Port, err)
	}
	defer l.Close()

	log.Printf("Server is listening on port %s", constants.Port)

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		for {
			conn, err := l.Accept()
			if err != nil {
				log.Printf("Error accepting connection: %v", err)
				continue
			}
			go handleConnection(conn)
		}
	}()

	<-sigChan
	log.Println("Shutting down the server...")
}

func handleConnection(conn net.Conn) {
	defer func() {
		if err := conn.Close(); err != nil {
			log.Printf("Error closing connection: %v", err)
		}
	}()

	buffer := make([]byte, 2048)
	for {
		n, err := conn.Read(buffer)
		if err != nil {
			if err.Error() == "EOF" {
				break
			}
			log.Printf("Error reading from connection: %v", err)
			break
		}

		rawRequest := string(buffer[:n])
		log.Printf("Received request: %s", rawRequest)
		response := handlers.HandleRequest(handlers.ParseRequest(rawRequest))

		_, err = conn.Write([]byte(response))
		if err != nil {
			log.Printf("Error writing to connection: %v", err)
			break
		}
	}
}

func parseCommandLineArguments() {
	workingDirectory := flag.String("directory", "/", "the directory to serve files from")
	flag.Parse()

	setWorkingDirectory(*workingDirectory)
}

func setWorkingDirectory(directory string) {
	os.Setenv("WORKING_DIRECTORY", directory)
}
