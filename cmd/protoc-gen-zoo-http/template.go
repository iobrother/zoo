package main

import (
	"embed"
	"io"
	"text/template"
)

//go:embed http.tpl
var Static embed.FS

var httpTemplate = template.Must(template.New("components").ParseFS(Static, "http.tpl")).
	Lookup("http.tpl")

type serviceDesc struct {
	ServiceType string // Greeter
	ServiceName string // helloworld.Greeter
	Metadata    string // api/v1/helloworld.proto
	Methods     []*methodDesc
	MethodSets  map[string]*methodDesc // unique because additional_bindings
}

type methodDesc struct {
	// method
	Name    string
	Num     int
	Request string
	Reply   string
	Comment string
	// http_rule
	Path         string
	Method       string
	HasParams    bool
	HasBody      bool
	Body         string
	ResponseBody string
}

func (s *serviceDesc) execute(w io.Writer) error {
	s.MethodSets = make(map[string]*methodDesc)
	for _, m := range s.Methods {
		s.MethodSets[m.Name] = m
	}
	return httpTemplate.Execute(w, s)
}
