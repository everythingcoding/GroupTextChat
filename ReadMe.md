install go
install npm
install wscat using npm -> npm install -g wscat
build -> go build
run -> ./grouptextchat

open 3 othe terminals and run wscat -c ws://localhost:8080/ws
send some message from each client and check if received or not by other clients
