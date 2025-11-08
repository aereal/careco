package test

import "github.com/99designs/gqlgen/graphql/handler"

func provideTestServer(srv *handler.Server) *TestHandler { return &TestHandler{Server: srv} }

type TestHandler struct {
	*handler.Server
}
