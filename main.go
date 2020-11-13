package main

import (
	"fmt"

	"github.com/giantswarm/valuemodifier/vault/decrypt"
	"github.com/giantswarm/valuemodifier/vault/encrypt"
	"github.com/hashicorp/vault/api"
)

func main() {
	var client *api.Client
	{
		config := api.DefaultConfig()
		err := config.ReadEnvironment()
		if err != nil {
			fmt.Printf("ERROR: %v\n", err)
			return
		}
		client, err = api.NewClient(config)
		if err != nil {
			fmt.Printf("ERROR: %v\n", err)
			return
		}
	}

	var encryptSvc *encrypt.Service
	{
		var err error
		config := encrypt.Config{
			VaultClient: client,
			Key:         "possum",
		}
		encryptSvc, err = encrypt.New(config)
		if err != nil {
			fmt.Printf("ERROR: %v\n", err)
			return
		}
	}

	var decryptSvc *decrypt.Service
	{
		var err error
		config := decrypt.Config{
			VaultClient: client,
			Key:         "possum",
		}
		decryptSvc, err = decrypt.New(config)
		if err != nil {
			fmt.Printf("ERROR: %v\n", err)
			return
		}
	}

	s := []byte("Hello there, general Kenobi")

	eByte, err := encryptSvc.Modify(s)
	if err != nil {
		panic(err)
	}
	fmt.Println(string(eByte))

	dByte, err := decryptSvc.Modify(eByte)
	if err != nil {
		panic(err)
	}

	fmt.Printf("input plaintext: %s\nencrypted: %s\ndecrypted: %s\n", s, string(eByte), string(dByte))
}
