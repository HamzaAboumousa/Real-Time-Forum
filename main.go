package main

import (
	"database/sql"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"net/url"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"
	"unicode"

	uuid "github.com/gofrs/uuid"
	_ "github.com/mattn/go-sqlite3"
	"golang.org/x/crypto/bcrypt"
)

var (
	Port      = ":8000"
	Msgdb     *Messagedb
	bunch     *Post
	guys      *User
	flood     *Comm
	templates *template.Template
	LD        *Likes
	LDcom     *Likescom
)

type Citem struct {
	CommentId   string
	PostId      string
	Author      string
	Content     string
	Creat_at    string
	ComIsAuthor bool
	L           int
	D           int
}

type Litem struct {
	PostID   string
	Username string
	Like     string
}

type Litemcom struct {
	ComID    string
	Username string
	Like     string
}

type Item struct {
	ID       string
	Author   string
	Content  string
	Thread   string
	Creat_at string
	L        int
	D        int
	IsAuthor bool
}

type Uitem struct {
	Email    string
	Username string
	Password string
	Age      int
	Fname    string
	Lname    string
	Sex      string
}

type User struct {
	DB *sql.DB
}

type Likescom struct {
	DB *sql.DB
}

type Likes struct {
	DB *sql.DB
}

type Comm struct {
	DB *sql.DB
}
type Post struct {
	DB *sql.DB
}

func (comm *Comm) Delete(id string) {
	stmt, _ := comm.DB.Prepare(`DELETE FROM "comments" WHERE "commentid" = ?`)
	stmt.Exec(id)
}

func (comm *Comm) Get(LD *Likescom, str string) []Citem {
	s := fmt.Sprintf("SELECT * FROM comments WHERE postid = '%v'", str)

	items := []Citem{}
	rows, _ := comm.DB.Query(s)
	var commentid string
	var postid string
	var author string
	var content string
	var creat_at string
	for rows.Next() {
		rows.Scan(&commentid, &postid, &author, &content, &creat_at)
		item := Citem{
			CommentId: commentid,
			PostId:    postid,
			Author:    author,
			Content:   content,
			Creat_at:  creat_at,
			L:         len(LD.Get(commentid, "l")),
			D:         len(LD.Get(commentid, "d")),
		}
		items = append(items, item)
	}
	rows.Close()
	return items
}

func (comm *Comm) Add(citem Citem) {
	stmt, _ := comm.DB.Prepare(`INSERT INTO "comments" (commentid, postid, author, content, creat_at) values(?, ?, ?, ?, ?)`)
	_, err := stmt.Exec(citem.CommentId, citem.PostId, citem.Author, citem.Content, citem.Creat_at)
	if err != nil {
		fmt.Println(err)
	}
}

func NewComm(db *sql.DB) *Comm {
	stmt, err := db.Prepare(`
		CREATE TABLE IF NOT EXISTS "comments" (
			"commentid" TEXT NOT NULL UNIQUE,
			"postid"	TEXT NOT NULL,
			"author"	TEXT NOT NULL COLLATE NOCASE,
			"content"	TEXT NOT NULL,
			"creat_at"  TEXT NOT NULL
		);
	`)
	if err != nil {
		println(err)
	}
	stmt.Exec()
	return &Comm{
		DB: db,
	}
}

func (likes *Likes) GetOne(id, user string) Litem {
	item := Litem{}

	s := fmt.Sprintf("SELECT * FROM likes WHERE postid = '%v' AND username = '%v'", id, user)

	rows, _ := likes.DB.Query(s)
	var postid string
	var author string
	var like string
	if rows.Next() {
		rows.Scan(&postid, &author, &like)
		item = Litem{
			PostID:   postid,
			Username: author,
			Like:     like,
		}
	}
	rows.Close()
	return item
}

func (likes *Likes) Add(litem Litem) {
	one := likes.GetOne(litem.PostID, litem.Username)
	var s string
	if one.Like == "" {
		s = "INSERT INTO likes (like, postid, username) values (?, ?, ?)"
	} else if litem.Like != one.Like {
		s = "UPDATE likes SET like = ? WHERE postid = ? AND username = ?"
	} else {
		s = "DELETE FROM likes WHERE like = ? AND postid = ? AND username = ?"
	}
	stmt, _ := likes.DB.Prepare(s)
	_, err := stmt.Exec(litem.Like, litem.PostID, litem.Username)
	if err != nil {
		fmt.Println(err)
	}
}

func (likes *Likes) Get(id, l string) []Litem {
	items := []Litem{}
	var s string
	if l == "all" {
		s = fmt.Sprintf("SELECT * FROM likes WHERE username = '%v'", id)
	} else {
		s = fmt.Sprintf("SELECT * FROM likes WHERE postid = '%v' AND like = '%v'", id, l)
	}

	rows, _ := likes.DB.Query(s)
	var postid string
	var author string
	var like string
	for rows.Next() {
		rows.Scan(&postid, &author, &like)
		item := Litem{
			PostID:   postid,
			Username: author,
			Like:     like,
		}
		items = append(items, item)
	}
	rows.Close()
	return items
}

func NewLD(db *sql.DB) *Likes {
	stmt, _ := db.Prepare(`
		CREATE TABLE IF NOT EXISTS "likes" (
			"postid"	TEXT NOT NULL,
			"username"	TEXT NOT NULL COLLATE NOCASE,
			"like"	TEXT
		);
	`)
	stmt.Exec()

	return &Likes{
		DB: db,
	}
}

func (likes *Likescom) GetOne(id, user string) Litemcom {
	item := Litemcom{}

	s := fmt.Sprintf("SELECT * FROM likescom WHERE comid = '%v' AND username = '%v'", id, user)

	rows, _ := likes.DB.Query(s)
	var comid string
	var author string
	var like string
	if rows.Next() {
		rows.Scan(&comid, &author, &like)
		item = Litemcom{
			ComID:    comid,
			Username: author,
			Like:     like,
		}
	}
	rows.Close()
	return item
}

func (likes *Likescom) Add(litemcom Litemcom) {
	one := likes.GetOne(litemcom.ComID, litemcom.Username)
	var s string
	if one.Like == "" {
		s = "INSERT INTO likescom (like, comid, username) values (?, ?, ?)"
	} else if litemcom.Like != one.Like {
		s = "UPDATE likescom SET like = ? WHERE comid = ? AND username = ?"
	} else {
		s = "DELETE FROM likescom WHERE like = ? AND comid = ? AND username = ?"
	}
	stmt, _ := likes.DB.Prepare(s)
	_, err := stmt.Exec(litemcom.Like, litemcom.ComID, litemcom.Username)
	if err != nil {
		fmt.Println(err)
	}
}

func (likes *Likescom) Get(id, l string) []Litemcom {
	items := []Litemcom{}
	var s string
	if l == "all" {
		s = fmt.Sprintf("SELECT * FROM likescom WHERE username = '%v'", id)
	} else {
		s = fmt.Sprintf("SELECT * FROM likescom WHERE comid = '%v' AND like = '%v'", id, l)
	}

	rows, _ := likes.DB.Query(s)
	var postid string
	var author string
	var like string
	for rows.Next() {
		rows.Scan(&postid, &author, &like)
		item := Litemcom{
			ComID:    postid,
			Username: author,
			Like:     like,
		}
		items = append(items, item)
	}
	rows.Close()
	return items
}

func NewLDcom(db *sql.DB) *Likescom {
	stmt, _ := db.Prepare(`
		CREATE TABLE IF NOT EXISTS "likescom" (
			"comid"	TEXT NOT NULL,
			"username"	TEXT NOT NULL COLLATE NOCASE,
			"like"	TEXT
		);
	`)
	stmt.Exec()

	return &Likescom{
		DB: db,
	}
}

func (post *Post) GetMyPosts(LD *Likes, str string) ([]Item, []Item) {
	s := fmt.Sprintf("SELECT * FROM posts WHERE author = '%v'", str)

	myitems := []Item{}
	mylikes := []Item{}
	likes := LD.Get(str, "all")
	rows, _ := post.DB.Query(s)
	var id string
	var author string
	var content string
	var thread string
	var creat_at string
	for rows.Next() {
		rows.Scan(&id, &author, &content, &thread, &creat_at)
		item := Item{
			ID:       id,
			Author:   author,
			Content:  content,
			Thread:   thread,
			Creat_at: creat_at,
			L:        len(LD.Get(id, "l")),
			D:        len(LD.Get(id, "d")),
		}
		myitems = append(myitems, item)
	}
	rows.Close()

	for _, v := range likes {
		s := fmt.Sprintf("SELECT * FROM posts WHERE id = '%v'", v.PostID)

		rows, _ := post.DB.Query(s)
		var id string
		var author string
		var content string
		var thread string
		var creat_at string
		var item Item
		if rows.Next() {
			rows.Scan(&id, &author, &content, &thread, &creat_at)
			item = Item{
				ID:       id,
				Author:   author,
				Content:  content,
				Thread:   thread,
				Creat_at: creat_at,
				L:        len(LD.Get(id, "l")),
				D:        len(LD.Get(id, "d")),
			}
			mylikes = append(mylikes, item)

		}

		rows.Close()
	}
	return myitems, mylikes
}

func (post *Post) Filter(LD *Likes, str string) []Item {
	s := fmt.Sprintf("SELECT * FROM posts WHERE thread LIKE '%v'", "%"+str+"%")

	items := []Item{}
	rows, _ := post.DB.Query(s)
	var id string
	var author string
	var content string
	var thread string
	var creat_at string
	for rows.Next() {
		rows.Scan(&id, &author, &content, &thread, &creat_at)
		item := Item{
			ID:       id,
			Author:   author,
			Content:  content,
			Thread:   thread,
			Creat_at: creat_at,
			L:        len(LD.Get(id, "l")),
			D:        len(LD.Get(id, "d")),
		}
		items = append(items, item)
	}
	rows.Close()
	return items
}

func (post *Post) Get(LD *Likes) []Item {
	items := []Item{}
	rows, _ := post.DB.Query(`
	SELECT * FROM "posts"
	`)
	var id string
	var author string
	var content string
	var thread string
	var creat_at string
	for rows.Next() {
		rows.Scan(&id, &author, &content, &thread, &creat_at)
		item := Item{
			ID:       id,
			Author:   author,
			Content:  content,
			Thread:   thread,
			Creat_at: creat_at,
			L:        len(LD.Get(id, "l")),
			D:        len(LD.Get(id, "d")),
		}
		items = append(items, item)
	}
	rows.Close()
	return items
}

func (post *Post) Add(item Item) {
	stmt, _ := post.DB.Prepare(`INSERT INTO "posts"(id, author, content, thread, creat_at) values(?, ?, ?, ?, ?)`)
	_, err := stmt.Exec(item.ID, item.Author, item.Content, item.Thread, item.Creat_at)
	if err != nil {
		fmt.Println(err)
	}
}

func NewPost(db *sql.DB) *Post {
	stmt, _ := db.Prepare(`
		CREATE TABLE IF NOT EXISTS "posts" (
			"id"	TEXT NOT NULL UNIQUE,
			"author"	TEXT NOT NULL COLLATE NOCASE,
			"content"	TEXT NOT NULL,
			"thread"	TEXT,
			"creat_at"  TEXT NOT NULL,
			PRIMARY KEY("id")
		);
	`)
	stmt.Exec()

	return &Post{
		DB: db,
	}
}

func NewMessage(db *sql.DB) *Messagedb {
	stmt, err := db.Prepare(`
		CREATE TABLE IF NOT EXISTS "messages" (
			"id"	INTEGER PRIMARY KEY AUTOINCREMENT,
			"from"	    TEXT NOT NULL COLLATE NOCASE,
			"to"	    TEXT NOT NULL COLLATE NOCASE,
			"content"	TEXT NOT NULL,
			"send_at"   Text NOT NULL
		);
	`)
	if err != nil {
		fmt.Println(err)
	}
	stmt.Exec()

	return &Messagedb{
		DB: db,
	}
}

func (message *Messagedb) Get(user string) []Message {
	items := []Message{}
	s := fmt.Sprintf("SELECT * FROM messages WHERE `to` LIKE '%v' OR `from` LIKE '%v'", user, user)
	rows, _ := message.DB.Query(s)
	var from string
	var to string
	var msg string
	var send_at string
	for rows.Next() {
		var id int
		rows.Scan(&id, &from, &to, &msg, &send_at)
		date := send_at
		item := Message{
			From:    from,
			To:      to,
			Content: msg,
			Date:    date,
		}
		items = append(items, item)
	}
	rows.Close()
	return items
}

func (message *Messagedb) Add(item Message) {
	stmt, er := message.DB.Prepare(`INSERT INTO messages ("from", "to", "content", send_at) values(?, ?, ?, ?);`)
	if er != nil {
		fmt.Println(er)
	}
	_, err := stmt.Exec(item.From, item.To, item.Content, item.Date)
	if err != nil {
		fmt.Println(err)
	}
}

func (post *Post) Update(item Item, id string) {
	stmt, _ := post.DB.Prepare(`UPDATE "posts" SET "content" = ?, "thread" = ? WHERE "id" = ?`)
	stmt.Exec(item.Content, item.Thread, id)
}

func (post *Post) Delete(id string) {
	stmt, _ := post.DB.Prepare(`DELETE FROM "posts" WHERE "id" = ?`)
	stmt.Exec(id)
}

func (user *User) Add(uitem Uitem) error {
	stmt, _ := user.DB.Prepare(`
	INSERT INTO "users" (email, username, password, age, fname, lastname, sex) values (?, ?, ?, ?, ?, ?, ?)
	`)
	_, err := stmt.Exec(uitem.Email, uitem.Username, uitem.Password, uitem.Age, uitem.Fname, uitem.Lname, uitem.Sex)
	if err != nil {
		return err
	}
	return nil
}

func NewUser(db *sql.DB) *User {
	stmt, err := db.Prepare(`
		CREATE TABLE IF NOT EXISTS "users" (
			"email"	TEXT UNIQUE COLLATE NOCASE,
			"username"	TEXT NOT NULL UNIQUE COLLATE NOCASE,
			"password"	TEXT NOT NULL,
			"age" INTEGER NOT NULL,
			"fname" TEXT NOT NULL COLLATE NOCASE,
			"lastname" TEXT NOT NULL COLLATE NOCASE,
			"sex" TEXT NOT NULL COLLATE NOCASE,
			PRIMARY KEY("username")
		);
	`)
	if err != nil {
		fmt.Println(err)
	}
	stmt.Exec()

	return &User{
		DB: db,
	}
}

func (user *User) Get() []Uitem {
	uitems := []Uitem{}
	rows, _ := user.DB.Query(`
	SELECT * FROM "users"
	`)
	var email string
	var username string
	var password string
	var age int
	var fname string
	var lname string
	var sex string
	for rows.Next() {
		rows.Scan(&email, &username, &password, &age, &fname, &lname, &sex)
		uitem := Uitem{
			Email:    email,
			Username: username,
			Password: password,
			Age:      age,
			Fname:    fname,
			Lname:    lname,
			Sex:      sex,
		}
		uitems = append(uitems, uitem)
	}
	rows.Close()
	return uitems
}

func (user *User) GetUser(str string) Uitem {
	s := fmt.Sprintf("SELECT * FROM users WHERE username = '%v'", str)
	rows, _ := user.DB.Query(s)
	var email string
	var username string
	var pass string
	var age int
	var fname string
	var lname string
	var sex string
	var uitem Uitem
	if rows.Next() {
		rows.Scan(&email, &username, &pass, &age, &fname, &lname, &sex)
		uitem = Uitem{
			Email:    email,
			Username: username,
			Password: pass,
			Age:      age,
			Fname:    fname,
			Lname:    lname,
			Sex:      sex,
		}
	}
	rows.Close()
	return uitem
}

func (user *User) GetUserbymail(str string) Uitem {
	s := fmt.Sprintf("SELECT * FROM users WHERE email = '%v'", str)
	rows, _ := user.DB.Query(s)
	var email string
	var username string
	var pass string
	var age int
	var fname string
	var lname string
	var sex string
	var uitem Uitem
	if rows.Next() {
		rows.Scan(&email, &username, &pass, &age, &fname, &lname, &sex)
		uitem = Uitem{
			Email:    email,
			Username: username,
			Password: pass,
			Age:      age,
			Fname:    fname,
			Lname:    lname,
			Sex:      sex,
		}
	}
	rows.Close()
	return uitem
}

func LikeDislikecom(w http.ResponseWriter, r *http.Request, s *Session) {
	if !s.IsAuthorized {
		http.Redirect(w, r, "/", http.StatusFound)
		return
	}
	var value string
	if r.URL.Path == "/likecom" {
		value = "l"
	} else if r.URL.Path == "/dislikecom" {
		value = "d"
	}
	values, _ := url.ParseQuery(r.URL.RawQuery)
	posid := values.Get("posid")
	if r.Method != http.MethodGet {
		http.Redirect(w, r, "/view?id="+posid, http.StatusFound)
	}

	comid := values.Get("coid")
	LDcom.Add(Litemcom{
		ComID:    comid,
		Username: s.Username,
		Like:     value,
	})
	http.Redirect(w, r, "/view?id="+posid, http.StatusFound)
}

func profile(w http.ResponseWriter, r *http.Request, s *Session) {
	if !s.IsAuthorized {
		http.Redirect(w, r, "/login", http.StatusFound)
	}
	myposts, mylikes := bunch.GetMyPosts(LD, s.Username)
	data := Info{
		Sess:       s,
		Posts:      myposts,
		LikedPosts: mylikes,
	}
	t, err := template.ParseFiles("template/profile.html", "template/header.html", "template/footer.html", "template/chat.html")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	t.ExecuteTemplate(w, "profile", data)
}

func Filter(w http.ResponseWriter, r *http.Request, s *Session) {
	if !s.IsAuthorized {
		http.Redirect(w, r, "/login", http.StatusFound)
	}
	thread := r.URL.Query().Get("thread")
	p := bunch.Filter(LD, thread)
	data := Info{
		Sess:  s,
		Posts: p,
	}
	t, err := template.ParseFiles("template/posts.html", "template/header.html", "template/footer.html", "template/chat.html")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	t.ExecuteTemplate(w, "posts", data)
}

func Logout(w http.ResponseWriter, r *http.Request, s *Session) {
	if !s.IsAuthorized {
		http.Redirect(w, r, "/login", http.StatusFound)
	}
	sessionStore.Delete(s)
	http.Redirect(w, r, "/", http.StatusFound)
}

func LikeDislike(w http.ResponseWriter, r *http.Request, s *Session) {
	if !s.IsAuthorized {
		http.Redirect(w, r, "/", http.StatusFound)
		return
	}
	var value string
	if r.URL.Path == "/like" {
		value = "l"
	} else if r.URL.Path == "/dislike" {
		value = "d"
	}
	if r.Method == http.MethodGet {
		LD.Add(Litem{
			PostID:   r.FormValue("id"),
			Username: s.Username,
			Like:     value,
		})
	}

	if r.FormValue("r") == "no" {
		http.Redirect(w, r, "/view?id="+r.FormValue("id"), http.StatusFound)
	} else {
		if r.FormValue("r") == "p" {
			http.Redirect(w, r, "/profile", http.StatusFound)
		}
		http.Redirect(w, r, "/#"+r.FormValue("id"), http.StatusFound)
	}
}

func DelComm(w http.ResponseWriter, r *http.Request, s *Session) {
	if !s.IsAuthorized {
		http.Redirect(w, r, "/", http.StatusFound)
		return
	}
	values, _ := url.ParseQuery(r.URL.RawQuery)
	comid := values.Get("coid")
	posid := values.Get("posid")
	if r.Method == http.MethodGet {
		flood.Delete(comid)
	}
	http.Redirect(w, r, "/view?id="+posid, http.StatusFound)
}

func SaveComm(w http.ResponseWriter, r *http.Request, s *Session) {
	if !s.IsAuthorized {
		http.Redirect(w, r, "/", http.StatusFound)
		return
	}
	if r.FormValue("content") != "" && verifyContent(r.FormValue("content")) && r.Method == http.MethodPost {
		flood.Add(Citem{
			CommentId: Generate(),
			PostId:    r.FormValue("id"),
			Author:    s.Username,
			Content:   r.FormValue("content"),
			Creat_at:  time.Now().Format("2006-01-02 15:04"),
		})
		http.Redirect(w, r, "/view?id="+r.FormValue("id"), http.StatusFound)
	} else {
		http.Redirect(w, r, "/comment?id="+r.FormValue("id")+"&err=1&content="+r.FormValue("content"), http.StatusFound)
	}
}

func View(w http.ResponseWriter, r *http.Request, s *Session) {
	var url string
	if !s.IsAuthorized {
		http.Redirect(w, r, "/login", http.StatusFound)
		return
	}
	if r.URL.Path != "/view" && s.Username != "" {
		url = "comment"
	} else {
		url = "view"
	}

	items := bunch.Get(LD)
	id := r.FormValue("id")
	eror := r.FormValue("err")
	var item Item

	for _, v := range items {
		if v.ID == id {
			item = v
		}
	}
	coms := flood.Get(LDcom, id)
	for i, v := range coms {
		if v.Author == s.Username {
			coms[i].ComIsAuthor = true
		}
	}

	data := Info{
		Sess:     s,
		Comments: coms,
		Post:     item,
		IsAuthor: item.Author == s.Username,
	}
	if eror == "1" {
		content := r.FormValue("content")
		data.Content = content
		data.Error = "You need to write something at least 8 symbols"
	}

	t, err := template.ParseFiles("template/comment.html", "template/view.html", "template/header.html", "template/footer.html", "template/chat.html")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	t.ExecuteTemplate(w, url, data)
}

func register(w http.ResponseWriter, r *http.Request, s *Session) {
	t, err := template.ParseFiles("template/auth.html", "template/header.html", "template/footer.html", "template/chat.html")
	if err != nil {
		return
	}

	if r.Method == "GET" {

		e := t.ExecuteTemplate(w, "register", nil)
		if e != nil {
			fmt.Fprint(w, e)
		}
	} else {
		type Data struct {
			Email      string
			Username   string
			Fname      string
			Lname      string
			Age        string
			Password   string
			PHemail    string
			PHusername string
			PHlname    string
			PHfname    string
			PHage      string
		}

		email := r.FormValue("email")
		username := r.FormValue("username")
		password := r.FormValue("password")
		password2 := r.FormValue("password2")
		fname := r.FormValue("fname")
		lname := r.FormValue("lname")
		age := r.FormValue("age")
		sex := r.FormValue("radio")
		data := Data{
			PHemail:    email,
			PHusername: username,
			PHlname:    lname,
			PHfname:    fname,
			PHage:      age,
		}
		if m, _ := regexp.MatchString(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`, email); !m {
			data.Email = "need an apropriate email address"
			t.ExecuteTemplate(w, "register", data)

			return
		}

		if guys.GetUser(username).Username == username {
			data.Username = "already taken"
			t.ExecuteTemplate(w, "register", data)
			return
		}
		fmt.Println(age)
		intage, ageerr := strconv.Atoi(age)
		if intage < 6 || ageerr != nil {
			if ageerr == nil {
				data.Age = "You still young go to play"
			} else {
				data.Age = "Enter a valid age"
			}
			t.ExecuteTemplate(w, "register", data)
			return
		}
		if len(strings.ReplaceAll(fname, " ", "")) < 1 {
			data.Fname = "Enter a valid first name"
			t.ExecuteTemplate(w, "register", data)
			return
		}
		if len(strings.ReplaceAll(lname, " ", "")) < 1 {
			data.Lname = "Enter a valid last name"
			t.ExecuteTemplate(w, "register", data)
			return
		}
		if m, _ := regexp.MatchString(`^\w{5,10}$`, username); !m {
			data.Username = "only eng letters from 5 to 10 symbols, no spaces"
			t.ExecuteTemplate(w, "register", data)
			return
		}
		if password != password2 {
			data.Password = "Password don't match"
			err := t.ExecuteTemplate(w, "register", data)
			if err != nil {
				http.Error(w, err.Error(), 500)
			}
			return
		}
		if m := verifyPassword(password); !m {
			data.Password = "must be eng letters at least 1 Upper case and special symbol, 6 to 10 long, no spaces"
			err := t.ExecuteTemplate(w, "register", data)
			if err != nil {
				http.Error(w, err.Error(), 500)
			}
			return
		}

		hashedPass, err := HashPassword(password)
		if err != nil {
			log.Fatal(err.Error())
		}

		err1 := guys.Add(Uitem{
			Email:    email,
			Username: username,
			Password: hashedPass,
			Age:      intage,
			Fname:    fname,
			Lname:    lname,
			Sex:      sex,
		})
		if err1 != nil {
			data.Email = "user with this email or username already exists"
			e := t.ExecuteTemplate(w, "register", data)
			if e != nil {
				http.Error(w, e.Error(), 500)
			}
			return
		}

		http.Redirect(w, r, "/login", http.StatusFound)
	}
}

func login(w http.ResponseWriter, r *http.Request, s *Session) {
	if AlreadyLoggedIn(r) {
		http.Redirect(w, r, "/", http.StatusFound)
	}
	t, err := template.ParseFiles("template/login.html", "template/header.html", "template/footer.html", "template/chat.html")
	if err != nil {
		return
	}

	type Data struct {
		Username string
		Password string
	}

	var data Data

	username := r.FormValue("username")
	if username == "" {
		t.ExecuteTemplate(w, "login", data)
		return
	}
	password := r.FormValue("password")
	user := guys.GetUser(username)
	mailuser := guys.GetUserbymail(username)
	if user.Username == username {
		if CheckPasswordHash(password, user.Password) {
			for k, v := range sessionStore.data {
				if v.Username == username {
					delete(sessionStore.data, k)
				}
			}
			s.IsAuthorized = true
			s.Username = user.Username

			http.Redirect(w, r, "/", http.StatusFound)
			return

		}
		data.Password = "this password is incorrect"
		t.ExecuteTemplate(w, "login", data)

		return
	} else {
		if mailuser.Email == username {
			if CheckPasswordHash(password, mailuser.Password) {
				for k, v := range sessionStore.data {
					if v.Username == mailuser.Username {
						delete(sessionStore.data, k)
					}
				}
				s.IsAuthorized = true
				s.Username = mailuser.Username

				http.Redirect(w, r, "/", http.StatusFound)
				return

			}
			data.Password = "this password is incorrect"
			t.ExecuteTemplate(w, "login", data)

			return
		}
	}
	data.Username = "no such user"
	t.ExecuteTemplate(w, "login", data)
}

func deletePost(w http.ResponseWriter, r *http.Request, s *Session) {
	if !s.IsAuthorized {
		http.Redirect(w, r, "/login", http.StatusFound)
		return
	}
	if r.Method == http.MethodGet {
		id := strings.TrimPrefix(r.URL.RequestURI(), "/delete?")
		bunch.Delete(id)
	}
	http.Redirect(w, r, "/", http.StatusFound)
}

func savePost(w http.ResponseWriter, r *http.Request, s *Session) {
	if !s.IsAuthorized {
		http.Redirect(w, r, "/", http.StatusFound)
	}

	fmt.Println(r.FormValue("content"), verifyContent(r.FormValue("content")))
	if r.FormValue("content") != "" && verifyContent(r.FormValue("content")) {
		if r.FormValue("id") != "" {
			id := r.FormValue("id")
			r.ParseForm()
			var thread string
			threadList := r.Form["thread"]
			if len(threadList) > 1 {
				for i, v := range threadList {
					thread += v
					if i != len(threadList)-1 {
						thread += "#"
					}

				}
			} else {
				thread = threadList[0]
			}
			bunch.Update(Item{
				Content: r.FormValue("content"),
				Thread:  thread,
			}, id)
		} else {
			r.ParseForm()
			var thread string
			threadList := r.Form["thread"]
			if len(threadList) > 1 {
				for i, v := range threadList {
					thread += v
					if i != len(threadList)-1 {
						thread += "#"
					}

				}
			} else {
				if len(threadList) > 1 {
					thread = threadList[0]
				} else {
					thread = "general"
				}
			}

			bunch.Add(Item{
				ID:       Generate(),
				Author:   s.Username,
				Content:  r.FormValue("content"),
				Thread:   thread,
				Creat_at: time.Now().Format("2006-01-02 à 15:04"),
			})

		}

		http.Redirect(w, r, "/", http.StatusFound)

	} else {
		t, err := template.ParseFiles("template/newpost.html", "template/header.html", "template/footer.html", "template/chat.html")
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		data := Info{
			Sess:  s,
			Error: "You need to write something at least 8 symbols",
			Post: Item{
				Content: r.FormValue("content"),
			},
		}
		t.ExecuteTemplate(w, "newpost", data)

	}
}

func newPost(w http.ResponseWriter, r *http.Request, s *Session) {
	if !s.IsAuthorized {
		http.Redirect(w, r, "/login", http.StatusFound)
	}
	items := bunch.Get(LD)
	id := strings.TrimPrefix(r.URL.RequestURI(), "/edit?")
	var item Item

	for _, v := range items {
		if v.ID == id {
			item = v
		}
	}
	t, err := template.ParseFiles("template/newpost.html", "template/header.html", "template/footer.html", "template/chat.html")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	data := Info{
		Sess: s,
		Post: item,
	}
	t.ExecuteTemplate(w, "newpost", data)
}

type Info struct {
	Sess       *Session
	Comments   []Citem
	Posts      []Item
	LikedPosts []Item
	Post       Item
	IsAuthor   bool
	Error      string
	Content    string
}

func Postsss(w http.ResponseWriter, r *http.Request, s *Session) {
	if r.URL.Path != "/" {
		w.WriteHeader(404)
		w.Write([]byte("Error 404 not found"))
		return
	}
	items := bunch.Get(LD)
	for i, v := range items {
		if v.Author == s.Username {
			items[i].IsAuthor = true
		}
	}
	data := Info{
		Sess:  s,
		Posts: items,
	}
	if !s.IsAuthorized {
		http.Redirect(w, r, "/login", http.StatusFound)
	}
	t, err := template.ParseFiles("template/posts.html", "template/header.html", "template/footer.html", "template/chat.html")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	t.ExecuteTemplate(w, "posts", data)
}

// utilities
func Generate() string {
	u2, err := uuid.NewV4()
	if err != nil {
		fmt.Println("error uuid")
	}
	return fmt.Sprintf("%x", u2)
}

func verifyPassword(s string) bool {
	letters := 0
	var sevenOrTen, number, upper, special bool
	for _, c := range s {
		switch {
		case unicode.IsNumber(c):
			number = true
			letters++
		case unicode.IsUpper(c):
			upper = true
			letters++
		case unicode.IsPunct(c) || unicode.IsSymbol(c):
			special = true
		case unicode.IsLetter(c):
			letters++
		default:
			// return false, false, false, false
		}
	}
	sevenOrTen = letters >= 8 && letters <= 10
	if sevenOrTen && number && upper && special {
		return true
	}
	return false
}

func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

const (
	COOKIE_NAME = "sessionId"
)

type Session struct {
	Id           string
	Username     string
	IsAuthorized bool
}

type SessionStore struct {
	data map[string]*Session
}

var sessionStore = NewSessionStore()

func NewSessionStore() *SessionStore {
	s := new(SessionStore)
	s.data = make(map[string]*Session)
	return s
}

func (store *SessionStore) Get(sessionId string) *Session {
	session := store.data[sessionId]
	if session == nil {
		return &Session{Id: sessionId}
	}
	return session
}

func (store *SessionStore) Set(session *Session) {
	store.data[session.Id] = session
}

func (store *SessionStore) Delete(session *Session) {
	delete(store.data, session.Id)
}

func ensureCookie(r *http.Request, w http.ResponseWriter) string {
	cookie, _ := r.Cookie(COOKIE_NAME)
	if cookie != nil {
		if cookie.Expires.Before(time.Now()) {
			cookie.Expires = time.Now().Add(5 * time.Hour)
			http.SetCookie(w, cookie)

		}
		return cookie.Value
	}
	sessionId := Generate()

	cookie = &http.Cookie{
		Name:     COOKIE_NAME,
		Value:    sessionId,
		Expires:  time.Now().Add(5 * time.Second),
		SameSite: http.SameSiteLaxMode,
	}
	http.SetCookie(w, cookie)

	return sessionId
}

func Convertfunc(next func(w http.ResponseWriter, r *http.Request, s *Session)) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		sessionId := ensureCookie(r, w)

		session := sessionStore.Get(sessionId)

		sessionStore.Set(session)
		next(w, r, session)
	}
}

func AlreadyLoggedIn(r *http.Request) bool {
	c, err := r.Cookie(COOKIE_NAME)
	if err != nil {
		return false
	}
	sess, _ := sessionStore.data[c.Value]
	return sess.Username != ""
}

func verifyContent(s string) bool {
	letters := 0

	var symbol, sevenOrTen bool
	for _, c := range s {
		if unicode.IsLetter(c) {
			symbol = true
			letters++
		}
	}
	sevenOrTen = letters >= 8
	if sevenOrTen && symbol {
		return true
	}
	return false
}

func main() {
	manager := newManager()
	db, _ := sql.Open("sqlite3", "Forum.db")
	defer db.Close()
	guys = NewUser(db)
	bunch = NewPost(db)
	flood = NewComm(db)
	LD = NewLD(db)
	LDcom = NewLDcom(db)
	Msgdb = NewMessage(db)
	fs := http.FileServer(http.Dir("template"))
	http.Handle("/template/", http.StripPrefix("/template/", fs))
	http.HandleFunc("/", Convertfunc(Postsss))
	http.HandleFunc("/write", Convertfunc(newPost))
	http.HandleFunc("/edit", Convertfunc(newPost))
	http.HandleFunc("/view", Convertfunc(View))
	http.HandleFunc("/comment", Convertfunc(View))
	http.HandleFunc("/delete", Convertfunc(deletePost))
	http.HandleFunc("/SavePost", Convertfunc(savePost))
	http.HandleFunc("/savecomm", Convertfunc(SaveComm))
	http.HandleFunc("/login", Convertfunc(login))
	http.HandleFunc("/register", Convertfunc(register))
	http.HandleFunc("/deleteCom", Convertfunc(DelComm))
	http.HandleFunc("/like", Convertfunc(LikeDislike))
	http.HandleFunc("/dislike", Convertfunc(LikeDislike))
	http.HandleFunc("/logout", Convertfunc(Logout))
	http.HandleFunc("/filter", Convertfunc(Filter))
	http.HandleFunc("/profile", Convertfunc(profile))
	http.HandleFunc("/likecom", Convertfunc(LikeDislikecom))
	http.HandleFunc("/dislikecom", Convertfunc(LikeDislikecom))
	http.HandleFunc("/ws", Convertfunc(manager.Servws))

	fmt.Printf("server start at https://localhost%s/\n", Port)
	if error := http.ListenAndServeTLS(Port, "cert.pem", "key.pem", nil); error != nil {
		log.Fatal(error)
		os.Exit(0)
	}
}
