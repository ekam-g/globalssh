# Welcome To Global SSH! 🎉️

## Why Should I Use It?

Global SSH offers an easy-to-use alternative to SSH, with a setup time of less than five minutes. Unlike SSH, which requires port forwarding and presents security risks, complexity, maintenance overhead, limited scalability, and problems with dynamic IPs, Global SSH uses a Redis pub/sub model to connect a host and client across servers, eliminating the need for port forwarding and network management.

SSH can be difficult to set up and maintain, requiring software installation and configuration on both local and remote systems, as well as firewall and access control configuration. Global SSH removes this complexity, making it accessible to non-technical users.

Moreover, Global SSH is more scalable than SSH, particularly in larger networks. With SSH, configuring and managing port forwarding settings becomes increasingly challenging as the network grows. Any changes or updates to port forwarding settings must be made on each device individually, which can be time-consuming and error-prone. In contrast, Global SSH requires only a Redis database to handle connections and can be easily scaled by creating new free accounts.

In summary, Global SSH offers an easier, more secure, and scalable alternative to SSH, with a simple setup process and minimal maintenance overhead.

## Features 🚀️

1. With the key file (current named redis_key.json) you can set it to the home directory (~) of any Unix computer to instatly set up a host or client
2. Not resource intensive and even work in poor wifi conditions
3. Extermely easy to use with a very low learning curve
4. With Global SSH set up you can easily switch between being a server or client by doing
   ````
   //server mode 
   $ global_ssh server
   //client mode
   $ global_ssh client
   ````
5. Does not require port forwarding to use across networks(henice the global in the name)
6. Client mode works on all platforms, including chips like 386 and arm
7. Server mode works on all unix based systems(so not windows)


## Setup
