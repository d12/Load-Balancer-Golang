package main

import (
    "fmt"
    "net/http"
    //"gopkg.in/yaml.v2"
)

var proxy = Proxy {
  hostOrigin: "http://localhost:9090",
}

func main() {
    http.HandleFunc("/", proxy.handler)
    err := http.ListenAndServe(":9090", nil)
    if err != nil {
        fmt.Print("ListenAndServe: ", err)
    }
}
