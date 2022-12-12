package main

import (
	"encoding/hex"
	"fmt"
	"github.com/status-im/doubleratchet"
	"syscall/js"
)

var web doubleratchet.Session

func initiateSession(this js.Value, input []js.Value) any {
	sk := doubleratchet.Key{
		0xeb, 0x8, 0x10, 0x7c, 0x33, 0x54, 0x0, 0x20, 0xe9, 0x4f, 0x6c, 0x84, 0xe4, 0x39, 0x50, 0x5a, 0x2f, 0x60, 0xbe, 0x81, 0xa, 0x78, 0x8b, 0xeb, 0x1e, 0x2c, 0x9, 0x8d, 0x4b, 0x4d, 0xc1, 0x40,
	}
	fmt.Println("SharedKey", sk)

	keyPair, err := doubleratchet.DefaultCrypto{}.GenerateDH()
	if err != nil {
		panic(err)
	}

	publicKey := keyPair.PublicKey()
	fmt.Println("PublicKey", publicKey)

	web, err = doubleratchet.New([]byte("web-session-id"), sk, keyPair, nil)
	if err != nil {
		panic(err)
	}

	return true
}

func initiateSessionRemote(this js.Value, input []js.Value) any {
	skString := fmt.Sprintf("%s", input[0])
	pkString := fmt.Sprintf("%s", input[1])

	pk, err := hex.DecodeString(pkString)
	if err != nil {
		panic(err)
	}

	sk, err := hex.DecodeString(skString)
	if err != nil {
		panic(err)
	}

	fmt.Printf("Decoded Byte Array: %v \n", pk)
	fmt.Printf("Decoded Byte Array: %v \n", sk)

	web, err = doubleratchet.NewWithRemoteKey([]byte("web-session-id"), sk, pk, nil)

	return true
}

func encryptMessage(this js.Value, input []js.Value) any {
	if web == nil {
		fmt.Println("Please instantiate Web client first")
		return false
	}

	messageString := fmt.Sprintf("%s", input[0])
	message, err := web.RatchetEncrypt([]byte(messageString), nil)
	if err != nil {
		panic(err)
	}

	fmt.Println(message)

	return true
}

// Some issues with returning the Message type. Will work on that.
//func decryptMessage(this js.Value, input []js.Value) any {
//	messageString := fmt.Sprintf("%s", input[0])
//	message, err := web.RatchetDecrypt(messageString, nil)
//
//	return true
//}

func registerCallbacks() {
	js.Global().Set("initiateSession", js.FuncOf(initiateSession))
	js.Global().Set("initiateSessionRemote", js.FuncOf(initiateSessionRemote))
	js.Global().Set("encryptMessage", js.FuncOf(encryptMessage))
}

func main() {
	c := make(chan struct{}, 0)

	println("WASM Go Initialized")
	// register functions
	registerCallbacks()
	<-c
}
