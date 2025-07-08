package main

import (
	"flag"
	"log"
	"path/filepath"
	"time"
)

func main() {

	start := time.Now()

	var dryRun = flag.Bool("dry-run", false, "Perform a dry run without copying files")
	var mirror = flag.Bool("mirror", false, "Mirror the source directory to the destination directory")
	var verbose = flag.Bool("verbose", false, "Enable verbose output")
	var remote = flag.String("remote", "", "Remote server in the form ip:port (e.g. 192.168.1.10:8080)")

	flag.Parse()
	args := flag.Args()
	if len(args) != 2 {
		log.Fatal("Usage: mysync [--dry-run] src dst")
	}

	src_dir := args[0]
	dst_dir := args[1]

	src_dir, err := filepath.Abs(src_dir)
	if err != nil {
		log.Fatal("Error resolving src path:", err)
	}

	dst_dir, err = filepath.Abs(dst_dir)
	if err != nil {
		log.Fatal("Error resolving dst path:", err)
	}

	if *mirror {
		mirrorDirs(src_dir, dst_dir, dryRun, verbose)
	}

	if *remote != "" {
		handleRemoteSync(src_dir, dst_dir, dryRun, mirror, verbose, remote, start)
	} else {
		handleLocalSync(src_dir, dst_dir, dryRun, mirror, verbose, start)
	}
}
