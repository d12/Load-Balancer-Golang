package main

import (
    "net/http"
    "net/url"
    "io/ioutil"
    "bytes"
    "io"
    "strings"
    "strconv"
    "math"
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
func (proxy Proxy)chooseServer(ignoreList []string) *Server {
  var min = -1
  var minIndex = 0
  for index,server := range proxy.Servers {
    var skip = false
    for _, ignore := range ignoreList {
      if(ignore == server.Name){
        skip = true
        break
      }
    }

    if skip {
      continue
    }

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

func (proxy Proxy)ReverseProxy(w http.ResponseWriter, r *http.Request, server Server) (int, error){
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
    LogErr("connection refused")
    return 0, err
  }
  LogInfo("Recieved response: " + strconv.Itoa(resp.StatusCode))

  bodyBytes, err := ioutil.ReadAll(resp.Body)
  if err != nil {
    LogErr("Proxy: Failed to read response body")
    http.NotFound(w, r)
    return 0, err
  }

  buffer := bytes.NewBuffer(bodyBytes)
  for k, v := range resp.Header {
    w.Header().Set(k, strings.Join(v, ";"))
  }

  w.WriteHeader(resp.StatusCode)

  io.Copy(w, buffer)
  return resp.StatusCode, nil
}

func (proxy Proxy)attemptServers(w http.ResponseWriter, r *http.Request, ignoreList []string) {
  if float64(len(ignoreList)) >= math.Min(float64(3), float64(len(proxy.Servers))) {
    LogErr("Failed to find server for request")
    http.NotFound(w, r)
    return
  }

  var server = proxy.chooseServer(ignoreList)
  LogInfo("Got request: " + r.RequestURI)
  LogInfo("Sending to server: " + server.Name)

  server.Connections += 1
  _, err := proxy.ReverseProxy(w, r, *server)
  server.Connections -= 1

  if err != nil && strings.Contains(err.Error(), "connection refused") {
    LogWarn("Server did not respond: " + server.Name)

    proxy.attemptServers(w, r, append(ignoreList, server.Name))
    return
  }

  LogInfo("Responded to request successfuly")
}

func (proxy Proxy)handler(w http.ResponseWriter, r *http.Request) {
  proxy.attemptServers(w, r, []string{})
}
