package main

import (
	"log"
	"time"
)

// StatsSummary is summary of stats
type StatsSummary struct {
	UpdatedTime time.Time
	SampleCount int
  DecileHistogram []int
	AveragePercent int
}

var statsSummaryCached = StatsSummary{}

// GetStatsSummary returns stats summary
func GetStatsSummary() (StatsSummary, error) {
	if time.Now().Sub(statsSummaryCached.UpdatedTime) > 1 * time.Hour {
		err := updateStatsSummary(&statsSummaryCached)
		if err != nil {
			log.Println(err)
			return StatsSummary{}, err
		}
	}

	return statsSummaryCached, nil
}

// ReloadStatsSummaryCache reloads the cache
func ReloadStatsSummaryCache() error {
	return updateStatsSummary(&statsSummaryCached)
}

func updateStatsSummary(summary *StatsSummary) error {
	log.Println("Updating stats summary")
	statsEntries, err := GetAllStatsEntries()
	if err != nil {
		log.Printf("Updating stats summary failed: %v", err)
		return err
	}

	summary.UpdatedTime = time.Now()
	summary.DecileHistogram = getDecilehistogramForEntries(statsEntries)
	summary.AveragePercent = getAveragePercentForEntries(statsEntries)
	summary.SampleCount = len(statsEntries)

	log.Printf("Done updating stats summary: %+v", summary)
	return nil
}


func getDecilehistogramForEntries(entries []StatsEntry) []int {
	decileHistogram := []int{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0}
	for _, entry := range entries {
		decile := 10 * entry.CorrectCount / entry.TotalCount
    decileHistogram[decile]++
	}
  return decileHistogram
}

func getAveragePercentForEntries(entries []StatsEntry) int {
	if len(entries) == 0 {
		return 0
	}

	totalSum := 0
	for _, entry := range entries {
		totalSum += 100 * entry.CorrectCount / entry.TotalCount
	}
  return totalSum / len(entries)
}
