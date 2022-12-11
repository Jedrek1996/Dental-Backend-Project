// Main File
package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"regexp"
	"strconv"
	"sync"

	"github.com/tawesoft/golib/v2/dialog"
)

// Templates (Tpl) 😎
// Login Page/ Sign Up Tpl
var tpl *template.Template
var signupTpl *template.Template

// Admin Tpl
var adminTpl *template.Template
var adminRemoveApptTpl *template.Template
var adminSessionTpl *template.Template
var adminUserTpl *template.Template

// User Tpl
var userLoggedInTpl *template.Template

// Store username for other methods
var currentUser string

// Error Logger
var Error *log.Logger

// Wait Group
var wg sync.WaitGroup
var muFile sync.Mutex

func init() {

	//Set Error Path (Change to this when submitting**)
	file, err := os.OpenFile("errors.txt",
		os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)

	if err != nil {
		log.Fatalln("Failed to open error log file:", err)
	}

	Error = log.New(io.MultiWriter(file, os.Stderr),
		"ERROR: ",
		log.Ldate|log.Ltime|log.Lshortfile)

	//Set custom error log
	log.SetPrefix("TRACE ERROR: ")
	log.SetFlags(log.Llongfile)

	//Login Template
	tpl = template.Must(template.ParseGlob("./LoginPage/LoginPageTemplates/LoginPage/*"))

	//Signup Template
	signupTpl = template.Must(template.ParseGlob("./LoginPage/LoginPageTemplates/SignUpPage/*"))

	//User Template
	userLoggedInTpl = template.Must(template.ParseGlob("./LoginPage/LoginPageTemplates/UserPage/*"))

	//Admin Templates
	adminTpl = template.Must(template.ParseGlob("./LoginPage/LoginPageTemplates/AdminPage/MainAdminPage/*"))
	adminRemoveApptTpl = template.Must(template.ParseGlob("./LoginPage/LoginPageTemplates/AdminPage/AppointmentRemovePage/*"))
	adminSessionTpl = template.Must(template.ParseGlob("./LoginPage/LoginPageTemplates/AdminPage/SessionAdminPage/*"))
	adminUserTpl = template.Must(template.ParseGlob("./LoginPage/LoginPageTemplates/AdminPage/UserAdminPage/*"))

}

func main() {
	//Link CSS file
	fs := http.FileServer(http.Dir("./LoginPage"))
	http.Handle("/LoginPage/", http.StripPrefix("/LoginPage/", fs))

	//Directory

	//Login Page
	http.HandleFunc("/", loginPageFunc)

	//Sign Up
	http.HandleFunc("/LoginPageSignUp", signUpPageFunc)

	//Admin Pages
	http.HandleFunc("/AdminPage", adminPageAddAppt)

	http.HandleFunc("/AdminRemoveAppointment", adminRemoveAppt)

	http.HandleFunc("/AdminSessionPage", adminSessionPage)

	http.HandleFunc("/AdminUserPage", userAdminPage)

	//User Pages
	http.HandleFunc("/UserLoggedIn", userLoggedInFunc)

	http.HandleFunc("/UserAppointment", userEditAppts)

	http.HandleFunc("/UserDeleteAppointment", userDeleteAppts)

	http.HandleFunc("/UpdateUsername", userUpdateUsername)

	http.HandleFunc("/UpdatePws", userUpdatePws)

	http.HandleFunc("/UserLogOut", userLogOut)

	http.Handle("/favicon.ico", http.NotFoundHandler())

	http.ListenAndServe(":8080", nil)

}

// Login Page ---✨
func loginPageFunc(res http.ResponseWriter, req *http.Request) {

	redirectAccessWithCookie(res, req, currentUser) //Used to redirect user/admin who are still logged in.

	userLogin := UserLoginStruct{LoggedIn: "Empty"}

	if req.Method == "GET" {
		err := tpl.ExecuteTemplate(res, "MainLoginPage.gohtml", userLogin)

		if err != nil {
			log.Fatalln(err)
		}
	} else if req.Method == "POST" {

		userName := req.FormValue("loginUsername")
		userPassword := req.FormValue("loginPassword")

		file, err := ioutil.ReadFile("usersData.json")
		if err != nil {
			log.Fatalln(err)
		}

		json.Unmarshal(file, &userData)

		for _, v := range userData {
			if v.LoginUsername == userName {
				//Username Exists
				match := CheckStringHash(userPassword, v.LoginPassword)
				if match {
					//Logged In
					currentUser = userName
					setCookie(res, req, currentUser)
					checkIsUserRedirect(res, req, currentUser, redirectAdminPage, redirectLoginPage)
					return

				} else {
					//Wrong Password
					userLogin.LoggedIn = "Wrong"
				}
			} else {
				userLogin.LoggedIn = "Wrong"
			}
		}

		tpl.ExecuteTemplate(res, "MainLoginPage.gohtml", userLogin)
		defer redirectLoginPage(res, req)
	}
}

// Sign Up Page ---✨
func signUpPageFunc(res http.ResponseWriter, req *http.Request) {
	fmt.Println("Signing Up Page")

	redirectAccessWithCookie(res, req, currentUser)

	err := signupTpl.ExecuteTemplate(res, "MainSignUpPage.gohtml", "")

	if err != nil {
		log.Fatalln(err)
	}

	newUserName := req.FormValue("signupUsername")
	newUserPassword := req.FormValue("signupPassword")

	//Regex
	userNameMatch, err := regexp.MatchString("^[a-zA-Z0-9]{4,20}$", newUserName)
	if err != nil {
		log.Fatalln(err)
	}

	passWordMatch, err := regexp.MatchString("^[a-zA-Z0-9]{4,20}$", newUserPassword)
	if err != nil {
		log.Fatalln(err)
	}

	hashedPassword, _ := HashString(newUserPassword)
	addUser := loginDetails{LoginUsername: newUserName, LoginPassword: hashedPassword}

	if len(newUserName) > 0 && len(newUserPassword) > 0 {

		if userNameMatch && passWordMatch {

			//Check if username exists
			checkUsername := addUsersToJson(addUser)

			if checkUsername {
				dialog.Alert("Account created!")
			} else if !checkUsername {
				dialog.Alert("Please select another username!")
			}

		} else if !userNameMatch || !passWordMatch {
			dialog.Alert("Username and Password should have more than 4 characters and no special characters!")
		}
	}
}

// User Logged In (Not admin) ---✨
func userLoggedInFunc(res http.ResponseWriter, req *http.Request) {

	appointmentsData := appDataStruct()

	if req.Method == "GET" {
		err := userLoggedInTpl.ExecuteTemplate(res, "MainUserPage.gohtml", &appointmentsData)
		checkLoggedIn(res, req, currentUser)

		if err != nil {
			log.Fatalln(err)
			return
		}

		getAvailableAppointments()

	} else if req.Method == "POST" {
		appointmentType := req.FormValue("appointmentType")
		dateSelect := req.FormValue("dateSelect")
		timeSelect := req.FormValue("timeSelect")

		userBookedApptDetails := AppointmentsDetail{AppointmentType: appointmentType, Date: dateSelect, Time: timeSelect, User: currentUser}

		wg.Add(1)
		duplicateAppt := addUserAppointment(userBookedApptDetails)
		wg.Wait()

		if duplicateAppt {
			dialog.Alert("Please choose another appointment date!")
			redirectUserPage(res, req)
			return
		} else if !duplicateAppt {
			dialog.Alert("Appt Added!")
			redirectUserPage(res, req)
		}
	}
}

// User Edit Appts Page ---✨
func userEditAppts(res http.ResponseWriter, req *http.Request) {

	checkLoggedIn(res, req, currentUser)

	appointmentsData := appDataStruct()

	if req.Method == "GET" {
		err := userLoggedInTpl.ExecuteTemplate(res, "MainUserEditPage.gohtml", &appointmentsData)

		if err != nil {
			log.Fatalln(err)
			return
		}

	} else if req.Method == "POST" {

		currentApptType := req.FormValue("currentApptType")
		currentDateSelect := req.FormValue("currentDateSelect")
		currentTimeSelect := req.FormValue("currentTimeSelect")

		newAppointmentType := req.FormValue("newAppointmentType")
		newDateSelect := req.FormValue("newDateSelect")
		newTimeSelect := req.FormValue("newTimeSelect")

		prevApptDetails := AppointmentsDetail{AppointmentType: currentApptType, Date: currentDateSelect, Time: currentTimeSelect, User: currentUser}

		newApptDetails := AppointmentsDetail{AppointmentType: newAppointmentType, Date: newDateSelect, Time: newTimeSelect, User: currentUser}

		//currentApptis to check whether the currentAppt is available in user's bookedAppointment
		//avaliability is not booked in user's bookedAppointment
		currentApptChannel := make(chan bool)
		avaliabilityChannel := make(chan bool)

		go checkCurrentAppointment(currentApptType, currentDateSelect, currentTimeSelect, currentUser, currentApptChannel)

		go checkForAvaliableAppt(newAppointmentType, newDateSelect, newTimeSelect, avaliabilityChannel)

		currentAppt := <-currentApptChannel
		avaliability := <-avaliabilityChannel

		if currentAppt && avaliability {

			wg.Add(2)
			go addUserAppointment(newApptDetails)
			go deleteAppointment(prevApptDetails)
			wg.Wait()

			redirectUserEditAppt(res, req)
			dialog.Alert("Appointment Updated")
		} else {
			if !currentAppt && !avaliability {
				dialog.Alert("Current appointment selection is not booked and slot is not avaliable!")
			} else if !currentAppt {
				dialog.Alert("Current appointment selection does not exist!")
			} else if !avaliability {
				dialog.Alert("New selected slot is not available please select another slot!")
			}
			redirectUserEditAppt(res, req)
		}

	}

}

// User Delete Apps Page ---✨
func userDeleteAppts(res http.ResponseWriter, req *http.Request) {

	checkLoggedIn(res, req, currentUser)

	appointmentsData := appDataStruct()

	if req.Method == "GET" {
		err := userLoggedInTpl.ExecuteTemplate(res, "MainUserDeletePage.gohtml", &appointmentsData)

		if err != nil {
			log.Fatalln(err)
			return
		}

	} else if req.Method == "POST" {
		fmt.Println("Deleting")
		deleteApptType := req.FormValue("deleteApptType")
		deleteDateSelect := req.FormValue("deleteDateSelect")
		deleteTimeSelect := req.FormValue("deleteTimeSelect")

		prevApptDetails := AppointmentsDetail{AppointmentType: deleteApptType, Date: deleteDateSelect, Time: deleteTimeSelect, User: currentUser}

		currentApptChannel2 := make(chan bool)

		go checkCurrentAppointment(deleteApptType, deleteDateSelect, deleteTimeSelect, currentUser, currentApptChannel2)
		currentAppt := <-currentApptChannel2

		fmt.Println(currentAppt)

		if currentAppt {
			muFile.Lock()
			deleteAppointment(prevApptDetails)
			defer muFile.Unlock()
			dialog.Alert("Your appointment has been deleted")
			defer redirectDeleteUserApp(res, req)
		} else if !currentAppt {
			dialog.Alert("Your appointment does not exist!")
			defer redirectDeleteUserApp(res, req)
		}
	}
}

// User Update Username ---✨
func userUpdateUsername(res http.ResponseWriter, req *http.Request) {

	checkLoggedIn(res, req, currentUser)

	if req.Method == "GET" {
		err := userLoggedInTpl.ExecuteTemplate(res, "MainUsernameUpdate.gohtml", "")

		if err != nil {
			log.Fatalln(err)
			return
		}
	} else if req.Method == "POST" {
		currPws := req.FormValue("currPws")
		newUserName := req.FormValue("newUserName")
		userList := getAllUsers()

		var updated bool

		for _, v := range userList {

			if v.LoginUsername == currentUser {
				match := CheckStringHash(currPws, v.LoginPassword)
				if match {
					userExist := updateUser(currentUser, match, newUserName) //Update username returns bool (if userExist will not update)
					userRegexMatch, err := regexp.MatchString("^[a-zA-Z0-9]{4,20}$", newUserName)

					if err != nil {
						log.Fatalln(err)
					}

					if userExist {
						dialog.Alert("Selected new Username Exists! Please choose another one!")
						redirectUpdateUsername(res, req)
						break
					} else if !userExist && userRegexMatch {
						dialog.Alert("Username Updated! Please login again!") //Updated
						updated = true
						userLogOut(res, req)
						break
					}
				}
			}
		}
		if !updated {
			dialog.Alert("Please re-enter your Username and Password! No special characters allowed")
		}
	}
}

// User Update Password ---✨
func userUpdatePws(res http.ResponseWriter, req *http.Request) {

	checkLoggedIn(res, req, currentUser)

	if req.Method == "GET" {
		err := userLoggedInTpl.ExecuteTemplate(res, "MainPwsUpdate.gohtml", "")

		if err != nil {
			log.Fatalln(err)
			return
		}

	} else if req.Method == "POST" {
		currUser := req.FormValue("currUser")
		currPws := req.FormValue("currPws")
		newPws := req.FormValue("newPws")

		userList := getAllUsers()

		var updated bool

		for _, v := range userList {
			if v.LoginUsername == currentUser {
				match := CheckStringHash(currPws, v.LoginPassword)
				if match {

					pwsRegexMatch, err := regexp.MatchString("^[a-zA-Z0-9]{4,20}$", newPws)

					if err != nil {
						log.Fatalln(err)
					}
					if pwsRegexMatch {
						updateUserPws(currUser, currPws, newPws)
						dialog.Alert("Password Updated! Please login again!") //Updated password
						updated = true
						userLogOut(res, req)
						break
					}
				}
			}
		}
		if !updated {
			dialog.Alert("Please re-enter your Username and Password! No special characters allowed")
		}
	}
}

// User Logout Delete Cookies and Json ---✨
func userLogOut(res http.ResponseWriter, req *http.Request) {
	cookies := getCookiesJSON()

	for _, v := range cookies {
		if v.Name == currentUser {
			deleteCookieFromJSON(res, req, v)
			deleteCookie(res, req, currentUser)
			fmt.Println("Cookie deleted")
			redirectLoginPage(res, req)
		}
	}
}

// --- ADMIN ---
// Main page to add appointments ---✨
func adminPageAddAppt(res http.ResponseWriter, req *http.Request) {

	checkLoggedIn(res, req, currentUser)
	redirectUser(res, req, currentUser)

	apptSettings := allAdminApptAndSettings()

	if req.Method == "GET" {
		err := adminTpl.ExecuteTemplate(res, "MainAdminPage.gohtml", &apptSettings)

		if err != nil {
			log.Fatalln(err)
			return
		}

	} else if req.Method == "POST" {
		setAppointmentType := req.FormValue("setAppointmentType")
		setDateSelect := req.FormValue("setDateSelect")
		setTimeSelect := req.FormValue("setTimeSelect")

		addAppt := AppointmentsDetail{AppointmentType: setAppointmentType, Date: setDateSelect, Time: setTimeSelect}

		duplicateAppt := addAppointmentsAdmin(addAppt)

		if duplicateAppt {
			dialog.Alert("Please choose another appointment date to add!")
			redirectAdminPage(res, req)
			return
		} else if !duplicateAppt {
			dialog.Alert("Appointment added!")
			redirectAdminPage(res, req)
		}
	}
}

// User remove appt ---✨
func adminRemoveAppt(res http.ResponseWriter, req *http.Request) {

	checkLoggedIn(res, req, currentUser)
	redirectUser(res, req, currentUser)

	apptSettings := allAdminApptAndSettings()

	if req.Method == "GET" {
		err := adminRemoveApptTpl.ExecuteTemplate(res, "MainAdminRemoveApptPage.gohtml", &apptSettings)

		if err != nil {
			log.Fatalln(err)
			return
		}

	} else if req.Method == "POST" {
		removeAppointmentType := req.FormValue("removeAppointmentType")
		removeDateSelect := req.FormValue("removeDateSelect")
		removeTimeSelect := req.FormValue("removeTimeSelect")

		removeAppt := AppointmentsDetail{AppointmentType: removeAppointmentType, Date: removeDateSelect, Time: removeTimeSelect}
		allAppts := getAppointments()

		for _, v := range allAppts {
			if v.AppointmentType == removeAppointmentType && v.Date == removeDateSelect && v.Time == removeTimeSelect {
				deleteAppointmentAdmin(removeAppt)
				redirectAdminRemoveApptPage(res, req)
				return
			}
		}

		dialog.Alert("Please select an available appointment!")
		redirectAdminRemoveApptPage(res, req)
	}
}

// User remove session ---✨
func adminSessionPage(res http.ResponseWriter, req *http.Request) {

	checkLoggedIn(res, req, currentUser)
	redirectUser(res, req, currentUser)

	cookies := getCookiesStruct()

	if req.Method == "GET" {
		err := adminSessionTpl.ExecuteTemplate(res, "MainSessionPage.gohtml", &cookies)

		if err != nil {
			log.Fatalln(err)
			return
		}

	} else if req.Method == "POST" {

		deleteSession := req.FormValue("deleteSession")

		cookiesData := getCookiesJSON()

		for _, v := range cookiesData {
			if deleteSession == v.Name {
				deleteCookieFromJSON(res, req, v)
				dialog.Alert("Cookie deleted!")
				redirectAdminSessionPage(res, req)
				break
			}
		}
	}
}

// Admin remove user ---✨
func userAdminPage(res http.ResponseWriter, req *http.Request) {

	checkLoggedIn(res, req, currentUser)
	redirectUser(res, req, currentUser)

	users := getAllUsers()

	type usersStruct struct {
		Users   []string
		UserNos []string
	}

	var num int
	var userList []string
	var userNo []string

	for _, v := range users {
		userList = append(userList, v.LoginUsername)
		num++
	}

	for i := 1; i < num+1; i++ {
		userNoString := "User " + strconv.Itoa(i)
		userNo = append(userNo, userNoString)
	}

	passUsers := usersStruct{
		Users:   userList,
		UserNos: userNo,
	}

	if req.Method == "GET" {
		err := adminUserTpl.ExecuteTemplate(res, "MainUserAdminPage.gohtml", &passUsers)

		if err != nil {
			log.Fatalln(err)
			return
		}

	} else if req.Method == "POST" {
		userList := getAllUsers()
		removeUser := req.FormValue("removeUser")

		for _, v := range userList {
			if v.LoginUsername == removeUser {
				deleteUserAdmin(removeUser)
				dialog.Alert("User removed!")
				redirectUsersAdminPage(res, req)
				return
			}
		}
	}
}
