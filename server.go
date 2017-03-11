package main

import (
    "fmt"
    "net/http"
    "net/url"
    "log"
    "io/ioutil"
    "strconv"
    "bytes"
    "io"
    "strings"
)

type Server struct {
    name string
    address string
    connections int
}

var servers = []Server {
    Server {
        name: "server A",
        address: "http://localhost:3000",
        connections: 0,
    },
    Server {
        name: "server B",
        address: "http://localhost:3000",
        connections: 0,
    },
}

func chooseServer() Server {
  var min = -1
  var minIndex = 0
  for index,server := range servers {
    var conn = server.connections
    if min == -1 {
      min = conn
      minIndex = index
    }else if(conn < min){
      min = conn
      minIndex = index
    }
  }

  return servers[minIndex]
}

func GetHandler(w http.ResponseWriter, r *http.Request, server Server) {
  u, err := url.Parse(server.address + r.RequestURI)
  if err != nil {
      panic(err)
  }
  r.URL = u
  r.Host = ""
  fmt.Println("HOST: " + r.Host)
  r.Header.Set("X-Forwarded-Host", "localhost:3000") // THIS SHOULD BE SET THE THE REQUESTED HOST
  r.Header.Set("Origin", server.address)
  r.Host = "localhost:3000"
  fmt.Println("HOST: " + r.Host)
  r.RequestURI = ""

  client := &http.Client{}

  resp, err := client.Do(r)
  if err != nil {
    log.Fatal(err)
    fmt.Fprint(w, "Internal server error sorry")
    defer resp.Body.Close()
    return
  }

  bodyBytes, err := ioutil.ReadAll(resp.Body)
  if err != nil {
    log.Fatal(err)
    fmt.Fprint(w, "Failed to read response body")
    defer resp.Body.Close()
    return
  }

  buffer := bytes.NewBuffer(bodyBytes)

  contentLength := strconv.Itoa(len(bodyBytes))

  w.Header().Set("Content-Length", contentLength)
  for k, v := range resp.Header {
    w.Header().Set(k, strings.Join(v, ";"))
  }

  io.Copy(w, buffer)
  defer resp.Body.Close()
}

func PostHandler(w http.ResponseWriter, r *http.Request, server Server) {
  u, err := url.Parse(server.address + r.RequestURI)
  if err != nil {
      panic(err)
  }
  r.URL = u
  r.Host = server.address
  r.Header.Set("X-Forwarded-Host", server.address)
  r.RequestURI = ""

  client := &http.Client{
  }

  resp, err := client.Do(r)
  if err != nil {
    log.Fatal(err)
  }/*
  b, err := ioutil.ReadAll(r.Body)
  buf := bytes.NewBuffer(b)

  resp, err := http.Post(server.address + r.RequestURI, r.Header.Get("Content-Type"), buf)*/

  fmt.Printf("WOW")

  fmt.Printf("RESPONSE: %+v\n", resp)
  fmt.Printf("URL: %+v\n",r.URL)

  if err != nil {
    log.Fatal(err)
    fmt.Fprint(w, "Internal server error sorry")
    defer resp.Body.Close()
    return
  }

  bodyBytes, err := ioutil.ReadAll(resp.Body)
  if err != nil {
    log.Fatal(err)
    fmt.Fprint(w, "Failed to read response body")
    defer resp.Body.Close()
    return
  }

  buffer := bytes.NewBuffer(bodyBytes)

  fmt.Println("BODY: " + string(bodyBytes))

  contentLength := strconv.Itoa(len(bodyBytes))

  w.Header().Set("Content-Length", contentLength)
  for k, v := range resp.Header {
    w.Header().Set(k, strings.Join(v, ";"))
  }

  io.Copy(w, buffer)
  defer resp.Body.Close()
}

func handler(w http.ResponseWriter, r *http.Request) {
    var server = chooseServer()

    server.connections += 1

    switch r.Method {
    case "GET":
      GetHandler(w, r, server)

    case "POST":
      GetHandler(w, r, server)
    }

    server.connections -= 1

    fmt.Println("Served a request wew")
}

func main() {
    http.HandleFunc("/", handler)
    err := http.ListenAndServe(":9090", nil)
    if err != nil {
        fmt.Print("ListenAndServe: ", err)
    }
}
