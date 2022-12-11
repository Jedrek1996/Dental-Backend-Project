package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"

	"github.com/tawesoft/golib/v2/dialog"
	// "github.com/tawesoft/golib/v2/dialog"
)

// Settings passed into Admin Gohtml files
type adminApptSettingsStruct struct {
	AppointmentType []string
	Date            []string
	Time            []string
	ApptVal         []string
	SettingApptType []string
	SettingDate     []string
	SettingTime     []string
	AvailableType   []string
	AvailableDate   []string
	AvailableTime   []string
	BookedType      []string
	BookedDate      []string
	BookedTime      []string
	BookedUser      []string
}

func allAdminApptAndSettings() adminApptSettingsStruct {

	allAppts := allAppDataStruct()
	apptSettings := getAppointmentSettings()
	availableAppts := getAvailableAppointments()
	bookedAppts := getBookedAppointmentsData()

	//For available appts
	var availableApptsType []string
	var availableApptsDate []string
	var availableApptsTime []string

	for _, v := range availableAppts {
		availableApptsType = append(availableApptsType, v.AppointmentType)
		availableApptsDate = append(availableApptsDate, v.Date)
		availableApptsTime = append(availableApptsTime, v.Time)
	}

	//For booked appts
	var bookedApptsType []string
	var bookedApptsDate []string
	var bookedApptsTime []string
	var bookedApptUser []string

	for _, v := range bookedAppts {
		bookedApptsType = append(bookedApptsType, v.AppointmentType)
		bookedApptsDate = append(bookedApptsDate, v.Date)
		bookedApptsTime = append(bookedApptsTime, v.Time)
		bookedApptUser = append(bookedApptUser, v.User)
	}

	passThis := adminApptSettingsStruct{
		AppointmentType: allAppts.AppointmentType,
		Date:            allAppts.Date,
		Time:            allAppts.Time,
		ApptVal:         allAppts.ApptVal,
		SettingApptType: apptSettings.AppointmentType,
		SettingDate:     apptSettings.Date,
		SettingTime:     apptSettings.Time,
		AvailableType:   availableApptsType,
		AvailableDate:   availableApptsDate,
		AvailableTime:   availableApptsTime,
		BookedType:      bookedApptsType,
		BookedDate:      bookedApptsDate,
		BookedTime:      bookedApptsTime,
		BookedUser:      bookedApptUser,
	}
	return passThis
}

// Admin Add Main Appointment---✨
func addAppointmentsAdmin(add AppointmentsDetail) bool {
	// defer wg.Done()
	appointmentList := getAppointments()

	for _, v := range appointmentList {
		if v.AppointmentType == add.AppointmentType && v.Date == add.Date && v.Time == add.Time {
			// wg.Done()
			return true
		}
	}
	appointmentList = append(appointmentList, add)

	bytes, err := json.MarshalIndent(appointmentList, "", "  ")

	if err != nil {
		log.Fatalln(err)
	}
	_ = ioutil.WriteFile("appointments.json", bytes, 0644)
	return false
}

// Admin Delete Main Appointment---✨
func deleteAppointmentAdmin(remove AppointmentsDetail) {
	appointmentList := getAppointments()
	for _, v := range appointmentList {
		if v.AppointmentType == remove.AppointmentType && v.Date == remove.Date && v.Time == remove.Time {
			removeApptAdmin(remove, appointmentList)
			break
		}
	}

}

// Admin Remove Main Appointment---✨
func removeApptAdmin(remove AppointmentsDetail, appointmentList []AppointmentsDetail) {
	isBooked := checkIfApptBookedAdmin(remove)
	if !isBooked {
		for i := 0; i < len(appointmentList); i++ {
			appt := appointmentList[i]
			if appt == remove {
				appointmentList = append(appointmentList[:i], appointmentList[i+1:]...)
				i--
				dialog.Alert("Appointment Removed")
				fmt.Println("Removed admin appt:", remove)
				break
			}
		}
	} else if isBooked {
		dialog.Alert("This appointment has already been booked!")
	}

	bytes, err := json.MarshalIndent(appointmentList, "", "  ")
	if err != nil {
		log.Fatalln(err)
	}
	_ = ioutil.WriteFile("appointments.json", bytes, 0644)
}

// Admin Check Appts (If booked by user return false)
func checkIfApptBookedAdmin(remove AppointmentsDetail) bool {
	bookedAppts := getBookedAppointmentsData()

	for _, v := range bookedAppts {
		if v.AppointmentType == remove.AppointmentType && v.Date == remove.Date && v.Time == remove.Time {
			return true
		}
	}
	return false
}

// Admin Delete User ---✨
func deleteUserAdmin(username string) {
	allUser := getAllUsers()
	for _, v := range allUser {
		if v.LoginUsername == username {
			removeUserAdmin(username, allUser)
			break
		}
	}
}

// Admin Remove User From JSON---✨
func removeUserAdmin(username string, allUser []loginDetails) {
	for i := 0; i < len(allUser); i++ {
		appt := allUser[i]
		if appt.LoginUsername == username {
			allUser = append(allUser[:i], allUser[i+1:]...)
			i--
			fmt.Println("Removed:", username)
			break
		}
	}

	bytes, err := json.MarshalIndent(allUser, "", "  ")
	if err != nil {
		log.Fatalln(err)
	}
	_ = ioutil.WriteFile("usersData.json", bytes, 0644)
}

// Admin Redirect User---✨
func redirectUser(res http.ResponseWriter, req *http.Request, userName string) {
	if !strings.EqualFold(userName, "Admin") {
		redirectUserPage(res, req)
		dialog.Alert("Please do not try and access the admin pages!")
	}
}
