package main

import (
	"fmt"

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

// GetByID fetches by id
func (repo *ResultRepositoryFirestore) GetByID(ID string) (Result, error) {
	docRef := repo.firestoreClient.Collection("result_abcd").Doc(ID)
	snapshot, err := docRef.Get(context.Background())
	if err != nil {
		return Result{}, err
	}
	if snapshot.Exists() != true {
		return Result{}, fmt.Errorf("Cannot find result with ID %s", ID)
	}
	var result Result
	err = snapshot.DataTo(&result)
	return result, err
}
