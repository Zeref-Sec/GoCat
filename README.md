# GoCat: Multi-Protocol Remote Administration Tool

## Overview

GoCat is a Go-based tool providing a flexible, multi-protocol solution for remote system administration. It supports both standard and reverse connection modes, making it versatile for various network configurations and firewall scenarios.

**IMPORTANT:** GoCat is intended for educational purposes and authorized system administration only. Misuse of this tool may violate laws and organizational policies. Always ensure you have explicit permission before using this on any system or network.

## Features

- Supports multiple protocols: HTTPS, HTTP, and TCP
- Standard and reverse connection modes
- Cross-platform compatibility (Windows, Linux, macOS)
- Interactive command-line interface for clients
- Flexible configuration through command-line flags

## Installation

### Prerequisites

- Go 1.15 or higher

### Building from Source

1. Clone the repository or download the source files (`main.go`, `server.go`, `client.go`).

2. Open a terminal and navigate to the directory containing the source files.

3. Build for your current platform:
   ```
   go build -o gocat
   ```

4. (Optional) Cross-compile for Windows and Linux:
   ```
   GOOS=windows GOARCH=amd64 go build -o gocat_windows_amd64.exe && GOOS=linux GOARCH=amd64 go build -o gocat_linux_amd64
   ```

## Usage

### Standard Mode

#### Server

```
./gocat -server -port <port> -proto <protocol> [-cert <cert_file> -key <key_file>]
```

Example (HTTPS):
```
./gocat -server -port 8443 -proto https -cert server.crt -key server.key
```

#### Client

```
./gocat -port <port> -proto <protocol> -host <host> [-ca <ca_file>]
```

Example (HTTPS):
```
./gocat -port 8443 -proto https -host example.com -ca ca.crt
```

### Reverse Connection Mode

#### Listener (Client)

```
./gocat -reverse -port <port>
```

Example:
```
./gocat -reverse -port 8080
```

#### Connector (Server)

```
./gocat -server -reverse -host <host> -port <port>
```

Example:
```
./gocat -server -reverse -host client_ip_address -port 8080
```

## Security Considerations

1. **Authorization:** Only use GoCat on systems and networks where you have explicit permission.
2. **Encryption:** Always use HTTPS when possible to ensure encrypted communication.
3. **Authentication:** Implement proper authentication mechanisms before deploying in any non-controlled environment.
4. **Firewalls:** The reverse connection feature can bypass firewalls. Use with extreme caution and only when necessary.
5. **Logging:** Consider implementing logging features for auditing purposes.
6. **Limited Access:** Restrict GoCat's capabilities to only what is necessary for the intended tasks.

## Disclaimer

GoCat is provided "as is" without warranty of any kind. The authors are not responsible for any misuse or damage caused by this tool. Use at your own risk and responsibility.

## Contributing

Contributions to improve GoCat's security, functionality, or documentation are welcome. Please submit pull requests or open issues on the project's repository.
