package utils

import "os"

type FileFilter func(q string, path string, info os.FileInfo) bool

var DefaultFilter FileFilter = func(q string, path string, info os.FileInfo) bool { return true }
