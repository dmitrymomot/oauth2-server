{{define "otp"}}
<label for="otp" class="block text-sm font-medium text-gray-700">Verification code</label>
<div class="relative mt-1">
  {{if .validation.otp}}
  <input id="otp" name="otp" type="text" autocomplete="otp" aria-invalid="true" aria-describedby="otp-error"
    class="block w-full rounded-md border-rose-300 py-3 px-4 shadow-sm text-rose-900 focus:border-rose-500 focus:ring-blue-500 focus:outline-none"
    value="{{if .form.OTP}}{{.form.OTP}}{{end}}">
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
  <input id="otp" name="otp" type="text" autocomplete="otp"
    class="block w-full rounded-md border-gray-300 py-3 px-4 shadow-sm focus:border-blue-500 focus:ring-blue-500"
    value="{{if .form.OTP}}{{.form.OTP}}{{end}}">
  {{end}}
</div>
{{ if .validation.otp }}
{{ range $key, $value := .validation.otp }}
<p class="mt-2 text-sm text-rose-600" id="otp-error-{{$key}}">{{$value}}</p>
{{ end }}
{{ end }}
{{end}}