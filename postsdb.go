package main

// import (
// 	"database/sql"
// 	"fmt"
// 	"time"
// )

// type Postdb struct {
// 	DB *sql.DB
// }

// type Post struct {
// 	Id       int
// 	User_id  int
// 	Creat_at time.Time
// 	Content  string
// 	Thread   string
// 	L        int
// 	D        int
// }

// func CreatPostTable(db *sql.DB) *Postdb {
// 	stmt, _ := db.Prepare(`
// 		CREATE TABLE IF NOT EXISTS "posts" (
// 			"post_id" INTEGER PRIMARY KEY AUTOINCREMENT,
// 			"user_id"	INTEGER NOT NUL,
// 			"creat_at" DATE NOT NULL,
// 			"content" TEXT NOT NULL,
// 			"thread" TEXT NOT NULL DEFAULT 'general',
// 			FOREIGN KEY(user_id) REFERENCES users(user_id) ON DELETE CASCADE
// 		);
// 	`)
// 	stmt.Exec()

// 	return &Postdb{
// 		DB: db,
// 	}
// }

// func (db *Postdb) addPost(post Post) error {
// 	stmt, er := db.DB.Prepare(`
// 	INSERT INTO "posts" (user_id, creat_at, content, thread) values (?, ?, ?, ?)
// 	`)
// 	if er != nil {
// 		return er
// 	}
// 	_, err := stmt.Exec(post.User_id, post.Creat_at.Format("2006-01-02"), post.Content, post.Thread)
// 	if err != nil {
// 		return err
// 	}
// 	return nil
// }

// func (db *Postdb) UpdatPost(post Post) error {
// 	stmt, er := db.DB.Prepare(`
// 	UPDATE "posts" SET "content" = ?, "thread" = ?, WHERE "post_id" = ?
// 	`)
// 	if er != nil {
// 		return er
// 	}
// 	_, err := stmt.Exec(post.Content, post.Thread, post.Id)
// 	if err != nil {
// 		return err
// 	}
// 	return nil
// }

// func (db *Postdb) getPostbyquery(s string, likes *Likesdb) Post {
// 	rows, _ := db.DB.Query(s)
// 	var id int
// 	var User_id int
// 	var Creat_at string
// 	var Content string
// 	var Thread string
// 	var post Post
// 	if rows.Next() {
// 		rows.Scan(&id, &User_id, &Creat_at, &Content, &Thread)
// 		age, _ := time.Parse("2006-01-02", Creat_at)
// 		post = Post{
// 			Id:       id,
// 			User_id:  User_id,
// 			Creat_at: age,
// 			Content:  Content,
// 			Thread:   Thread,
// 			L:        len(likes.GetAllPostLikesOrDislike(id, "l")),
// 			D:        len(likes.GetAllPostLikesOrDislike(id, "d")),
// 		}
// 	}
// 	rows.Close()
// 	return post
// }

// func (db *Postdb) GetAllPostsByQuery(s string, likes *Likesdb) []Post {
// 	Posts := []Post{}
// 	rows, _ := db.DB.Query(s)
// 	var id int
// 	var User_id int
// 	var Creat_at string
// 	var Content string
// 	var Thread string
// 	for rows.Next() {
// 		rows.Scan(&id, &User_id, &Creat_at, &Content, &Thread)
// 		age, _ := time.Parse("2006-01-02", Creat_at)
// 		post := Post{
// 			Id:       id,
// 			User_id:  User_id,
// 			Creat_at: age,
// 			Content:  Content,
// 			Thread:   Thread,
// 			L:        len(likes.GetAllPostLikesOrDislike(id, "l")),
// 			D:        len(likes.GetAllPostLikesOrDislike(id, "d")),
// 		}
// 		Posts = append(Posts, post)
// 	}
// 	rows.Close()
// 	return Posts
// }

// func (db *Postdb) GetAllPosts(like *Likesdb) []Post {
// 	s := fmt.Sprintf("SELECT * FROM posts ORDER BY creat_at DESC")
// 	return db.GetAllPostsByQuery(s, like)
// }

// func (db *Postdb) GetPostById(post_id int, like *Likesdb) Post {
// 	post := Post{}
// 	s := fmt.Sprintf("SELECT * FROM posts WHERE post_id = %d", post_id)
// 	rows, _ := db.DB.Query(s)
// 	var id int
// 	var User_id int
// 	var Creat_at string
// 	var Content string
// 	var Thread string
// 	if rows.Next() {
// 		rows.Scan(&id, &User_id, &Creat_at, &Content, &Thread)
// 		age, _ := time.Parse("2006-01-02", Creat_at)
// 		post = Post{
// 			Id:       id,
// 			User_id:  User_id,
// 			Creat_at: age,
// 			Content:  Content,
// 			Thread:   Thread,
// 			L:        len(like.GetAllPostLikesOrDislike(id, "l")),
// 			D:        len(like.GetAllPostLikesOrDislike(id, "d")),
// 		}
// 	}
// 	rows.Close()
// 	return post
// }

// func (db *Postdb) GetAllPostsByThread(thread string, like *Likesdb) []Post {
// 	s := fmt.Sprintf("SELECT * FROM posts WHERE thread LIKE '%v' ORDER BY creat_at DESC", "%"+thread+"%")
// 	return db.GetAllPostsByQuery(s, like)
// }

// func (db *Postdb) GetAllMyLikedPost(user_id int, like *Likesdb) []Post {
// 	user_likes := like.GetAllUserLikes(user_id)
// 	var result []Post
// 	for _, v := range user_likes {
// 		result = append(result, db.GetPostById(v.Post_id, like))
// 	}
// 	return result
// }

// func (db *Postdb) GetAllMyCreatedPosts(user_id int, like *Likesdb) []Post {
// 	s := fmt.Sprintf("SELECT * FROM posts WHERE user_id = %v ORDER BY creat_at DESC", user_id)
// 	return db.GetAllPostsByQuery(s, like)
// }

// func (db *Postdb) Delete(post_id int) {
// 	stmt, _ := db.DB.Prepare(`DELETE FROM "posts" WHERE "post_id" = ?`)
// 	stmt.Exec(post_id)
// }
