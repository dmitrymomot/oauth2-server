{{define "email"}}
<label for="email" class="block text-sm font-medium text-gray-700">Email</label>
<div class="relative mt-1">
  {{if .validation.email}}
  <input id="email" name="email" type="email" autocomplete="email" aria-invalid="true" aria-describedby="email-error"
    class="block w-full rounded-md border-rose-300 py-3 px-4 shadow-sm text-rose-900 focus:border-rose-500 focus:ring-blue-500 focus:outline-none"
    value="{{if .form.Email}}{{.form.Email}}{{end}}">
  <div class="pointer-events-none absolute inset-y-0 right-0 flex items-center pr-3">
    <!-- Heroicon name: mini/exclamation-circle -->
    <svg class="h-5 w-5 text-rose-500" xmlns="http://www.w3.org/2000/svg" viewBox="0 0 20 20" fill="currentColor"
      aria-hidden="true">
      <path fill-rule="evenodd"
        d="M18 10a8 8 0 11-16 0 8 8 0 0116 0zm-8-5a.75.75 0 01.75.75v4.5a.75.75 0 01-1.5 0v-4.5A.75.75 0 0110 5zm0 10a1 1 0 100-2 1 1 0 000 2z"
        clip-rule="evenodd" />
    </svg>
  </div>
  {{else}}
  <input id="email" name="email" type="email" autocomplete="email"
    class="block w-full rounded-md border-gray-300 py-3 px-4 shadow-sm focus:border-blue-500 focus:ring-blue-500"
    value="{{if .form.Email}}{{.form.Email}}{{end}}">
  {{end}}
</div>
{{ if .validation.email }}
{{ range $key, $value := .validation.email }}
<p class="mt-2 text-sm text-rose-600" id="email-error-{{$key}}">{{$value}}</p>
{{ end }}
{{ end }}
{{end}}