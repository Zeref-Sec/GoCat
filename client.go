package main

import (
	"bufio"
	"crypto/tls"
	"crypto/x509"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"os"
	"strings"
)

type Command struct {
	Cmd  string   `json:"cmd"`
	Args []string `json:"args"`
}

type Response struct {
	Output string `json:"output"`
	Error  string `json:"error"`
}

func runClient(config *Config) {
	switch config.Protocol {
	case "https", "http":
		runHTTPClient(config)
	case "tcp":
		runTCPClient(config)
	default:
		log.Fatalf("Unsupported protocol: %s", config.Protocol)
	}
}

func runHTTPClient(config *Config) {
	var client *http.Client
	var protocol string

	if config.Protocol == "https" {
		caCert, err := ioutil.ReadFile(config.CAFile)
		if err != nil {
			log.Fatal(err)
		}
		caCertPool := x509.NewCertPool()
		caCertPool.AppendCertsFromPEM(caCert)

		tlsConfig := &tls.Config{
			RootCAs: caCertPool,
		}
		transport := &http.Transport{TLSClientConfig: tlsConfig}
		client = &http.Client{Transport: transport}
		protocol = "https"
	} else {
		client = &http.Client{}
		protocol = "http"
	}

	fmt.Printf("Connected to %s://%s:%d\n", protocol, config.Host, config.Port)
	fmt.Println("Type 'exit' to quit the session.")

	reader := bufio.NewReader(os.Stdin)
	for {
		fmt.Print("> ")
		input, _ := reader.ReadString('\n')
		input = strings.TrimSpace(input)

		if input == "exit" {
			fmt.Println("Exiting session.")
			return
		}

		parts := strings.SplitN(input, " ", 2)
		cmd := Command{Cmd: parts[0]}
		if len(parts) > 1 {
			cmd.Args = strings.Split(parts[1], " ")
		}

		jsonData, err := json.Marshal(cmd)
		if err != nil {
			fmt.Println("Error:", err)
			continue
		}

		resp, err := client.Post(fmt.Sprintf("%s://%s:%d/execute", protocol, config.Host, config.Port), "application/json", strings.NewReader(string(jsonData)))
		if err != nil {
			fmt.Println("Error:", err)
			continue
		}

		var response Response
		err = json.NewDecoder(resp.Body).Decode(&response)
		resp.Body.Close()
		if err != nil {
			fmt.Println("Error:", err)
			continue
		}

		if response.Error != "" {
			fmt.Println("Error:", response.Error)
		}
		fmt.Print(response.Output)
	}
}

func runTCPClient(config *Config) {
	conn, err := net.Dial("tcp", fmt.Sprintf("%s:%d", config.Host, config.Port))
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	go func() {
		scanner := bufio.NewScanner(conn)
		for scanner.Scan() {
			fmt.Println(scanner.Text())
		}
	}()

	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		fmt.Fprintf(conn, "%s\n", scanner.Text())
	}
}

func runReverseClient(config *Config) {
	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", config.Port))
	if err != nil {
		log.Fatal(err)
	}
	defer listener.Close()

	log.Printf("Listening on :%d", config.Port)

	conn, err := listener.Accept()
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	log.Printf("Server connected from %s", conn.RemoteAddr())

	go io.Copy(os.Stdout, conn)

	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		fmt.Fprintf(conn, "%s\n", scanner.Text())
	}
}
