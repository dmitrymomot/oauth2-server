{{define "terms_checkbox"}}
<div class="flex items-start">
  <input type="hidden" name="terms" :value="terms">
  <div class="flex-shrink-0">
    <!-- Enabled: "bg-blue-600", Not Enabled: "bg-gray-200" -->
    <button type="button" @click="terms = !terms; submitDisabled = !terms"
      :class="[ terms ? 'bg-blue-600' : 'bg-gray-200' ]"
      class="relative inline-flex h-6 w-11 flex-shrink-0 cursor-pointer rounded-full border-2 border-transparent transition-colors duration-200 ease-in-out focus:outline-none focus:ring-2 focus:ring-blue-500 focus:ring-offset-2"
      role="switch" aria-checked="false">
      <span class="sr-only">Agree to
        <a target="_blank" href="{{.termsOfUseURL}}" class="text-blue-600 hover:text-blue-500">terms of use</a>,
        <a target="_blank" href="{{.privacyPolicyURL}}" class="text-blue-600 hover:text-blue-500">privacy policy</a>
        and
        <a target="_blank" href="{{.cookiesPolicyURL}}" class="text-blue-600 hover:text-blue-500">cookie policy</a>
      </span>
      <!-- Enabled: "translate-x-5", Not Enabled: "translate-x-0" -->
      <span aria-hidden="true" :class="[ terms ? 'translate-x-5' : 'translate-x-0' ]"
        class="translate-x-0 inline-block h-5 w-5 transform rounded-full bg-white shadow ring-0 transition duration-200 ease-in-out"></span>
    </button>
  </div>
  <div class="ml-3">
    <p class="text-base text-gray-500">
      By selecting this, you agree to the
      <a target="_blank" href="{{.termsOfUseURL}}" class="font-medium text-gray-700 underline underline-offset-4">Terms
        of Use</a>,
      <a target="_blank" href="{{.privacyPolicyURL}}"
        class="font-medium text-gray-700 underline underline-offset-4">Privacy Policy</a>
      and
      <a target="_blank" href="{{.cookiesPolicyURL}}"
        class="font-medium text-gray-700 underline underline-offset-4">Cookie Policy</a>.
    </p>
  </div>
</div>
{{ if .validation.terms }}
{{ range $key, $value := .validation.terms }}
<p class="mt-2 text-sm text-rose-600" id="terms-error-{{$key}}">{{$value}}</p>
{{ end }}
{{ end }}
{{end}}