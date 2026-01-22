package analyzer

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/alexanderritik/loghive/internal/domain"
)

type Service struct {
	InputCh <-chan domain.LogEntry
}

func NewService(input <-chan domain.LogEntry) *Service {
	return &Service{InputCh: input}
}

func (s *Service) Start() {
	go s.loop()
}

func (s *Service) loop() {

	var currentFile *os.File
	var currentHour int = -1
	defer func() {
		if currentFile != nil {
			currentFile.Close()
		}
	}()

	ticker := time.NewTicker(5 * time.Second)
	defer ticker.Stop()

	var errorCount int
	var totalBytes int

	for {
		select {
		case entry, ok := <-s.InputCh:
			if !ok {
				log.Println("Analyzer stopping...")
				return
			}

			now := time.Now()
			if now.Hour() != currentHour {
				if currentFile != nil {
					currentFile.Close()
				}

				filename := fmt.Sprintf("logs_%s_%02d.txt", now.Format("2006-01-02"), now.Hour())

				f, err := os.Create(filename)
				if err != nil {
					log.Printf("CRITICAL: Cannot create log file: %v", err)
					continue
				}

				currentFile = f
				currentHour = now.Hour()
				log.Printf("Rotated log file to: %s", filename)
			}
			n, _ := currentFile.WriteString(fmt.Sprintf("%s [%s] %s\n",
				entry.Timestamp.Format(time.RFC3339),
				entry.Level,
				entry.Message,
			))

			totalBytes += n
			if entry.Level == "ERROR" {
				errorCount++
			}

		case <-ticker.C:
			log.Printf("--- STATS (Last 5s) ---")
			log.Printf("Errors: %d", errorCount)
			log.Printf("Bytes Written: %d", totalBytes)
			log.Println("-----------------------")

			// Reset counters for the next 5-second window
			errorCount = 0
			totalBytes = 0
		}
	}
}
