package lib

import (
	"context"
	"fmt"
	"io"
	"net"
	"strings"
	"sync"
)

type WriterFunc func(w io.Writer) *Writer

type Server struct {
	mu          sync.RWMutex
	ListenAddr  string
	ln          net.Listener
	quitChan    chan struct{}
	aof         *AppendOnlyFile
	queue       []Value
	tranMode    bool
	spawnWriter WriterFunc
	writer      *Writer
}

func NewServer(addr string) Server {
	return Server{
		mu:          sync.RWMutex{},
		ListenAddr:  addr,
		quitChan:    make(chan struct{}),
		aof:         nil,
		queue:       make([]Value, 0),
		tranMode:    false, //the server starts in normal mode.
		spawnWriter: NewWriter,
	}
}

func (s *Server) Start() error {
	//create a tcp ln on port 6379
	ln, err := net.Listen("tcp", s.ListenAddr)
	if err != nil {
		fmt.Println(err)
		return err
	}

	s.ln = ln

	s.aof = s.createAOF("minired.aof")

	go s.acceptConn()

	<-s.quitChan

	defer ln.Close()

	return nil
}

func (s *Server) acceptConn() {
	for {
		conn, err := s.ln.Accept()
		if err != nil {
			fmt.Println("CONNECTION_ERROR", err)
			continue
		}
		go s.readConn(conn)
	}
}

func (s *Server) readConn(conn net.Conn) {
	defer conn.Close()

	for {
		resp := NewResp(conn)
		value, err := resp.Read()
		if err != nil {
			if err == io.EOF {
				break
			}
			fmt.Println("READ_ERROR", err)
			continue
		}

		// if the request sent is not an array type ignore it
		// or if the request is empty.
		if value.Typ != "array" || len(value.Array) == 0 {
			continue
		}

		s.writer = s.spawnWriter(conn)
		result := s.handleCommandExecution(value)
		s.writer.Write(result)
	}
}

func (s *Server) handleCommandExecution(value Value) Value {
	command := strings.ToLower(value.Array[0].Bulk)

	if command == "multi" {
		s.turnOnTranMode()
		s.clearQueue() //clear queue at every transaction initiaition
		return Value{Typ: "string", Str: "Ok"}
	}

	if command == "exec" {
		results := s.executeQueuedCommands()
		s.turnOffTranMode()
		s.clearQueue()
		return results
	}

	if command == "discard" {
		s.turnOffTranMode()
		s.clearQueue()
		return Value{Typ: "string", Str: "Ok"}
	}

	if s.isInTranscMode() {
		s.enqeueCommand(value)
		return Value{Typ: "string", Str: "Queued"}
	}

	if command == "set" || command == "hset" {
		s.aof.Write(value)
	}

	result := s.execCommand(value)
	return result
}

func (s *Server) clearQueue() {
	s.mu.RLock()
	defer s.mu.RUnlock()
	s.queue = s.queue[:0] //clear the queue
}

func (s *Server) executeQueuedCommands() Value {
	results := Value{Typ: "array"}

	for _, value := range s.queue {

		result := s.execCommand(value)

		command := strings.ToLower(value.Array[0].Bulk)
		if command == "set" || command == "hset" {
			s.aof.Write(value)
		}

		results.Array = append(results.Array, result)
	}

	return results
}

func (s *Server) execCommand(value Value) Value {
	command := strings.ToLower(value.Array[0].Bulk)
	args := value.Array[1:]

	// get the handler for the command .
	handler, ok := CommandHandlers[command]
	if !ok {
		fmt.Printf("Command not supported [%s]\n", command)
		return Value{Typ: "string", Str: ""}
	}

	// and feed it the arguements
	result := handler(context.Background(), args)
	return result
}

func (s *Server) isInTranscMode() bool {
	return s.tranMode
}

func (s *Server) turnOnTranMode() bool {
	s.mu.Lock()
	defer s.mu.Unlock()

	if s.isInTranscMode() {
		return true
	}
	s.tranMode = true

	return true
}

func (s *Server) turnOffTranMode() {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.tranMode = false
}

func (s *Server) enqeueCommand(args Value) {
	s.mu.Lock()
	s.queue = append(s.queue, args)
	s.mu.Unlock()
}

func (s *Server) createAOF(path string) *AppendOnlyFile {
	aof, err := NewAppendOnlyFile(path)
	if err != nil {
		fmt.Println("Error: ", err)
		return nil
	}

	// execute every write entry in the log file to populate the
	// store with the data before the server was shutdown.
	aof.Read(func(value Value) {
		s.execCommand(value)
	})

	return aof
}
