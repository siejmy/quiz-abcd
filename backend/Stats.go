package main

import (
	"log"

	"cloud.google.com/go/firestore"
	"golang.org/x/net/context"
	"gopkg.in/validator.v2"
)

// StatsEntry — unit of statistics
type StatsEntry struct {
		CorrectCount int `json:"correctCount" validate:"min=0"`
		TotalCount int `json:"totalCount" validate:"min=0"`
}

// Validate validates
func (entry StatsEntry) Validate() error {
	return validator.Validate(entry);
}

// WriteStats writes the stats
func WriteStats(quiz Quiz, result Result) error {
	entry := GetStatsEntryForResult(quiz, result)
	docRef := GetFirestoreCollectionRef("stats_entry").NewDoc()
	_, err := docRef.Create(context.Background(), entry)
	if err != nil {
		return err
	}
	return ReloadStatsSummaryCache()
}

// GetAllStatsEntries returns all stats entries
func GetAllStatsEntries() ([]StatsEntry, error) {
	collRef := GetFirestoreCollectionRef("stats_entry")
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

// GetStatsEntryForResult gets stats for result
func GetStatsEntryForResult(quiz Quiz, result Result) StatsEntry {
	TotalCount := len(quiz.Questions)
	CorrectCount := 0

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
