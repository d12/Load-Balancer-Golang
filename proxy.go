package main

import (
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
  Scheme string
  Servers []Server
}

func (proxy Proxy) origin() string {
  return (proxy.Scheme + "://" + proxy.Host + ":" + strconv.Itoa(proxy.Port));
}

// TODO: This crashes if we define no servers in our config
func (proxy Proxy)chooseServer() *Server {
  var min = -1
  var minIndex = 0
  for index,server := range proxy.Servers {
    var conn = server.Connections
    if min == -1 {
      min = conn
      minIndex = index
    }else if(conn < min){
      min = conn
      minIndex = index
    }
  }

  return &proxy.Servers[minIndex]
}

func (proxy Proxy)ReverseProxy(w http.ResponseWriter, r *http.Request, server Server) {
  u, err := url.Parse(server.Url() + r.RequestURI)
  if err != nil {
      LogErrAndCrash(err.Error())
  }

  r.URL = u
  r.Header.Set("X-Forwarded-Host", r.Host)
  r.Header.Set("Origin", proxy.origin())
  r.Host = server.Url()
  r.RequestURI = ""


  client := &http.Client{
    CheckRedirect: func(req *http.Request, via []*http.Request) error {
        return http.ErrUseLastResponse
    },
  }

  // TODO: If the server doesn't respond, try a new web server
  // We could return a status code from this function and let the handler try passing the request to a new server.
  resp, err := client.Do(r)
  if err != nil {
    // For now, this is a fatal error
    // When we can fail to another webserver, this should only be a warning.
    LogErr(err.Error())
    http.NotFound(w, r)
    return
  }
  LogInfo("Recieved response: " + strconv.Itoa(resp.StatusCode))

  bodyBytes, err := ioutil.ReadAll(resp.Body)
  if err != nil {
    LogErr("Proxy: Failed to read response body")
    http.NotFound(w, r)
    return
  }

  buffer := bytes.NewBuffer(bodyBytes)
  for k, v := range resp.Header {
    w.Header().Set(k, strings.Join(v, ";"))
  }

  w.WriteHeader(resp.StatusCode)

  io.Copy(w, buffer)
  defer resp.Body.Close()
}

func (proxy Proxy)handler(w http.ResponseWriter, r *http.Request) {
    var server = proxy.chooseServer()
    LogInfo("Got request: " + r.RequestURI)
    LogInfo("Sending to server: " + server.Name)

    server.Connections += 1

    proxy.ReverseProxy(w, r, *server)

    server.Connections -= 1

    LogInfo("Responded to request successfuly")
}
