package main

import (
	"cloud.google.com/go/firestore"
	"golang.org/x/net/context"
)

// ResultRepositoryFirestore is a repository for results
type ResultRepositoryFirestore struct {
	firestoreClient *firestore.Client
}

// Save saves result
func (repo *ResultRepositoryFirestore) Save(ID string, result Result) error {
	docRef := repo.firestoreClient.Collection("result_abcd").Doc(ID)
	_, err := docRef.Create(context.Background(), result)
	return err
}
