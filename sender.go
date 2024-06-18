package log

import "log"

type SenderI interface {
	Start()
	LogChan() chan<- Entry
	Stop()
}

type Sender struct {
	LogClient ClientI
	logChan   chan Entry
}

func NewLogSender(logClient ClientI) SenderI {
	return &Sender{
		LogClient: logClient,
		logChan:   make(chan Entry, 50),
	}
}

func (s *Sender) Start() {
	go func() {
		for logMessage := range s.logChan {
			err := s.LogClient.Send(logMessage)
			if err != nil {
				log.Printf("Failed to send log: %v", err)
			} else {
				log.Printf("Log sent successfully: %v", logMessage)
			}
		}
	}()
}

func (s *Sender) LogChan() chan<- Entry {
	return s.logChan
}

func (s *Sender) Stop() {
	close(s.logChan)
}
