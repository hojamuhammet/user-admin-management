package server

import "net/http"

type Server struct {
	HttpServer *http.Server
}
