/*
Global SSH offers an easier more secure, and scalable alternative to SSH,
with a simple setup process and minimal maintenance overhead.

Usage:
	globalssh [mode]

The Modes are:
	client:
		connect to the server listed in your redis_key.json
	server:
		creates a server with the name listed in the redis_key.json

Features:

1. Easy Setup: Set up Global SSH in less than 5 minutes by placing the key file (redis_key.json) in
   the home directory (~) of any Unix computer. This streamlined process allows for quick configuration of hosts or clients.
2. Network Scalability: Global SSH eliminates the need for port forwarding, making server-client connections across networks effortless and efficient. Say goodbye to the complexities of managing port forwarding settings as your network grows.
3. Low Resource Consumption: Global SSH is designed to be lightweight and not resource-intensive, ensuring optimal performance even in poor Wi-Fi conditions.
4. User-Friendly Interface: We have prioritized creating an intuitive and user-friendly experience with Global SSH. The platform offers a minimal learning curve, enabling users of all technical expertise levels to navigate and utilize it effectively.
5. Seamless Mode Switching: Transitioning between server and client modes is a breeze with Global SSH. By issuing simple commands in the terminal, you can easily switch between being a server or a client:

   * Server Mode: $ `globalssh server`
   * Client Mode: $ `globalssh client`
6. No Port Forwarding Needed: Global SSH removes the requirement for port forwarding, simplifying the process of connecting hosts and clients across networks. Say goodbye to the complexities of managing port forwarding settings.
7. Platform Agnostic: Global SSH's client mode is compatible with all platforms, including chips like 386 and ARM.
8. Unix Compatibility: Global SSH's server mode works seamlessly on all Unix-based systems, ensuring wide compatibility across various platforms. Please note that server mode is not supported on Windows.
9. Collaboration Feature: Global SSH enables server and client collaboration, allowing them to work together in a shared shell environment.
10. Enhanced Security:

* Global SSH provides robust security measures to safeguard your server and client connections.
* Enjoy an additional layer of anonymity as Global SSH cannot be easily discovered or traced back to a specific location.
* The Redis server used in Global SSH acts as a proxy, effectively concealing your server's IP address and making it significantly harder for potential attackers to target your system.
* AES encryption is integrated into Global SSH, ensuring the confidentiality of data transmission between the server and client.
* Global SSH's design ensures that no data is stored through the Redis server, minimizing the risk of an attack even if the Redis server is compromised.

*/

package main

import (
	"fmt"
	"globalssh/net"
	"os"
	"strings"

	"globalssh/client"
	"globalssh/server"
)

const title = "_______________       ______          ______       ______________________  __\n__  ____/___  /______ ___  /_ ______ ____  /       __  ___/__  ___/___  / / /\n_  / __  __  / _  __ \\__  __ \\_  __ `/__  /        _____ \\ _____ \\ __  /_/ / \n/ /_/ /  _  /  / /_/ /_  /_/ // /_/ / _  /         ____/ / ____/ / _  __  /  \n\\____/   /_/   \\____/ /_.___/ \\__,_/  /_/          /____/  /____/  /_/ /_/   \n                                                                             "

const help = "\nExample 'globalssh client', 'globalssh client {servername}', 'globalssh server', 'globalssh server {servername}', globalssh update'"

func main() {
	if !(len(os.Args) >= 2) {
		fmt.Println(title)
		fmt.Println("Please give an arg like client or server or update" + help)
		os.Exit(1)
	}
	switch strings.Trim(os.Args[1], " ") {
	case "client":
		if len(os.Args) == 5 {
			client.CommandSend(os.Args[2], os.Args[3], os.Args[4])
		} else if len(os.Args) == 3 {
			client.Run(os.Args[2])
		}
		client.Run("")
	case "server":
		if len(os.Args) == 3 {
			server.Start(os.Args[2])
		}
		server.Start("")
	case "update":
		net.Update()
	default:
		fmt.Println(title)
		fmt.Println("Bad Arg Given, Please Put in server or client or update" + help)
		os.Exit(1)

	}
}
