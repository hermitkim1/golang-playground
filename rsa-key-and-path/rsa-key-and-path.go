package main

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
)

func Exists(filename string) bool {
	info, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return false
	}
	return !info.IsDir()
}

// var publicKeys = make(map[string]string)

func main() {
	ex, err := os.Executable()
	if err != nil {
		panic(err)
	}
	exePath := filepath.Dir(ex)
	fmt.Printf("exePath: %v\n", exePath)

	// Load cb-log config from the current directory (usually for the production)
	secretPath := filepath.Join(exePath, "secret")

	hostID := "aaiai123-123-11-aas"
	rsaPrivateKeyPassword := ""

	privateKeyFile := hostID + ".pem"
	privateKeyPath := filepath.Join(secretPath, privateKeyFile)

	publicKeyFile := hostID + ".pub"
	publicKeyPath := filepath.Join(secretPath, publicKeyFile)

	var privateKey *rsa.PrivateKey
	var publicKey *rsa.PublicKey

	if !Exists(privateKeyPath) {
		fmt.Println("RSA key doesn't exist")

		// Create directory or folder if not exist
		_, err := os.Stat(secretPath)

		if os.IsNotExist(err) {
			errDir := os.MkdirAll(secretPath, 0755)
			if errDir != nil {
				log.Fatal(err)
			}

		}

		// generate key
		privateKey, err = rsa.GenerateKey(rand.Reader, 2048)
		if err != nil {
			fmt.Printf("Cannot generate RSA key\n")
			os.Exit(1)
		}
		publicKey = &privateKey.PublicKey

		// dump private key to file
		var privateKeyBytes []byte = x509.MarshalPKCS1PrivateKey(privateKey)
		privateKeyBlock := &pem.Block{
			Type:  "RSA PRIVATE KEY",
			Bytes: privateKeyBytes,
		}
		privatePem, err := os.Create(privateKeyPath)
		if err != nil {
			fmt.Printf("error when create private.pem: %s \n", err)
			os.Exit(1)
		}
		err = pem.Encode(privatePem, privateKeyBlock)
		if err != nil {
			fmt.Printf("error when encode private pem: %s \n", err)
			os.Exit(1)
		}

		// dump public key to file
		publicKeyBytes, err := x509.MarshalPKIXPublicKey(publicKey)
		if err != nil {
			fmt.Printf("error when dumping publickey: %s \n", err)
			os.Exit(1)
		}
		publicKeyBlock := &pem.Block{
			Type:  "RSA PUBLIC KEY",
			Bytes: publicKeyBytes,
		}
		publicPem, err := os.Create(publicKeyPath)
		if err != nil {
			fmt.Printf("error when create public.pem: %s \n", err)
			os.Exit(1)
		}

		err = pem.Encode(publicPem, publicKeyBlock)
		if err != nil {
			fmt.Printf("error when encode public pem: %s \n", err)
			os.Exit(1)
		}
	} else {
		fmt.Println("RSA key exists")

		privateKeyBytes, err := ioutil.ReadFile(privateKeyPath)
		if err != nil {
			fmt.Printf("No RSA private key found, generating temp one: %s \n", err)
		}

		privateKeyPem, _ := pem.Decode(privateKeyBytes)
		var privPemBytes []byte
		if privateKeyPem.Type != "RSA PRIVATE KEY" {
			fmt.Println("RSA private key is of the wrong type")
		}

		if rsaPrivateKeyPassword != "" {
			privPemBytes, err = x509.DecryptPEMBlock(privateKeyPem, []byte(rsaPrivateKeyPassword))
		} else {
			privPemBytes = privateKeyPem.Bytes
		}

		var parsedKey interface{}
		if parsedKey, err = x509.ParsePKCS1PrivateKey(privPemBytes); err != nil {
			if parsedKey, err = x509.ParsePKCS8PrivateKey(privPemBytes); err != nil {
				fmt.Printf("Unable to parse RSA private key, generating a temp one: %s \n", err)
			}
		}
		fmt.Println("parsedKey")
		fmt.Printf("%#v\n", parsedKey)

		var ok bool
		privateKey, ok = parsedKey.(*rsa.PrivateKey)
		if !ok {
			fmt.Printf("Unable to parse RSA private key, generating a temp one: %s \n", err)
		}

		publicKeyBytes, err := ioutil.ReadFile(publicKeyPath)
		fmt.Println("publicKeyBytes")
		fmt.Printf("%#v\n", publicKeyBytes)
		if err != nil {
			fmt.Printf("No RSA public key found, generating temp one: %s \n", err)
		}
		publicKeyPem, _ := pem.Decode(publicKeyBytes)
		fmt.Println("publicKeyPem")
		fmt.Printf("%#v\n", publicKeyPem)
		if publicKeyPem == nil {
			fmt.Printf("Use `ssh-keygen -f id_rsa.pub -e -m pem > id_rsa.pem` to generate the pem encoding of your RSA public key: %s \n",
				errors.New("rsa public key not in pem format"))
		}

		if publicKeyPem.Type != "RSA PUBLIC KEY" {
			fmt.Println("RSA public key is of the wrong type")
		}

		if parsedKey, err = x509.ParsePKIXPublicKey(publicKeyPem.Bytes); err != nil {
			fmt.Printf("Unable to parse RSA public key, generating a temp one: %s \n", err)
		}
		fmt.Println("parsedKey")
		fmt.Printf("%#v\n", parsedKey)

		if publicKey, ok = parsedKey.(*rsa.PublicKey); !ok {
			fmt.Printf("Unable to parse RSA public key, generating a temp one: %s \n", err)
		}

		privateKey.PublicKey = *publicKey

	}

	fmt.Println("Private Key")
	fmt.Printf("%#v\n", privateKey)
	fmt.Println("Public Key")
	fmt.Printf("%#v\n", publicKey)
}
