package output

import (
	_ "embed"
	"fmt"
	"io"
	"net/http"
	"strings"
	"text/template"

	"github.com/barasher/dep-carto/internal/model"
)

//go:embed dotTemplate.tpl
var dotTpl string

type DotGraph struct {
	Servers         []model.Server
	ExternalServers []string
	Dependencies    []DotGraphDep
}

type DotGraphDep struct {
	From  string
	To    string
	Label string
}

type DotFormatter struct{}

func NewDotFormatter() DotFormatter {
	return DotFormatter{}
}

func (DotFormatter) FormatToWriter(s []model.Server, w io.Writer) error {
	dg := newDotGraph(s)
	funcs := template.FuncMap{"join": strings.Join}
	tpl, err := template.New("test").Funcs(funcs).Parse(dotTpl)
	if err != nil {
		return fmt.Errorf("error while parsing template: %w", err)
	}
	if err := tpl.Execute(w, dg); err != nil {
		return fmt.Errorf("error while applying template: %w", err)
	}
	return nil
}

func (d DotFormatter) Format(s []model.Server, w http.ResponseWriter) {
	w.Header().Add("Content-Type", "text/plain")
	w.WriteHeader(http.StatusOK)
	if err := d.FormatToWriter(s, w); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
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
			if ref, found := identifierToHost[curD.Resource]; found {
				dg.Dependencies = append(dg.Dependencies, DotGraphDep{curS.Hostname, ref, curD.Label})
			} else {
				dg.ExternalServers = append(dg.ExternalServers, curD.Resource)
				dg.Dependencies = append(dg.Dependencies, DotGraphDep{curS.Hostname, curD.Resource, curD.Label})
			}
		}
	}

	return dg
}
