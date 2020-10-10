package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
)

var db *sql.DB
var err error

// Customer struct (Model) ...
type Customer struct {
	CustomerID   string `json:"CustomerID"`
	CompanyName  string `json:"CompanyName"`
	ContactName  string `json:"ContactName"`
	ContactTitle string `json:"ContactTitle"`
	Address      string `json:"Address"`
	City         string `json:"City"`
	Country      string `json:"Country"`
	Phone        string `json:"Phone"`
	PostalCode   string `json:"PostalCode"`
}
type account struct {
	Accountid string `json:"account_id"`
	Name      string `json:"name"`
	Accnum    string `json:"acc_num"`
	Debit     string `json:"debit"`
	Credit    string `json:"credit"`
	Balance   string `json:"balance"`
	Parentid  string `json:"parent_id"`
}
type adjtable struct {
	Empid  string `json:"emp_id"`
	Name   string `json:"name"`
	Salary string `json:"salary"`
	Bossid string `json:"boss_id"`
}

type book struct {
	Bookid        string `json:"BookID"`
	Title         string `json:"Title"`
	Author        string `json:"Author"`
	DatePublished string `json:"DatePublished"`
	Publisher     string `json:"Publisher"`
	Edition       string `json:"Edition"`
	City          string `json:"City"`
}

type categorie struct {
	Categoryid   string `json:"CategoryID"`
	Categoryname string `json:"CategoryName"`
	Description  string `json:"Description"`
	Picture      string `json:"Picture"`
}
type city struct {
	Cityid     string `json:"CityID"`
	Provinceid string `json:"ProvinceID"`
	Cityname   string `json:"CityNAme"`
}
type countrie struct {
	Countryid   string `json:"CountryID"`
	Countryname string `json:"CountryName"`
}

// Get all orders

func getCustomers(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var customers []Customer

	sql := `SELECT
				CustomerID,
				IFNULL(CompanyName,''),
				IFNULL(ContactName,'') ContactName,
				IFNULL(ContactTitle,'') ContactTitle,
				IFNULL(Address,'') Address,
				IFNULL(City,'') City,
				IFNULL(Country,'') Country,
				IFNULL(Phone,'') Phone ,
				IFNULL(PostalCode,'') PostalCode
			FROM customers`

	result, err := db.Query(sql)

	defer result.Close()

	if err != nil {
		panic(err.Error())
	}

	for result.Next() {

		var customer Customer
		err := result.Scan(&customer.CustomerID, &customer.CompanyName, &customer.ContactName,
			&customer.ContactTitle, &customer.Address, &customer.City, &customer.Country,
			&customer.Phone, &customer.PostalCode)

		if err != nil {
			panic(err.Error())
		}
		customers = append(customers, customer)
	}

	json.NewEncoder(w).Encode(customers)
}

func createCustomer(w http.ResponseWriter, r *http.Request) {

	if r.Method == "POST" {

		CustomerID := r.FormValue("CustomerID")
		CompanyName := r.FormValue("CompanyName")
		ContactName := r.FormValue("ContactName")
		ContactTitle := r.FormValue("ContactTitle")
		Address := r.FormValue("Address")
		City := r.FormValue("City")
		Region := r.FormValue("Region")
		PostalCode := r.FormValue("PostalCode")
		Country := r.FormValue("Country")
		Phone := r.FormValue("Phone")
		Fax := r.FormValue("Fax")

		stmt, err := db.Prepare("INSERT INTO customers (CustomerID,CompanyName,ContactName,ContactTitle,Address,City,Region,PostalCode,Country,Phone,Fax) VALUES (?,?,?,?,?,?,?,?,?,?,?)")

		_, err = stmt.Exec(CustomerID, CompanyName, ContactName, ContactTitle, Address, City, Region, PostalCode, Country, Phone, Fax)

		if err != nil {
			fmt.Fprintf(w, "Data Duplicate")
		} else {
			fmt.Fprintf(w, "Data Created")
		}

	}
}

func getCustomer(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var customers []Customer
	params := mux.Vars(r)

	sql := `SELECT
				CustomerID,
				IFNULL(CompanyName,''),
				IFNULL(ContactName,'') ContactName,
				IFNULL(ContactTitle,'') ContactTitle,
				IFNULL(Address,'') Address,
				IFNULL(City,'') City,
				IFNULL(Country,'') Country,
				IFNULL(Phone,'') Phone ,
				IFNULL(PostalCode,'') PostalCode
			FROM customers WHERE CustomerID = ?`

	result, err := db.Query(sql, params["id"])

	if err != nil {
		panic(err.Error())
	}

	defer result.Close()

	var customer Customer

	for result.Next() {

		err := result.Scan(&customer.CustomerID, &customer.CompanyName, &customer.ContactName,
			&customer.ContactTitle, &customer.Address, &customer.City, &customer.Country,
			&customer.Phone, &customer.PostalCode)

		if err != nil {
			panic(err.Error())
		}

		customers = append(customers, customer)
	}

	json.NewEncoder(w).Encode(customers)
}

func updateCustomer(w http.ResponseWriter, r *http.Request) {

	if r.Method == "PUT" {

		params := mux.Vars(r)

		newCompanyName := r.FormValue("CompanyName")
		newContactName := r.FormValue("ContactName")

		stmt, err := db.Prepare("UPDATE customers SET CompanyName = ?,ContactName = ? WHERE CustomerID = ?")

		_, err = stmt.Exec(newCompanyName, newContactName, params["id"])

		if err != nil {
			fmt.Fprintf(w, "Data not found or Request error")
		}

		fmt.Fprintf(w, "Customer with CustomerID = %s was updated", params["id"])
	}
}

func deleteCustomer(w http.ResponseWriter, r *http.Request) {

	params := mux.Vars(r)
	stmt, err := db.Prepare("DELETE FROM customers WHERE CustomerID = ?")

	_, err = stmt.Exec(params["id"])

	if err != nil {
		fmt.Fprintf(w, "delete failed")
	}

	fmt.Fprintf(w, "Customer with ID = %s was deleted", params["id"])
}

func getPost(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "application/json")

	var customers []Customer

	CustomerID := r.FormValue("CustomerID")
	CompanyName := r.FormValue("CompanyName")

	sql := `SELECT
				CustomerID,
				IFNULL(CompanyName,''),
				IFNULL(ContactName,'') ContactName,
				IFNULL(ContactTitle,'') ContactTitle,
				IFNULL(Address,'') Address,
				IFNULL(City,'') City,
				IFNULL(Country,'') Country,
				IFNULL(Phone,'') Phone ,
				IFNULL(PostalCode,'') PostalCode
			FROM customers WHERE CustomerID = ? AND CompanyName = ?`

	result, err := db.Query(sql, CustomerID, CompanyName)

	if err != nil {
		panic(err.Error())
	}

	defer result.Close()

	var customer Customer

	for result.Next() {

		err := result.Scan(&customer.CustomerID, &customer.CompanyName, &customer.ContactName,
			&customer.ContactTitle, &customer.Address, &customer.City, &customer.Country,
			&customer.Phone, &customer.PostalCode)

		if err != nil {
			panic(err.Error())
		}

		customers = append(customers, customer)
	}

	json.NewEncoder(w).Encode(customers)

}

func getAccounts(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var Accounts []account

	sql := `SELECT
				account_id,
				IFNULL(name,''),
				IFNULL(acc_num,'') acc_num,
				IFNULL(debit,'') debit,
				IFNULL(credit,'') credit,
				IFNULL(balance,'') balance,
				IFNULL(parent_id,'') parent_id
				
			FROM accounts`

	result, err := db.Query(sql)

	defer result.Close()

	if err != nil {
		panic(err.Error())
	}

	for result.Next() {

		var account account
		err := result.Scan(&account.Accountid, &account.Name, &account.Accnum,
			&account.Debit, &account.Credit, &account.Balance, &account.Parentid)

		if err != nil {
			panic(err.Error())
		}
		Accounts = append(Accounts, account)
	}

	json.NewEncoder(w).Encode(Accounts)
}
func createAccount(w http.ResponseWriter, r *http.Request) {

	if r.Method == "POST" {

		// Accountid := r.FormValue("account_id")
		Name := r.FormValue("name")
		Accnum := r.FormValue("acc_num")
		Debit := r.FormValue("debit")
		Credit := r.FormValue("credit")
		Balance := r.FormValue("balance")
		Parentid := r.FormValue("parent_id")

		stmt, err := db.Prepare("INSERT INTO accounts VALUES (NULL,?,?,?,?,?,?)")

		_, err = stmt.Exec(Name, Accnum, Debit, Credit, Balance, Parentid)

		if err != nil {
			fmt.Fprintf(w, "Data Duplicate")
		} else {
			fmt.Fprintf(w, "Data Created")
		}

	}
}
func getAccount(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	var Accounts []account

	sql := `SELECT
				account_id,
				IFNULL(name,''),
				IFNULL(acc_num,'') acc_num,
				IFNULL(debit,'') debit,
				IFNULL(credit,'') credit,
				IFNULL(balance,'') balance,
				IFNULL(parent_id,'') parent_id
				
			FROM accounts where account_id=?`

	result, err := db.Query(sql, params["id"])

	defer result.Close()

	if err != nil {
		panic(err.Error())
	}

	for result.Next() {

		var account account
		err := result.Scan(&account.Accountid, &account.Name, &account.Accnum,
			&account.Debit, &account.Credit, &account.Balance, &account.Parentid)

		if err != nil {
			panic(err.Error())
		}
		Accounts = append(Accounts, account)
	}

	json.NewEncoder(w).Encode(Accounts)
}
func updateAccount(w http.ResponseWriter, r *http.Request) {

	if r.Method == "PUT" {

		params := mux.Vars(r)

		newAccountname := r.FormValue("name")

		stmt, err := db.Prepare("UPDATE accounts SET name = ? WHERE account_id = ?")

		_, err = stmt.Exec(newAccountname, params["id"])

		if err != nil {
			fmt.Fprintf(w, "Data not found or Request error")
		}

		fmt.Fprintf(w, "Customer with AccountID = %s was updated", params["id"])
	}
}
func deleteAcount(w http.ResponseWriter, r *http.Request) {

	params := mux.Vars(r)
	stmt, err := db.Prepare("DELETE FROM accounts WHERE account_id = ?")

	_, err = stmt.Exec(params["id"])

	if err != nil {
		fmt.Fprintf(w, "delete failed")
	}

	fmt.Fprintf(w, "Account with ID = %s was deleted", params["id"])
}
func deleteAcountPost(w http.ResponseWriter, r *http.Request) {
	Accountid := r.FormValue("account_id")
	stmt, err := db.Prepare("DELETE FROM accounts WHERE account_id = ?")

	_, err = stmt.Exec(Accountid)

	if err != nil {
		fmt.Fprintf(w, "delete failed")
	}

	fmt.Fprintf(w, "Account with ID = "+Accountid+" was deleted")
}
func getPostAccount(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	Accountid := r.FormValue("account_id")
	Name := r.FormValue("name")

	var Accounts []account

	sql := `SELECT
				account_id,
				IFNULL(name,''),
				IFNULL(acc_num,'') acc_num,
				IFNULL(debit,'') debit,
				IFNULL(credit,'') credit,
				IFNULL(balance,'') balance,
				IFNULL(parent_id,'') parent_id
				
			FROM accounts where account_id = ? and name = ?`

	result, err := db.Query(sql, Accountid, Name)

	defer result.Close()

	if err != nil {
		panic(err.Error())
	}

	for result.Next() {

		var account account
		err := result.Scan(&account.Accountid, &account.Name, &account.Accnum,
			&account.Debit, &account.Credit, &account.Balance, &account.Parentid)

		if err != nil {
			panic(err.Error())
		}
		Accounts = append(Accounts, account)
	}

	json.NewEncoder(w).Encode(Accounts)
}

func getAdjtables(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var Adjtablearray []adjtable

	sql := `SELECT
				emp_id,
				IFNULL(name,''),
				IFNULL(salary,'') salary,
				IFNULL(boss_id,'') boss_id
				
			FROM adj_table`

	result, err := db.Query(sql)

	defer result.Close()

	if err != nil {
		panic(err.Error())
	}

	for result.Next() {

		var adjtables adjtable
		err := result.Scan(&adjtables.Empid, &adjtables.Name, &adjtables.Salary,
			&adjtables.Bossid)

		if err != nil {
			panic(err.Error())
		}
		Adjtablearray = append(Adjtablearray, adjtables)
	}

	json.NewEncoder(w).Encode(Adjtablearray)
}
func createAdjtable(w http.ResponseWriter, r *http.Request) {

	if r.Method == "POST" {

		// Empid := r.FormValue("emp_id")
		Name := r.FormValue("name")
		Salary := r.FormValue("salary")
		Bossid := r.FormValue("boss_id")

		stmt, err := db.Prepare("INSERT INTO adj_table VALUES (NULL,?,?,?)")

		_, err = stmt.Exec(Name, Salary, Bossid)

		if err != nil {
			fmt.Fprintf(w, "Data Duplicate")
		} else {
			fmt.Fprintf(w, "Data Created")
		}

	}
}
func getAdjTable(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	var Adjtablearray []adjtable

	sql := `SELECT
			emp_id,
			IFNULL(name,''),
			IFNULL(salary,'') salary,
			IFNULL(boss_id,'') boss_id
			
		FROM adj_table where emp_id=?`

	result, err := db.Query(sql, params["id"])

	defer result.Close()

	if err != nil {
		panic(err.Error())
	}

	for result.Next() {

		var adjtables adjtable
		err := result.Scan(&adjtables.Empid, &adjtables.Name, &adjtables.Salary,
			&adjtables.Bossid)

		if err != nil {
			panic(err.Error())
		}
		Adjtablearray = append(Adjtablearray, adjtables)
	}

	json.NewEncoder(w).Encode(Adjtablearray)
}
func updateAdjtable(w http.ResponseWriter, r *http.Request) {

	if r.Method == "PUT" {

		params := mux.Vars(r)

		newAdjtablename := r.FormValue("name")

		stmt, err := db.Prepare("UPDATE adj_table SET name = ? WHERE emp_id = ?")

		_, err = stmt.Exec(newAdjtablename, params["id"])

		if err != nil {
			fmt.Fprintf(w, "Data not found or Request error")
		}

		fmt.Fprintf(w, "Customer with AccountID = %s was updated", params["id"])
	}
}
func deleteAdjtable(w http.ResponseWriter, r *http.Request) {

	params := mux.Vars(r)
	stmt, err := db.Prepare("DELETE FROM adj_table WHERE emp_id = ?")

	_, err = stmt.Exec(params["id"])

	if err != nil {
		fmt.Fprintf(w, "delete failed")
	}

	fmt.Fprintf(w, "Adj with ID = %s was deleted", params["id"])
}
func deleteAdjTablePost(w http.ResponseWriter, r *http.Request) {
	Empid := r.FormValue("emp_id")
	stmt, err := db.Prepare("DELETE FROM adj_table WHERE emp_id = ?")

	_, err = stmt.Exec(Empid)

	if err != nil {
		fmt.Fprintf(w, "delete failed")
	}

	fmt.Fprintf(w, "Adj with ID = "+Empid+" was deleted")
}
func getPostAdj(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	Empid := r.FormValue("emp_id")
	Name := r.FormValue("name")

	var Accounts []account

	sql := `SELECT
				emp_id,
			IFNULL(name,''),
			IFNULL(salary,'') salary,
			IFNULL(boss_id,'') boss_id
			
		FROM adj_table where emp_id=? and name ?`

	result, err := db.Query(sql, Empid, Name)

	defer result.Close()

	if err != nil {
		panic(err.Error())
	}

	for result.Next() {

		var account account
		err := result.Scan(&account.Accountid, &account.Name, &account.Accnum,
			&account.Debit, &account.Credit, &account.Balance, &account.Parentid)

		if err != nil {
			panic(err.Error())
		}
		Accounts = append(Accounts, account)
	}

	json.NewEncoder(w).Encode(Accounts)
}

func getBooks(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var Bookarray []book

	sql := `SELECT
				BookID,
				IFNULL(Title,''),
				IFNULL(Author,'') Author,
				IFNULL(DatePublished,'') DatePublished,
				IFNULL(Publisher,'') Publisher,
				IFNULL(Edition,'') Edition,
				IFNULL(City,'') City
				
			FROM Book`

	result, err := db.Query(sql)

	defer result.Close()

	if err != nil {
		panic(err.Error())
	}

	for result.Next() {

		var books book
		err := result.Scan(&books.Bookid, &books.Title, &books.Author,
			&books.DatePublished, &books.Publisher,
			&books.Edition, &books.City)

		if err != nil {
			panic(err.Error())
		}
		Bookarray = append(Bookarray, books)
	}

	json.NewEncoder(w).Encode(Bookarray)
}
func createBook(w http.ResponseWriter, r *http.Request) {

	if r.Method == "POST" {

		Title := r.FormValue("Title")
		Author := r.FormValue("Author")
		DatePublished := r.FormValue("DatePublished")
		Publisher := r.FormValue("Publisher")
		Edition := r.FormValue("Edition")
		City := r.FormValue("City")

		stmt, err := db.Prepare("INSERT INTO Book VALUES (NULL,?,?,?,?,?,?)")

		_, err = stmt.Exec(Title, Author, DatePublished, Publisher, Edition, City)

		if err != nil {
			fmt.Fprintf(w, "Data Duplicate")
		} else {
			fmt.Fprintf(w, "Data Created")
		}

	}
}
func getBook(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	var Bookarray []book

	sql := `SELECT
			BookID,
			IFNULL(Title,''),
			IFNULL(Author,'') Author,
			IFNULL(DatePublished,'') DatePublished,
			IFNULL(Publisher,'') Publisher,
			IFNULL(Edition,'') Edition,
			IFNULL(City,'') City

			
		FROM Book where BookID=?`

	result, err := db.Query(sql, params["id"])

	defer result.Close()

	if err != nil {
		panic(err.Error())
	}

	for result.Next() {

		var books book
		err := result.Scan(&books.Bookid, &books.Title, &books.Author,
			&books.DatePublished, &books.Publisher,
			&books.Edition, &books.City)

		if err != nil {
			panic(err.Error())
		}
		Bookarray = append(Bookarray, books)
	}

	json.NewEncoder(w).Encode(Bookarray)
}
func updateBook(w http.ResponseWriter, r *http.Request) {

	if r.Method == "PUT" {

		params := mux.Vars(r)

		newTitle := r.FormValue("Title")

		stmt, err := db.Prepare("UPDATE Book SET Title = ? WHERE BookID = ?")

		_, err = stmt.Exec(newTitle, params["id"])

		if err != nil {
			fmt.Fprintf(w, "Data not found or Request error")
		}

		fmt.Fprintf(w, "Customer with BookID = %s was updated", params["id"])
	}
}
func deleteBook(w http.ResponseWriter, r *http.Request) {

	params := mux.Vars(r)
	stmt, err := db.Prepare("DELETE FROM Book WHERE BookID = ?")

	_, err = stmt.Exec(params["id"])

	if err != nil {
		fmt.Fprintf(w, "delete failed")
	}

	fmt.Fprintf(w, "Book with ID = %s was deleted", params["id"])
}
func deleteBookPost(w http.ResponseWriter, r *http.Request) {
	Bookid := r.FormValue("BookID")
	stmt, err := db.Prepare("DELETE FROM Book WHERE BookID = ?")

	_, err = stmt.Exec(Bookid)

	if err != nil {
		fmt.Fprintf(w, "delete failed")
	}

	fmt.Fprintf(w, "Adj with ID = "+Bookid+" was deleted")
}
func getPostBook(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	Bookid := r.FormValue("BookID")
	Title := r.FormValue("Title")

	var Bookarray []book

	sql := `SELECT
			BookID,
			IFNULL(Title,''),
			IFNULL(Author,'') Author,
			IFNULL(DatePublished,'') DatePublished,
			IFNULL(Publisher,'') Publisher,
			IFNULL(Edition,'') Edition,
			IFNULL(City,'') City

		FROM Book where BookID=? and Title =?`

	result, err := db.Query(sql, Bookid, Title)

	defer result.Close()

	if err != nil {
		panic(err.Error())
	}

	for result.Next() {

		var books book
		err := result.Scan(&books.Bookid, &books.Title, &books.Author,
			&books.DatePublished, &books.Publisher,
			&books.Edition, &books.City)

		if err != nil {
			panic(err.Error())
		}
		Bookarray = append(Bookarray, books)
	}

	json.NewEncoder(w).Encode(Bookarray)
}

func getCategories(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var Categoriearray []categorie

	sql := `SELECT
				CategoryID,
				IFNULL(CategoryName,''),
				IFNULL(Description,'') Description,
				IFNULL(Picture,'') Picture
				
			FROM Categories`

	result, err := db.Query(sql)

	defer result.Close()

	if err != nil {
		panic(err.Error())
	}

	for result.Next() {

		var categories categorie
		err := result.Scan(&categories.Categoryid, &categories.Categoryname, &categories.Description,
			&categories.Picture)

		if err != nil {
			panic(err.Error())
		}
		Categoriearray = append(Categoriearray, categories)
	}

	json.NewEncoder(w).Encode(Categoriearray)
}

// func createCategorie(w http.ResponseWriter, r *http.Request) {

// 	if r.Method == "POST" {

// 		CategoryName := r.FormValue("CategoryName")
// 		Description := r.FormValue("Description")
// 		Picture, handler, err := r.FormFile("Picture")

// 		stmt, err := db.Prepare("INSERT INTO Categories VALUES (NULL,?,?,?)")

// 		_, err = stmt.Exec(CategoryName, Description, Picture)

// 		if err != nil {
// 			fmt.Fprintf(w, "Data Duplicate")
// 		} else {
// 			fmt.Fprintf(w, "Data Created")
// 		}

// 	}
// }

// func getBook(w http.ResponseWriter, r *http.Request) {
// 	w.Header().Set("Content-Type", "application/json")
// 	params := mux.Vars(r)
// 	var Bookarray []book

// 	sql := `SELECT
// 			BookID,
// 			IFNULL(Title,''),
// 			IFNULL(Author,'') Author,
// 			IFNULL(DatePublished,'') DatePublished,
// 			IFNULL(Publisher,'') Publisher,
// 			IFNULL(Edition,'') Edition,
// 			IFNULL(City,'') City

// 		FROM Book where BookID=?`

// 	result, err := db.Query(sql, params["id"])

// 	defer result.Close()

// 	if err != nil {
// 		panic(err.Error())
// 	}

// 	for result.Next() {

// 		var books book
// 		err := result.Scan(&books.Bookid, &books.Title, &books.Author,
// 			&books.DatePublished, &books.Publisher,
// 			&books.Edition, &books.City)

// 		if err != nil {
// 			panic(err.Error())
// 		}
// 		Bookarray = append(Bookarray, books)
// 	}

// 	json.NewEncoder(w).Encode(Bookarray)
// }

// func updateBook(w http.ResponseWriter, r *http.Request) {

// 	if r.Method == "PUT" {

// 		params := mux.Vars(r)

// 		newTitle := r.FormValue("Title")

// 		stmt, err := db.Prepare("UPDATE Book SET Title = ? WHERE BookID = ?")

// 		_, err = stmt.Exec(newTitle, params["id"])

// 		if err != nil {
// 			fmt.Fprintf(w, "Data not found or Request error")
// 		}

// 		fmt.Fprintf(w, "Customer with BookID = %s was updated", params["id"])
// 	}
// }
// func deleteBook(w http.ResponseWriter, r *http.Request) {

// 	params := mux.Vars(r)
// 	stmt, err := db.Prepare("DELETE FROM Book WHERE BookID = ?")

// 	_, err = stmt.Exec(params["id"])

// 	if err != nil {
// 		fmt.Fprintf(w, "delete failed")
// 	}

// 	fmt.Fprintf(w, "Book with ID = %s was deleted", params["id"])
// }
// func deleteBookPost(w http.ResponseWriter, r *http.Request) {
// 	Bookid := r.FormValue("BookID")
// 	stmt, err := db.Prepare("DELETE FROM Book WHERE BookID = ?")

// 	_, err = stmt.Exec(Bookid)

// 	if err != nil {
// 		fmt.Fprintf(w, "delete failed")
// 	}

// 	fmt.Fprintf(w, "Adj with ID = "+Bookid+" was deleted")
// }
// func getPostBook(w http.ResponseWriter, r *http.Request) {
// 	w.Header().Set("Content-Type", "application/json")
// 	Bookid := r.FormValue("BookID")
// 	Title := r.FormValue("Title")

// 	var Bookarray []book

// 	sql := `SELECT
// 			BookID,
// 			IFNULL(Title,''),
// 			IFNULL(Author,'') Author,
// 			IFNULL(DatePublished,'') DatePublished,
// 			IFNULL(Publisher,'') Publisher,
// 			IFNULL(Edition,'') Edition,
// 			IFNULL(City,'') City

// 		FROM Book where BookID=? and Title =?`

// 	result, err := db.Query(sql, Bookid, Title)

// 	defer result.Close()

// 	if err != nil {
// 		panic(err.Error())
// 	}

// 	for result.Next() {

// 		var books book
// 		err := result.Scan(&books.Bookid, &books.Title, &books.Author,
// 			&books.DatePublished, &books.Publisher,
// 			&books.Edition, &books.City)

// 		if err != nil {
// 			panic(err.Error())
// 		}
// 		Bookarray = append(Bookarray, books)
// 	}

// 	json.NewEncoder(w).Encode(Bookarray)
// }

func getCitys(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var Cityarray []city

	sql := `SELECT
				CityID,
				IFNULL(ProvinceID,''),
				IFNULL(CityName,'') CityName
				
			FROM City`

	result, err := db.Query(sql)

	defer result.Close()

	if err != nil {
		panic(err.Error())
	}

	for result.Next() {

		var citys city
		err := result.Scan(&citys.Cityid, &citys.Provinceid, &citys.Cityname)

		if err != nil {
			panic(err.Error())
		}
		Cityarray = append(Cityarray, citys)
	}

	json.NewEncoder(w).Encode(Cityarray)
}
func createCity(w http.ResponseWriter, r *http.Request) {

	if r.Method == "POST" {

		ProvinceID := r.FormValue("ProvinceID")
		CityName := r.FormValue("CityName")

		stmt, err := db.Prepare("INSERT INTO city VALUES (NULL,?,?)")

		_, err = stmt.Exec(ProvinceID, CityName)

		if err != nil {
			fmt.Fprintf(w, "Data Duplicate")
		} else {
			fmt.Fprintf(w, "Data Created")
		}

	}
}
func getCity(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	var Cityarray []city

	sql := `SELECT
			CityID,
			IFNULL(ProvinceID,''),
			IFNULL(CityName,'') CityName
			

			
		FROM Book where CityID=?`

	result, err := db.Query(sql, params["id"])

	defer result.Close()

	if err != nil {
		panic(err.Error())
	}

	for result.Next() {

		var citys city
		err := result.Scan(&citys.Cityid, &citys.Provinceid, &citys.Cityname)

		if err != nil {
			panic(err.Error())
		}
		Cityarray = append(Cityarray, citys)
	}

	json.NewEncoder(w).Encode(Cityarray)
}
func updateCity(w http.ResponseWriter, r *http.Request) {

	if r.Method == "PUT" {

		params := mux.Vars(r)

		newCityName := r.FormValue("CityName")

		stmt, err := db.Prepare("UPDATE City SET CityName = ? WHERE CityID = ?")

		_, err = stmt.Exec(newCityName, params["id"])

		if err != nil {
			fmt.Fprintf(w, "Data not found or Request error")
		}

		fmt.Fprintf(w, "Customer with CityID = %s was updated", params["id"])
	}
}
func deleteCity(w http.ResponseWriter, r *http.Request) {

	params := mux.Vars(r)
	stmt, err := db.Prepare("DELETE FROM City WHERE CityID = ?")

	_, err = stmt.Exec(params["id"])

	if err != nil {
		fmt.Fprintf(w, "delete failed")
	}

	fmt.Fprintf(w, "City with ID = %s was deleted", params["id"])
}
func deleteCityPost(w http.ResponseWriter, r *http.Request) {
	Cityid := r.FormValue("CityID")
	stmt, err := db.Prepare("DELETE FROM City WHERE CityID = ?")

	_, err = stmt.Exec(Cityid)

	if err != nil {
		fmt.Fprintf(w, "delete failed")
	}

	fmt.Fprintf(w, "City with ID = "+Cityid+" was deleted")
}
func getPostCity(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	CityID := r.FormValue("CityID")
	ProvinceID := r.FormValue("ProvinceID")

	var Cityarray []city

	sql := `SELECT
			CityID,
			IFNULL(ProvinceID,''),
			IFNULL(CityName,'') CityName
		FROM city where CityID=? and ProvinceID =?`

	result, err := db.Query(sql, CityID, ProvinceID)

	defer result.Close()

	if err != nil {
		panic(err.Error())
	}

	for result.Next() {

		var citys city
		err := result.Scan(&citys.Cityid, &citys.Provinceid, &citys.Cityname)

		if err != nil {
			panic(err.Error())
		}
		Cityarray = append(Cityarray, citys)
	}

	json.NewEncoder(w).Encode(Cityarray)
}

func getCountries(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var Countriearray []countrie

	sql := `SELECT
				CountryID,
				IFNULL(CountryName,'')
				
			FROM Countries`

	result, err := db.Query(sql)

	defer result.Close()

	if err != nil {
		panic(err.Error())
	}

	for result.Next() {

		var countries countrie
		err := result.Scan(&countries.Countryid, &countries.Countryname)

		if err != nil {
			panic(err.Error())
		}
		Countriearray = append(Countriearray, countries)
	}

	json.NewEncoder(w).Encode(Countriearray)
}
func createCountrie(w http.ResponseWriter, r *http.Request) {

	if r.Method == "POST" {

		CountryName := r.FormValue("CountryName")

		stmt, err := db.Prepare("INSERT INTO Countries VALUES (NULL,?)")

		_, err = stmt.Exec(CountryName)

		if err != nil {
			fmt.Fprintf(w, "Data Duplicate")
		} else {
			fmt.Fprintf(w, "Data Created")
		}

	}
}
func getCountrie(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	var Countriearray []countrie

	sql := `SELECT
			CountryID,
			IFNULL(CountryName,'')
			
		FROM Countries where CountryID=?`

	result, err := db.Query(sql, params["id"])

	defer result.Close()

	if err != nil {
		panic(err.Error())
	}

	for result.Next() {

		var countries countrie
		err := result.Scan(&countries.Countryid, &countries.Countryname)

		if err != nil {
			panic(err.Error())
		}
		Countriearray = append(Countriearray, countries)
	}

	json.NewEncoder(w).Encode(Countriearray)
}
func updateCountrie(w http.ResponseWriter, r *http.Request) {

	if r.Method == "PUT" {

		params := mux.Vars(r)

		newCountryName := r.FormValue("CountryName")

		stmt, err := db.Prepare("UPDATE Countries SET CountryName = ? WHERE CountryID = ?")

		_, err = stmt.Exec(newCountryName, params["id"])

		if err != nil {
			fmt.Fprintf(w, "Data not found or Request error")
		}

		fmt.Fprintf(w, "Country with CountryID = %s was updated", params["id"])
	}
}
func deleteCountrie(w http.ResponseWriter, r *http.Request) {

	params := mux.Vars(r)
	stmt, err := db.Prepare("DELETE FROM Contriesty WHERE CountryID = ?")

	_, err = stmt.Exec(params["id"])

	if err != nil {
		fmt.Fprintf(w, "delete failed")
	}

	fmt.Fprintf(w, "City with ID = %s was deleted", params["id"])
}
func deleteCountriePost(w http.ResponseWriter, r *http.Request) {
	Countryid := r.FormValue("CountryID")
	stmt, err := db.Prepare("DELETE FROM Country WHERE CountryID = ?")

	_, err = stmt.Exec(Countryid)

	if err != nil {
		fmt.Fprintf(w, "delete failed")
	}

	fmt.Fprintf(w, "Country with ID = "+Countryid+" was deleted")
}
func getPostContrie(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	CountryID := r.FormValue("CountryID")
	var Countriearray []countrie

	sql := `SELECT
			CountryID,
			IFNULL(CountryName,'')
			
		FROM Countries where CountryID=?`

	result, err := db.Query(sql, CountryID)

	defer result.Close()

	if err != nil {
		panic(err.Error())
	}

	for result.Next() {

		var countries countrie
		err := result.Scan(&countries.Countryid, &countries.Countryname)

		if err != nil {
			panic(err.Error())
		}
		Countriearray = append(Countriearray, countries)
	}

	json.NewEncoder(w).Encode(Countriearray)
}

type employee struct {
	EmployeeID      string `json:"EmployeeID"`
	LastName        string `json:"LastName"`
	FirstName       string `json:"FirstName"`
	Title           string `json:"Title"`
	TitleOfCourtesy string `json:"TitleOfCourtesy"`
	BirthDate       string `json:"BirthDate"`
	HireDate        string `json:"HireDate"`
	Address         string `json:"Address"`
	City            string `json:"City"`
	Region          string `json:"Region"`
	PostalCode      string `json:"PostalCode"`
	Country         string `json:"Country"`
	HomePhone       string `json:"HomePhone"`
	Extension       string `json:"Extension"`
	Photo           string `json:"Photo"`
	Notes           string `json:"Notes"`
	ReportsTo       string `json:"ReportsTo"`
	ProvinceName    string `json:"ProvinceName"`
}

func getEmployees(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var Employeearray []employee

	sql := `SELECT
				EmployeeID,
				IFNULL(LastName,''),
				IFNULL(FirstName,'') FirstName,
				IFNULL(Title,'') Title,
				IFNULL(TitleOfCourtesy,'') TitleOfCourtesy,
				IFNULL(BirthDate,'') BirthDate,
				IFNULL(HireDate,'') HireDate,
				IFNULL(Address,'') Address,
				IFNULL(City,'') City,
				IFNULL(Region,'') Region,
				IFNULL(PostalCode,'') PostalCode,
				IFNULL(Country,'') Country,
				IFNULL(HomePhone,'') HomePhone,
				IFNULL(Extension,'') Extension,
				IFNULL(Photo,'') Photo,
				IFNULL(Notes,'') Notes,
				IFNULL(ReportsTo,'') ReportsTo,
				IFNULL(ProvinceName,'') ProvinceName

				
			FROM employees`

	result, err := db.Query(sql)

	defer result.Close()

	if err != nil {
		panic(err.Error())
	}

	for result.Next() {

		var employees employee
		err := result.Scan(&employees.EmployeeID, &employees.LastName,
			&employees.FirstName, &employees.Title, &employees.TitleOfCourtesy,
			&employees.BirthDate, &employees.HireDate, &employees.Address,
			&employees.City, &employees.Region, &employees.PostalCode,
			&employees.Country, &employees.HomePhone, &employees.Extension,
			&employees.Photo, &employees.Notes, &employees.ReportsTo,
			&employees.ProvinceName)

		if err != nil {
			panic(err.Error())
		}
		Employeearray = append(Employeearray, employees)
	}

	json.NewEncoder(w).Encode(Employeearray)
}
func createEmployee(w http.ResponseWriter, r *http.Request) {

	if r.Method == "POST" {

		LastName := r.FormValue("LastName")
		FirstName := r.FormValue("FirstName")
		Title := r.FormValue("Title")
		TitleOfCourtesy := r.FormValue("TitleOfCourtesy")
		BirthDate := r.FormValue("BirthDate")
		HireDate := r.FormValue("HireDate")
		Address := r.FormValue("Address")
		City := r.FormValue("City")
		Region := r.FormValue("Region")
		PostalCode := r.FormValue("PostalCode")
		Country := r.FormValue("Country")
		HomePhone := r.FormValue("HomePhone")
		Extension := r.FormValue("Extension")
		Photo := r.FormValue("Photo")
		Notes := r.FormValue("Notes")
		ReportsTo := r.FormValue("ReportsTo")
		ProvinceName := r.FormValue("ProvinceName")

		stmt, err := db.Prepare("INSERT INTO employees VALUES (NULL,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?)")

		_, err = stmt.Exec(LastName, FirstName, Title,
			TitleOfCourtesy, BirthDate, HireDate, Address,
			City, Region, PostalCode, Country,
			HomePhone, Extension, Photo, Notes, ReportsTo,
			ProvinceName)

		if err != nil {
			fmt.Fprintf(w, "Data Duplicate")
		} else {
			fmt.Fprintf(w, "Data Created")
		}

	}
}
func getEmployee(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	var Employeearray []employee

	sql := `SELECT
			EmployeeID,
			IFNULL(LastName,''),
			IFNULL(FirstName,'') FirstName,
			IFNULL(Title,'') Title,
			IFNULL(TitleOfCourtesy,'') TitleOfCourtesy,
			IFNULL(BirthDate,'') BirthDate,
			IFNULL(HireDate,'') HireDate,
			IFNULL(Address,'') Address,
			IFNULL(City,'') City,
			IFNULL(Region,'') Region,
			IFNULL(PostalCode,'') PostalCode,
			IFNULL(Country,'') Country,
			IFNULL(HomePhone,'') HomePhone,
			IFNULL(Extension,'') Extension,
			IFNULL(Photo,'') Photo,
			IFNULL(Notes,'') Notes,
			IFNULL(ReportsTo,'') ReportsTo,
			IFNULL(ProvinceName,'') ProvinceName
					
		FROM employees where EmployeeID=?`

	result, err := db.Query(sql, params["id"])

	defer result.Close()

	if err != nil {
		panic(err.Error())
	}

	for result.Next() {

		var employees employee
		err := result.Scan(&employees.EmployeeID, &employees.LastName,
			&employees.FirstName, &employees.Title, &employees.TitleOfCourtesy,
			&employees.BirthDate, &employees.HireDate, &employees.Address,
			&employees.City, &employees.Region, &employees.PostalCode,
			&employees.Country, &employees.HomePhone, &employees.Extension,
			&employees.Photo, &employees.Notes, &employees.ReportsTo,
			&employees.ProvinceName)

		if err != nil {
			panic(err.Error())
		}
		Employeearray = append(Employeearray, employees)
	}

	json.NewEncoder(w).Encode(Employeearray)
}
func updateEmployee(w http.ResponseWriter, r *http.Request) {

	if r.Method == "PUT" {

		params := mux.Vars(r)

		newFirstName := r.FormValue("FirstName")
		newLastName := r.FormValue("LastName")

		stmt, err := db.Prepare("UPDATE employees SET FirstName = ?, LastName =? WHERE EmployeeID = ?")

		_, err = stmt.Exec(newFirstName, newLastName, params["id"])

		if err != nil {
			fmt.Fprintf(w, "Data not found or Request error")
		}

		fmt.Fprintf(w, "Employee with CityID = %s was updated", params["id"])
	}
}
func deleteEmployee(w http.ResponseWriter, r *http.Request) {

	params := mux.Vars(r)
	stmt, err := db.Prepare("DELETE FROM employees WHERE EmployeeID = ?")

	_, err = stmt.Exec(params["id"])

	if err != nil {
		fmt.Fprintf(w, "delete failed")
	}

	fmt.Fprintf(w, "Employee with ID = %s was deleted", params["id"])
}
func deleteEmployeePost(w http.ResponseWriter, r *http.Request) {
	EmployeeID := r.FormValue("EmployeeID")
	stmt, err := db.Prepare("DELETE FROM employees WHERE EmployeeID = ?")

	_, err = stmt.Exec(EmployeeID)

	if err != nil {
		fmt.Fprintf(w, "delete failed")
	}

	fmt.Fprintf(w, "Employee with ID = "+EmployeeID+" was deleted")
}
func getPostEmployee(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	EmployeeID := r.FormValue("EmployeeID")
	FirstName := r.FormValue("FirstName")

	var Employeearray []employee

	sql := `SELECT
			EmployeeID,
			IFNULL(LastName,''),
			IFNULL(FirstName,'') FirstName,
			IFNULL(Title,'') Title,
			IFNULL(TitleOfCourtesy,'') TitleOfCourtesy,
			IFNULL(BirthDate,'') BirthDate,
			IFNULL(HireDate,'') HireDate,
			IFNULL(Address,'') Address,
			IFNULL(City,'') City,
			IFNULL(Region,'') Region,
			IFNULL(PostalCode,'') PostalCode,
			IFNULL(Country,'') Country,
			IFNULL(HomePhone,'') HomePhone,
			IFNULL(Extension,'') Extension,
			IFNULL(Photo,'') Photo,
			IFNULL(Notes,'') Notes,
			IFNULL(ReportsTo,'') ReportsTo,
			IFNULL(ProvinceName,'') ProvinceName
					
		FROM employees where EmployeeID=? and FirstName = ?`

	result, err := db.Query(sql, EmployeeID, FirstName)

	defer result.Close()

	if err != nil {
		panic(err.Error())
	}

	for result.Next() {

		var employees employee
		err := result.Scan(&employees.EmployeeID, &employees.LastName,
			&employees.FirstName, &employees.Title, &employees.TitleOfCourtesy,
			&employees.BirthDate, &employees.HireDate, &employees.Address,
			&employees.City, &employees.Region, &employees.PostalCode,
			&employees.Country, &employees.HomePhone, &employees.Extension,
			&employees.Photo, &employees.Notes, &employees.ReportsTo,
			&employees.ProvinceName)

		if err != nil {
			panic(err.Error())
		}
		Employeearray = append(Employeearray, employees)
	}

	json.NewEncoder(w).Encode(Employeearray)
}

type load struct {
	Loadid     string `json:"id"`
	MemberID   string `json:"MemberID"`
	BookID     string `json:"BookID"`
	LoanDate   string `json:"LoanDate"`
	DueDate    string `json:"DueDate"`
	ReturnDate string `json:"ReturnDate"`
}

func getLoads(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var Loadarray []load

	sql := `SELECT
				id,
				IFNULL(MemberID,''),
				IFNULL(BookID,'') BookID,
				IFNULL(LoanDate,'') LoanDate,
				IFNULL(DueDate,'') DueDate,
				IFNULL(ReturnDate,'') ReturnDate

				
			FROM loads`

	result, err := db.Query(sql)

	defer result.Close()

	if err != nil {
		panic(err.Error())
	}

	for result.Next() {

		var loads load
		err := result.Scan(&loads.Loadid, &loads.MemberID,
			&loads.BookID, &loads.LoanDate, &loads.DueDate,
			&loads.ReturnDate)

		if err != nil {
			panic(err.Error())
		}
		Loadarray = append(Loadarray, loads)
	}

	json.NewEncoder(w).Encode(Loadarray)
}
func createload(w http.ResponseWriter, r *http.Request) {

	if r.Method == "POST" {

		MemberID := r.FormValue("MemberID")
		BookID := r.FormValue("BookID")
		LoanDate := r.FormValue("LoanDate")
		DueDate := r.FormValue("DueDate")
		ReturnDate := r.FormValue("ReturnDate")

		stmt, err := db.Prepare("INSERT INTO loads VALUES (NULL,?,?,?,?,?)")

		_, err = stmt.Exec(MemberID, BookID, LoanDate,
			DueDate, ReturnDate)

		if err != nil {
			fmt.Fprintf(w, "Data Duplicate")
		} else {
			fmt.Fprintf(w, "Data Created")
		}

	}
}
func getLoad(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	var Loadarray []load

	sql := `SELECT
			id,
			IFNULL(MemberID,''),
			IFNULL(BookID,'') BookID,
			IFNULL(LoanDate,'') LoanDate,
			IFNULL(DueDate,'') DueDate,
			IFNULL(ReturnDate,'') ReturnDate
					
		FROM loads where id=?`

	result, err := db.Query(sql, params["id"])

	defer result.Close()

	if err != nil {
		panic(err.Error())
	}

	for result.Next() {

		var loads load
		err := result.Scan(&loads.Loadid, &loads.MemberID,
			&loads.BookID, &loads.LoanDate, &loads.DueDate,
			&loads.ReturnDate)

		if err != nil {
			panic(err.Error())
		}
		Loadarray = append(Loadarray, loads)
	}

	json.NewEncoder(w).Encode(Loadarray)
}
func updateLoad(w http.ResponseWriter, r *http.Request) {

	if r.Method == "PUT" {

		params := mux.Vars(r)

		newMemberID := r.FormValue("MemberID")
		newBookID := r.FormValue("BookID")

		stmt, err := db.Prepare("UPDATE loads SET MemberID = ?, BookID =? WHERE id = ?")

		_, err = stmt.Exec(newMemberID, newBookID, params["id"])

		if err != nil {
			fmt.Fprintf(w, "Data not found or Request error")
		}

		fmt.Fprintf(w, "Load with ID = %s was updated", params["id"])
	}
}
func deleteLoad(w http.ResponseWriter, r *http.Request) {

	params := mux.Vars(r)
	stmt, err := db.Prepare("DELETE FROM loads WHERE id = ?")

	_, err = stmt.Exec(params["id"])

	if err != nil {
		fmt.Fprintf(w, "delete failed")
	}

	fmt.Fprintf(w, "Load with ID = %s was deleted", params["id"])
}
func deleteLoadPost(w http.ResponseWriter, r *http.Request) {
	id := r.FormValue("id")
	stmt, err := db.Prepare("DELETE FROM loads WHERE id = ?")

	_, err = stmt.Exec(id)

	if err != nil {
		fmt.Fprintf(w, "delete failed")
	}

	fmt.Fprintf(w, "Load with ID = "+id+" was deleted")
}
func getPostLoad(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	id := r.FormValue("id")
	MemberID := r.FormValue("MemberID")

	var Loadarray []load

	sql := `SELECT
			id,
			IFNULL(MemberID,''),
			IFNULL(BookID,'') BookID,
			IFNULL(LoanDate,'') LoanDate,
			IFNULL(DueDate,'') DueDate,
			IFNULL(ReturnDate,'') ReturnDate
					
		FROM loads where id=? and MemberID = ?`

	result, err := db.Query(sql, id, MemberID)

	defer result.Close()

	if err != nil {
		panic(err.Error())
	}

	for result.Next() {

		var loads load
		err := result.Scan(&loads.Loadid, &loads.MemberID,
			&loads.BookID, &loads.LoanDate, &loads.DueDate,
			&loads.ReturnDate)

		if err != nil {
			panic(err.Error())
		}
		Loadarray = append(Loadarray, loads)
	}

	json.NewEncoder(w).Encode(Loadarray)
}

type loan struct {
	Loanid     string `json:"id"`
	MemberID   string `json:"MemberID"`
	BookID     string `json:"BookID"`
	LoanDate   string `json:"LoanDate"`
	DueDate    string `json:"DueDate"`
	ReturnDate string `json:"ReturnDate"`
}

func getLoans(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var Loanarray []loan

	sql := `SELECT
				id,
				IFNULL(MemberID,''),
				IFNULL(BookID,'') BookID,
				IFNULL(LoanDate,'') LoanDate,
				IFNULL(DueDate,'') DueDate,
				IFNULL(ReturnDate,'') ReturnDate

				
			FROM loan`

	result, err := db.Query(sql)

	defer result.Close()

	if err != nil {
		panic(err.Error())
	}

	for result.Next() {

		var loans loan
		err := result.Scan(&loans.Loanid, &loans.MemberID,
			&loans.BookID, &loans.LoanDate, &loans.DueDate,
			&loans.ReturnDate)

		if err != nil {
			panic(err.Error())
		}
		Loanarray = append(Loanarray, loans)
	}

	json.NewEncoder(w).Encode(Loanarray)
}
func createloan(w http.ResponseWriter, r *http.Request) {

	if r.Method == "POST" {

		MemberID := r.FormValue("MemberID")
		BookID := r.FormValue("BookID")
		LoanDate := r.FormValue("LoanDate")
		DueDate := r.FormValue("DueDate")
		ReturnDate := r.FormValue("ReturnDate")

		stmt, err := db.Prepare("INSERT INTO loan VALUES (NULL,?,?,?,?,?)")

		_, err = stmt.Exec(MemberID, BookID, LoanDate,
			DueDate, ReturnDate)

		if err != nil {
			fmt.Fprintf(w, "Data Duplicate")
		} else {
			fmt.Fprintf(w, "Data Created")
		}

	}
}
func getLoan(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	var Loanarray []loan

	sql := `SELECT
			id,
			IFNULL(MemberID,''),
			IFNULL(BookID,'') BookID,
			IFNULL(LoanDate,'') LoanDate,
			IFNULL(DueDate,'') DueDate,
			IFNULL(ReturnDate,'') ReturnDate
					
		FROM loan where id=?`

	result, err := db.Query(sql, params["id"])

	defer result.Close()

	if err != nil {
		panic(err.Error())
	}

	for result.Next() {

		var loans loan
		err := result.Scan(&loans.Loanid, &loans.MemberID,
			&loans.BookID, &loans.LoanDate, &loans.DueDate,
			&loans.ReturnDate)

		if err != nil {
			panic(err.Error())
		}
		Loanarray = append(Loanarray, loans)
	}

	json.NewEncoder(w).Encode(Loanarray)
}
func updateLoan(w http.ResponseWriter, r *http.Request) {

	if r.Method == "PUT" {

		params := mux.Vars(r)

		newMemberID := r.FormValue("MemberID")
		newBookID := r.FormValue("BookID")

		stmt, err := db.Prepare("UPDATE loan SET MemberID = ?, BookID =? WHERE id = ?")

		_, err = stmt.Exec(newMemberID, newBookID, params["id"])

		if err != nil {
			fmt.Fprintf(w, "Data not found or Request error")
		}

		fmt.Fprintf(w, "Loan with ID = %s was updated", params["id"])
	}
}
func deleteLoan(w http.ResponseWriter, r *http.Request) {

	params := mux.Vars(r)
	stmt, err := db.Prepare("DELETE FROM loan WHERE id = ?")

	_, err = stmt.Exec(params["id"])

	if err != nil {
		fmt.Fprintf(w, "delete failed")
	}

	fmt.Fprintf(w, "Loan with ID = %s was deleted", params["id"])
}
func deleteLoanPost(w http.ResponseWriter, r *http.Request) {
	id := r.FormValue("id")
	stmt, err := db.Prepare("DELETE FROM loan WHERE id = ?")

	_, err = stmt.Exec(id)

	if err != nil {
		fmt.Fprintf(w, "delete failed")
	}

	fmt.Fprintf(w, "Loan with ID = "+id+" was deleted")
}
func getPostLoan(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	id := r.FormValue("id")
	MemberID := r.FormValue("MemberID")

	var Loanarray []loan

	sql := `SELECT
			id,
			IFNULL(MemberID,''),
			IFNULL(BookID,'') BookID,
			IFNULL(LoanDate,'') LoanDate,
			IFNULL(DueDate,'') DueDate,
			IFNULL(ReturnDate,'') ReturnDate
					
		FROM loan where id=? and MemberID = ?`

	result, err := db.Query(sql, id, MemberID)

	defer result.Close()

	if err != nil {
		panic(err.Error())
	}

	for result.Next() {

		var loans loan
		err := result.Scan(&loans.Loanid, &loans.MemberID,
			&loans.BookID, &loans.LoanDate, &loans.DueDate,
			&loans.ReturnDate)

		if err != nil {
			panic(err.Error())
		}
		Loanarray = append(Loanarray, loans)
	}

	json.NewEncoder(w).Encode(Loanarray)
}

type longorder struct {
	OrderID        string `json:"OrderID"`
	CustomerID     string `json:"CustomerID"`
	EmployeeID     string `json:"EmployeeID"`
	OrderDate      string `json:"OrderDate"`
	RequiredDate   string `json:"RequiredDate"`
	ShippedDate    string `json:"ShippedDate"`
	ShipVia        string `json:"ShipVia"`
	Freight        string `json:"Freight"`
	ShipName       string `json:"ShipName"`
	ShipAddress    string `json:"ShipAddress"`
	ShipCity       string `json:"ShipCity"`
	ShipRegion     string `json:"ShipRegion"`
	ShipPostalCode string `json:"ShipPostalCode"`
	ShipCountry    string `json:"ShipCountry"`
}

func getLongorders(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var Longorderarray []longorder

	sql := `SELECT
			OrderID,
				IFNULL(CustomerID,''),
				IFNULL(EmployeeID,'') EmployeeID,
				IFNULL(OrderDate,'') OrderDate,
				IFNULL(RequiredDate,'') RequiredDate,
				IFNULL(ShippedDate,'') ShippedDate,
				IFNULL(ShipVia,'') ShipVia,
				IFNULL(Freight,'') Freight,
				IFNULL(ShipName,'') ShipName,
				IFNULL(ShipAddress,'') ShipAddress,
				IFNULL(ShipCity,'') ShipCity,
				IFNULL(ShipRegion,'') ShipRegion,
				IFNULL(ShipPostalCode,'') ShipPostalCode,
				IFNULL(ShipCountry,'') ShipCountry

				
			FROM longorders`

	result, err := db.Query(sql)

	defer result.Close()

	if err != nil {
		panic(err.Error())
	}

	for result.Next() {

		var longorders longorder
		err := result.Scan(&longorders.OrderID, &longorders.CustomerID,
			&longorders.EmployeeID, &longorders.OrderDate, &longorders.RequiredDate,
			&longorders.ShippedDate, &longorders.ShipVia, &longorders.Freight,
			&longorders.ShipName, &longorders.ShipAddress, &longorders.ShipCity,
			&longorders.ShipRegion, &longorders.ShipPostalCode, &longorders.ShipCountry)

		if err != nil {
			panic(err.Error())
		}
		Longorderarray = append(Longorderarray, longorders)
	}

	json.NewEncoder(w).Encode(Longorderarray)
}
func createLongorder(w http.ResponseWriter, r *http.Request) {

	if r.Method == "POST" {

		CustomerID := r.FormValue("CustomerID")
		EmployeeID := r.FormValue("EmployeeID")
		OrderDate := r.FormValue("OrderDate")
		RequiredDate := r.FormValue("RequiredDate")
		ShippedDate := r.FormValue("ShippedDate")
		ShipVia := r.FormValue("ShipVia")
		Freight := r.FormValue("Freight")
		ShipName := r.FormValue("ShipName")
		ShipAddress := r.FormValue("ShipAddress")
		ShipCity := r.FormValue("ShipCity")
		ShipRegion := r.FormValue("ShipRegion")
		ShipPostalCode := r.FormValue("ShipPostalCode")
		ShipCountry := r.FormValue("ShipCountry")

		stmt, err := db.Prepare("INSERT INTO loan VALUES (NULL,?,?,?,?,?,?,?,?,?,?,?,?,?)")

		_, err = stmt.Exec(CustomerID, EmployeeID, OrderDate,
			RequiredDate, ShippedDate, ShipVia, Freight, ShipName,
			ShipAddress, ShipCity, ShipRegion, ShipPostalCode, ShipCountry)

		if err != nil {
			fmt.Fprintf(w, "Data Duplicate")
		} else {
			fmt.Fprintf(w, "Data Created")
		}

	}
}
func getLongorder(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	var Longorderarray []longorder

	sql := `SELECT
				OrderID,
				IFNULL(CustomerID,''),
				IFNULL(EmployeeID,'') EmployeeID,
				IFNULL(OrderDate,'') OrderDate,
				IFNULL(RequiredDate,'') RequiredDate,
				IFNULL(ShippedDate,'') ShippedDate,
				IFNULL(ShipVia,'') ShipVia,
				IFNULL(Freight,'') Freight,
				IFNULL(ShipName,'') ShipName,
				IFNULL(ShipAddress,'') ShipAddress,
				IFNULL(ShipCity,'') ShipCity,
				IFNULL(ShipRegion,'') ShipRegion,
				IFNULL(ShipPostalCode,'') ShipPostalCode,
				IFNULL(ShipCountry,'') ShipCountry

				
			FROM longorders where OrderID = ?`

	result, err := db.Query(sql, params["id"])

	defer result.Close()

	if err != nil {
		panic(err.Error())
	}

	for result.Next() {

		var longorders longorder
		err := result.Scan(&longorders.OrderID, &longorders.CustomerID,
			&longorders.EmployeeID, &longorders.OrderDate, &longorders.RequiredDate,
			&longorders.ShippedDate, &longorders.ShipVia, &longorders.Freight,
			&longorders.ShipName, &longorders.ShipAddress, &longorders.ShipCity,
			&longorders.ShipRegion, &longorders.ShipPostalCode, &longorders.ShipCountry)

		if err != nil {
			panic(err.Error())
		}
		Longorderarray = append(Longorderarray, longorders)
	}

	json.NewEncoder(w).Encode(Longorderarray)
}
func updateLongorder(w http.ResponseWriter, r *http.Request) {

	if r.Method == "PUT" {

		params := mux.Vars(r)

		newCustomerID := r.FormValue("CustomerID")
		newEmployeeID := r.FormValue("EmployeeID")

		stmt, err := db.Prepare("UPDATE longorders SET EmployeeID = ?, EmployeeID =? WHERE CityID = ?")

		_, err = stmt.Exec(newCustomerID, newEmployeeID, params["id"])

		if err != nil {
			fmt.Fprintf(w, "Data not found or Request error")
		}

		fmt.Fprintf(w, "Longorder with ID = %s was updated", params["id"])
	}
}
func deleteLongorder(w http.ResponseWriter, r *http.Request) {

	params := mux.Vars(r)
	stmt, err := db.Prepare("DELETE FROM longorders WHERE OrderID = ?")

	_, err = stmt.Exec(params["id"])

	if err != nil {
		fmt.Fprintf(w, "delete failed")
	}

	fmt.Fprintf(w, "Longorder with ID = %s was deleted", params["id"])
}
func deleteLongorderPost(w http.ResponseWriter, r *http.Request) {
	OrderID := r.FormValue("OrderID")
	stmt, err := db.Prepare("DELETE FROM longorders WHERE OrderID = ?")

	_, err = stmt.Exec(OrderID)

	if err != nil {
		fmt.Fprintf(w, "delete failed")
	}

	fmt.Fprintf(w, "Longorder with ID = "+OrderID+" was deleted")
}
func getPostLongorder(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	OrderID := r.FormValue("OrderID")
	CustomerID := r.FormValue("CustomerID")

	var Longorderarray []longorder

	sql := `SELECT
				OrderID,
				IFNULL(CustomerID,''),
				IFNULL(EmployeeID,'') EmployeeID,
				IFNULL(OrderDate,'') OrderDate,
				IFNULL(RequiredDate,'') RequiredDate,
				IFNULL(ShippedDate,'') ShippedDate,
				IFNULL(ShipVia,'') ShipVia,
				IFNULL(Freight,'') Freight,
				IFNULL(ShipName,'') ShipName,
				IFNULL(ShipAddress,'') ShipAddress,
				IFNULL(ShipCity,'') ShipCity,
				IFNULL(ShipRegion,'') ShipRegion,
				IFNULL(ShipPostalCode,'') ShipPostalCode,
				IFNULL(ShipCountry,'') ShipCountry

				
			FROM longorders where OrderID = ? and CustomerID = ?`

	result, err := db.Query(sql, OrderID, CustomerID)

	defer result.Close()

	if err != nil {
		panic(err.Error())
	}

	for result.Next() {

		var longorders longorder
		err := result.Scan(&longorders.OrderID, &longorders.CustomerID,
			&longorders.EmployeeID, &longorders.OrderDate, &longorders.RequiredDate,
			&longorders.ShippedDate, &longorders.ShipVia, &longorders.Freight,
			&longorders.ShipName, &longorders.ShipAddress, &longorders.ShipCity,
			&longorders.ShipRegion, &longorders.ShipPostalCode, &longorders.ShipCountry)

		if err != nil {
			panic(err.Error())
		}
		Longorderarray = append(Longorderarray, longorders)
	}

	json.NewEncoder(w).Encode(Longorderarray)
}

type mahasiswa struct {
	MahasiswaID string `json:"MahasiswaID"`
	Nama        string `json:"Nama"`
	ProvinceID  string `json:"ProvinceID"`
}

func getMahasiswas(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var Mahasiswaarray []mahasiswa

	sql := `SELECT
			MahasiswaID,
				IFNULL(Nama,''),
				IFNULL(ProvinceID,'') ProvinceID
				
			FROM mahasiswa`

	result, err := db.Query(sql)

	defer result.Close()

	if err != nil {
		panic(err.Error())
	}

	for result.Next() {

		var mahasiswas mahasiswa
		err := result.Scan(&mahasiswas.MahasiswaID, &mahasiswas.Nama,
			&mahasiswas.ProvinceID)

		if err != nil {
			panic(err.Error())
		}
		Mahasiswaarray = append(Mahasiswaarray, mahasiswas)
	}

	json.NewEncoder(w).Encode(Mahasiswaarray)
}
func createMahasiswa(w http.ResponseWriter, r *http.Request) {

	if r.Method == "POST" {

		Nama := r.FormValue("Nama")
		ProvinceID := r.FormValue("ProvinceID")

		stmt, err := db.Prepare("INSERT INTO loan VALUES (NULL,?,?)")

		_, err = stmt.Exec(Nama, ProvinceID)

		if err != nil {
			fmt.Fprintf(w, "Data Duplicate")
		} else {
			fmt.Fprintf(w, "Data Created")
		}

	}
}
func getMahasiswa(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	var Mahasiswaarray []mahasiswa

	sql := `SELECT
			MahasiswaID,
				IFNULL(Nama,''),
				IFNULL(ProvinceID,'') ProvinceID
				
			FROM mahasiswa where MahasiswaID = ?`

	result, err := db.Query(sql, params["id"])

	defer result.Close()

	if err != nil {
		panic(err.Error())
	}

	for result.Next() {

		var mahasiswas mahasiswa
		err := result.Scan(&mahasiswas.MahasiswaID, &mahasiswas.Nama,
			&mahasiswas.ProvinceID)

		if err != nil {
			panic(err.Error())
		}
		Mahasiswaarray = append(Mahasiswaarray, mahasiswas)
	}

	json.NewEncoder(w).Encode(Mahasiswaarray)
}
func updateMahasiswa(w http.ResponseWriter, r *http.Request) {

	if r.Method == "PUT" {

		params := mux.Vars(r)

		newNama := r.FormValue("Nama")
		newProvinceID := r.FormValue("ProvinceID")

		stmt, err := db.Prepare("UPDATE mahasiswa SET Nama = ?, ProvinceID =? WHERE MahasiswaID = ?")

		_, err = stmt.Exec(newNama, newProvinceID, params["id"])

		if err != nil {
			fmt.Fprintf(w, "Data not found or Request error")
		}

		fmt.Fprintf(w, "Mahasiswa with ID = %s was updated", params["id"])
	}
}
func deleteMahasiswa(w http.ResponseWriter, r *http.Request) {

	params := mux.Vars(r)
	stmt, err := db.Prepare("DELETE FROM mahasiswa WHERE MahasiswaID = ?")

	_, err = stmt.Exec(params["id"])

	if err != nil {
		fmt.Fprintf(w, "delete failed")
	}

	fmt.Fprintf(w, "Mahasiswa with ID = %s was deleted", params["id"])
}
func deleteMahasiswaPost(w http.ResponseWriter, r *http.Request) {
	MahasiswaID := r.FormValue("MahasiswaID")
	stmt, err := db.Prepare("DELETE FROM mahasiswa WHERE MahasiswaID = ?")

	_, err = stmt.Exec(MahasiswaID)

	if err != nil {
		fmt.Fprintf(w, "delete failed")
	}

	fmt.Fprintf(w, "Mahasiswa with ID = "+MahasiswaID+" was deleted")
}
func getPostMahasiswa(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	MahasiswaID := r.FormValue("MahasiswaID")
	// CustomerID := r.FormValue("CustomerID")

	var Mahasiswaarray []mahasiswa

	sql := `SELECT
			MahasiswaID,
				IFNULL(Nama,''),
				IFNULL(ProvinceID,'') ProvinceID
				
			FROM mahasiswa where MahasiswaID = ?`

	result, err := db.Query(sql, MahasiswaID)

	defer result.Close()

	if err != nil {
		panic(err.Error())
	}

	for result.Next() {

		var mahasiswas mahasiswa
		err := result.Scan(&mahasiswas.MahasiswaID, &mahasiswas.Nama,
			&mahasiswas.ProvinceID)

		if err != nil {
			panic(err.Error())
		}
		Mahasiswaarray = append(Mahasiswaarray, mahasiswas)
	}

	json.NewEncoder(w).Encode(Mahasiswaarray)
}

type member struct {
	MemberID        string `json:"MemberID"`
	CardID          string `json:"CardID"`
	LastName        string `json:"LastName"`
	FirstName       string `json:"FirstName"`
	Sex             string `json:"Sex"`
	Title           string `json:"Title"`
	TitleOfCourtesy string `json:"TitleOfCourtesy"`
	BirthDate       string `json:"BirthDate"`
	HireDate        string `json:"HireDate"`
	Address         string `json:"Address"`
	City            string `json:"City"`
	Region          string `json:"Region"`
	PostalCode      string `json:"PostalCode"`
	Country         string `json:"Country"`
	HomePhone       string `json:"HomePhone"`
	Extension       string `json:"Extension"`
	Photo           string `json:"Photo"`
	Notes           string `json:"Notes"`
	ReportsTo       string `json:"ReportsTo"`
	AdmisionFee     string `json:"AdmisionFee"`
}

func getMembers(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var Memberarray []member

	sql := `SELECT
				MemberID,
				IFNULL(CardID,''),
				IFNULL(LastName,'') LastName,
				IFNULL(FirstName,'') FirstName,
				IFNULL(Sex,'') Sex,
				IFNULL(Title,'') Title,
				IFNULL(TitleOfCourtesy,'') TitleOfCourtesy,
				IFNULL(BirthDate,'') BirthDate,
				IFNULL(HireDate,'') HireDate,
				IFNULL(Address,'') Address,
				IFNULL(City,'') City,
				IFNULL(Region,'') Region,
				IFNULL(PostalCode,'') PostalCode,
				IFNULL(Country,'') Country,
				IFNULL(HomePhone,'') HomePhone,
				IFNULL(Extension,'') Extension,
				IFNULL(Photo,'') Photo,
				IFNULL(Notes,'') Notes,
				IFNULL(ReportsTo,'') ReportsTo,
				IFNULL(AdmisionFee,'') AdmisionFee
				
			FROM member`

	result, err := db.Query(sql)

	defer result.Close()

	if err != nil {
		panic(err.Error())
	}

	for result.Next() {

		var members member
		err := result.Scan(&members.MemberID, &members.CardID,
			&members.LastName, &members.FirstName, &members.Sex,
			&members.Title, &members.TitleOfCourtesy, &members.BirthDate,
			&members.HireDate, &members.Address, &members.City,
			&members.Region, &members.PostalCode, &members.Country,
			&members.HomePhone, &members.Extension, &members.Photo,
			&members.Notes,
			&members.ReportsTo, &members.AdmisionFee)

		if err != nil {
			panic(err.Error())
		}
		Memberarray = append(Memberarray, members)
	}

	json.NewEncoder(w).Encode(Memberarray)
}
func createMember(w http.ResponseWriter, r *http.Request) {

	if r.Method == "POST" {

		// MemberID := r.FormValue("MemberID")
		CardID := r.FormValue("CardID")
		LastName := r.FormValue("LastName")
		FirstName := r.FormValue("FirstName")
		Sex := r.FormValue("Sex")
		Title := r.FormValue("Title")
		TitleOfCourtesy := r.FormValue("TitleOfCourtesy")
		BirthDate := r.FormValue("BirthDate")
		HireDate := r.FormValue("HireDate")
		Address := r.FormValue("Address")
		City := r.FormValue("City")
		Region := r.FormValue("Region")
		PostalCode := r.FormValue("PostalCode")
		Country := r.FormValue("Country")
		HomePhone := r.FormValue("HomePhone")
		Extension := r.FormValue("Extension")
		Photo := r.FormValue("Photo")
		Notes := r.FormValue("Notes")
		ReportsTo := r.FormValue("ReportsTo")
		AdmisionFee := r.FormValue("AdmisionFee")

		stmt, err := db.Prepare("INSERT INTO member VALUES (NULL,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?)")

		_, err = stmt.Exec(CardID, LastName, FirstName, Sex, Title, TitleOfCourtesy, BirthDate, HireDate, Address, City, Region, PostalCode, Country, HomePhone, Extension, Photo, Notes, ReportsTo, AdmisionFee)

		if err != nil {
			fmt.Fprintf(w, "Data Duplicate")
		} else {
			fmt.Fprintf(w, "Data Created")
		}

	}
}
func getMember(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	var Memberarray []member

	sql := `SELECT
			MemberID,
			IFNULL(CardID,''),
			IFNULL(LastName,'') LastName,
			IFNULL(FirstName,'') FirstName,
			IFNULL(Sex,'') Sex,
			IFNULL(Title,'') Title,
			IFNULL(TitleOfCourtesy,'') TitleOfCourtesy,
			IFNULL(BirthDate,'') BirthDate,
			IFNULL(HireDate,'') HireDate,
			IFNULL(Address,'') Address,
			IFNULL(City,'') City,
			IFNULL(Region,'') Region,
			IFNULL(PostalCode,'') PostalCode,
			IFNULL(Country,'') Country,
			IFNULL(HomePhone,'') HomePhone,
			IFNULL(Extension,'') Extension,
			IFNULL(Photo,'') Photo,
			IFNULL(Notes,'') Notes,
			IFNULL(ReportsTo,'') ReportsTo,
			IFNULL(AdmisionFee,'') AdmisionFee
				
			FROM member where MemberID = ?`

	result, err := db.Query(sql, params["id"])

	defer result.Close()

	if err != nil {
		panic(err.Error())
	}

	for result.Next() {

		var members member
		err := result.Scan(&members.MemberID, &members.CardID,
			&members.LastName, &members.FirstName, &members.Sex,
			&members.Title, &members.TitleOfCourtesy, &members.BirthDate,
			&members.HireDate, &members.Address, &members.City,
			&members.Region, &members.PostalCode, &members.Country,
			&members.HomePhone, &members.Extension, &members.Photo,
			&members.Notes,
			&members.ReportsTo, &members.AdmisionFee)

		if err != nil {
			panic(err.Error())
		}
		Memberarray = append(Memberarray, members)
	}

	json.NewEncoder(w).Encode(Memberarray)
}
func updateMember(w http.ResponseWriter, r *http.Request) {

	if r.Method == "PUT" {

		params := mux.Vars(r)

		newCardID := r.FormValue("CardID")

		stmt, err := db.Prepare("UPDATE member SET CardID = ? WHERE MemberID = ?")

		_, err = stmt.Exec(newCardID, params["id"])

		if err != nil {
			fmt.Fprintf(w, "Data not found or Request error")
		}

		fmt.Fprintf(w, "Member with ID = %s was updated", params["id"])
	}
}
func deleteMember(w http.ResponseWriter, r *http.Request) {

	params := mux.Vars(r)
	stmt, err := db.Prepare("DELETE FROM member WHERE MemberID = ?")

	_, err = stmt.Exec(params["id"])

	if err != nil {
		fmt.Fprintf(w, "delete failed")
	}

	fmt.Fprintf(w, "Member with ID = %s was deleted", params["id"])
}
func deleteMemberPost(w http.ResponseWriter, r *http.Request) {
	MemberID := r.FormValue("MemberID")
	stmt, err := db.Prepare("DELETE FROM member WHERE MemberID = ?")

	_, err = stmt.Exec(MemberID)

	if err != nil {
		fmt.Fprintf(w, "delete failed")
	}

	fmt.Fprintf(w, "Member with ID = "+MemberID+" was deleted")
}
func getPostMember(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	MemberID := r.FormValue("MemberID")
	// CustomerID := r.FormValue("CustomerID")
	var Memberarray []member

	sql := `SELECT
			MemberID,
			IFNULL(CardID,''),
			IFNULL(LastName,'') LastName,
			IFNULL(FirstName,'') FirstName,
			IFNULL(Sex,'') Sex,
			IFNULL(Title,'') Title,
			IFNULL(TitleOfCourtesy,'') TitleOfCourtesy,
			IFNULL(BirthDate,'') BirthDate,
			IFNULL(HireDate,'') HireDate,
			IFNULL(Address,'') Address,
			IFNULL(City,'') City,
			IFNULL(Region,'') Region,
			IFNULL(PostalCode,'') PostalCode,
			IFNULL(Country,'') Country,
			IFNULL(HomePhone,'') HomePhone,
			IFNULL(Extension,'') Extension,
			IFNULL(Photo,'') Photo,
			IFNULL(Notes,'') Notes,
			IFNULL(ReportsTo,'') ReportsTo,
			IFNULL(AdmisionFee,'') AdmisionFee
				
			FROM member where MemberID = ?`

	result, err := db.Query(sql, MemberID)

	defer result.Close()

	if err != nil {
		panic(err.Error())
	}

	for result.Next() {

		var members member
		err := result.Scan(&members.MemberID, &members.CardID,
			&members.LastName, &members.FirstName, &members.Sex,
			&members.Title, &members.TitleOfCourtesy, &members.BirthDate,
			&members.HireDate, &members.Address, &members.City,
			&members.Region, &members.PostalCode, &members.Country,
			&members.HomePhone, &members.Extension, &members.Photo,
			&members.Notes,
			&members.ReportsTo, &members.AdmisionFee)

		if err != nil {
			panic(err.Error())
		}
		Memberarray = append(Memberarray, members)
	}

	json.NewEncoder(w).Encode(Memberarray)
}

//TABEL BARU
type tabelbaru struct {
	Tabelid     string `json:"id"`
	Namalengkap string `json:"nama_lengkap"`
	Alamat      string `json:"alamat"`
	Tgllahir    string `json:"tgl_lahir"`
	Gender      string `json:"gender"`
	Hobi        string `json:"hobi"`
	Pddtrakhir  string `json:"pdd_trakhir"`
}

func gettbs(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var tbarray []tabelbaru

	sql := `SELECT
				id,
				IFNULL(nama_lengkap,''),
				IFNULL(alamat,'') alamat,
				IFNULL(tgl_lahir,'') tgl_lahir,
				IFNULL(gender,'') gender,
				IFNULL(hobi,'') hobi,
				IFNULL(pdd_trakhir,'') pdd_trakhir
				
			FROM tabelbaru`

	result, err := db.Query(sql)

	defer result.Close()

	if err != nil {
		panic(err.Error())
	}

	for result.Next() {

		var tbs tabelbaru
		err := result.Scan(&tbs.Tabelid, &tbs.Namalengkap,
			&tbs.Alamat, &tbs.Tgllahir, &tbs.Gender,
			&tbs.Hobi, &tbs.Pddtrakhir)

		if err != nil {
			panic(err.Error())
		}
		tbarray = append(tbarray, tbs)
	}

	json.NewEncoder(w).Encode(tbarray)
}
func createtb(w http.ResponseWriter, r *http.Request) {

	if r.Method == "POST" {

		// MemberID := r.FormValue("MemberID")
		namalengkap := r.FormValue("nama_lengkap")
		alamat := r.FormValue("alamat")
		tgllahir := r.FormValue("tgl_lahir")
		gender := r.FormValue("gender")
		hobi := r.FormValue("hobi")
		pddtrakhir := r.FormValue("pdd_trakhir")

		stmt, err := db.Prepare("INSERT INTO tabelbaru VALUES (NULL,?,?,?,?,?,?)")

		_, err = stmt.Exec(namalengkap, alamat, tgllahir, gender, hobi, pddtrakhir)

		if err != nil {
			fmt.Fprintf(w, "Data Duplicate")
		} else {
			fmt.Fprintf(w, "Data Created")
		}

	}
}
func gettb(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	var tbarray []tabelbaru

	sql := `SELECT
				id,
				IFNULL(nama_lengkap,''),
				IFNULL(alamat,'') alamat,
				IFNULL(tgl_lahir,'') tgl_lahir,
				IFNULL(gender,'') gender,
				IFNULL(hobi,'') hobi,
				IFNULL(pdd_trakhir,'') pdd_trakhir
				
			FROM tabelbaru where id = ?`

	result, err := db.Query(sql, params["id"])

	defer result.Close()

	if err != nil {
		panic(err.Error())
	}

	for result.Next() {

		var tbs tabelbaru
		err := result.Scan(&tbs.Tabelid, &tbs.Namalengkap,
			&tbs.Alamat, &tbs.Tgllahir, &tbs.Gender,
			&tbs.Hobi, &tbs.Pddtrakhir)

		if err != nil {
			panic(err.Error())
		}
		tbarray = append(tbarray, tbs)
	}

	json.NewEncoder(w).Encode(tbarray)
}
func updatetb(w http.ResponseWriter, r *http.Request) {

	if r.Method == "PUT" {

		params := mux.Vars(r)

		newCardID := r.FormValue("CardID")

		stmt, err := db.Prepare("UPDATE tabelbaru SET nama_lengkap = ? WHERE id = ?")

		_, err = stmt.Exec(newCardID, params["id"])

		if err != nil {
			fmt.Fprintf(w, "Data not found or Request error")
		}

		fmt.Fprintf(w, "Data Tabrlbaru with ID = %s was updated", params["id"])
	}
}
func deletetb(w http.ResponseWriter, r *http.Request) {

	params := mux.Vars(r)
	stmt, err := db.Prepare("DELETE FROM tabelbaru WHERE id = ?")

	_, err = stmt.Exec(params["id"])

	if err != nil {
		fmt.Fprintf(w, "delete failed")
	}

	fmt.Fprintf(w, "Data Tabelbaru with ID = %s was deleted", params["id"])
}
func deleteTBPost(w http.ResponseWriter, r *http.Request) {
	id := r.FormValue("id")
	stmt, err := db.Prepare("DELETE FROM tabelbaru WHERE id = ?")

	_, err = stmt.Exec(id)

	if err != nil {
		fmt.Fprintf(w, "delete failed")
	}

	fmt.Fprintf(w, "Data Tabelbaru with ID = "+id+" was deleted")
}
func getPosttb(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	id := r.FormValue("id")
	var tbarray []tabelbaru

	sql := `SELECT
				id,
				IFNULL(nama_lengkap,''),
				IFNULL(alamat,'') alamat,
				IFNULL(tgl_lahir,'') tgl_lahir,
				IFNULL(gender,'') gender,
				IFNULL(hobi,'') hobi,
				IFNULL(pdd_trakhir,'') pdd_trakhir
				
			FROM tabelbaru where id = ?`

	result, err := db.Query(sql, id)

	defer result.Close()

	if err != nil {
		panic(err.Error())
	}

	for result.Next() {

		var tbs tabelbaru
		err := result.Scan(&tbs.Tabelid, &tbs.Namalengkap,
			&tbs.Alamat, &tbs.Tgllahir, &tbs.Gender,
			&tbs.Hobi, &tbs.Pddtrakhir)
		if err != nil {
			panic(err.Error())
		}
		tbarray = append(tbarray, tbs)
	}

	json.NewEncoder(w).Encode(tbarray)
}

// Main function
func main() {

	db, err = sql.Open("mysql", "root:@tcp(127.0.0.1:3306)/northwind")
	if err != nil {
		panic(err.Error())
	}

	defer db.Close()

	// Init router
	r := mux.NewRouter()

	// Route handles & endpoints

	r.HandleFunc("/customers", getCustomers).Methods("GET")
	r.HandleFunc("/customers/{id}", getCustomer).Methods("GET")
	r.HandleFunc("/customers", createCustomer).Methods("POST")
	r.HandleFunc("/customers/{id}", updateCustomer).Methods("PUT")
	r.HandleFunc("/customers/{id}", deleteCustomer).Methods("DELETE")

	//New
	r.HandleFunc("/getcustomer", getPost).Methods("POST")

	//account
	r.HandleFunc("/accounts", getAccounts).Methods("GET")
	r.HandleFunc("/accounts", createAccount).Methods("POST")
	r.HandleFunc("/accounts/{id}", getAccount).Methods("GET")
	r.HandleFunc("/accounts/{id}", updateAccount).Methods("PUT")
	r.HandleFunc("/accounts/{id}", deleteAcount).Methods("DELETE")

	r.HandleFunc("/delaccounts", deleteAcountPost).Methods("POST")

	r.HandleFunc("/getaccount", getPostAccount).Methods("POST")

	//adjtable
	r.HandleFunc("/adjtable", getAdjtables).Methods("GET")
	r.HandleFunc("/adjtable", createAdjtable).Methods("POST")
	r.HandleFunc("/adjtable/{id}", getAdjTable).Methods("GET")
	r.HandleFunc("/adjtable/{id}", updateAdjtable).Methods("PUT")
	r.HandleFunc("/adjtable/{id}", deleteAdjtable).Methods("DELETE")

	r.HandleFunc("/deladjtable", deleteAdjTablePost).Methods("POST")

	r.HandleFunc("/getadjtable", getPostAdj).Methods("POST")

	//Book
	r.HandleFunc("/book", getBooks).Methods("GET")
	r.HandleFunc("/book", createBook).Methods("POST")
	r.HandleFunc("/book/{id}", getBook).Methods("GET")
	r.HandleFunc("/book/{id}", updateBook).Methods("PUT")
	r.HandleFunc("/book/{id}", deleteBook).Methods("DELETE")

	r.HandleFunc("/delbook", deleteBookPost).Methods("POST")

	r.HandleFunc("/getbook", getPostBook).Methods("POST")

	//City
	r.HandleFunc("/city", getCitys).Methods("GET")
	r.HandleFunc("/city", createCity).Methods("POST")
	r.HandleFunc("/city/{id}", getCity).Methods("GET")
	r.HandleFunc("/city/{id}", updateCity).Methods("PUT")
	r.HandleFunc("/city/{id}", deleteCity).Methods("DELETE")

	r.HandleFunc("/delcity", deleteCityPost).Methods("POST")

	r.HandleFunc("/getcity", getPostCity).Methods("POST")

	//Countries
	r.HandleFunc("/country", getCountries).Methods("GET")
	r.HandleFunc("/country", createCountrie).Methods("POST")
	r.HandleFunc("/country/{id}", getCountrie).Methods("GET")
	r.HandleFunc("/country/{id}", updateCountrie).Methods("PUT")
	r.HandleFunc("/country/{id}", deleteCountrie).Methods("DELETE")

	r.HandleFunc("/delcountry", deleteCountriePost).Methods("POST")

	r.HandleFunc("/getcountry", getPostContrie).Methods("POST")

	//Employee
	r.HandleFunc("/employee", getEmployees).Methods("GET")
	r.HandleFunc("/employee", createEmployee).Methods("POST")
	r.HandleFunc("/employee/{id}", getEmployee).Methods("GET")
	r.HandleFunc("/employee/{id}", updateEmployee).Methods("PUT")
	r.HandleFunc("/employee/{id}", deleteEmployee).Methods("DELETE")

	r.HandleFunc("/delemployee", deleteEmployeePost).Methods("POST")

	r.HandleFunc("/getemployee", getPostEmployee).Methods("POST")

	//load
	r.HandleFunc("/load", getLoads).Methods("GET")
	r.HandleFunc("/load", createload).Methods("POST")
	r.HandleFunc("/load/{id}", getLoad).Methods("GET")
	r.HandleFunc("/load/{id}", updateLoad).Methods("PUT")
	r.HandleFunc("/load/{id}", deleteLoad).Methods("DELETE")

	r.HandleFunc("/delload", deleteLoadPost).Methods("POST")

	r.HandleFunc("/getload", getPostLoad).Methods("POST")

	//loan
	r.HandleFunc("/loan", getLoans).Methods("GET")
	r.HandleFunc("/loan", createloan).Methods("POST")
	r.HandleFunc("/loan/{id}", getLoan).Methods("GET")
	r.HandleFunc("/loan/{id}", updateLoan).Methods("PUT")
	r.HandleFunc("/loan/{id}", deleteLoan).Methods("DELETE")

	r.HandleFunc("/delloan", deleteLoanPost).Methods("POST")

	r.HandleFunc("/getloan", getPostLoan).Methods("POST")

	//longorders
	r.HandleFunc("/longorder", getLongorders).Methods("GET")
	r.HandleFunc("/longorder", createLongorder).Methods("POST")
	r.HandleFunc("/longorder/{id}", getLongorder).Methods("GET")
	r.HandleFunc("/longorder/{id}", updateLongorder).Methods("PUT")
	r.HandleFunc("/longorder/{id}", deleteLongorder).Methods("DELETE")

	r.HandleFunc("/dellongorder", deleteLongorderPost).Methods("POST")

	r.HandleFunc("/getlongorder", getPostLongorder).Methods("POST")

	//mahasiswa
	r.HandleFunc("/mahasiswa", getMahasiswas).Methods("GET")
	r.HandleFunc("/mahasiswa", createMahasiswa).Methods("POST")
	r.HandleFunc("/mahasiswa/{id}", getMahasiswa).Methods("GET")
	r.HandleFunc("/mahasiswa/{id}", updateMahasiswa).Methods("PUT")
	r.HandleFunc("/mahasiswa/{id}", deleteMahasiswa).Methods("DELETE")

	r.HandleFunc("/delmahasiswa", deleteMahasiswaPost).Methods("POST")

	r.HandleFunc("/getmahasiswa", getPostMahasiswa).Methods("POST")

	//member
	r.HandleFunc("/member", getMembers).Methods("GET")
	r.HandleFunc("/member", createMember).Methods("POST")
	r.HandleFunc("/member/{id}", getMember).Methods("GET")
	r.HandleFunc("/member/{id}", updateMember).Methods("PUT")
	r.HandleFunc("/member/{id}", deleteMember).Methods("DELETE")

	r.HandleFunc("/delmember", deleteMemberPost).Methods("POST")

	r.HandleFunc("/getmember", getPostMember).Methods("POST")

	//Tabel Baru
	r.HandleFunc("/tabelbaru", gettbs).Methods("GET")
	r.HandleFunc("/tabelbaru", createtb).Methods("POST")
	r.HandleFunc("/tabelbaru/{id}", gettb).Methods("GET")
	r.HandleFunc("/tabelbaru/{id}", updatetb).Methods("PUT")
	r.HandleFunc("/tabelbaru/{id}", deletetb).Methods("DELETE")

	r.HandleFunc("/deltabelbaru", deleteTBPost).Methods("POST")

	r.HandleFunc("/gettabelbaru", getPosttb).Methods("POST")

	// Start server
	log.Fatal(http.ListenAndServe(":8080", r))
}
