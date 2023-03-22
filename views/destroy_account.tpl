{{ define "content"}}
<div class="text-center">
  {{include "partials/logo"}}
  {{if .form.Email}}
  <h2 class="text-3xl font-bold tracking-tight text-gray-900 sm:text-4xl">Confirm destroy account</h2>
  <p class="mt-4 text-lg leading-6 text-gray-500">Verification code was sent to your email.</p>
  {{else}}
  <h2 class="text-3xl font-bold tracking-tight text-gray-900 sm:text-4xl">Request destroy account</h2>
  <p class="mt-4 text-lg leading-6 text-gray-500">
    Enter your email address to send email with a destroy account confirmation code.
  </p>
  {{end}}
</div>
<div class="mt-8">

  {{if .form.Email}}

  <form action="/auth/account/destroy/verify" method="GET" role="form" id="form-destroy-account-link"
    class="grid grid-cols-1 gap-y-6 sm:grid-cols-2 sm:gap-x-8">

    {{template "messages" .}}

    <input type="hidden" name="email" value="{{.form.Email}}">
    <div class="sm:col-span-2"> {{template "otp" .}} </div>
    <div class="sm:col-span-2"> {{template "submit_button" "Destroy account"}} </div>
  </form>

  <div class="mt-4 text-center items-center">
    <span class="text-base text-gray-500"> Didn't receive the code? </span>
    <form action="/auth/account/destroy/request" method="POST" class="inline">
      <input type="hidden" name="email" value="{{.form.Email}}">
      <button type="submit" class="inline font-medium text-gray-700 underline underline-offset-4">Resend</button>
    </form>
  </div>

  {{else}}

  <form action="/auth/account/destroy/request" method="POST" role="form" id="form-destroy-account-request"
    class="grid grid-cols-1 gap-y-6 sm:grid-cols-2 sm:gap-x-8">

    {{template "messages" .}}

    <div class="sm:col-span-2"> {{template "email" .}} </div>
    <div class="sm:col-span-2"> {{template "submit_button" "Send confirmation code"}} </div>
  </form>

  {{end}}

</div>
{{end}}