{{ define "content"}}
<main class="flex-grow flex flex-col justify-center max-w-7xl w-full mx-auto mt-12 px-4 sm:px-6 lg:px-8">
  <div class="flex-shrink-0 flex justify-center">
    <svg xmlns="http://www.w3.org/2000/svg" class="h-24 w-24 text-green-500" fill="none" viewBox="0 0 24 24"
      stroke="currentColor" stroke-width="2">
      <path stroke-linecap="round" stroke-linejoin="round" d="M9 12l2 2 4-4m6 2a9 9 0 11-18 0 9 9 0 0118 0z" />
    </svg>
  </div>
  <div class="py-8">
    <div class="text-center">
      <p class="text-sm font-semibold text-gray-400 uppercase tracking-wide">Success</p>
      <h1 class="mt-2 text-2xl font-extrabold text-gray-900 tracking-tight sm:text-4xl">
        Recovery email sent
      </h1>
      <p class="mt-2 text-base text-gray-500">
        Password recovery instruction has been sent to <span class="text-gray-700 underline">{{if
          .form.Email}}{{.form.Email}}{{end}}</span>. Please check your email
        and follow the instructions inside.
      </p>
      <div class="mt-6 text-gray-600">
        Do not receive the email?
        <a href="/auth/password/recovery" class="text-base font-medium text-blue-600 hover:text-blue-500">Resend
          recovery email</a>
      </div>
    </div>
  </div>
</main>
{{end}}