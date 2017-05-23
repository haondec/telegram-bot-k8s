package main

import (
    "encoding/json"
    "fmt"
    "io/ioutil"
    "os"
)

type Claim struct {
    UserName string `json:"username"`
    Role string `json:"role"`
}

func (p Claim) toString() string {
    return toJson(p)
}

func toJson(p interface{}) string {
    bytes, err := json.Marshal(p)
    if err != nil {
        fmt.Println(err.Error())
        os.Exit(1)
    }

    return string(bytes)
}

func maintesting() {

    claims := getClaims()
    for _, p := range claims {
        fmt.Println(p.toString())
    }

    fmt.Println(toJson(claims))
}

func getClaims() []Claim {
    raw, err := ioutil.ReadFile("./Claims.json")
    if err != nil {
        fmt.Println(err.Error())
        os.Exit(1)
    }

    var c []Claim
    json.Unmarshal(raw, &c)
    return c
}

