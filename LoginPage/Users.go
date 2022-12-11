package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"

	"github.com/tawesoft/golib/v2/dialog"
)

type loginDetails struct {
	LoginUsername string `json:"Username"`
	LoginPassword string `json:"Password"`
	Rights        string `json:"Rights"`
}

var userData []loginDetails

// Get all users from usersData Json ---✨
func getAllUsers() []loginDetails {

	file, err := ioutil.ReadFile("usersData.json")
	if err != nil {
		log.Fatal(err)
	}

	json.Unmarshal(file, &userData)
	fmt.Println(userData)
	return userData
}

// Update Users into usersData Json ---✨
func updateUser(oldUserName string, matchPws bool, newUserName string) bool {

	userList := getAllUsers()

	for i, v := range userList {
		if strings.EqualFold(newUserName, v.LoginUsername) {
			return true
		} else if !strings.EqualFold(newUserName, v.LoginUsername) {
			if oldUserName == v.LoginUsername && matchPws {
				updateApptUser(oldUserName, newUserName)

				userList[i].LoginUsername = newUserName

				bytes, err := json.MarshalIndent(userList, "", "  ")

				if err != nil {
					log.Fatalln(err)
				}
				_ = ioutil.WriteFile("usersData.json", bytes, 0644)

				return false
			}
		}
	}
	return true
}

// Update username appt into bookedAppointments json ---✨
func updateApptUser(oldUserName string, newUserName string) {

	bookedAppointmentsData := getBookedAppointmentsData()

	for i, v := range bookedAppointmentsData {
		if oldUserName == v.User {
			bookedAppointmentsData[i].User = newUserName //Update username
			dialog.Alert("Your username has been updated!")
		}
	}

	bytes, err := json.MarshalIndent(bookedAppointmentsData, "", "  ")

	if err != nil {
		log.Fatalln(err)
	}

	_ = ioutil.WriteFile("bookedAppointments.json", bytes, 0644)
}

// Update password into usersData Json ---✨
func updateUserPws(oldUserName string, oldPws string, newPws string) {
	userList := getAllUsers()

	//Check if the new password is the same as the old password
	hashedPassword, _ := HashString(newPws)
	match := CheckStringHash(oldPws, hashedPassword)

	for i, v := range userList {
		if oldUserName == v.LoginUsername && !match {
			hashedPassword, _ := HashString(newPws)
			userList[i].LoginPassword = hashedPassword
			dialog.Alert("Password updated!") //Updated password
			break
		} else {
			dialog.Alert("New Password should not be the same!")
			break
		}
	}

	bytes, err := json.MarshalIndent(userList, "", "  ")

	if err != nil {
		log.Fatalln(err)
	}
	_ = ioutil.WriteFile("usersData.json", bytes, 0644)
}

// Adding new users to usersData.json ---✨
func addUsersToJson(newUser loginDetails) bool {
	file, err := ioutil.ReadFile("usersData.json")
	if err != nil {
		log.Fatal(err)
	}

	json.Unmarshal(file, &userData)

	for _, v := range userData {
		fmt.Println(v.LoginUsername, newUser.LoginUsername)
		if strings.EqualFold(v.LoginUsername, newUser.LoginUsername) {

			return false

		} else if !strings.EqualFold(v.LoginUsername, newUser.LoginUsername) {

			userData = append(userData, newUser)

			bytes, err := json.MarshalIndent(userData, "", "  ")

			if err != nil {
				log.Fatalln(err)
			}
			_ = ioutil.WriteFile("usersData.json", bytes, 0644)

			return true
		}
	}
	return false
}

// Check if user/admin redirects accordingly ---✨
func checkIsUserRedirect(res http.ResponseWriter, req *http.Request, userName string, adminPageRedirect redirectFn, userPageRedirect redirectFn) {

	if strings.EqualFold(userName, "Admin") {
		adminPageRedirect(res, req)
	} else {
		userPageRedirect(res, req)
	}
}
