package main

import (
	"context"
	"fmt"
	"log"

	vault "github.com/hashicorp/vault/api"
)

func main() {
	config := vault.DefaultConfig()
	config.Address = "http://127.0.0.1:8200"
	client, err := vault.NewClient(config)
	if err != nil {
		log.Fatalf("unable to initialize Vault client: %v", err)
	}

	client.SetToken("dev-only-token")

	secretData := map[string]interface{}{
		"pg_host":     "localhost",
		"pg_port":     "5532",
		"pg_user":     "footballiot",
		"pg_password": "footballiot",
		"pg_dbname":   "footballiot",

		"mongo_host":     "localhost",
		"mongo_port":     "27018",
		"mongo_user":     "footballiot",
		"mongo_password": "footballiot",
	}
	fmt.Println(secretData)
	_, err = client.KVv2("secret").Put(context.Background(), "football-iot-secret", secretData)
	if err != nil {
		log.Fatalf("unable to write secret: %v", err)
	}

	fmt.Println("Secret written successfully.")

	// Read a secret
	secret, err := client.KVv2("secret").Get(context.Background(), "football-iot-secret")
	if err != nil {
		log.Fatalf("unable to read secret: %v", err)
	}

	value, _ := secret.Data["pg_port"].(string)
	fmt.Println(value)

}
