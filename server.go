package main

type Server struct {
    name string
    scheme string
    host string
    port string
    connections int
}

func (server Server) Url() string {
  return server.scheme + "://" + server.host + ":" + server.port;
}
