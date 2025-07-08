package main

import (
	"encoding/binary"
	"encoding/json"
	"fmt"
	"io"
	"net"
	"os"
	"path/filepath"
	"time"
)

type FileMeta struct {
	Path    string `json:"path"`
	Size    int64  `json:"size"`
	ModTime int64  `json:"mod_time"`
	Dir     bool   `json:"dir"`
}

func main() {
	ln, err := net.Listen("tcp", ":8080")
	if err != nil {
		panic(err)
	}
	fmt.Println("Server is listening on port 8080...")
	for {
		conn, err := ln.Accept()
		if err != nil {
			fmt.Println("Error accepting connection:", err)
			continue
		}
		go handleConnection(conn)
	}
}

func handleConnection(conn net.Conn) {
	defer conn.Close()
	fmt.Println("Handling connection from:", conn.RemoteAddr())

	lenBuf := make([]byte, 4)
	_, err := io.ReadFull(conn, lenBuf)
	if err != nil {
		fmt.Println("Error reading length:", err)
		return
	}
	jsonLen := binary.BigEndian.Uint32(lenBuf)
	jsonBuf := make([]byte, jsonLen)
	_, err = io.ReadFull(conn, jsonBuf)
	if err != nil {
		fmt.Println("Error reading JSON data:", err)
		return
	}
	fmt.Println("Received JSON data:", string(jsonBuf))

	var meta FileMeta
	err = json.Unmarshal(jsonBuf, &meta)
	if err != nil {
		fmt.Println("Error unmarshalling JSON:", err)
		return
	}

	fmt.Printf("Received file metadata: %+v\n", meta)

	f, err := os.Stat(meta.Path)

	if err != nil {
		if os.IsNotExist(err) {
			if meta.Dir {
				err := os.MkdirAll(meta.Path, 0755)
				if err != nil {
					fmt.Println("Failed to create directory:", err)
				}
				skipFile(conn)
				return
			} else {
				fmt.Println("File does not exist in remote, copying:", meta.Path)
				getFile(conn, meta)
			}
		}
		fmt.Println("Error checking file:", err)
		return
	}

	if f.Size() != meta.Size || f.ModTime().Unix() != meta.ModTime {
		fmt.Println("File sizes or modification times differ, copying:", meta.Path)
		getFile(conn, meta)
	} else {
		skipFile(conn)
		fmt.Println("File already exists and is up to date:", meta.Path)
	}

}

func skipFile(conn net.Conn) {
	conn.Write([]byte("SKIP\n"))
}

func getFile(conn net.Conn, meta FileMeta) {
	conn.Write([]byte("SEND\n"))

	err := os.MkdirAll(filepath.Dir(meta.Path), 0755)
	if err != nil {
		fmt.Println("Error creating parent dirs:", err)
		return
	}

	dstFile, err := os.Create(meta.Path)
	if err != nil {
		fmt.Println("Error creating file:", err)
		return
	}
	defer dstFile.Close()

	_, err = io.CopyN(dstFile, conn, meta.Size)
	if err != nil {
		fmt.Println("Error reading file:", err)
		return
	}

	modTime := time.Unix(meta.ModTime, 0)
	os.Chtimes(meta.Path, modTime, modTime)

}
