# OAuth2 server

[![GitHub tag (latest SemVer)](https://img.shields.io/github/tag/dmitrymomot/oauth2-server)](https://github.com/dmitrymomot/oauth2-server)
[![Tests](https://github.com/dmitrymomot/oauth2-server/actions/workflows/tests.yml/badge.svg)](https://github.com/dmitrymomot/oauth2-server/actions/workflows/tests.yml)
[![CodeQL Analysis](https://github.com/dmitrymomot/oauth2-server/actions/workflows/codeql-analysis.yml/badge.svg)](https://github.com/dmitrymomot/oauth2-server/actions/workflows/codeql-analysis.yml)
[![Go Report Card](https://goreportcard.com/badge/github.com/dmitrymomot/oauth2-server)](https://goreportcard.com/report/github.com/dmitrymomot/oauth2-server)
[![Go Reference](https://pkg.go.dev/badge/github.com/dmitrymomot/oauth2-server.svg)](https://pkg.go.dev/github.com/dmitrymomot/oauth2-server)
[![License](https://img.shields.io/github/license/dmitrymomot/oauth2-server)](https://github.com/dmitrymomot/oauth2-server/blob/main/LICENSE)

This is a simple OAuth2 server implementation in Go.

## Features

- [x] Implements the [OAuth2 Authorization Framework](http://tools.ietf.org/html/rfc6749)
- [x] Implements the [OAuth2 Token Revocation](http://tools.ietf.org/html/rfc7009) extension
- [x] Implements the [OAuth2 Token Introspection](http://tools.ietf.org/html/rfc7662) extension
- [x] Signin/Signup pages
- [x] Reset password flow
- [x] API to create and manage clients
- [x] API to manage user data