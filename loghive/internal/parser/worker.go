package parser

import (
	"log"
	"strings"
	"sync"
	"time"

	"github.com/alexanderritik/loghive/internal/domain"
)

type Service struct {
	InputCh   <-chan []byte
	ResultsCh chan domain.LogEntry
	wg        sync.WaitGroup
}

func NewService(input <-chan []byte) *Service {
	return &Service{
		InputCh:   input,
		ResultsCh: make(chan domain.LogEntry, 100), // Buffered to prevent blocking workers
	}
}

func (s *Service) StartWorkers(numWorkers int) {

	// fan-IN
	for i := 0; i < numWorkers; i++ {
		s.wg.Add(1)
		go s.worker(i)
	}

	go func() {
		s.wg.Wait()
		close(s.ResultsCh)
		log.Println("All parser workers stopped, ResultsCh closed.")
	}()
}

func (s *Service) worker(id int) {
	defer s.wg.Done()

	for rawMsg := range s.InputCh {
		entry, err := s.parse(rawMsg)
		if err != nil {
			log.Printf("Worker %d: Failed to parse log: %v", id, err)
			continue
		}

		s.ResultsCh <- entry
	}
}

func (s *Service) parse(msg []byte) (domain.LogEntry, error) {
	line := string(msg)

	parts := strings.SplitN(line, "|", 3)

	if len(parts) < 3 {
		// In a real app, define a custom error
		return domain.LogEntry{}, nil
	}

	ts, err := time.Parse(time.RFC3339, parts[0])
	if err != nil {
		ts = time.Now() // Fallback
	}

	entry := domain.LogEntry{
		Timestamp: ts,
		Level: parts[1],
		Message: parts[2],
	}

	return  entry, nil

}
