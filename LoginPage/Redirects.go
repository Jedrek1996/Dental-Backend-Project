package main

import "net/http"

// Redirect Methods

type redirectFn func(res http.ResponseWriter, req *http.Request)

func redirectLoginPage(res http.ResponseWriter, req *http.Request) {
	http.Redirect(res, req, "/", http.StatusSeeOther)
}
func redirectUserPage(res http.ResponseWriter, req *http.Request) {
	http.Redirect(res, req, "/UserLoggedIn", http.StatusSeeOther)
}
func redirectLogout(res http.ResponseWriter, req *http.Request) {
	http.Redirect(res, req, "/UserLogOut", http.StatusSeeOther)
}
func redirectUserEditAppt(res http.ResponseWriter, req *http.Request) {
	http.Redirect(res, req, "/UserAppointment", http.StatusSeeOther)
}
func redirectDeleteUserApp(res http.ResponseWriter, req *http.Request) {
	http.Redirect(res, req, "/UserDeleteAppointment", http.StatusSeeOther)
}
func redirectUpdateUsername(res http.ResponseWriter, req *http.Request) {
	http.Redirect(res, req, "/UpdateUsername", http.StatusSeeOther)
}
func redirectUpdatePws(res http.ResponseWriter, req *http.Request) {
	http.Redirect(res, req, "/UpdatePws", http.StatusSeeOther)
}
func redirectAdminPage(res http.ResponseWriter, req *http.Request) {
	http.Redirect(res, req, "/AdminPage", http.StatusSeeOther)
}
func redirectAdminRemoveApptPage(res http.ResponseWriter, req *http.Request) {
	http.Redirect(res, req, "/AdminRemoveAppointment", http.StatusSeeOther)
}
func redirectAdminSessionPage(res http.ResponseWriter, req *http.Request) {
	http.Redirect(res, req, "/AdminSessionPage", http.StatusSeeOther)
}
func redirectUsersAdminPage(res http.ResponseWriter, req *http.Request) {
	http.Redirect(res, req, "/AdminUserPage", http.StatusSeeOther)
}
