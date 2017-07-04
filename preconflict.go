package main

import (
        "errors"
        "os"
        "strings"
)

const (
        flock string = ".lock"        
)

func check_lock(fn string) bool {
        fInfo, err := os.Stat(fn)
        if err == nil {
                if fInfo.IsDir() == false {
                        return true
                }
        }
        return false
}

func make_lock(path string) {
        if strings.Contains(path, flock) == false {
                fn := validatePath(path) + flock
        }
        if check_lock(fn) == false {
                f, err := os.Create(fn)
        }
}

func un_lock(path string) {
        if strings.Contains(path, flock) == false {
                fn := validatePath(path) + flock
        }
        if check_lock(fn) == false {
                f, err := os.Remove(fn)
        }
}
