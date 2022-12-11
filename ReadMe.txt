🌭🍿--- Read Me File --- 🍕🍔
------------------------------------------------------------------------------------------------
This project is a part of the GoTrack assignment assigned to the participants.

A brief summary about this project, it is a backend project that uses Golang as the primary language. Accompanied with CSS and GoHTML for the front end formatting.
The project has several external Go libraries that has to be imported before running. (Do a go mod tidy followed by go get etc..)
This project is also using JSON files as the primary or "mock" database; as database were not allowed to be used. JSON and Data structures were accepted. Plese do not that there are popups with Dialog Alert 😡, please click on them to proceed on~
Do not that this web application is not fully customized yet and will be going through more developments in the future. Any feedback will be greatly appreciated - Jed 😊 (https://www.linkedin.com/in/jedrek-koh/)

The command that can be used is : go run -v ./Loginpage (Connect to localhost 8080)
The http url: http://localhost:8080/   (Login Page)
------------------------------------------------------------------------------------------------
Credentials for admin has been created as tasked in the assignment.
Admin Account 👮‍♂️
Username: Admin
Password: Password

Normal User Account 🤦‍♂️
Username: Test
Password: Test123
------------------------------------------------------------------------------------------------
-Files- 📖
1.GO files- The Main.go is the main file which excutes the http connection, it is the back bone of the entire web application. 
The other files are functions that serves different purposes, comments have been labelled above each function with a ---✨ at the end of the comment.

2. CSS - Contains the Style sheet needed for the GoTemplates (CSS file connected through the fileserver in the main func for GO and link href for the GoTemplates)

3. LoginPageTemplates - In this folder lies 4 other sub folders which holds the different templates based on the (folder).

4. JSON files - The primary data storage used in this web application. There are 5 JSON files. 
- Appointments : Stores all the appointments
- BookedAppoinments : Stores booked appoinments from users
- AppointmentsSetting : Store the appointment values that is used to render the settings. (HARD CODED) ⚠️ Will update this so that the admin can input their own values in the web application
- Cookies: Session cookie data is being stored when the user is logged in. Remove when logged out.
- UsersData: The credentials of users are stored here (Username and Hashed Password)

5. GoMod/GoSum - Directory and specify versions of dependencies for each module/ The checksum of direct and indirect dependency required along with the version.
(Please do a go mod tidy if necessary prior to running this web application.)
------------------------------------------------------------------------------------------------