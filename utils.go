package main

import "strings"

type FileJob struct {
	srcPath string
	dstPath string
}

type FileMeta struct {
	Path    string `json:"path"`
	Size    int64  `json:"size"`
	ModTime int64  `json:"mod_time"`
	Dir     bool   `json:"dir"`
}

func isRemote(dst_dir string) bool {
	return strings.Contains(dst_dir, "@") && strings.Contains(dst_dir, ":")
}
