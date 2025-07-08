package main

import (
	"encoding/binary"
	"encoding/json"
	"fmt"
	"net"
	"time"
)

type FileMeta struct {
	Path    string `json:"path"`
	Size    int64  `json:"size"`
	ModTime int64  `json:"mod_time"`
}

func main() {
	conn, err := net.Dial("tcp", "localhost:8080")
	if err != nil {
		panic(err)
	}
	defer conn.Close()

	meta := FileMeta{
		Path:    "example/file.txt",
		Size:    1024,
		ModTime: time.Now().Unix(),
	}

	// Convert to JSON
	jsonBytes, err := json.Marshal(meta)
	if err != nil {
		panic(err)
	}

	// First send 4 bytes for JSON length
	lengthBuf := make([]byte, 4)
	binary.BigEndian.PutUint32(lengthBuf, uint32(len(jsonBytes)))
	conn.Write(lengthBuf)

	// Then send the JSON itself
	conn.Write(jsonBytes)

	fmt.Println("Sent metadata:", string(jsonBytes))
}
