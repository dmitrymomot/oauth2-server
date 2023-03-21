{{ define "content"}}
<div class="text-center">
  {{include "partials/logo"}}
  {{if .form.Email}}
  <h2 class="text-3xl font-bold tracking-tight text-gray-900 sm:text-4xl">Set new password</h2>
  {{else}}
  <h2 class="text-3xl font-bold tracking-tight text-gray-900 sm:text-4xl">Password recovery</h2>
  {{end}}
  <p class="mt-4 text-lg leading-6 text-gray-500">Or back to <a
      class="font-medium text-gray-700 underline underline-offset-4" href="/auth/login">sign in</a> page
  </p>
</div>
<div class="mt-8">

  {{if .form.Email}}

  <form action="/auth/password/reset" method="POST" role="form" id="form-verification"
    class="grid grid-cols-1 gap-y-6 sm:grid-cols-2 sm:gap-x-8">

    {{template "messages" .}}

    <input type="hidden" name="email" value="{{.form.Email}}">
    <input type="hidden" name="otp" value="{{.form.OTP}}">
    <div> {{template "password" .}} </div>
    <div> {{template "password_confirmation" .}} </div>
    <div class="sm:col-span-2"> {{template "submit_button" "Reset password"}} </div>

    <div class="sm:col-span-2 text-center items-center">
      <p class="text-base text-gray-500">
        Do not receive the email? <a href="/auth/password/recovery"
          class="font-medium text-gray-700 underline underline-offset-4">Resend recovery email</a></p>
    </div>
  </form>

  {{else}}

  <form action="/auth/password/recovery" method="POST" role="form" id="form-verification"
    class="grid grid-cols-1 gap-y-6 sm:grid-cols-2 sm:gap-x-8">

    {{template "messages" .}}

    <div class="sm:col-span-2"> {{template "email" .}} </div>
    <div class="sm:col-span-2"> {{template "submit_button" "Send password reset instruction"}} </div>

    <div class="sm:col-span-2 text-center items-center">
      <p class="text-base text-gray-500">
        Or <a href="/auth/login" class="font-medium text-gray-700 underline underline-offset-4">go back</a>
        to the login page
      </p>
    </div>
  </form>

  {{end}}

</div>
{{end}}