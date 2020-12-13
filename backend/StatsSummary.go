package main

import (
	"log"
	"time"
)

// StatsSummary is summary of stats
type StatsSummary struct {
	UpdatedTime time.Time
	SampleCount uint
  DecileHistogram []uint
  AveragePercent uint
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
	statsEntries, err := GetAllStatsEntries()
	if err != nil {
		return err
	}

	summary.UpdatedTime = time.Now()
	summary.DecileHistogram = getDecilehistogramForEntries(statsEntries)
	summary.AveragePercent = getAveragePercentForEntries(statsEntries)
	return nil
}


func getDecilehistogramForEntries(entries []StatsEntry) []uint {
	decileHistogram := []uint{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0}
	for _, entry := range entries {
		decile := 10 * entry.CorrectCount / entry.TotalCount
    decileHistogram[decile]++
	}
  return decileHistogram
}

func getAveragePercentForEntries(entries []StatsEntry) uint {
	totalSum := uint(0)
	for _, entry := range entries {
		totalSum += 100 * entry.CorrectCount / entry.TotalCount
	}
  return totalSum / uint(len(entries))
}
