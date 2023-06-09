package main

// import (
// 	"database/sql"
// 	"fmt"
// )

// type Likesdb struct {
// 	DB *sql.DB
// }

// type Like struct {
// 	Id      int
// 	User_id int
// 	Post_id int
// 	Value   string
// }

// func CreatLikesTable(db *sql.DB) *Likesdb {
// 	stmt, _ := db.Prepare(`
// 		CREATE TABLE IF NOT EXISTS "likes" (
// 			"like_id" INTEGER PRIMARY KEY AUTOINCREMENT,
// 			"user_id"	INTEGER NOT NUL,
// 			"post_id" INTEGER NOT NULL,
// 			"value" TEXT NOT NULL,
// 			FOREIGN KEY(user_id) REFERENCES users(user_id) ON DELETE CASCADE,
// 			FOREIGN KEY(post_id) REFERENCES posts(post_id) ON DELETE CASCADE
// 		);
// 	`)
// 	stmt.Exec()

// 	return &Likesdb{
// 		DB: db,
// 	}
// }

// func (likes *Likesdb) GetLikeByPostAndUser(like_id, user_id int) Like {
// 	like := Like{}

// 	s := fmt.Sprintf("SELECT * FROM likes WHERE post_id = %v AND user_id = %v", like_id, user_id)

// 	rows, _ := likes.DB.Query(s)
// 	var postid int
// 	var author int
// 	var value string
// 	if rows.Next() {
// 		rows.Scan(&postid, &author, &like)
// 		like = Like{
// 			Post_id: postid,
// 			User_id: author,
// 			Value:   value,
// 		}
// 	}
// 	rows.Close()
// 	return like
// }

// func (likes *Likesdb) AddLike(like Like) {
// 	one := likes.GetLikeByPostAndUser(like.Post_id, like.User_id)
// 	var s string
// 	if one.Value == "" {
// 		s = "INSERT INTO likes (value, post_id, user_id) values (?, ?, ?)"
// 	} else if like.Value != one.Value {
// 		s = "UPDATE likes SET value = ? WHERE post_id = ? AND user_id = ?"
// 	} else {
// 		s = "DELETE FROM likes WHERE like = ? AND post_id = ? AND user_id = ?"
// 	}
// 	stmt, _ := likes.DB.Prepare(s)
// 	_, err := stmt.Exec(like.Value, like.Post_id, like.User_id)
// 	fmt.Println("this is likes db", err)
// }

// func (likes *Likesdb) GetAllPostLikesOrDislike(post_id int, like_value string) []Like {
// 	items := []Like{}
// 	var s string
// 	s = fmt.Sprintf("SELECT * FROM likes WHERE post_id = %v AND value = '%v'", post_id, like_value)

// 	rows, _ := likes.DB.Query(s)
// 	var idi int
// 	var postid int
// 	var author int
// 	var like string
// 	for rows.Next() {
// 		rows.Scan(&idi, &postid, &author, &like)
// 		item := Like{
// 			Id:      idi,
// 			Post_id: postid,
// 			User_id: author,
// 			Value:   like,
// 		}
// 		items = append(items, item)
// 	}
// 	rows.Close()
// 	return items
// }

// func (likes *Likesdb) GetAllUserLikes(user_id int) []Like {
// 	items := []Like{}
// 	var s string
// 	s = fmt.Sprintf("SELECT * FROM likes WHERE user_id = %v AND value = 'l'", user_id)

// 	rows, _ := likes.DB.Query(s)
// 	var idi int
// 	var postid int
// 	var author int
// 	var like string
// 	for rows.Next() {
// 		rows.Scan(&idi, &postid, &author, &like)
// 		item := Like{
// 			Id:      idi,
// 			Post_id: postid,
// 			User_id: author,
// 			Value:   like,
// 		}
// 		items = append(items, item)
// 	}
// 	rows.Close()
// 	return items
// }
