<!DOCTYPE html>
<html lang="en" class="h-full">

<head>
  <title>{{if .page_title}}{{.page_title}}{{else}}OAuth2 Service{{end}}</title>

  {{include "partials/head"}}
</head>

<body class="h-full">
  <div class="overflow-hidden bg-white py-16 px-4 sm:px-6 lg:px-8 lg:py-24">
    <div class="relative mx-auto max-w-lg">
      {{include "partials/bgsq"}}
      {{template "content" .}}
    </div>
  </div>
</body>

</html>