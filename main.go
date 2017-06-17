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
    proxy, err := readConfig()
    if err != nil {
      log.Fatal(err)
    }

    fmt.Println("Main: Forwarding requests to proxy...")
    fmt.Println("Main: Listening to requests on port: " + strconv.Itoa(proxy.Port))
    http.HandleFunc("/", proxy.handler)
    err = http.ListenAndServe(":" + strconv.Itoa(proxy.Port), nil)
    if err != nil {
        log.Fatal(err)
    }
}
