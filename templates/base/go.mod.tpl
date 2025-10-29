module {{.ProjectName}}

go {{.GoVersion}}

require (
{{- range .GetAllImports}}
	{{.}} latest
{{- end}}
)
