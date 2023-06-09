package main

// import (
// 	"database/sql"
// 	"fmt"
// 	"time"
// )

// type Userdb struct {
// 	DB *sql.DB
// }

// type User struct {
// 	Id        int
// 	Email     string
// 	Pseudo    string
// 	Password  string
// 	Age       time.Time
// 	FirstName string
// 	LastName  string
// 	Picture   string
// }

// func CreatUserTable(db *sql.DB) *Userdb {
// 	stmt, _ := db.Prepare(`
// 		CREATE TABLE IF NOT EXISTS "users" (
// 			"user_id" INTEGER PRIMARY KEY AUTOINCREMENT,
// 			"email"	TEXT NOT NULL UNIQUE COLLATE NOCASE,
// 			"pseudo"	TEXT NOT NULL UNIQUE COLLATE NOCASE,
// 			"password"	TEXT NOT NULL,
// 			"age" DATE NOT NULL,
// 			"first_name" TEXT NOT NULL COLLATE NOCASE,
// 			"last_name" TEXT NOT NULL COLLATE NOCASE,
// 			"picture" TEXT NOT NULL
// 		);
// 	`)
// 	stmt.Exec()

// 	return &Userdb{
// 		DB: db,
// 	}
// }

// func (db *Userdb) addUserr(user User) error {
// 	stmt, er := db.DB.Prepare(`
// 	INSERT INTO "users" (email, pseudo, password, age, first_name, last_name, picture) values (?, ?, ?, ?, ?, ?, ?)
// 	`)
// 	if er != nil {
// 		return er
// 	}
// 	_, err := stmt.Exec(user.Email, user.Pseudo, user.Password, user.Age.Format("2006-01-02"), user.FirstName, user.LastName, user.Picture)
// 	if err != nil {
// 		return err
// 	}
// 	return nil
// }

// func (db *Userdb) UpdatePic(user User) error {
// 	stmt, er := db.DB.Prepare(`
// 	UPDATE "users" SET "picture" = ? WHERE "pseudo" = ?
// 	`)
// 	if er != nil {
// 		return er
// 	}
// 	_, err := stmt.Exec(user.Picture, user.Pseudo)
// 	if err != nil {
// 		return err
// 	}
// 	return nil
// }

// func (db *Userdb) GetAllUsers() []User {
// 	Users := []User{}
// 	rows, _ := db.DB.Query(`
// 	SELECT * FROM "users"
// 	`)
// 	var id int
// 	var email string
// 	var pseudo string
// 	var password string
// 	var date string
// 	var first_name string
// 	var last_name string
// 	var picture string
// 	for rows.Next() {
// 		rows.Scan(&id, &email, &pseudo, &password, &date, &first_name, &last_name, &picture)
// 		age, _ := time.Parse("2006-01-02", date)
// 		User := User{
// 			Id:        id,
// 			Email:     email,
// 			Pseudo:    pseudo,
// 			Password:  password,
// 			Age:       age,
// 			FirstName: first_name,
// 			LastName:  last_name,
// 			Picture:   picture,
// 		}
// 		Users = append(Users, User)
// 	}
// 	rows.Close()
// 	return Users
// }

// func (db *Userdb) GetUserByPseudo(str string) User {
// 	s := fmt.Sprintf("SELECT * FROM users WHERE username = '%v'", str)
// 	return db.getuserbyquery(s)
// }

// func (db *Userdb) GetUserByEmail(str string) User {
// 	s := fmt.Sprintf("SELECT * FROM users WHERE email = '%v'", str)
// 	return db.getuserbyquery(s)
// }

// func (db *Userdb) GetUserById(str int) User {
// 	s := fmt.Sprintf("SELECT * FROM users WHERE id = %d", str)
// 	return db.getuserbyquery(s)
// }

// func (db *Userdb) GetUsersByFirst_name(str string) []User {
// 	s := fmt.Sprintf("SELECT * FROM users WHERE first_name = '%v'", str)
// 	Users := []User{}
// 	rows, _ := db.DB.Query(s)
// 	var id int
// 	var email string
// 	var pseudo string
// 	var password string
// 	var date string
// 	var first_name string
// 	var last_name string
// 	var picture string
// 	for rows.Next() {
// 		rows.Scan(&id, &email, &pseudo, &password, &date, &first_name, &last_name, &picture)
// 		age, _ := time.Parse("2006-01-02", date)
// 		User := User{
// 			Id:        id,
// 			Email:     email,
// 			Pseudo:    pseudo,
// 			Password:  password,
// 			Age:       age,
// 			FirstName: first_name,
// 			LastName:  last_name,
// 			Picture:   picture,
// 		}
// 		Users = append(Users, User)
// 	}
// 	rows.Close()
// 	return Users
// }

// func (db *Userdb) GetUsersByLast_name(str string) []User {
// 	s := fmt.Sprintf("SELECT * FROM users WHERE last_name = '%v'", str)
// 	Users := []User{}
// 	rows, _ := db.DB.Query(s)
// 	var id int
// 	var email string
// 	var pseudo string
// 	var password string
// 	var date string
// 	var first_name string
// 	var last_name string
// 	var picture string
// 	for rows.Next() {
// 		rows.Scan(&id, &email, &pseudo, &password, &date, &first_name, &last_name, &picture)
// 		age, _ := time.Parse("2006-01-02", date)
// 		User := User{
// 			Id:        id,
// 			Email:     email,
// 			Pseudo:    pseudo,
// 			Password:  password,
// 			Age:       age,
// 			FirstName: first_name,
// 			LastName:  last_name,
// 			Picture:   picture,
// 		}
// 		Users = append(Users, User)
// 	}
// 	rows.Close()
// 	return Users
// }

// func (db *Userdb) getuserbyquery(s string) User {
// 	rows, _ := db.DB.Query(s)
// 	var id int
// 	var email string
// 	var pseudo string
// 	var password string
// 	var date string
// 	var first_name string
// 	var last_name string
// 	var picture string
// 	var user User
// 	if rows.Next() {
// 		rows.Scan(&id, &email, &pseudo, &password, &date, &first_name, &last_name, &picture)
// 		age, _ := time.Parse("2006-01-02", date)
// 		user = User{
// 			Id:        id,
// 			Email:     email,
// 			Pseudo:    pseudo,
// 			Password:  password,
// 			Age:       age,
// 			FirstName: first_name,
// 			LastName:  last_name,
// 			Picture:   picture,
// 		}
// 	}
// 	rows.Close()
// 	return user
// }
