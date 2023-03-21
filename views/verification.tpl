{{ define "content"}}
<div class="text-center">
  {{include "partials/logo"}}
  {{if .form.Email}}
  <h2 class="text-3xl font-bold tracking-tight text-gray-900 sm:text-4xl">Verify your email</h2>
  <p class="mt-4 text-lg leading-6 text-gray-500">Verification code was sent to your email.</p>
  {{else}}
  <h2 class="text-3xl font-bold tracking-tight text-gray-900 sm:text-4xl">Resend verification code</h2>
  <p class="mt-4 text-lg leading-6 text-gray-500">
    Enter your email address to resend email with a new verification code
  </p>
  {{end}}
</div>
<div class="mt-8">

  {{if .form.Email}}

  <form action="/verification/link" method="GET" role="form" id="form-verification"
    class="grid grid-cols-1 gap-y-6 sm:grid-cols-2 sm:gap-x-8">

    {{template "messages" .}}

    <input type="hidden" name="email" value="{{.form.Email}}">
    <div class="sm:col-span-2"> {{template "otp" .}} </div>
    <div class="sm:col-span-2"> {{template "submit_button" "Verify email"}} </div>
  </form>

  <div class="mt-4 text-center items-center">
    <span class="text-base text-gray-500"> Didn't receive the code? </span>
    <form action="/verification/resend" method="POST" class="inline">
      <input type="hidden" name="email" value="{{.form.Email}}">
      <button type="submit" class="inline font-medium text-gray-700 underline underline-offset-4">Resend</button>
    </form>
  </div>

  {{else}}

  <form action="/verification/resend" method="POST" role="form" id="form-verification"
    class="grid grid-cols-1 gap-y-6 sm:grid-cols-2 sm:gap-x-8">

    {{template "messages" .}}

    <div class="sm:col-span-2"> {{template "email" .}} </div>
    <div class="sm:col-span-2"> {{template "submit_button" "Resend verification code"}} </div>
  </form>

  {{end}}

</div>
{{end}}