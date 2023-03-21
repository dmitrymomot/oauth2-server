{{define "messages"}}
{{if alLeastOneExists .errors .success}}
<div class="sm:col-span-2 mt-3">
  <div class="grid grid-cols-1 gap-y-3">
    {{template "error_messages" .}}
    {{template "success_messages" .}}
  </div>
</div>
{{end}}
{{end}}