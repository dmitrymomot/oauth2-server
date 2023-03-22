{{ if appLogo }}
<img class="mx-auto h-12 w-auto mb-4" src="{{appLogo}}" alt="{{if appName}}{{appName}}{{else}}logo{{end}}">
{{ else }}
<div class="py-6"></div>
{{ end }}