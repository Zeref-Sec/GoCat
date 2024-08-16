package main

import (
	"flag"
	"fmt"
	"os"
)

type Config struct {
	IsServer  bool
	Port      int
	CertFile  string
	KeyFile   string
	CAFile    string
	Protocol  string
	Host      string
	IsReverse bool
}

func main() {
	config := Config{}

	flag.BoolVar(&config.IsServer, "server", false, "Run in server mode")
	flag.IntVar(&config.Port, "port", 8080, "Port to use")
	flag.StringVar(&config.CertFile, "cert", "server.crt", "Certificate file for HTTPS")
	flag.StringVar(&config.KeyFile, "key", "server.key", "Key file for HTTPS")
	flag.StringVar(&config.CAFile, "ca", "ca.crt", "CA certificate file for HTTPS client")
	flag.StringVar(&config.Protocol, "proto", "https", "Protocol to use (https, http, or tcp)")
	flag.StringVar(&config.Host, "host", "localhost", "Host to connect to")
	flag.BoolVar(&config.IsReverse, "reverse", false, "Use reverse connection (server connects to client)")

	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage of %s:\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "Standard mode:\n")
		fmt.Fprintf(os.Stderr, "  Server: %s -server -port <port> -proto <protocol> [-cert <cert_file> -key <key_file>]\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "  Client: %s -port <port> -proto <protocol> -host <host> [-ca <ca_file>]\n\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "Reverse connection mode:\n")
		fmt.Fprintf(os.Stderr, "  Listener (Client): %s -reverse -port <port>\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "  Connector (Server): %s -server -reverse -host <host> -port <port>\n\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "Options:\n")
		flag.PrintDefaults()
		fmt.Fprintf(os.Stderr, "\nWARNING: The reverse connection feature (-reverse) bypasses typical firewall rules and can be easily misused.\n")
		fmt.Fprintf(os.Stderr, "Use it only when absolutely necessary, in controlled environments, and with proper authorization.\n")
	}

	flag.Parse()

	if len(os.Args) == 1 {
		flag.Usage()
		os.Exit(1)
	}

	if config.IsReverse {
		if config.IsServer {
			runReverseServer(&config)
		} else {
			runReverseClient(&config)
		}
	} else if config.IsServer {
		runServer(&config)
	} else {
		runClient(&config)
	}
}
