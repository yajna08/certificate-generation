package main

import (
	"fmt"
	"os"
)

func main() {

	// domain
	// number of peers
	// number of users

	dirs := []string{
		"org1.example.com/ca/admincerts/",
		"org1.example.com/ca/cacerts/",
		"org1.example.com/ca/tlscacerts",
		"org1.example.com/msp",
		"org1.example.com/peers/peer0.org1.example.com/msp/admincerts",
		"org1.example.com/peers/peer0.org1.example.com/msp/cacerts",
		"org1.example.com/peers/peer0.org1.example.com/msp/keystore",
		"org1.example.com/peers/peer0.org1.example.com/msp/signcerts",
		"org1.example.com/peers/peer0.org1.example.com/msp/tlscacerts",
		"org1.example.com/peers/peer0.org1.example.com/tls/",
		"org1.example.com/peers/peer1.org1.example.com/msp/admincerts",
		"org1.example.com/peers/peer1.org1.example.com/msp/cacerts",
		"org1.example.com/peers/peer1.org1.example.com/msp/keystore",
		"org1.example.com/peers/peer1.org1.example.com/msp/signcerts",
		"org1.example.com/peers/peer1.org1.example.com/msp/tlscacerts",
		"org1.example.com/peers/peer1.org1.example.com/tls/",
		"org1.example.com/tlsca",
		"org1.example.com/users/Admin@org1.example.com/msp/admincerts",
		"org1.example.com/users/Admin@org1.example.com/msp/cacerts",
		"org1.example.com/users/Admin@org1.example.com/msp/keystore",
		"org1.example.com/users/Admin@org1.example.com/msp/signcerts",
		"org1.example.com/users/Admin@org1.example.com/msp/tlscacerts",
		"org1.example.com/users/Admin@org1.example.com/tls",
		"org1.example.com/users/User1@org1.example.com/msp/admincerts",
		"org1.example.com/users/User1@org1.example.com/msp/cacerts",
		"org1.example.com/users/User1@org1.example.com/msp/keystore",
		"org1.example.com/users/User1@org1.example.com/msp/signcerts",
		"org1.example.com/users/User1@org1.example.com/msp/tlscacerts",
		"org1.example.com/users/User1@org1.example.com/tls",
		"org1.example.com/users/User2@org1.example.com/msp/admincerts",
		"org1.example.com/users/User2@org1.example.com/msp/cacerts",
		"org1.example.com/users/User2@org1.example.com/msp/keystore",
		"org1.example.com/users/User2@org1.example.com/msp/signcerts",
		"org1.example.com/users/User2@org1.example.com/msp/tlscacerts",
		"org1.example.com/users/User2@org1.example.com/tls",
	}

	for _, v := range dirs {
		err := os.MkdirAll(v, os.ModePerm)
		if err != nil {
			fmt.Print("Directory Created")
		}
	}

	fmt.Println("Dir structure created")
}
