package main

import (
    "fmt"
    "net/http"
    "net/url"
    "log"
    "io/ioutil"
    "bytes"
    "io"
    "strings"
)

var servers = []Server {
    Server {
        name: "server A",
        scheme: "http",
        host: "localhost",
        port: "3000",
    },
    Server {
        name: "server B",
        scheme: "http",
        host: "localhost",
        port: "3000",
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

func ReverseProxy(w http.ResponseWriter, r *http.Request, server Server) {
  u, err := url.Parse(server.Url() + r.RequestURI)
  fmt.Println("GOT REQUEST: " + r.RequestURI);
  if err != nil {
      panic(err)
  }

  r.URL = u
  r.Header.Set("X-Forwarded-Host", r.Host)
  r.Header.Set("Origin", "http://localhost:9090")
  r.Host = server.Url()
  r.RequestURI = ""

  client := &http.Client{
    CheckRedirect: func(req *http.Request, via []*http.Request) error {
        return http.ErrUseLastResponse
    },
  }

  resp, err := client.Do(r)
  if err != nil {
    fmt.Println(err)
    fmt.Println(w, "Internal server error sorry")
    http.NotFound(w, r)
    return
  }

  bodyBytes, err := ioutil.ReadAll(resp.Body)
  if err != nil {
    fmt.Println(w, "Failed to read response body")
    http.notFound(w, r)
    return
  }

  buffer := bytes.NewBuffer(bodyBytes)

  for k, v := range resp.Header {
    w.Header().Set(k, strings.Join(v, ";"))
  }

  w.WriteHeader(resp.StatusCode)
  fmt.Println(resp.StatusCode)

  io.Copy(w, buffer)
  defer resp.Body.Close()
}

func handler(w http.ResponseWriter, r *http.Request) {
    var server = chooseServer()

    server.connections += 1

    ReverseProxy(w, r, server)

    server.connections -= 1

    fmt.Println("Served a request")
}

func main() {
    http.HandleFunc("/", handler)
    err := http.ListenAndServe(":9090", nil)
    if err != nil {
        fmt.Print("ListenAndServe: ", err)
    }
}
