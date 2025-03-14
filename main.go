package main

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/thedevsaddam/renderer"
	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

// Global renderer instance for JSON responses
var rnd *renderer.Render

// Global MongoDB database reference
var db *mgo.Database

// Constants for MongoDB and server configuration
const (
	hostName       string = "localhost:27017" // MongoDB host
	dbName         string = "demo_todo"       // Database name
	collectionName string = "todo"            // Collection name
	port           string = ":9000"           // Server port
)

// Struct for MongoDB document representation
type todoModel struct {
	ID        bson.ObjectId `bson:"_id,omitempty"`
	Title     string        `bson:"title"`
	Completed bool          `bson:"completed"`
	CreatedAt time.Time     `bson:"createdAt"`
}

// Struct for API response
type todo struct {
	ID        string    `json:"id"`
	Title     string    `json:"title"`
	Completed bool      `json:"completed"`
	CreatedAt time.Time `json:"createdAt"`
}

// Initialize the renderer and database connection
func init() {
	rnd = renderer.New()

	// Establish connection to MongoDB
	sess, err := mgo.Dial(hostName)
	checkErr(err)

	// Set session mode for consistency
	sess.SetMode(mgo.Monotonic, true)

	// Select the database
	db = sess.DB(dbName)
}

// Error handling function
func checkErr(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

// Home page handler
func homeHandler(w http.ResponseWriter, r *http.Request) {
	err := rnd.Template(w, http.StatusOK, []string{"static/home.tpl"}, nil)
	checkErr(err)
}

// Fetch all To-Do items from the database
func fetchTodos(w http.ResponseWriter, r *http.Request) {
	var todos []todoModel

	// Retrieve all todos from the database
	if err := db.C(collectionName).Find(bson.M{}).All(&todos); err != nil {
		rnd.JSON(w, http.StatusInternalServerError, renderer.M{"message": "Failed to fetch todos", "error": err})
		return
	}

	// Convert to API response format
	var todoList []todo
	for _, t := range todos {
		todoList = append(todoList, todo{
			ID:        t.ID.Hex(),
			Title:     t.Title,
			Completed: t.Completed,
			CreatedAt: t.CreatedAt,
		})
	}

	// Return the JSON response
	rnd.JSON(w, http.StatusOK, renderer.M{"data": todoList})
}

// Create a new To-Do item
func createTodo(w http.ResponseWriter, r *http.Request) {
	var t todo

	// Decode request body
	if err := json.NewDecoder(r.Body).Decode(&t); err != nil {
		rnd.JSON(w, http.StatusBadRequest, renderer.M{"message": "Invalid request payload", "error": err})
		return
	}

	// Validate title
	if t.Title == "" {
		rnd.JSON(w, http.StatusBadRequest, renderer.M{"message": "Title is required"})
		return
	}

	// Create a new todo document
	tm := todoModel{
		ID:        bson.NewObjectId(),
		Title:     t.Title,
		Completed: false,
		CreatedAt: time.Now(),
	}

	// Insert into database
	if err := db.C(collectionName).Insert(&tm); err != nil {
		rnd.JSON(w, http.StatusInternalServerError, renderer.M{"message": "Failed to save todo", "error": err})
		return
	}

	rnd.JSON(w, http.StatusCreated, renderer.M{"message": "Todo created successfully", "todo_id": tm.ID.Hex()})
}

// Delete a To-Do item
func deleteTodo(w http.ResponseWriter, r *http.Request) {
	id := strings.TrimSpace(chi.URLParam(r, "id"))

	// Validate ObjectId
	if !bson.IsObjectIdHex(id) {
		rnd.JSON(w, http.StatusBadRequest, renderer.M{"message": "Invalid ID"})
		return
	}

	// Remove document
	if err := db.C(collectionName).RemoveId(bson.ObjectIdHex(id)); err != nil {
		rnd.JSON(w, http.StatusInternalServerError, renderer.M{"message": "Failed to remove todo", "error": err})
		return
	}

	rnd.JSON(w, http.StatusOK, renderer.M{"message": "Todo deleted successfully"})
}

// Update a To-Do item
func updateTodo(w http.ResponseWriter, r *http.Request) {
	id := strings.TrimSpace(chi.URLParam(r, "id"))

	// Validate ObjectId
	if !bson.IsObjectIdHex(id) {
		rnd.JSON(w, http.StatusBadRequest, renderer.M{"message": "Invalid ID"})
		return
	}

	var t todo

	// Decode request body
	if err := json.NewDecoder(r.Body).Decode(&t); err != nil {
		rnd.JSON(w, http.StatusBadRequest, renderer.M{"message": "Invalid request payload", "error": err})
		return
	}

	// Update document
	if err := db.C(collectionName).Update(
		bson.M{"_id": bson.ObjectIdHex(id)},
		bson.M{"$set": bson.M{"title": t.Title, "completed": t.Completed}},
	); err != nil {
		rnd.JSON(w, http.StatusInternalServerError, renderer.M{"message": "Failed to update todo", "error": err})
		return
	}

	rnd.JSON(w, http.StatusOK, renderer.M{"message": "Todo updated successfully"})
}

// Set up and start the HTTP server
func main() {
	r := chi.NewRouter()

	// Middleware
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	// Routes
	r.Get("/", homeHandler)
	r.Get("/todos", fetchTodos)
	r.Post("/todos", createTodo)
	r.Put("/todos/{id}", updateTodo)
	r.Delete("/todos/{id}", deleteTodo)

	// Start server
	log.Println("Server started on port", port)
	http.ListenAndServe(port, r)
}
