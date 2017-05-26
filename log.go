package main

import (
	"encoding/json"
	"strings"
	"time"
	"os"
	"crypto/md5"
	"fmt"
	"io"
)

type LogConfig struct {
	log_count	int
	log_count_max	int
	log_sub		int

	log_file_name	string
}

type LogJson struct {
	ID		string `json:"id"`
	Userid		string `json:"userid"` 
	Date		string `json:"date"`
	Message		string `json:"message"`
}

// stringToMap splits `s string` using `sep` as separator and
// set every word as a key in a new map. The value of
// all keys is set to `true`
func stringToMap(s string, sep string) map[string]bool {
	ss := strings.Split(s, sep)
	m := make(map[string]bool)
	for _, word := range ss {
		m[word] = true
	}

	return m

}

// Format filename
func fmFileName(fn string) string {
	fn = strings.Replace(fn, " ", "_", -1)
	return strings.ToLower(fn)
}

// Format json message
func fmJsonMessage(str string) string {
	str = strings.Replace(str, "\"id\":", "\n\t\"id\":", -1)
	str = strings.Replace(str, "\"userid\":", "\n\t\"userid\":", -1)
	str = strings.Replace(str, "\"date\":", "\n\t\"date\":", -1)
	str = strings.Replace(str, "\"message\":", "\n\t\"message\":", -1)
	str = strings.Replace(str, "\"}", "\"\n}", -1)
	return str
}

// Append file
func appendFile(str string) {
	if lcfg_main.log_count >= lcfg_main.log_count_max {
		lcfg_main.log_sub++
	}
	fn := validatePath(os.Getenv(telegramLogLabel)) +
		lcfg_main.log_file_name +
		fmt.Sprintf("_%d", lcfg_main.log_sub) +
		".log"
	
	// Open append
	f, err := os.OpenFile(fn, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0664)
	if err != nil {
		panic(err)
		return
	}

	defer f.Close()

	// Add comma
	fInfo, err := os.Stat(fn)
	commaCheck := false
	if err != nil {
		panic(err)
		return
	}
	if fInfo.Size() > 0 {
		commaCheck = true
	}
	
	if commaCheck {
		if _, err = f.WriteString(",\n" + fmJsonMessage(str)); err != nil {
			panic(err)
			f.Close()
			return
		}
	} else {
		if _, err = f.WriteString(fmJsonMessage(str)); err != nil {
                        panic(err)
			f.Close()
                        return
                }
	}

	f.Close()
}

// Init log name
func initLogName() string {
	t := time.Now()
	time := t.Format(timeFM)
	return fmFileName(time)	
}

// Write log
func writeLog(uid string, str string) {
	t := time.Now()
	time := t.Format(timeFM)
	logjs := LogJson{}
	h := md5.New()
	io.WriteString(h, str)

	logjs.ID 	= fmt.Sprintf("%x", h.Sum(nil))
	logjs.Userid 	= uid
	logjs.Date 	= time
	logjs.Message 	= str

	mesjs, _ := json.Marshal(logjs)
	appendFile(string(mesjs))
}
