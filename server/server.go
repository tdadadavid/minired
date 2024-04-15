package server

import (
	"context"
	"fmt"
	"io"
	"minired/lib"
	"net"
	"strings"
)

type Server struct {
	ListenAddr string
	ln         net.Listener
	quitChan   chan struct{}
	aof        *lib.AppendOnlyFile
}

func NewServer(addr string) Server {
	return Server{
		ListenAddr: addr,
		quitChan:   make(chan struct{}),
		aof:        createAOF("minired.aof"),
	}
}

func (s Server) Start() error {
	//create a tcp ln on port 6379
	ln, err := net.Listen("tcp", s.ListenAddr)
	if err != nil {
		fmt.Println(err)
		return err
	}

	s.ln = ln

	go s.acceptConn()

	<-s.quitChan

	defer ln.Close()

	return nil
}

func (s Server) acceptConn() {
	for {
		conn, err := s.ln.Accept()
		if err != nil {
			fmt.Println("CONNECTION_ERROR", err)
			continue
		}
		go s.readConn(conn)
	}
}

func (s Server) readConn(conn net.Conn) {
	defer conn.Close()

	for {
		resp := lib.NewResp(conn)
		value, err := resp.Read()
		if err != nil {
			if err == io.EOF {
				break
			}
			fmt.Println("READ_ERROR", err)
			continue
		}

		// if the request sent is not an array type ignore it
		if value.Typ != "array" {
			continue
		}

		// ignore empty request
		if len(value.Array) == 0 {
			continue
		}

		writer := lib.NewWriter(conn)

		command := strings.ToLower(value.Array[0].Bulk)
		args := value.Array[1:]

		if command == "set" || command == "hset" {
			s.aof.Write(value)
		}

		// get the handler for the command .
		handler, ok := lib.CommandHandlers[command]
		if !ok {
			fmt.Println("Command not supported: [", command, "]")
			writer.Write(lib.Value{Typ: "string", Str: ""})
			continue
		}
		// and feed it the arguements
		result := handler(context.Background(), args)

		writer.Write(result)
	}
}

func createAOF(path string) *lib.AppendOnlyFile {
	aof, err := lib.NewAppendOnlyFile(path)
	if err != nil {
		fmt.Println("Error: ", err)
		return nil
	}

	// execute every write entry in the log file to populate the
	// store with the data before the server was shutdown.
	aof.Read(func(value lib.Value) {
		command := strings.ToLower(value.Array[0].Bulk)
		args := value.Array[1:]

		handler, ok := lib.CommandHandlers[command]
		if !ok {
			fmt.Println("Command not supported: [", command, "]")
		}
		handler(context.Background(), args)
	})

	return aof
}
