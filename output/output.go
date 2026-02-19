package output

import (
	"encoding/json"
	"os"
	"time"
)

type Handler struct {
	OutputFile string
}

type ScanInfo struct {
	StartTime      string  `json:"start_time"`
	EndTime        string  `json:"end_time"`
	DurationSeconds float64 `json:"duration_seconds"`
	TotalRequests  int     `json:"total_requests"`
	Found          int     `json:"found"`
}

type JSONOutput struct {
	ScanInfo ScanInfo               `json:"scan_info"`
	Results  []map[string]interface{} `json:"results"`
}

func NewHandler(outputFile string) *Handler {
	return &Handler{
		OutputFile: outputFile,
	}
}

func (h *Handler) Save(results []map[string]interface{}, stats interface{}) error {
	if h.OutputFile == "" {
		return nil // No output file specified
	}

	
	scanStats := stats.(struct {
		TotalRequests   int
		DiscoveredCount int
		Duration        time.Duration
	})

	output := JSONOutput{
		ScanInfo: ScanInfo{
			StartTime:      time.Now().Add(-scanStats.Duration).Format(time.RFC3339),
			EndTime:        time.Now().Format(time.RFC3339),
			DurationSeconds: scanStats.Duration.Seconds(),
			TotalRequests:  scanStats.TotalRequests,
			Found:          scanStats.DiscoveredCount,
		},
		Results: results,
	}

	file, err := os.Create(h.OutputFile)
	if err != nil {
		return err
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")
	return encoder.Encode(output)
}
