package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/a-h/templ"
	_ "github.com/jackc/pgx/v5/stdlib"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
	"sync"
)

type Post struct {
	ID   int    `json:"id"`
	Body string `json:"body"`
}

type EnvVars struct {
	DBPass string
}

type User struct {
	ID           int
	FirstName    string
	LastName     string
	EmailAddress string
	LocationID   int
}

type UserQuery struct {
	ID           int
	FirstName    string
	LastName     string
	EmailAddress string
	Location     string
}

const (
	host   = "localhost"
	port   = 5432
	user   = "go_user"
	dbname = "go_database"
)

var (
	posts   = make(map[int]Post)
	nextID  = 1
	postsMu sync.Mutex
)

func main() {
	envVars := EnvVars{}
	envVars.getEnv()

	dbUrl := generateDBConnectionURL(envVars.DBPass)
	dbConn, err := sql.Open("pgx", dbUrl)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
		os.Exit(1)
	}
	fmt.Println("Connected to database")

	var users []UserQuery
	rows, err := dbConn.Query("SELECT u.id, first_name, last_name, email_address, l.name FROM users u JOIN public.locations l on l.id = u.location_id;")

	if err != nil {
		fmt.Fprintf(os.Stderr, "Query failed: %v\n", err)
		os.Exit(1)
	}
	defer rows.Close()

	for rows.Next() {
		var user UserQuery
		if err := rows.Scan(&user.ID, &user.FirstName, &user.LastName, &user.EmailAddress, &user.Location); err != nil {
			fmt.Fprintf(os.Stderr, "Unable to scan row: %v\n", err)
			os.Exit(1)
		}
		users = append(users, user)
	}

	fmt.Println(users)

	http.HandleFunc("/posts", postsHandler)
	http.HandleFunc("/posts/", postHandler)

	fmt.Println("Server is running at http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func (envvars *EnvVars) getEnv() {
	envvars.DBPass = os.Getenv("DB_PASS")
	if envvars.DBPass == "" {
		log.Fatal("DB_PASS environment variable not set")
	}
}

func generateDBConnectionURL(dbpassword string) string {
	// postgres://username:password@localhost:5432/database_name
	return fmt.Sprintf("postgres://%s:%s@%s:%d/%s", user, dbpassword, host, port, dbname)
}

//func execQuery(query string) {
//	db, err := sql.Open("sqlite3", "./posts.db")
//}

func postsHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		handleGetPosts(w, r)
	case "POST":
		handlePostPosts(w, r)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func postHandler(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.URL.Path[len("/posts/"):])
	if err != nil {
		http.Error(w, "Invalid post ID", http.StatusBadRequest)
		return
	}

	switch r.Method {
	case "GET":
		handleGetPost(w, r, id)
	case "DELETE":
		handleDeletePost(w, r, id)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func handleGetPosts(w http.ResponseWriter, r *http.Request) {
	postsMu.Lock()
	defer postsMu.Unlock()

	ps := make([]Post, 0, len(posts))
	for _, p := range posts {
		ps = append(ps, p)
	}

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(ps)
}

func handlePostPosts(w http.ResponseWriter, r *http.Request) {
	var p Post

	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Error reading request body", http.StatusInternalServerError)
		return
	}

	if err := json.Unmarshal(body, &p); err != nil {
		http.Error(w, "Error parsing request body", http.StatusBadRequest)
		return
	}

	postsMu.Lock()
	defer postsMu.Unlock()

	p.ID = nextID
	nextID++
	posts[p.ID] = p

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	_ = json.NewEncoder(w).Encode(p)
}

func handleGetPost(w http.ResponseWriter, r *http.Request, id int) {
	postsMu.Lock()
	defer postsMu.Unlock()

	p, ok := posts[id]
	if !ok {
		http.Error(w, "post not found", http.StatusNotFound)
	}

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(p)
}

func handleDeletePost(w http.ResponseWriter, r *http.Request, id int) {
	postsMu.Lock()
	defer postsMu.Unlock()

	_, ok := posts[id]
	if !ok {
		http.Error(w, "post not found", http.StatusNotFound)
	}

	delete(posts, id)
	w.WriteHeader(http.StatusOK)
}
