package main

import (
    "fmt"
    "net/http"
    "net/url"
    "io/ioutil"
    "bytes"
    "io"
    "strings"
    "strconv"
)

type Proxy struct {
  Host string
  Port int
}

// TODO: Optional ports, different schemes
func (proxy Proxy) origin() string {
  return ("http://" + proxy.Host + ":" + strconv.Itoa(proxy.Port));
}

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

func (proxy Proxy)chooseServer() Server {
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

func (proxy Proxy)ReverseProxy(w http.ResponseWriter, r *http.Request, server Server) {
  fmt.Println("Proxy: Parsing URL...")
  u, err := url.Parse(server.Url() + r.RequestURI)
  if err != nil {
      panic(err)
  }

  fmt.Println("Proxy: Requested resource: " + r.RequestURI)
  fmt.Println("Proxy: Re-assigning request headers...")

  r.URL = u
  r.Header.Set("X-Forwarded-Host", r.Host)
  r.Header.Set("Origin", proxy.origin())
  r.Host = server.Url()
  r.RequestURI = ""

  fmt.Println("Proxy: Headers re-assigned. Sending new request to web server...")

  client := &http.Client{
    CheckRedirect: func(req *http.Request, via []*http.Request) error {
        return http.ErrUseLastResponse
    },
  }

  resp, err := client.Do(r)
  if err != nil {
    fmt.Println(err)
    fmt.Println(w, "Proxy: Internal server error sorry")
    http.NotFound(w, r)
    return
  }

  fmt.Println("Proxy: Request successful, recieved response.")
  fmt.Println("Proxy: Parsing response...")

  bodyBytes, err := ioutil.ReadAll(resp.Body)
  if err != nil {
    fmt.Println(w, "Proxy: Failed to read response body")
    http.NotFound(w, r)
    return
  }

  buffer := bytes.NewBuffer(bodyBytes)

  fmt.Println("Proxy: Response parsed.")
  fmt.Println("Proxy: Rewriting response headers for client...")

  for k, v := range resp.Header {
    w.Header().Set(k, strings.Join(v, ";"))
  }

  fmt.Println("Proxy: Headers rewritten. Sending response back to client...")

  w.WriteHeader(resp.StatusCode)
  fmt.Println(resp.StatusCode)

  io.Copy(w, buffer)
  defer resp.Body.Close()
}

func (proxy Proxy)handler(w http.ResponseWriter, r *http.Request) {
    fmt.Println("Proxy: Recieved a request, assigning a web server...");
    var server = proxy.chooseServer()

    server.connections += 1

    proxy.ReverseProxy(w, r, server)

    server.connections -= 1

    fmt.Println("Proxy: Responded to request successfuly!")
}
