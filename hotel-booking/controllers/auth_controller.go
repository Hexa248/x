package controllers

import "net/http"

func (a *App) LoginPage(w http.ResponseWriter, _ *http.Request) {
	render(w, "login.html", map[string]any{"Title": "Login Multi Role"})
}

func (a *App) RegisterPage(w http.ResponseWriter, _ *http.Request) {
	render(w, "register.html", map[string]any{"Title": "Register"})
}
