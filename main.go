package main

import (
	"database/sql"
	"fmt"
	"io/ioutil"
	"log"
	"mux"
	_ "mysql"
	"net/http"
)

var uname, email, password string

func dbConn() (db *sql.DB) {
	dbDriver := "mysql"
	dbUser := "root"
	dbPass := "root"
	dbName := "Neta"
	db, err := sql.Open(dbDriver, dbUser+":"+dbPass+"@/"+dbName)
	if err != nil {
		panic(err.Error())
	}
	//fmt.Println("Connected to db")
	return db

}

func loginpagehandler(w http.ResponseWriter, r *http.Request) {
	//fmt.Println("handler executed")
	content, err := ioutil.ReadFile("login.html")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Fprint(w, string(content))
}
func loginhandler(w http.ResponseWriter, r *http.Request) {
  db := dbConn()
  var tempuname, temppass string
	uname := r.FormValue("uname")
	password := r.FormValue("password")
	if len(uname) != 0 && len(password) != 0 {
		check , err := db.Query("SELECT uname, password from user where uname=?",uname)
		if err != nil {
			panic(err.Error())
		}
		for check.Next() {
			//fmt.Println("Enter username and password")
			err = check.Scan(&tempuname,&temppass)
			if err != nil {
				panic(err.Error())
			}
		}
		if tempuname == uname && temppass == password {
			fmt.Println("The user is authenticated")
			http.Redirect(w, r, "/home", http.StatusTemporaryRedirect)
			//fmt.Fprint(w, "logged in")
		}

	}
	if uname != tempuname && len(uname) != 0 && len(password) != 0 {

		http.Redirect(w, r, "/register", http.StatusMovedPermanently)
	} else {
		http.Redirect(w, r, "/", http.StatusMovedPermanently)
	}
	defer db.Close()

}
func registerpagehandler(w http.ResponseWriter, r *http.Request) {
	content, err := ioutil.ReadFile("register.html")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Fprint(w, string(content))

}

func registerhandler(w http.ResponseWriter, r *http.Request) {

	uname := r.FormValue("uname")
	email := r.FormValue("email")
	password := r.FormValue("password")
	if len(uname) != 0 && len(email) != 0 && len(password) != 0 {
		//fmt.Println("registered")
		db := dbConn()
		if r.Method == "POST" {
			insForm, err := db.Prepare("INSERT INTO user(uname, email, password) VALUES(?,?,?)")
			if err != nil {
				panic(err.Error())
			}
			insForm.Exec(uname, email, password)
			fmt.Println("INSERT: Name: " + uname + " | E-Mail: " + email)
			http.Redirect(w, r, "/", http.StatusMovedPermanently)
		}
		defer db.Close()

	} else {
		http.Redirect(w, r, "/register", http.StatusMovedPermanently)
	}
}

/*func home(w http.ResponseWriter, r *http.Request){
	content, err := ioutil.ReadFile("home.html")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Fprint(w, string(content))

}*/

func read(w http.ResponseWriter, r *http.Request) {
//fmt.Println("Displaying records")
	var tempuname, tempemail, temppassword string
	db := dbConn()
	uname := r.FormValue("uname")
	email := r.FormValue("email")
	password := r.FormValue("password")

	if len(uname) != 0 && len(password) != 0 {
		check, err := db.Query("SELECT * from user order by uname")
		if err != nil {
			panic(err.Error())
		}
		res := []string{}
		for check.Next() {

			err = check.Scan(&tempuname, &tempemail, &temppassword)
			if err != nil {
				panic(err.Error())
			}
			uname = tempuname
			email = tempemail
			password = temppassword
			res = append(res, uname, email, password)
			fmt.Fprintln(w, tempuname, tempemail, temppassword)
		}

		defer db.Close()

	}

}
func updatepagehandler(w http.ResponseWriter, r *http.Request){
	content, err := ioutil.ReadFile("update.html")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Fprint(w, string(content))

}

func update(w http.ResponseWriter, r *http.Request ) {
	db := dbConn()
	if r.Method == "POST" {
		uname := r.FormValue("uname")
		email := r.FormValue("email")
		password := r.FormValue("password")
		insForm, err := db.Prepare("UPDATE user SET uname=?, password=? WHERE email=?")
		if err != nil {
			panic(err.Error())
		}
		insForm.Exec(uname , password,email)
		log.Println("UPDATE: Name: " + uname + " | Password: " + password)
		fmt.Println("Updated successfully")
	}
	defer db.Close()
	http.Redirect(w, r, "/home", 307)
}

func deletepagehandler(w http.ResponseWriter, r *http.Request){
	content, err := ioutil.ReadFile("delete.html")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Fprint(w, string(content))

}
func delete(w http.ResponseWriter, r *http.Request ){
	db := dbConn()
	email := r.FormValue("email")
	delForm, err := db.Prepare("DELETE FROM user WHERE email=?")
	if err != nil {
		panic(err.Error())
	}
	delForm.Exec(email)
	//log.Println("DELETE")
	defer db.Close()
	//http.Redirect(w, r, "/home", 301)
	fmt.Fprint(w,"Deleted successfully")
}

func main() {
	//fmt.Println("Main executed")
	router := mux.NewRouter()
	//fmt.Println("Router executed")
	router.HandleFunc("/", loginpagehandler).Methods("GET")
	router.HandleFunc("/login", loginhandler).Methods("POST")
	router.HandleFunc("/register", registerpagehandler).Methods("GET")
	router.HandleFunc("/register", registerhandler).Methods("POST")
	//router.HandleFunc("/home",home).Methods("GET")
	router.HandleFunc("/home",read).Methods("POST")
	router.HandleFunc("/update",updatepagehandler).Methods("GET")
	router.HandleFunc("/update",update).Methods("POST")
	router.HandleFunc("/delete",deletepagehandler).Methods("GET")
	router.HandleFunc("/delete",delete).Methods("POST")
	http.ListenAndServe(":9090", router)
}

