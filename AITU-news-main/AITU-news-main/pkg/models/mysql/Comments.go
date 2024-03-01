package mysql

import (
	"database/sql"
	_ "database/sql"
	"encoding/json"
	_ "errors"
	"net/http"

	// Import the models package that we just created. You need to prefix this with
	// whatever module path you set up back in chapter 02.02 (Project Setup and Enabling // Modules) so that the import statement looks like this:
	// "{your-module-path}/pkg/models".
	_ "xximsz.net/snippetbox/pkg/models"
)

type CommentModel struct {
	DB *sql.DB
}

type Comment struct {
	ID     int    `json:"id"`
	UserID int    `json:"user_id"`
	NewsID int    `json:"news_id"`
	Text   string `json:"text"`
}

func PostCommentHandler(w http.ResponseWriter, r *http.Request) {
	// Parse request body to get comment data
	var comment Comment
	err := json.NewDecoder(r.Body).Decode(&comment)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Save comment to database
	db := sql.DB{}
	_, err = db.Exec("INSERT INTO comments (user_id, news_id, text) VALUES ($1, $2, $3)",
		comment.UserID, comment.NewsID, comment.Text)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Respond with success message
	w.WriteHeader(http.StatusCreated)
	w.Write([]byte("Comment posted successfully"))
}

func DeleteCommentHandler(w http.ResponseWriter, r *http.Request) {
	// Parse comment ID from request parameters or URL
	commentID := r.URL.Query().Get("comment_id")

	// Implement your authentication logic here to verify if the user is the owner or an admin

	// Delete comment from database
	db := sql.DB{}
	_, err := db.Exec("DELETE FROM comments WHERE id = $1", commentID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Respond with success message
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Comment deleted successfully"))
}

func ShowCommentsHandler(w http.ResponseWriter, r *http.Request) {
	// Parse news ID from request parameters or URL
	newsID := r.URL.Query().Get("news_id")

	// Query comments for the specified news ID
	db := sql.DB{}
	rows, err := db.Query("SELECT id, user_id, text FROM comments WHERE news_id = $1", newsID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	// Fetch comments and construct response
	var comments []Comment
	for rows.Next() {
		var comment Comment
		if err := rows.Scan(&comment.ID, &comment.UserID, &comment.Text); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		comments = append(comments, comment)
	}

	// Convert comments to JSON and respond
	response, err := json.Marshal(comments)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(response)
}
