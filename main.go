package main

import (
    "net/http"
    "strconv"
)

func main() {
    LogInfo("Spinning up load balancer...")
    LogInfo("Reading Config.yml...")
    proxy, err := ReadConfig()
    if err != nil {
      LogErr("An error occurred while trying to parse config.yml")
      LogErrAndCrash(err.Error())
    }

    http.HandleFunc("/", proxy.handler)
    err = http.ListenAndServe(":" + strconv.Itoa(proxy.Port), nil)
    if err != nil {
        LogErr("Failed to bind to port " + strconv.Itoa(proxy.Port))
        LogErrAndCrash("Make sure the port is available and not reserved")
    }
    LogInfo("Listening to requests on port: " + strconv.Itoa(proxy.Port))
}
