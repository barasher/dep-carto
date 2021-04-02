package output

import (
	_ "embed"
	"fmt"
	"github.com/barasher/dep-carto/internal/model"
	"net/http"
	"strings"
	"text/template"
)

//go:embed dot.tpl
var dotTpl string

type DotGraph struct {
	Servers         []model.Server
	ExternalServers []string
	Dependencies    []DotGraphDep
}

type DotGraphDep struct {
	From string
	To   string
}

type DotFormatter struct{}

func NewDotFormatter() DotFormatter {
	return DotFormatter{}
}

func (DotFormatter) Format(s []model.Server, w http.ResponseWriter) {
	dg := newDotGraph(s)
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	funcs := template.FuncMap{"join": strings.Join}
	tpl, err := template.New("test").Funcs(funcs).Parse(dotTpl)
	if err != nil {
		http.Error(w, fmt.Sprintf("error while parsing template: %v", err), http.StatusInternalServerError)
		return
	}
	if err := tpl.Execute(w, dg); err != nil {
		http.Error(w, fmt.Sprintf("error while applying template: %v", err), http.StatusInternalServerError)
		return
	}
	return
}

func newDotGraph(s []model.Server) DotGraph {
	dg := DotGraph{
		Dependencies:    []DotGraphDep{},
		ExternalServers: []string{},
		Servers:         []model.Server{},
	}

	hostnameToServer := make(map[string]model.Server)
	identifierToHost := make(map[string]string)
	for _, curS := range s {
		hostnameToServer[curS.Hostname] = model.Server{
			Hostname: curS.Hostname,
			IPs:      curS.IPs,
		}
		identifierToHost[curS.Hostname] = curS.Hostname
		for _, curD := range curS.IPs {
			identifierToHost[curD] = curS.Hostname
		}
	}

	for _, curS := range hostnameToServer {
		dg.Servers = append(dg.Servers, curS)
	}

	for _, curS := range s {
		for _, curD := range curS.Dependencies {
			if ref, found := identifierToHost[curD]; found {
				dg.Dependencies = append(dg.Dependencies, DotGraphDep{curS.Hostname, ref})
			} else {
				dg.ExternalServers = append(dg.ExternalServers, curD)
				dg.Dependencies = append(dg.Dependencies, DotGraphDep{curS.Hostname, curD})
			}
		}
	}

	return dg
}
