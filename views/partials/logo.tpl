{{ if .AppLogo }}
<img class="mx-auto h-12 w-auto mb-4" src="{{.AppLogo}}" alt="{{if .AppName}}{{.AppName}}{{else}}logo{{end}}">
{{ else }}
<div class="py-6"></div>
{{ end }}