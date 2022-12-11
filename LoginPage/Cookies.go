package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"

	uuid "github.com/satori/go.uuid"
)

type cookieDetails struct {
	Name  string
	Value string
}

type cookieDataStruct struct {
	Name  []string
	Value []string
}

var cookiesData []cookieDetails

//--- JSON COOKIES FUNCTION ---

// Get new session cookies from Cookies JSON ---✨
func getCookiesJSON() []cookieDetails {

	file, err := ioutil.ReadFile("cookies.json")
	if err != nil {
		log.Fatal(err)
	}

	json.Unmarshal(file, &cookiesData)

	return cookiesData
}

// Add new session cookies from Cookies JSON ---✨
func addCookieToJSON(newCookie cookieDetails) {

	cookies := getCookiesJSON()
	cookies = append(cookies, newCookie)
	bytes, err := json.MarshalIndent(cookies, "", "  ")

	if err != nil {
		log.Fatalln(err)
	}
	_ = ioutil.WriteFile("cookies.json", bytes, 0644)

}

// Delete session cookies (logout) ---✨
func deleteCookieFromJSON(res http.ResponseWriter, req *http.Request, cookie cookieDetails) {

	cookies := getCookiesJSON()
	for _, v := range cookies { //&& v.Expires == cookie.Expires
		if v.Name == cookie.Name && v.Value == cookie.Value {
			removeCookieFromJson(res, req, cookie.Name)
			break
		}
	}
	// wg.Done()
}

// Remove upon session cookie Logout from cookies JSON (logout) ---✨
func removeCookieFromJson(res http.ResponseWriter, req *http.Request, remove string) {

	cookies := getCookiesJSON()

	myCookie, err := req.Cookie(remove)
	if err != nil {
		log.Fatalln(err)
	}

	for i := 0; i < len(cookies); i++ {
		cookie := cookies[i]
		if cookie.Name == myCookie.Name {
			cookies = append(cookies[:i], cookies[i+1:]...)
			i--
			fmt.Println("Removed:", remove)
			break
		}
	}

	bytes, err := json.MarshalIndent(cookies, "", "  ")
	if err != nil {
		log.Fatalln(err)
	}
	_ = ioutil.WriteFile("cookies.json", bytes, 0644)
}

// --- COOKIES FUNCTIONS ---

// Set cookie when logged in ---✨
func setCookie(res http.ResponseWriter, req *http.Request, currentUser string) {
	fmt.Println("Setting Cookie!")
	id := uuid.NewV4()

	http.SetCookie(res, &http.Cookie{
		Name:  currentUser,
		Value: id.String(),
	})

	cookieTest := cookieDetails{Name: currentUser, Value: id.String()}
	addCookieToJSON(cookieTest)

}

// Delete cookie upon logout ---✨
func deleteCookie(res http.ResponseWriter, req *http.Request, currentUser string) {
	fmt.Println("Cookie Deleted")

	myCookie, err := req.Cookie(currentUser)
	fmt.Println(myCookie)
	if err != nil {
		log.Fatalln(err)
		return
	}
	myCookie.MaxAge = -1
	http.SetCookie(res, myCookie)
}

// Check if user is logged in else redirect ---✨
func checkLoggedIn(res http.ResponseWriter, req *http.Request, user string) {
	loggedIn := checkCookieExists(req, user)

	if loggedIn {
		return
	} else if !loggedIn {
		redirectLoginPage(res, req)
		return
	}
}

// Used to redirect user when they acces login/signup page when logged in ---✨
func redirectAccessWithCookie(res http.ResponseWriter, req *http.Request, user string) {

	loggedIn := checkCookieExists(req, user)

	if loggedIn {
		if !strings.EqualFold(user, "Admin") {
			redirectUserPage(res, req)
		} else if strings.EqualFold(user, "Admin") {
			redirectAdminPage(res, req)
		}
	}
}

// Check if cookie exists ---✨
func checkCookieExists(req *http.Request, user string) bool {

	myCookie, err := req.Cookie(user)
	if err != nil {
		return false
	}

	cookiesData := getCookiesJSON()
	for _, v := range cookiesData {
		if v.Name == user && myCookie.Name == user {
			return true
		}
	}
	return false
}

// Returns the cookie struct to pass into Admin Session Page gohtml ---✨
func getCookiesStruct() cookieDataStruct {

	var cookiesName []string
	var cookiesValue []string

	cookiesData := getCookiesJSON()

	for _, v := range cookiesData {
		cookiesName = append(cookiesName, v.Name)
		cookiesValue = append(cookiesValue, v.Value)
	}

	cookies := cookieDataStruct{
		Name:  cookiesName,
		Value: cookiesValue,
	}

	return cookies
}
