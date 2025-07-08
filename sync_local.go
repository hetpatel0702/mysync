package main

import (
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"sync"
	"time"
)

func handleLocalSync(src_dir string, dst_dir string, dryRun *bool, mirror *bool, verbose *bool, start time.Time) {

	jobs := make(chan FileJob, 100)
	wg := &sync.WaitGroup{}

	for i := 0; i < 20; i++ {
		go func() {
			for job := range jobs {
				copyFile(job.srcPath, job.dstPath, *dryRun)
				wg.Done()
			}
		}()
	}

	filepath.WalkDir(src_dir, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			// Handle error accessing the path
			return err
		}

		rel_path, _ := filepath.Rel(src_dir, path)
		dst_p := filepath.Join(dst_dir, rel_path)
		src_p := filepath.Join(src_dir, rel_path)

		if d.IsDir() {

			_, e := os.Stat(dst_p)
			if e != nil {
				if os.IsNotExist(e) {
					if *verbose {
						fmt.Println("Directory does not exist, creating:", dst_p)
					}
					err := os.MkdirAll(dst_p, 0755)
					if err != nil {
						fmt.Println(err)
					}
					return nil
				}
				return e
			}

			return nil // Continue walking
		}

		dst_f, e := os.Stat(dst_p)
		src_f, _ := os.Stat(src_p)

		if e != nil {
			if os.IsNotExist(e) {
				if *verbose {
					fmt.Println("File does not exist, copying:", src_p, "to", dst_p)
				}
				wg.Add(1)
				jobs <- FileJob{srcPath: src_p, dstPath: dst_p}
				return nil
			}
			return e
		}

		if dst_f.Size() != src_f.Size() || !dst_f.ModTime().Equal(src_f.ModTime()) {
			if *verbose {
				fmt.Println("File sizes or modification times differ, copying:", src_p, "to", dst_p)
			}
			wg.Add(1)
			jobs <- FileJob{srcPath: src_p, dstPath: dst_p}
			return nil
		}

		return nil
	})

	close(jobs)
	wg.Wait()

	fmt.Println("Sync completed in", time.Since(start))
}

func copyFile(src_p string, dst_p string, dryRun bool) {

	if dryRun {
		return
	}

	srcFile, err := os.Open(src_p)
	if err != nil {
		fmt.Println(err)
	}
	defer srcFile.Close()

	dstFile, _ := os.Create(dst_p)
	io.Copy(dstFile, srcFile)
	defer dstFile.Close()

	srcInfo, _ := os.Stat(src_p)
	os.Chtimes(dst_p, srcInfo.ModTime(), srcInfo.ModTime())

}

func mirrorDirs(src_dir string, dst_dir string, dryRun *bool, verbose *bool) {
	var toDelete []string
	filepath.WalkDir(dst_dir, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			// Handle error accessing the path
			return err
		}

		rel_path, _ := filepath.Rel(dst_dir, path)
		src_p := filepath.Join(src_dir, rel_path)

		if path == dst_dir {
			return nil // skip root
		}

		_, err = os.Stat(src_p)
		if os.IsNotExist(err) {
			toDelete = append(toDelete, path)
		}

		return nil
	})
	for i := len(toDelete) - 1; i >= 0; i-- {
		if *verbose {
			fmt.Println("Deleting:", toDelete[i])
		}
		if !*dryRun {
			os.RemoveAll(toDelete[i])
		}
	}
}
