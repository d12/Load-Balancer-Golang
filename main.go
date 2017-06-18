package main

import (
    "net/http"
    "strconv"
)

func main() {
    LogInfo("Spinning up load balancer...")
    LogInfo("Reading Config.yml...")
    proxy, err := readConfig()
    if err != nil {
      LogErr("Failed to read config.yml")
      LogErrAndCrash(err.Error())
    }

    LogInfo("Listening to requests on port: " + strconv.Itoa(proxy.Port))
    http.HandleFunc("/", proxy.handler)
    err = http.ListenAndServe(":" + strconv.Itoa(proxy.Port), nil)
    if err != nil {
        LogErrAndCrash(err.Error())
    }
}
