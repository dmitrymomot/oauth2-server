{{define "content"}}
<main class="flex-grow flex flex-col justify-center max-w-7xl w-full mx-auto sm:mt-12 px-4 sm:px-6 lg:px-8">
  <div class="flex-shrink-0 flex justify-center">
    <svg xmlns="http://www.w3.org/2000/svg" class="h-24 w-24 text-green-500" fill="none" viewBox="0 0 24 24"
      stroke="currentColor" stroke-width="2">
      <path stroke-linecap="round" stroke-linejoin="round" d="M9 12l2 2 4-4m6 2a9 9 0 11-18 0 9 9 0 0118 0z" />
    </svg>
  </div>
  <div class="py-8">
    <div class="text-center">
      <p class="text-sm font-semibold text-gray-400 uppercase tracking-wide">Success</p>
      <h1 class="mt-2 text-3xl font-extrabold text-gray-900 tracking-tight sm:text-4xl">Password has been changed.
      </h1>
      <p class="mt-2 text-base text-gray-500">
        The new password has been set. Use it to login to your account.
      </p>
      <div class="mt-6">
        <a href="/auth/login" class="text-base font-medium text-blue-600 hover:text-blue-500">Go to login
          page<span aria-hidden="true"> &rarr;</span></a>
      </div>
    </div>
  </div>
</main>
{{end}}