package main

import (
	"log"

	"cloud.google.com/go/firestore"
	"golang.org/x/net/context"
	"gopkg.in/validator.v2"
)

// StatsEntry — unit of statistics
type StatsEntry struct {
		CorrectCount uint `json:"correctCount" validate:"min=0"`
		TotalCount uint `json:"totalCount" validate:"min=0"`
}

// Validate validates
func (entry StatsEntry) Validate() error {
	return validator.Validate(entry);
}

// WriteStats writes the stats
func WriteStats(firestoreClient *firestore.Client, quiz Quiz, result Result) error {
	entry := getEntryForResult(quiz, result)
	docRef := firestoreClient.Collection("stats_entry_abcd").NewDoc()
	_, err := docRef.Create(context.Background(), entry)
	if err != nil {
		return err
	}
	return ReloadStatsSummaryCache()
}

// GetAllStatsEntries returns all stats entries
func GetAllStatsEntries() ([]StatsEntry, error) {
	collRef := firestoreClient.Collection("stats_entry_abcd")
	snapshots, err := collRef.Documents(context.Background()).GetAll()
	if err != nil {
		return []StatsEntry{}, err
	}

	results := make([]StatsEntry, 0, len(snapshots))

	for _, snapshot := range snapshots {
		observer, err := entryFromSnapshot(snapshot)
		if err != nil {
			log.Printf("Invalid event observer fetched: %v", err)
		} else {
			results = append(results, observer)
		}
	}
	return results, nil
}

func getEntryForResult(quiz Quiz, result Result) StatsEntry {
	TotalCount := uint8(len(quiz.Questions))
	CorrectCount := uint8(0)

	for i, question := range quiz.Questions {
		if i >= len(result.Answers) {
			break
		}

		answer := result.Answers[i]
		if answer == question.CorrectNo {
			CorrectCount++
		}
	}

	return StatsEntry{
		CorrectCount,
		TotalCount,
	}
}

func entryFromSnapshot(snapshot *firestore.DocumentSnapshot) (StatsEntry, error) {
	var entry StatsEntry
	err := snapshot.DataTo(&entry)
	if err != nil {
		return StatsEntry{}, err
	}
	return entry, nil
}
