package main

import (
	"context"
	"fmt"
	"log"

	"google.golang.org/api/iterator"
	"google.golang.org/api/option"

	"cloud.google.com/go/firestore"
	"github.com/kiselev-nikolay/direct-to-me/pkg/conf"
)

func main() {
	conf := conf.ReadConfig("./conf.yaml")
	projectID := "decent-genius-311507"
	ctx := context.Background()
	client, err := firestore.NewClient(ctx, projectID, option.WithCredentialsFile(conf.Google.Application.Credentials.Storage))
	if err != nil {
		log.Fatalf("Failed to create client: %v", err)
	}
	// Close client when done with
	defer client.Close()
	_, _, err = client.Collection("users").Add(ctx, map[string]interface{}{
		"first": "Ada",
		"last":  "Lovelace",
		"born":  1815,
	})
	if err != nil {
		log.Fatalf("Failed adding alovelace: %v", err)
	}
	_, _, err = client.Collection("users").Add(ctx, map[string]interface{}{
		"first":  "Alan",
		"middle": "Mathison",
		"last":   "Turing",
		"born":   1912,
	})
	if err != nil {
		log.Fatalf("Failed adding aturing: %v", err)
	}
	iter := client.Collection("users").Documents(ctx)
	for {
		doc, err := iter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			log.Fatalf("Failed to iterate: %v", err)
		}
		fmt.Println(doc.Data())
	}
}
