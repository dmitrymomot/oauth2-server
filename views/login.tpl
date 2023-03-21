{{ define "content"}}
<div class="text-center">
  {{include "partials/logo"}}
  <h2 class="text-3xl font-bold tracking-tight text-gray-900 sm:text-4xl">Sign in to your account</h2>
  <p class="mt-4 text-lg leading-6 text-gray-500">Don't have an account? <a
      class="font-medium text-gray-700 underline underline-offset-4" href="/auth/register">Create a new one</a></p>
</div>
<div class="mt-12">
  <form action="" method="POST" role="form" id="form-login" class="grid grid-cols-1 gap-y-6 sm:grid-cols-2 sm:gap-x-8">

    {{template "messages" .}}

    <div class="sm:col-span-2"> {{template "email" .}} </div>
    <div class="sm:col-span-2"> {{template "password" .}} </div>
    <div class="sm:col-span-2"> {{template "submit_button" "Sign in"}} </div>

    <div class="sm:col-span-2 text-center items-center">
      <div class="text-base">
        <a href="/auth/password/recovery" class="font-medium text-gray-700 underline underline-offset-4">
          Forgot your password?
        </a>
      </div>
    </div>
  </form>
</div>
{{end}}