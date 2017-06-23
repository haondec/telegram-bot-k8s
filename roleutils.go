package main

import (
    "encoding/json"
    "fmt"
    "io/ioutil"
    "os"
)

type Claim struct {
    UserName 	string `json:"username"`
    Role 	string `json:"role"`
}

func (p Claim) toString() string {
    return toJson(p)
}

func toJson(p interface{}) string {
    bytes, err := json.MarshalIndent(&p, "", "\t")
    if err != nil {
        fmt.Println(err.Error())
//        os.Exit(1)
    }

    return string(bytes)
}

func maintesting() {
    
    claims := getClaims("pathtofile")
    for _, p := range claims {
        fmt.Println(p.toString())
    }

    fmt.Println(toJson(claims))
}

func getClaims(fn string) []Claim {
    raw, err := ioutil.ReadFile(fn)
    if err != nil {
        fmt.Println(err.Error())
        os.Exit(1)
    }

    var c []Claim
    json.Unmarshal(raw, &c)
    return c
}

// Func map file roles of user (file .json) 
func rolemap(fn string) map[string]string {
	claims := getClaims(fn)
	var rm map[string]string
	rm = make(map[string]string)
	for _, p:= range claims {
		rm[p.UserName] = p.Role
	}
	return rm
}
