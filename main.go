package main

import (
    "fmt"
    "net/http"
    "log"
)



func main() {
    config, err := readConfig()
    if err != nil {
      log.Fatal(err)
    }

    var proxy = Proxy {
      HostOrigin: config.HostOrigin,
    }

    http.HandleFunc("/", proxy.handler)
    err = http.ListenAndServe(":9090", nil)
    if err != nil {
        fmt.Print("ListenAndServe: ", err)
    }
}
