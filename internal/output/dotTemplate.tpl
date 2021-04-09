digraph g {

graph [
rankdir = "LR"
];

{{range .Servers}}
"{{.Hostname}}"[
    label="{{.Hostname}}|{{join .IPs "\\n"}}"
    shape = "record"
]
{{end}}

{{range .ExternalServers}}
"{{.}}"
{{end}}

{{range .Dependencies}}
"{{.From}}" -> "{{.To}}" {{if ne .Label ""}}[label="{{.Label}}"]{{end}}
{{end}}

}