package main

import (
	"html/template"
	"net/http"
	"strconv"
	"sync"

	"html-go/models"
)

var templates = template.Must(
	template.ParseGlob("templates/*.html"),
)

var mu sync.Mutex

func main() {

	// Home
	http.HandleFunc("/", homeHandler)

	// Get Post By ID
	http.HandleFunc("/post", getPostHandler)

	// Create
	http.HandleFunc("/create", createPostHandler)

	// Update
	http.HandleFunc("/update", updatePostHandler)

	// Delete
	http.HandleFunc("/delete", deletePostHandler)

	// Static files
	http.Handle(
		"/style/",
		http.StripPrefix(
			"/style/",
			http.FileServer(http.Dir("style")),
		),
	)

	println("Server running on http://localhost:8080")

	http.ListenAndServe(":8080", nil)
}

// =========================
// HOME
// =========================

func homeHandler(w http.ResponseWriter, r *http.Request) {

	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}

	mu.Lock()
	posts := append([]models.Post(nil), models.Posts...)
	mu.Unlock()

	templates.ExecuteTemplate(w, "index.html", posts)
}

// =========================
// GET POST BY ID
// =========================

func getPostHandler(w http.ResponseWriter, r *http.Request) {

	id, err := strconv.Atoi(r.URL.Query().Get("id"))

	if err != nil {
		http.NotFound(w, r)
		return
	}

	mu.Lock()
	defer mu.Unlock()

	for _, post := range models.Posts {

		if post.ID == id {

			templates.ExecuteTemplate(
				w,
				"post.html",
				post,
			)

			return
		}
	}

	http.NotFound(w, r)
}

// =========================
// CREATE
// =========================

func createPostHandler(w http.ResponseWriter, r *http.Request) {

	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	title := r.FormValue("title")
	content := r.FormValue("content")

	mu.Lock()

	newID := 1

	if len(models.Posts) > 0 {
		newID = models.Posts[len(models.Posts)-1].ID + 1
	}

	post := models.Post{
		ID:      newID,
		Title:   title,
		Content: content,
	}

	models.Posts = append(models.Posts, post)

	mu.Unlock()

	http.Redirect(w, r, "/", http.StatusSeeOther)
}

// =========================
// UPDATE
// =========================

func updatePostHandler(w http.ResponseWriter, r *http.Request) {

	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	id, err := strconv.Atoi(r.FormValue("id"))

	if err != nil {
		http.NotFound(w, r)
		return
	}

	title := r.FormValue("title")
	content := r.FormValue("content")

	mu.Lock()
	defer mu.Unlock()

	for i := range models.Posts {

		if models.Posts[i].ID == id {

			models.Posts[i].Title = title
			models.Posts[i].Content = content

			break
		}
	}

	http.Redirect(w, r, "/", http.StatusSeeOther)
}

// =========================
// DELETE
// =========================

func deletePostHandler(w http.ResponseWriter, r *http.Request) {

	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	id, err := strconv.Atoi(r.FormValue("id"))

	if err != nil {
		http.NotFound(w, r)
		return
	}

	mu.Lock()
	defer mu.Unlock()

	for i, post := range models.Posts {

		if post.ID == id {

			models.Posts = append(
				models.Posts[:i],
				models.Posts[i+1:]...,
			)

			break
		}
	}

	http.Redirect(w, r, "/", http.StatusSeeOther)
}