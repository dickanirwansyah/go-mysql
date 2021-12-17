package main

import (
	"database/sql"
	"log"
	"net/http"
	"text/template"

	_ "github.com/go-sql-driver/mysql"
)

type VendorCategory struct {
	VendorCategoryId int
	Name             string
	Activated        string
}

/** initial template **/
var tmpl = template.Must(template.ParseGlob("form/*"))

/** function koneksi ke database **/
func dbConnection() (db *sql.DB) {
	dbDriver := "mysql"
	dbUser := "root"
	dbPass := "rootroot"
	dbName := "db_app_backend"
	db, err := sql.Open(dbDriver, dbUser+":"+dbPass+"@/"+dbName)
	if err != nil {
		panic(err.Error())
	}
	return db
}

/** function show data list vendor **/
func Index(w http.ResponseWriter, r *http.Request) {
	db := dbConnection()
	selDB, err := db.Query("SELECT * FROM VENDOR_CATEGORY ORDER BY VENDOR_CATEGORY_ID DESC")
	if err != nil {
		panic(err.Error())
	}
	vendCategory := VendorCategory{}
	res := []VendorCategory{}
	for selDB.Next() {
		var vendor_category_id int
		var name, activated string
		err := selDB.Scan(&vendor_category_id, &name, &activated)
		if err != nil {
			panic(err.Error())
		}
		vendCategory.VendorCategoryId = vendor_category_id
		vendCategory.Name = name
		vendCategory.Activated = activated
		res = append(res, vendCategory)
	}
	log.Println("Data : ", res)
	tmpl.ExecuteTemplate(w, "Index", res)
	defer db.Close()
}

/** template insert **/
func NewTmpl(w http.ResponseWriter, r *http.Request) {
	log.Println("template insert..")
	tmpl.ExecuteTemplate(w, "New", nil)
}

/** function action insert **/
func Insert(w http.ResponseWriter, r *http.Request) {
	log.Print("process save vendor..")
	db := dbConnection()
	if r.Method == "POST" {
		name := r.FormValue("name")
		activeted := r.FormValue("activated")
		insertForm, err := db.Prepare("INSERT INTO VENDOR_CATEGORY (name, activated) VALUES (?,?)")
		if err != nil {
			panic(err.Error())
		}
		insertForm.Exec(name, activeted)
		log.Println("process insert vendor category : " + "name : " + name + " & activated : " + activeted)
	}
	defer db.Close()
	http.Redirect(w, r, "/", 301)
}

/** func main golang **/
func main() {
	log.Print("Server started on : http://localhost:9999")
	http.HandleFunc("/", Index)
	http.HandleFunc("/new", NewTmpl)
	http.HandleFunc("/insert", Insert)
	http.ListenAndServe(":9999", nil)
}
