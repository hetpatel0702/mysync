package main

import (
	"encoding/binary"
	"encoding/json"
	"fmt"
	"io"
	"io/fs"
	"net"
	"os"
	"path/filepath"
	"time"
)

func handleRemoteSync(src_dir string, dst_dir string, dryRun *bool, mirror *bool, verbose *bool, remote *string, start time.Time) {

	filepath.WalkDir(src_dir, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			// Handle error accessing the path
			return err
		}

		rel_path, _ := filepath.Rel(src_dir, path)
		dst_p := filepath.Join(dst_dir, rel_path)
		src_p := filepath.Join(src_dir, rel_path)

		if d.IsDir() {
			syncRemote(src_p, dst_p, *remote, d.IsDir())
			return nil // Continue walking
		}

		syncRemote(src_p, dst_p, *remote, d.IsDir())
		return nil
	})

	fmt.Println("Sync completed in", time.Since(start))
}

func syncRemote(src_p string, dst_p string, remoteAddr string, isDir bool) {
	src_f, _ := os.Stat(src_p)

	metaData := FileMeta{
		Path:    dst_p,
		Size:    src_f.Size(),
		ModTime: src_f.ModTime().Unix(),
		Dir:     isDir,
	}

	jsonMetaData, err := json.Marshal(metaData)
	if err != nil {
		panic(err)
	}
	fmt.Println(string(jsonMetaData))

	conn, err := net.Dial("tcp", remoteAddr)
	if err != nil {
		panic(err)
	}
	defer conn.Close()

	// Send length of JSON
	lengthBuf := make([]byte, 4)
	binary.BigEndian.PutUint32(lengthBuf, uint32(len(jsonMetaData)))
	conn.Write(lengthBuf)

	// Send JSON
	conn.Write(jsonMetaData)

	resp := make([]byte, 10)
	n, _ := conn.Read(resp)
	fmt.Println("server response: ", string(resp[:n]))
	if string(resp[:n]) == "SEND\n" {

		// if it's a directory
		if isDir {
			return
		}

		file, err := os.Open(src_p)
		if err != nil {
			panic(err)
		}
		defer file.Close()

		_, err = io.Copy(conn, file) // conn is net.Conn
		if err != nil {
			fmt.Println("Error sending file:", err)
		}
	}
}
