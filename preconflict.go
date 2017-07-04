package main

import (
        "os"
        "strings"
	"fmt"
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

func subCL_getLast(path string, sep string) string {
        s_split := strings.Split(path, sep)
	// Always s_split have at least an member.
	return s_split[len(s_split) - 1]
}

func check_lock_v2(path string) bool {
        fn := subCL_getLast(path, "/")
        if fn != flock {
                fn = validatePath(path) + flock
        } else {
		fn = path
	}
        fInfo, err := os.Stat(fn)
        if err == nil {
                if fInfo.IsDir() == false {
                        return true
                }
        }
        return false
}

func make_lock(path string) {
        fn := subCL_getLast(path, "/")
        if fn != flock {
                fn = validatePath(path) + flock
	} else {
		fn = path
	}
	fmt.Println(fn)
        if check_lock(fn) == false {
                _, _ = os.Create(fn)
        }
}

func un_lock(path string) {
        fn := subCL_getLast(path, "/")
        if fn != flock {
                fn = validatePath(path) + flock
        } else {
		fn = path
	}
        if check_lock(fn) == true {
                _ = os.Remove(fn)
        }
}
