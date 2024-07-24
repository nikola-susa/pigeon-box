package serverevent

import (
	"bytes"
	"fmt"
	"net/http"
	"sync"
	"sync/atomic"
	"time"
)

type Event struct {
	Data      []byte
	Timestamp time.Time
}

type Subscriber struct {
	Events chan *Event
}

type Stream struct {
	Events      chan *Event
	Subscribers []*Subscriber
	SubCount    int32
	mu          sync.RWMutex
}

type Server struct {
	Streams map[string]*Stream
	mu      sync.RWMutex
}

func New() *Server {
	return &Server{
		Streams: make(map[string]*Stream),
	}
}

func NewStream() *Stream {
	s := &Stream{
		Events:      make(chan *Event),
		Subscribers: make([]*Subscriber, 0),
	}

	go s.run()
	return s
}

func (s *Server) AddStream(name string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.Streams[name] = NewStream()
}

func (s *Server) RemoveStream(name string) {
	s.mu.Lock()
	defer s.mu.Unlock()

	stream, ok := s.Streams[name]
	if !ok {
		return
	}
	stream.close()
	delete(s.Streams, name)
}

func (s *Server) GetStream(name string) *Stream {
	s.mu.RLock()
	defer s.mu.RUnlock()

	stream, ok := s.Streams[name]
	if !ok {
		stream = NewStream()
		s.Streams[name] = stream
	}

	return stream
}

func NewSubscriber() *Subscriber {
	return &Subscriber{
		Events: make(chan *Event),
	}
}

func (s *Stream) subscribe(sub *Subscriber) {
	s.Subscribers = append(s.Subscribers, sub)
	atomic.AddInt32(&s.SubCount, 1)
}

func (s *Stream) unsubscribe(sub *Subscriber) {
	s.mu.Lock()
	defer s.mu.Unlock()
	for i, subscriber := range s.Subscribers {
		if subscriber == sub {
			s.Subscribers = append(s.Subscribers[:i], s.Subscribers[i+1:]...)
			atomic.AddInt32(&s.SubCount, -1)
			return
		}
	}
}

func (s *Stream) publish(event *Event) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	for _, subscriber := range s.Subscribers {
		subscriber.Events <- event
	}
}

func (s *Stream) run() {
	for {
		select {
		case event := <-s.Events:
			s.publish(event)
		}
	}
}

func (s *Stream) close() {
	s.mu.Lock()
	defer s.mu.Unlock()
	for _, subscriber := range s.Subscribers {
		close(subscriber.Events)
	}
	close(s.Events)
}

func (s *Stream) getSubCount() int32 {
	return atomic.LoadInt32(&s.SubCount)
}

func (s *Server) Broadcast(streamName string, data []byte) {

	stream := s.GetStream(streamName)
	if stream == nil {
		fmt.Println("Stream not found", streamName)
		return
	}

	lines := bytes.Split(data, []byte("\n"))
	eventData := make([]byte, 0)
	for _, line := range lines {
		if len(line) == 0 {
			continue
		}
		eventData = append(eventData, []byte("data: "+string(line)+"\n")...)
	}

	event := &Event{
		Data:      eventData,
		Timestamp: time.Now(),
	}
	stream.Events <- event
}

func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	streamName := r.PathValue("stream")
	if streamName == "" {
		http.Error(w, "stream query parameter is required", http.StatusBadRequest)
		return
	}

	stream := s.GetStream(streamName)

	if stream == nil {
		s.AddStream(streamName)
	}

	sub := NewSubscriber()
	stream.subscribe(sub)

	defer stream.unsubscribe(sub)

	flusher, ok := w.(http.Flusher)
	if !ok {
		http.Error(w, "streaming unsupported", http.StatusBadRequest)
		return
	}

	go func() {
		<-r.Context().Done()
		stream.unsubscribe(sub)

		if stream.getSubCount() == 0 {
			s.RemoveStream(streamName)
		}

	}()

	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")

	for event := range sub.Events {
		_, err := fmt.Fprintf(w, "event: %s\n%s\n\n", streamName, event.Data)
		if err != nil {
			fmt.Println("Error writing to stream", err)
			return
		}
		flusher.Flush()
	}
}
