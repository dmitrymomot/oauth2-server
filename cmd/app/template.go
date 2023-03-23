package main

import (
	"html/template"
	"time"

	"github.com/foolin/goview"
)

func initTemplateEngine() {
	gv := goview.New(goview.Config{
		Root:      "views",
		Extension: ".tpl",
		Master:    "layouts/master",
		Partials: []string{
			"partials/input/email",
			"partials/input/password",
			"partials/input/password_confirmation",
			"partials/input/otp",
			"partials/input/submit",
			"partials/input/terms",
			"partials/messages/errors",
			"partials/messages/success",
			"partials/messages",
		},
		Funcs: template.FuncMap{
			"copy": func() string {
				return time.Now().Format("2006")
			},
			"alLeastOneExists": func(in ...interface{}) bool {
				for _, i := range in {
					if i != nil {
						return true
					}
				}
				return false
			},
			"appLogo": func() string {
				return appLogoURL
			},
			"appName": func() string {
				return productName
			},
		},
		DisableCache: true,
	})

	// Set new instance
	goview.Use(gv)
}
