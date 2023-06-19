# Websocks
### What is Websocks
A lightweight network obfuscation proxy based on **socks5** in **Golang**
### Usage
First clone the project in your server and local

In the server do a
```
go run cmd/server/main.go
```
It will generate a .Websocks json file and take note for the port & password the program is using.

Then in the local do
```
go run cmd/local/main.go
```
It will reports an error and follow the error to the generated .Websocks json file and paste the port (with the server's IP address) & password

Then use [SwitchyOmega](https://github.com/FelisCatus/SwitchyOmega), a Chrome extension. Adding the proxy profile based on .Websocks. After that local can take request from Chrome and send to server
