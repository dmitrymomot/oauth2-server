{{ define "content"}}
<div class="text-center">
  {{include "partials/logo"}}
  <h2 class="text-3xl font-bold tracking-tight text-gray-900 sm:text-4xl">Create a new account</h2>
  <p class="mt-4 text-lg leading-6 text-gray-500">Already have an account? Just <a
      class="font-medium text-gray-700 underline underline-offset-4" href="/auth/login">sign in</a></p>
</div>
<div class="mt-12">
  <form action="/auth/register" method="POST" role="form" id="form-register"
    x-data="{submitDisabled: true, terms: false}" class="grid grid-cols-1 gap-y-6 sm:grid-cols-2 sm:gap-x-8">

    {{template "messages" .}}

    <div class="sm:col-span-2"> {{template "email" .}} </div>
    <div> {{template "password" .}} </div>
    <div> {{template "password_confirmation" .}} </div>
    <div class="sm:col-span-2"> {{template "terms_checkbox" .}} </div>
    <div class="sm:col-span-2"> {{template "submit_button" "Sign up"}} </div>

    {{if .return_uri}}
    <div class="sm:col-span-2 text-center items-center">
      <p class="text-base text-gray-500">
        Or <a href="{{.return_uri}}" class="font-medium text-gray-700 underline underline-offset-4">go back</a> to
        previous page.
      </p>
    </div>
    {{end}}
  </form>
</div>
{{end}}