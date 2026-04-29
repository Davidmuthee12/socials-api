package db

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"math/rand"

	"github.com/Davidmuthee12/socials/internal/store"
)

var username = []string{
	"kwame", "jabari", "amara", "zuri", "imani",
	"tariq", "sekou", "chini", "malik", "ayodele",
	"biko", "sade", "ngozi", "amarae", "zola",
	"nala", "kofi", "abasi", "folami", "adisa",
	"themba", "thabo", "lerato", "kgosi", "tumelo",
	"chiamaka", "chinwe", "olamide", "adetola", "ayanna",
	"makena", "baraka", "rafiki", "jabali", "simba",
	"nyasha", "tendai", "rutendo", "tafadzwa", "munashe",
	"kato", "mutiso", "mwangi", "otieno", "kamau",
	"nyong", "abebe", "selam", "desta", "kebede",
}

var titles = []string{
	"firstpost", "myjourney", "techlife", "golangfun", "dailythoughts",
	"webdevtips", "codingvibes", "startupdream", "devdiary", "learnfast",
	"buildinpublic", "debuglife", "cleanCode", "backendmagic", "frontendflow",
	"systemdesign", "cloudnotes", "apilife", "devmindset", "scalablesystems",
}

var content = []string{
	"learning go is fun and powerful",
	"building scalable systems step by step",
	"debugging teaches patience and logic",
	"writing clean code improves teamwork",
	"backend development is all about structure",
	"frontend brings ideas to life visually",
	"consistency beats motivation in coding",
	"small improvements daily lead to mastery",
	"understanding basics is everything",
	"practice makes coding feel natural",
	"mistakes are part of the dev journey",
	"reading docs saves hours of confusion",
	"projects build real experience",
	"focus on solving real problems",
	"code readability matters a lot",
	"performance optimization is key",
	"learning never really stops in tech",
	"building APIs is a core backend skill",
	"good architecture scales better",
	"testing ensures reliability",
}

var tags = []string{
	"go", "backend", "frontend", "api", "database",
	"postgres", "docker", "kubernetes", "cloud", "devops",
	"javascript", "react", "nextjs", "tailwind", "css",
	"html", "systemdesign", "testing", "debugging", "performance",
}

var comments = []string{
	"great post",
	"very helpful",
	"learned something new",
	"this is insightful",
	"thanks for sharing",
	"nice explanation",
	"well written",
	"this helped me a lot",
	"clear and simple",
	"awesome work",
	"keep it up",
	"interesting perspective",
	"good breakdown",
	"makes sense",
	"really useful",
	"solid content",
	"straight to the point",
	"helpful tips",
	"appreciate this",
	"good read",
}

func Seed(store store.Storage, db *sql.DB) {
	ctx := context.Background()

	users := generateUsers(100)
	tx, _ := db.BeginTx(ctx, nil)
	for _, user := range users {
		_ = tx.Rollback()
		if err := store.Users.Create(ctx, tx, user); err != nil {
			log.Println("Error creating user: ", err)
			return
		}
	}

	tx.Commit()
	posts := generatePosts(200, users)
	for _, post := range posts {
		if err := store.Posts.Create(ctx, post); err != nil {
			log.Println("Error creating post: ", err)
			return
		}
	}

	comments := generateComments(500, users, posts)
	for _, comment := range comments {
		if err := store.Comments.Create(ctx, comment); err != nil {
			log.Println("Error creating comment: ", err)
			return
		}
	}

	log.Println("Seeding complete")
}

func generateUsers(num int) []*store.User {
	users := make([]*store.User, num)

	for i := 0; i < num; i++ {
		users[i] = &store.User{
			Username: username[i%len(username)] + fmt.Sprintf("%d", i),
			Email:    username[i%len(username)] + fmt.Sprintf("%d", i) + "@example.com",
			Role: store.Role{
				Name: "user",
			},
		}
	}

	return users
}

func generatePosts(num int, users []*store.User) []*store.Post {
	posts := make([]*store.Post, num)
	for i := 0; i < num; i++ {
		user := users[rand.Intn(len(users))]

		posts[i] = &store.Post{
			UserID:  user.ID,
			Title:   titles[rand.Intn(len(titles))],
			Content: content[rand.Intn(len(content))],
			Tags: []string{
				tags[rand.Intn(len(tags))],
				tags[rand.Intn(len(tags))],
			},
		}
	}

	return posts
}

func generateComments(num int, users []*store.User, posts []*store.Post) []*store.Comment {
	cms := make([]*store.Comment, num)
	for i := 0; i < num; i++ {
		cms[i] = &store.Comment{
			PostID:  posts[rand.Intn(len(posts))].ID,
			UserID:  users[rand.Intn(len(users))].ID,
			Content: comments[rand.Intn(len(comments))],
		}
	}

	return cms
}
