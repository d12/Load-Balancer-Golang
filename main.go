package main

import (
    "fmt"
    "net/http"
    "log"
    "strconv"
)



func main() {
    fmt.Println("Main: Spinning up load balancer...")
    fmt.Println("Main: Reading Config.yml...")
    config, err := readConfig()
    if err != nil {
      log.Fatal(err)
    }

    var proxy = Proxy {
      Host: config.Host,
      Port: config.Port,
    }

    fmt.Println("Main: Forwarding requests to proxy...")
    fmt.Println("Main: Listening to requests on port: " + proxy.Port)
    http.HandleFunc("/", proxy.handler)
    err = http.ListenAndServe(":" + strconv.Itoa(proxy.Port), nil)
    if err != nil {
        log.Fatal(err)
    }
}
