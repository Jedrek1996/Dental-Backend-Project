package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"strconv"
	"strings"
)

var appointmentsData []AppointmentsDetail

var bookedAppointmentsData []AppointmentsDetail

type allAppointmentsMap map[string]interface{}

// Struct for the hard coded values for Admin to edit
// (Admin will have to access the JSON file to edit the values)
type AppointmentSetting struct {
	AppointmentType []string
	Date            []string
	Time            []string
}

// Appointment Detail used to pass into add/remove functions
type AppointmentsDetail struct {
	AppointmentType string `json:"AppointmentType"`
	Date            string `json:"Date"`
	Time            string `json:"Time"`
	User            string
}

// Used to display all appt (Admin)
type AllAppointmentsDataStruct struct {
	AppointmentType []string
	Date            []string
	Time            []string
	ApptVal         []string
}

// Used to display users appt
type AppointmentDataStruct struct {
	AppointmentsDate []string
	AppointmentsTime []string
	UserAppType      []string
	UserDate         []string
	UserTime         []string
	ApptVal          []string
	AvailableType    []string
	AvailableDate    []string
	AvailableTime    []string
}

// To hold if user is logged in
type UserLoginStruct struct {
	LoggedIn string
}

// Returns ApptDataStrcut (User) for Gohtml template ---✨
func appDataStruct() AppointmentDataStruct {

	apptDate, apptTime := getAppointmentDateTime()

	//For available appts
	var availableApptsType []string
	var availableApptsDate []string
	var availableApptsTime []string

	userAppType, userDate, userTime, userUser := getUserAppointments(currentUser)

	//To display Appointment No.
	var apptNo []string
	for i := 1; i < len(userUser)+1; i++ {
		apptValue := "Appointment " + strconv.Itoa(i)
		apptNo = append(apptNo, apptValue)
	}
	availableAppts := getAvailableAppointments()

	for _, v := range availableAppts {
		availableApptsType = append(availableApptsType, v.AppointmentType)
		availableApptsDate = append(availableApptsDate, v.Date)
		availableApptsTime = append(availableApptsTime, v.Time)
	}

	appointmentsData := AppointmentDataStruct{
		AppointmentsDate: apptDate,
		AppointmentsTime: apptTime,
		UserAppType:      userAppType,
		UserDate:         userDate,
		UserTime:         userTime,
		ApptVal:          apptNo,
		AvailableType:    availableApptsType,
		AvailableDate:    availableApptsDate,
		AvailableTime:    availableApptsTime,
	}

	return appointmentsData
}
 
// Returns AppDataStrcut (Admin) to pass into allAdminApptAndSettings Struct for Gohtml template ---✨
func allAppDataStruct() AllAppointmentsDataStruct {

	allAppt := getAppointments()

	var apptType []string
	var apptDate []string
	var apptTime []string
	var apptNo []string

	for _, v := range allAppt {
		apptType = append(apptType, v.AppointmentType)
		apptDate = append(apptDate, v.Date)
		apptTime = append(apptTime, v.Time)
	}

	for i := 1; i < len(apptType)+1; i++ {
		apptValue := "Appointment " + strconv.Itoa(i)
		apptNo = append(apptNo, apptValue)
	}

	allApptsData := AllAppointmentsDataStruct{
		AppointmentType: apptType,
		Date:            apptDate,
		Time:            apptTime,
		ApptVal:         apptNo,
	}

	return allApptsData
}

// Returns AppointmentSettings (Admin) to pass into allAdminApptAndSettings Struct for Gohtml template ---✨
func getAppointmentSettings() AppointmentSetting {

	var apptType []string
	var apptDate []string
	var apptTime []string

	var apptTypeReturn []string
	var apptDateReturn []string
	var apptTimeReturn []string

	appts := allAppointmentsMap{}

	file, err := ioutil.ReadFile("appointmentsSetting.json")
	if err != nil {
		log.Fatal(err)
	}

	json.Unmarshal(file, &appts)

	//Get the values from JSON and convert from interface

	//Appt Types
	apptTypes := appts["AppointmentType"]
	apptTypesStr := fmt.Sprintf("%v", apptTypes)
	apptType = append(apptType, apptTypesStr)

	for _, v := range apptType {
		v = v[:len(v)-1]
		v = v[1:]
		s := strings.Split(v, " ")
		for i := 0; i < len(s); i++ {
			apptTypeReturn = append(apptTypeReturn, s[i])
		}
	}

	//Appt Dates
	apptDates := appts["Date"]
	apptDatesStr := fmt.Sprintf("%v", apptDates)
	apptDate = append(apptDate, apptDatesStr)

	for _, v := range apptDate {
		v = v[:len(v)-1]
		v = v[1:]
		s := strings.Split(v, " ")
		for i := 0; i < len(s); i++ {
			apptDateReturn = append(apptDateReturn, s[i])
		}
	}

	//Appt Time
	apptTimes := appts["Time"]
	apptTimesStr := fmt.Sprintf("%v", apptTimes)
	apptTime = append(apptTime, apptTimesStr)

	for _, v := range apptTime {
		v = v[:len(v)-1]
		v = v[1:]
		s := strings.Split(v, " ")
		for i := 0; i < len(s); i++ {
			apptTimeReturn = append(apptTimeReturn, s[i])
		}
	}

	apptSettings := AppointmentSetting{
		AppointmentType: apptTypeReturn,
		Date:            apptDateReturn,
		Time:            apptTimeReturn,
	}

	return apptSettings
}

// According to the documentation ioutil.Readfile has a auto close feature.
// Get Appointment from Appointments JSON ---✨
func getAppointments() []AppointmentsDetail {

	file, err := ioutil.ReadFile("appointments.json")
	if err != nil {
		log.Fatal(err)
	}

	json.Unmarshal(file, &appointmentsData)

	return appointmentsData
}

// Get Booked Appointments from BookedAppointments JSON---✨
func getBookedAppointmentsData() []AppointmentsDetail {

	file, err := ioutil.ReadFile("bookedAppointments.json")
	if err != nil {
		log.Fatal(err)
	}

	json.Unmarshal(file, &bookedAppointmentsData)

	return bookedAppointmentsData
}

// Compares the BookedAppointments and Appointments JSON (Returns available Appointments) ---✨
// Result: Available Appointment for Users and Admin structs
func getAvailableAppointments() []AppointmentsDetail {

	allAppt := getAppointments()
	bookedAppt := getBookedAppointmentsData()

	var diff []AppointmentsDetail

	//Loop twice to swap
	for i := 0; i < 2; i++ {
		for _, s1 := range bookedAppt {
			found := false
			for _, s2 := range allAppt {
				if s1.Date == s2.Date && s1.Time == s2.Time && s1.AppointmentType == s2.AppointmentType {
					found = true
					break
				}
			}
			if !found {
				diff = append(diff, s1)
			}
		}
		// Swap the slices, only if it was the first loop
		if i == 0 {
			bookedAppt, allAppt = allAppt, bookedAppt
		}
	}
	fmt.Println(diff)
	return diff
}

// Appointment date time for users dropdown select ---✨
func getAppointmentDateTime() ([]string, []string) {

	appointmentData := getAppointments()

	var apptDates []string
	var apptTimes []string

	for _, v := range appointmentData {
		apptDates = append(apptDates, v.Date)
		apptTimes = append(apptTimes, v.Time)
	}

	apptDates = removeDuplicateStr(apptDates)
	apptTimes = removeDuplicateStr(apptTimes)

	return apptDates, apptTimes
}

// Removes any duplicates in slice ---✨
func removeDuplicateStr(strSlice []string) []string {
	allKeys := make(map[string]bool)

	list := []string{}

	for _, item := range strSlice {
		if _, value := allKeys[item]; !value {
			allKeys[item] = true //If not value
			list = append(list, item)
		}
	}
	return list //Returns non duplicate slice
}

// Comapres new appt with existing appt ---✨
func addUserAppointment(details AppointmentsDetail) bool {

	defer wg.Done()

	fmt.Println("Added:", details)

	bookedAppointmentsData := getBookedAppointmentsData()

	for _, v := range bookedAppointmentsData {
		if v.AppointmentType == details.AppointmentType && v.Date == details.Date && v.Time == details.Time {
			wg.Done()
			return true
		}
	}
	bookedAppointmentsData = append(bookedAppointmentsData, details)

	bytes, err := json.MarshalIndent(bookedAppointmentsData, "", "  ")

	if err != nil {
		log.Fatalln(err)
	}
	_ = ioutil.WriteFile("bookedAppointments.json", bytes, 0644)

	return false

}

// Get users appointment info ---✨
func getUserAppointments(username string) ([]string, []string, []string, []string) {

	var userAppts []AppointmentsDetail

	var userApptType []string
	var userDate []string
	var userTime []string
	var userUser []string

	bookedAppointmentsData := getBookedAppointmentsData()

	for _, v := range bookedAppointmentsData {
		if v.User == username {
			userAppts = append(userAppts, v)
		}
	}

	for _, v := range userAppts {
		userApptType = append(userApptType, v.AppointmentType)
		userDate = append(userDate, v.Date)
		userTime = append(userTime, v.Time)
		userUser = append(userUser, v.User)
	}

	return userApptType, userDate, userTime, userUser
}

// To check whether the currentAppt is available in user's bookedAppointment ---✨
func checkCurrentAppointment(appointmentType string, date string, time string, user string, channel chan bool) {
	bookedAppointmentsData := getBookedAppointmentsData()

	for _, v := range bookedAppointmentsData {
		if user == v.User {
			if appointmentType == v.AppointmentType && date == v.Date && time == v.Time {
				fmt.Println("Current appointment available/booked")
				channel <- true
				break
			}
		}
	}
	channel <- false

}

// To check whether the appt is not booked or not in user's bookedAppointment ---✨
func checkForAvaliableAppt(appointmentType string, date string, time string, channel chan bool) {

	bookedAppointmentsData := getBookedAppointmentsData()

	for _, v := range bookedAppointmentsData {
		if appointmentType == v.AppointmentType && date == v.Date && time == v.Time {
			fmt.Println("Appointment UNavailable ")
			channel <- false
			break
		}
	}
	channel <- true

}

// Delete appointment (User) ---✨
func deleteAppointment(appt AppointmentsDetail) {
	bookedAppointmentsData := getBookedAppointmentsData()
	for _, v := range bookedAppointmentsData {
		if v.AppointmentType == appt.AppointmentType && v.Date == appt.Date && v.Time == appt.Time {
			removeAppt(appt)
			break
		}
	}
}

// Remove appointment from booked appointment JSON (User) ---✨
func removeAppt(remove AppointmentsDetail) {
	bookedAppointmentsData := getBookedAppointmentsData()

	for i := 0; i < len(bookedAppointmentsData); i++ {
		bookedAppt := bookedAppointmentsData[i]
		if bookedAppt == remove {
			bookedAppointmentsData = append(bookedAppointmentsData[:i], bookedAppointmentsData[i+1:]...)
			i--
			fmt.Println("Removed user appt:", remove)
			bytes, err := json.MarshalIndent(bookedAppointmentsData, "", "  ")
			if err != nil {
				log.Fatalln(err)
			}
			_ = ioutil.WriteFile("bookedAppointments.json", bytes, 0644)
			break
		}
	}

}