package main

import (
	"fmt"
	"log"

	"cloud.google.com/go/firestore"
	firebase "firebase.google.com/go/v4"
	"golang.org/x/net/context"
)

var firebaseClient = initializeFirebase()
var firestoreClient = initializeFirestore()

func initializeFirebase() *firebase.App {
	app, err := firebase.NewApp(context.Background(), nil)
	if err != nil {
			log.Fatalf("error initializing app: %v\n", err)
	}

	return app
}

func initializeFirestore() *firestore.Client {
	client, err := firebaseClient.Firestore(context.Background())
	if err != nil {
			log.Fatalf("error initializing firestore: %v\n", err)
	}
	return client
}

// GetFirestoreCollectionRef returns ref to the collection relative to this quiz
func GetFirestoreCollectionRef(name string) *firestore.CollectionRef {
	collName := fmt.Sprintf("/quiz/%s/%s", quiz.ID, name)
	return firestoreClient.Collection(collName)
}
