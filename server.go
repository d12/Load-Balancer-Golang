package main

type Server struct {
    Name string
    Scheme string
    Host string
    Port string
    Connections int
}

func (server Server) Url() string {
  return server.Scheme + "://" + server.Host + ":" + server.Port;
}
