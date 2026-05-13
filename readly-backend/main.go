package main

import (
	"database/sql"
	"encoding/base64"
	"encoding/json"
	"io"
	"log"
	"net/http"

	"golang.org/x/crypto/bcrypt"
	_ "modernc.org/sqlite"
)

// ! MODELLER
type User struct {
	FirstName string `json:"firstname"`
	LastName  string `json:"lastname"`
	Email     string `json:"email"`
	Username  string `json:"username"`
	Password  string `json:"password"`
	DOB       string `json:"dob"`
	Gender    string `json:"gender"`
	Bio       string `json:"bio"`
	Photo     string `json:"photo"`
	Privacy   string `json:"privacy"`
}
type Book struct {
	ID          int    `json:"id"`
	Title       string `json:"title"`
	Description string `json:"description"`
	Image       string `json:"image"`
	Username    string `json:"username"`
}
type Like struct {
	ID       int    `json:"id"`
	Username string `json:"username"`
	BookID   int    `json:"book_id"`
}
type Comment struct {
	ID       int    `json:"id"`
	Username string `json:"username"`
	BookID   int    `json:"book_id"`
	Content  string `json:"content"`
}

//! MODELLER END ----------------------------------------------------------------------------------

var db *sql.DB

func initDB() {
	var err error
	db, err = sql.Open("sqlite", "./users.db")
	if err != nil {
		log.Fatal("Veritabanı bağlantı hatası:", err)
	}

	//! TABLOLARIN OLUŞTURULMASI
	// Users tablosu
	createUsersTable := `CREATE TABLE IF NOT EXISTS users (
		"id" INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
		"firstname" TEXT,
		"lastname" TEXT,
		"email" TEXT UNIQUE,
		"username" TEXT UNIQUE,
		"password" TEXT,
		"dob" TEXT,
		"gender" TEXT,
		"bio" TEXT,
		"photo" TEXT
	);`
	_, err = db.Exec(createUsersTable)
	if err != nil {
		log.Fatal("Users tablosunu oluşturma hatası:", err)
	}

	// Books tablosu
	createBooksTable := `CREATE TABLE IF NOT EXISTS books (
		"id" INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
		"title" TEXT,
		"description" TEXT,
		"image" TEXT,
		"username" TEXT
	);`
	_, err = db.Exec(createBooksTable)
	if err != nil {
		log.Fatal("Books tablosunu oluşturma hatası:", err)
	}
	//username sütunu yoksa ekle
	_, err = db.Exec(`ALTER TABLE books ADD COLUMN username TEXT`)
	if err != nil {
		log.Println("Books tablosuna username eklenemedi (muhtemelen zaten var):", err)
	}

	// Follow tablosu
	createFollowTable := `CREATE TABLE IF NOT EXISTS follow (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		follower TEXT NOT NULL,
		following TEXT NOT NULL,
		status TEXT NOT NULL DEFAULT 'accepted',
		UNIQUE(follower, following)
	);`
	_, err = db.Exec(createFollowTable)
	if err != nil {
		log.Fatal("Follow tablosunu oluşturma hatası:", err)
	}
	_, _ = db.Exec(`ALTER TABLE follow ADD COLUMN status TEXT DEFAULT 'accepted'`)
	_, _ = db.Exec(`UPDATE follow SET status = 'accepted' WHERE status IS NULL`)

	//photo sütunu yoksa ekle
	_, _ = db.Exec(`ALTER TABLE users ADD COLUMN photo TEXT`)
	_, _ = db.Exec(`ALTER TABLE users ADD COLUMN privacy TEXT`)
	_, _ = db.Exec(`UPDATE users SET privacy = 'Herkes' WHERE privacy = 'Sadece Takipçiler'`)
	log.Println("Veritabanı başarıyla hazırlandı.")

	// Like tablosu
	createLikesTable := `CREATE TABLE IF NOT EXISTS likes(
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		username TEXT NOT NULL,
		book_id INTEGER NOT NULL, 
		UNIQUE(username, book_id)        
	);`
	_, err = db.Exec(createLikesTable)
	if err != nil {
		log.Fatal("Likes tablosunu oluşturma hatası:", err)
	}

	//Yorum tablosu
	createCommentTable := `CREATE TABLE IF NOT EXISTS comments(
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		username TEXT NOT NULL,        
		book_id INTEGER NOT NULL,
		content TEXT NOT NULL
	);`
	_, err = db.Exec(createCommentTable)
	if err != nil {
		log.Fatal("Comment tablosunu oluşturma hatası:", err)
	}

	// Mesajlaşma Tabloları
	createConversationsTable := `CREATE TABLE IF NOT EXISTS conversations(
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		type TEXT NOT NULL, 
		name TEXT,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP
	);`
	_, err = db.Exec(createConversationsTable)
	if err != nil {
		log.Fatal("Conversations tablosunu oluşturma hatası:", err)
	}

	createConversationMembersTable := `CREATE TABLE IF NOT EXISTS conversation_members(
		conversation_id INTEGER NOT NULL,
		username TEXT NOT NULL,
		last_read_time DATETIME DEFAULT CURRENT_TIMESTAMP,
		UNIQUE(conversation_id, username)
	);`
	_, err = db.Exec(createConversationMembersTable)
	if err != nil {
		log.Fatal("Conversation members tablosunu oluşturma hatası:", err)
	}

	createMessagesTable := `CREATE TABLE IF NOT EXISTS messages(
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		conversation_id INTEGER NOT NULL,
		sender_username TEXT NOT NULL,
		content TEXT NOT NULL,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP
	);`
	_, err = db.Exec(createMessagesTable)
	if err != nil {
		log.Fatal("Messages tablosunu oluşturma hatası:", err)
	}
}

//! TABLOLARIN OLUŞTURULMASI END ----------------------------------------------------------------------------------

// --- Yardımcı Fonksiyonlar ---
func respondJSON(w http.ResponseWriter, status int, payload interface{}) {
	response, err := json.Marshal(payload)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	w.Write(response)
}

func enableCORS(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}
		next.ServeHTTP(w, r)
	})
}

func canViewProfile(targetUsername, currentUsername string) (bool, error) {
	if targetUsername == currentUsername {
		return true, nil
	}

	var privacy string
	err := db.QueryRow("SELECT COALESCE(privacy, 'Herkes') FROM users WHERE username = ?", targetUsername).Scan(&privacy)
	if err != nil {
		return false, err
	}

	switch privacy {
	case "Herkes":
		return true, nil
	case "Gizli Profil":
		var count int
		err := db.QueryRow("SELECT COUNT(*) FROM follow WHERE follower = ? AND following = ? AND status = 'accepted'", currentUsername, targetUsername).Scan(&count)
		if err != nil {
			return false, err
		}
		return count > 0, nil
	default:
		return false, nil
	}
}

// ! REGISTER - LOGIN
func registerHandler(w http.ResponseWriter, r *http.Request) {
	var user User
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		respondJSON(w, http.StatusBadRequest, map[string]string{"error": "Geçersiz istek formatı"})
		return
	}
	defer r.Body.Close()

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		respondJSON(w, http.StatusInternalServerError, map[string]string{"error": "Şifre oluşturulamadı"})
		return
	}

	insertSQL := `INSERT INTO users(firstname, lastname, email, username, password, dob, gender, bio) VALUES (?, ?, ?, ?, ?, ?, ?, ?)`
	_, err = db.Exec(insertSQL, user.FirstName, user.LastName, user.Email, user.Username, string(hashedPassword), user.DOB, user.Gender, user.Bio)
	if err != nil {
		respondJSON(w, http.StatusBadRequest, map[string]string{"error": "Bu kullanıcı adı veya e-posta zaten kullanılıyor."})
		return
	}

	respondJSON(w, http.StatusCreated, map[string]string{"message": "Kayıt başarıyla tamamlandı!"})
}

type LoginCredentials struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func loginHandler(w http.ResponseWriter, r *http.Request) {
	var creds LoginCredentials
	err := json.NewDecoder(r.Body).Decode(&creds)
	if err != nil {
		respondJSON(w, http.StatusBadRequest, map[string]string{"error": "Geçersiz istek"})
		return
	}
	defer r.Body.Close()

	var storedPasswordHash string
	err = db.QueryRow("SELECT password FROM users WHERE username = ?", creds.Username).Scan(&storedPasswordHash)
	if err != nil {
		if err == sql.ErrNoRows {
			respondJSON(w, http.StatusUnauthorized, map[string]string{"error": "Kullanıcı adı veya şifre hatalı"})
			return
		}
		respondJSON(w, http.StatusInternalServerError, map[string]string{"error": "Veritabanı hatası"})
		return
	}

	err = bcrypt.CompareHashAndPassword([]byte(storedPasswordHash), []byte(creds.Password))
	if err != nil {
		respondJSON(w, http.StatusUnauthorized, map[string]string{"error": "Kullanıcı adı veya şifre hatalı"})
		return
	}

	respondJSON(w, http.StatusOK, map[string]string{"message": "Giriş başarılı!"})
}

//! REGISTER - LOGIN END ----------------------------------------------------------------------------------

// ! PROFILE GET - UPDATE
func profileHandler(w http.ResponseWriter, r *http.Request) {
	username := r.URL.Query().Get("username")
	currentUser := r.URL.Query().Get("currentUser")
	log.Println("Requested username:", username, "currentUser:", currentUser)

	allowed, err := canViewProfile(username, currentUser)
	if err != nil {
		respondJSON(w, http.StatusInternalServerError, map[string]string{"error": "Veritabanı hatası"})
		return
	}
	if !allowed {
		respondJSON(w, http.StatusForbidden, map[string]string{"error": "Bu profili görüntüleme yetkiniz yok"})
		return
	}

	var user User
	err = db.QueryRow(`SELECT firstname, lastname, username, dob, gender, bio, photo, COALESCE(privacy, 'Herkes') FROM users WHERE username=?`, username).
		Scan(&user.FirstName, &user.LastName, &user.Username, &user.DOB, &user.Gender, &user.Bio, &user.Photo, &user.Privacy)
	if err != nil {
		respondJSON(w, http.StatusNotFound, map[string]string{"error": "Kullanıcı bulunamadı"})
		return
	}

	if user.Bio == "" {
		user.Bio = ""
	}

	log.Println("Fetched user:", user.Username, " - ", user.FirstName)
	respondJSON(w, http.StatusOK, user)
}

func updateProfileHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPut {
		http.Error(w, "only PUT allowed", http.StatusMethodNotAllowed)
		return
	}
	var user User
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	_, err := db.Exec(
		"UPDATE users SET firstname = ?, lastname = ?, bio = ?, photo = ?, privacy = ? WHERE username = ?",
		user.FirstName, user.LastName, user.Bio, user.Photo, user.Privacy, user.Username,
	)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	respondJSON(w, http.StatusOK, map[string]string{"message": "Profil başarıyla güncellendi!"})
}

func changeAccountSettingsHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPut {
		http.Error(w, "only PUT allowed", http.StatusMethodNotAllowed)
		return
	}

	var req struct {
		CurrentUsername string `json:"currentUsername"`
		NewUsername     string `json:"newUsername"`
		OldPassword     string `json:"oldPassword"`
		NewPassword     string `json:"newPassword"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// First verify old password if password is being changed
	if req.NewPassword != "" {
		if req.OldPassword == "" {
			http.Error(w, "Eski şifre zorunludur", http.StatusBadRequest)
			return
		}
		var storedHash string
		err := db.QueryRow("SELECT password FROM users WHERE username = ?", req.CurrentUsername).Scan(&storedHash)
		if err != nil {
			http.Error(w, "Kullanıcı bulunamadı", http.StatusInternalServerError)
			return
		}
		if err := bcrypt.CompareHashAndPassword([]byte(storedHash), []byte(req.OldPassword)); err != nil {
			http.Error(w, "Eski şifre yanlış", http.StatusUnauthorized)
			return
		}
		
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.NewPassword), bcrypt.DefaultCost)
		if err != nil {
			http.Error(w, "Şifre oluşturulamadı", http.StatusInternalServerError)
			return
		}
		_, err = db.Exec("UPDATE users SET password = ? WHERE username = ?", string(hashedPassword), req.CurrentUsername)
		if err != nil {
			http.Error(w, "Şifre güncellenemedi", http.StatusInternalServerError)
			return
		}
	}

	// Change username if requested
	if req.NewUsername != "" && req.NewUsername != req.CurrentUsername {
		var exists int
		db.QueryRow("SELECT COUNT(*) FROM users WHERE username = ?", req.NewUsername).Scan(&exists)
		if exists > 0 {
			http.Error(w, "Bu kullanıcı adı zaten alınmış", http.StatusConflict)
			return
		}

		tx, err := db.Begin()
		if err != nil {
			http.Error(w, "Transaction hatası", http.StatusInternalServerError)
			return
		}

		queries := []struct {
			Query string
			Args  []interface{}
		}{
			{"UPDATE users SET username = ? WHERE username = ?", []interface{}{req.NewUsername, req.CurrentUsername}},
			{"UPDATE books SET username = ? WHERE username = ?", []interface{}{req.NewUsername, req.CurrentUsername}},
			{"UPDATE likes SET username = ? WHERE username = ?", []interface{}{req.NewUsername, req.CurrentUsername}},
			{"UPDATE comments SET username = ? WHERE username = ?", []interface{}{req.NewUsername, req.CurrentUsername}},
			{"UPDATE follow SET follower = ? WHERE follower = ?", []interface{}{req.NewUsername, req.CurrentUsername}},
			{"UPDATE follow SET following = ? WHERE following = ?", []interface{}{req.NewUsername, req.CurrentUsername}},
		}

		for _, q := range queries {
			if _, err := tx.Exec(q.Query, q.Args...); err != nil {
				tx.Rollback()
				http.Error(w, "Güncelleme sırasında hata: "+err.Error(), http.StatusInternalServerError)
				return
			}
		}

		if err := tx.Commit(); err != nil {
			http.Error(w, "Transaction commit hatası", http.StatusInternalServerError)
			return
		}
	}

	respondJSON(w, http.StatusOK, map[string]string{"message": "Hesap ayarları başarıyla güncellendi!"})
}

//! PROFILE GET - UPDATE END ----------------------------------------------------------------------------------

// ! BOOKS ADD - GET - DELETE - UPDATE
func addBookHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "only POST allowed", http.StatusMethodNotAllowed)
		return
	}

	var b Book
	if err := json.NewDecoder(r.Body).Decode(&b); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	res, err := db.Exec("INSERT INTO books (title, description, image, username) VALUES (?, ?, ?, ?)", b.Title, b.Description, b.Image, b.Username)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	lastID, _ := res.LastInsertId()
	b.ID = int(lastID)
	respondJSON(w, http.StatusCreated, b)
}

func getBooksHandler(w http.ResponseWriter, r *http.Request) {
	username := r.URL.Query().Get("username")
	currentUser := r.URL.Query().Get("currentUser")
	if username == "" {
		http.Error(w, "username is required", http.StatusBadRequest)
		return
	}

	allowed, err := canViewProfile(username, currentUser)
	if err != nil {
		http.Error(w, "Veritabanı hatası", http.StatusInternalServerError)
		return
	}
	if !allowed {
		http.Error(w, "Bu profili görüntüleme yetkiniz yok", http.StatusForbidden)
		return
	}

	rows, err := db.Query("SELECT id, title, description, image FROM books WHERE username = ?", username)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var books []Book
	for rows.Next() {
		var b Book
		if err := rows.Scan(&b.ID, &b.Title, &b.Description, &b.Image); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		books = append(books, b)
	}
	if err := rows.Err(); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	respondJSON(w, http.StatusOK, books)
}

func deleteBookHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		http.Error(w, "only DELETE allowed", http.StatusMethodNotAllowed)
		return
	}

	var b Book
	if err := json.NewDecoder(r.Body).Decode(&b); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	_, err := db.Exec("DELETE FROM books WHERE id = ? AND username = ?", b.ID, b.Username)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	respondJSON(w, http.StatusOK, map[string]string{"message": "Kitap silindi"})
}
func updateBookHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPut {
		http.Error(w, "only PUT allowed", http.StatusMethodNotAllowed)
		return
	}

	var b Book
	if err := json.NewDecoder(r.Body).Decode(&b); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	_, err := db.Exec("UPDATE books SET title = ?, description = ?, image = ? WHERE id = ? AND username = ?", b.Title, b.Description, b.Image, b.ID, b.Username)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	respondJSON(w, http.StatusOK, map[string]string{"message": "Kitap güncellendi"})
}

//! BOOKS ADD - GET - DELETE - UPDATE END ----------------------------------------------------------------------------------

func findUsersHandler(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query().Get("q")
	currentUser := r.URL.Query().Get("current")

	rows, err := db.Query(`SELECT username, firstname, photo, COALESCE(privacy, 'Herkes') FROM users WHERE username LIKE ?`, "%"+query+"%")
	if err != nil {
		respondJSON(w, http.StatusInternalServerError, map[string]string{"error": "DB error"})
		return
	}
	defer rows.Close()

	var users []map[string]interface{}
	for rows.Next() {
		var username, firstname, photo, privacy string
		rows.Scan(&username, &firstname, &photo, &privacy)
		var acceptedCount, pendingCount int
		db.QueryRow(`SELECT COUNT(*) FROM follow WHERE follower=? AND following=? AND status='accepted'`, currentUser, username).Scan(&acceptedCount)
		db.QueryRow(`SELECT COUNT(*) FROM follow WHERE follower=? AND following=? AND status='pending'`, currentUser, username).Scan(&pendingCount)

		users = append(users, map[string]interface{}{
			"username":    username,
			"firstname":   firstname,
			"photo":       photo,
			"privacy":     privacy,
			"isFollowing": acceptedCount > 0,
			"isPending":   pendingCount > 0,
		})
	}
	respondJSON(w, http.StatusOK, users)
}

// ! FOLLOW - UNFOLLOW
func followUserHandler(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Follower  string `json:"follower"`
		Following string `json:"following"`
		Action    string `json:"action"`
	}
	json.NewDecoder(r.Body).Decode(&req)

	if req.Action == "follow" {
		var privacy string
		err := db.QueryRow(`SELECT COALESCE(privacy, 'Herkes') FROM users WHERE username = ?`, req.Following).Scan(&privacy)
		if err != nil {
			respondJSON(w, 500, map[string]string{"error": "DB error"})
			return
		}

		if privacy == "Gizli Profil" {
			_, err = db.Exec(`INSERT OR IGNORE INTO follow(follower, following, status) VALUES(?,?,?)`, req.Follower, req.Following, "pending")
		} else {
			_, err = db.Exec(`INSERT OR IGNORE INTO follow(follower, following, status) VALUES(?,?,?)`, req.Follower, req.Following, "accepted")
		}
		if err != nil {
			respondJSON(w, 500, map[string]string{"error": "DB error"})
			return
		}
	} else if req.Action == "approve" {
		_, err := db.Exec(`UPDATE follow SET status = 'accepted' WHERE follower = ? AND following = ? AND status = 'pending'`, req.Follower, req.Following)
		if err != nil {
			respondJSON(w, 500, map[string]string{"error": "DB error"})
			return
		}
	} else if req.Action == "reject" {
		_, err := db.Exec(`DELETE FROM follow WHERE follower = ? AND following = ? AND status = 'pending'`, req.Follower, req.Following)
		if err != nil {
			respondJSON(w, 500, map[string]string{"error": "DB error"})
			return
		}
	} else {
		_, err := db.Exec(`DELETE FROM follow WHERE follower=? AND following=?`, req.Follower, req.Following)
		if err != nil {
			respondJSON(w, 500, map[string]string{"error": "DB error"})
			return
		}
	}
	respondJSON(w, 200, map[string]string{"message": "success"})
}

func getFollowersHandler(w http.ResponseWriter, r *http.Request) {
	username := r.URL.Query().Get("username")
	currentUser := r.URL.Query().Get("currentUser")
	allowed, err := canViewProfile(username, currentUser)
	if err != nil {
		http.Error(w, "Veritabanı hatası", http.StatusInternalServerError)
		return
	}
	if !allowed {
		http.Error(w, "Bu profili görüntüleme yetkiniz yok", http.StatusForbidden)
		return
	}

	var followersCount, followingCount int
	db.QueryRow("SELECT COUNT(*) FROM follow WHERE following = ? AND follower != ? AND status = 'accepted'", username, username).Scan(&followersCount)
	db.QueryRow("SELECT COUNT(*) FROM follow WHERE follower = ? AND following != ? AND status = 'accepted'", username, username).Scan(&followingCount)

	respondJSON(w, http.StatusOK, map[string]int{
		"followers": followersCount,
		"following": followingCount,
	})
}

//! FOLLOW - UNFOLLOW END ----------------------------------------------------------------------------------

func getAllBooksHandler(w http.ResponseWriter, r *http.Request) {
	currentUser := r.URL.Query().Get("currentUser")

	query := `
        SELECT 
		b.id, b.title, b.description, b.image, b.username, 
		COALESCE(u.photo, ''),
		(SELECT COUNT(*) FROM likes WHERE book_id=b.id) as like_count,
		(SELECT COUNT(*) FROM likes WHERE book_id=b.id AND username=?) as is_liked,
		(SELECT COUNT(*) FROM comments WHERE book_id=b.id) as comment_count
        FROM books b 
        LEFT JOIN users u ON b.username = u.username
        LEFT JOIN follow f ON f.follower = ? AND f.following = b.username AND f.status = 'accepted'
        WHERE COALESCE(u.privacy, 'Herkes') = 'Herkes' OR b.username = ? OR f.id IS NOT NULL`

	rows, err := db.Query(query, currentUser, currentUser, currentUser)
	if err != nil {
		respondJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}
	defer rows.Close()

	var books []map[string]interface{}
	for rows.Next() {
		var id, likeCount, isLiked, commentCount int
		var title, description, image, username, userPhoto string
		if err := rows.Scan(&id, &title, &description, &image, &username, &userPhoto, &likeCount, &isLiked, &commentCount); err != nil {
			continue
		}
		books = append(books, map[string]interface{}{
			"id":            id,
			"title":         title,
			"description":   description,
			"image":         image,
			"username":      username,
			"user_photo":    userPhoto,
			"like_count":    likeCount,
			"is_liked":      isLiked > 0,
			"comment_count": commentCount,
		})
	}
	respondJSON(w, http.StatusOK, books)
}

func getNotificationsHandler(w http.ResponseWriter, r *http.Request) {
	currentUser := r.URL.Query().Get("username")

	rows, err := db.Query(`SELECT follower, status FROM follow WHERE following = ?`, currentUser)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var pending []map[string]string
	var accepted []map[string]string
	for rows.Next() {
		var follower, status string
		if err := rows.Scan(&follower, &status); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		if status == "pending" {
			pending = append(pending, map[string]string{"follower": follower, "message": follower + " seni takip etmek istiyor."})
		} else {
			accepted = append(accepted, map[string]string{"follower": follower, "message": follower + " seni takip etti."})
		}
	}

	respondJSON(w, http.StatusOK, map[string]interface{}{
		"pendingRequests": pending,
		"acceptedFollows": accepted,
	})
}

// ! GET FOLLOWERS - FOLLOWING
func checkFollowStatusHandler(w http.ResponseWriter, r *http.Request) {
	follower := r.URL.Query().Get("follower")
	following := r.URL.Query().Get("following")

	var status string
	err := db.QueryRow("SELECT status FROM follow WHERE follower = ? AND following = ?", follower, following).Scan(&status)
	if err != nil && err != sql.ErrNoRows {
		respondJSON(w, http.StatusInternalServerError, map[string]string{"error": "DB error"})
		return
	}

	respondJSON(w, http.StatusOK, map[string]bool{
		"isFollowing": status == "accepted",
		"isPending":   status == "pending",
	})
}

func getFollowersListHandler(w http.ResponseWriter, r *http.Request) {
	username := r.URL.Query().Get("username")
	currentUser := r.URL.Query().Get("currentUser")
	allowed, err := canViewProfile(username, currentUser)
	if err != nil {
		http.Error(w, "Veritabanı hatası", http.StatusInternalServerError)
		return
	}
	if !allowed {
		http.Error(w, "Bu profili görüntüleme yetkiniz yok", http.StatusForbidden)
		return
	}

	rows, err := db.Query(`
        SELECT f.follower, COALESCE(u.photo, '') 
        FROM follow f
        LEFT JOIN users u ON f.follower = u.username
        WHERE f.following = ? AND f.status = 'accepted'`, username)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	defer rows.Close()

	var followers = []map[string]string{}
	for rows.Next() {
		var uname, photo string
		if err := rows.Scan(&uname, &photo); err != nil {
			http.Error(w, err.Error(), 500)
			return
		}
		followers = append(followers, map[string]string{
			"username": uname,
			"photo":    photo,
		})
	}
	respondJSON(w, 200, followers)
}

func getFollowingListHandler(w http.ResponseWriter, r *http.Request) {
	username := r.URL.Query().Get("username")
	currentUser := r.URL.Query().Get("currentUser")
	allowed, err := canViewProfile(username, currentUser)
	if err != nil {
		http.Error(w, "Veritabanı hatası", http.StatusInternalServerError)
		return
	}
	if !allowed {
		http.Error(w, "Bu profili görüntüleme yetkiniz yok", http.StatusForbidden)
		return
	}

	rows, err := db.Query(`
        SELECT f.following, COALESCE(u.photo, '') 
        FROM follow f
        LEFT JOIN users u ON f.following = u.username
        WHERE f.follower = ? AND f.status = 'accepted'`, username)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	defer rows.Close()

	var following = []map[string]string{}
	for rows.Next() {
		var uname, photo string
		if err := rows.Scan(&uname, &photo); err != nil {
			http.Error(w, err.Error(), 500)
			return
		}
		following = append(following, map[string]string{
			"username": uname,
			"photo":    photo,
		})
	}
	respondJSON(w, 200, following)
}

//! GET FOLLOWERS - FOLLOWING END ----------------------------------------------------------------------------------

// ! PROFILE PHOTO
func uploadProfilePhotoHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "only POST allowed", http.StatusMethodNotAllowed)
		return
	}

	// limit parsed form size
	if err := r.ParseMultipartForm(10 << 20); err != nil {
		http.Error(w, "failed to parse multipart form", http.StatusBadRequest)
		return
	}

	username := r.FormValue("username")
	file, _, err := r.FormFile("photo")
	if err != nil {
		http.Error(w, "missing photo file", http.StatusBadRequest)
		return
	}
	defer file.Close()

	data, err := io.ReadAll(file)
	if err != nil {
		http.Error(w, "failed to read file", http.StatusInternalServerError)
		return
	}

	// detect content type
	contentType := http.DetectContentType(data)
	encoded := base64.StdEncoding.EncodeToString(data)
	dataURL := "data:" + contentType + ";base64," + encoded

	_, err = db.Exec("UPDATE users SET photo = ? WHERE username = ?", dataURL, username)
	if err != nil {
		http.Error(w, "failed to save photo", http.StatusInternalServerError)
		return
	}

	respondJSON(w, http.StatusOK, map[string]string{"photo": dataURL})
}

// removeProfilePhotoHandler clears the photo field for the given username.
func removeProfilePhotoHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "only POST allowed", http.StatusMethodNotAllowed)
		return
	}
	var req struct {
		Username string `json:"username"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request", http.StatusBadRequest)
		return
	}

	_, err := db.Exec("UPDATE users SET photo = '' WHERE username = ?", req.Username)
	if err != nil {
		http.Error(w, "failed to remove photo", http.StatusInternalServerError)
		return
	}

	respondJSON(w, http.StatusOK, map[string]string{"message": "photo removed"})
}

//! PROFILE PHOTO END ----------------------------------------------------------------------------------

// ! LIKES
func toggleLikeHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "only POST allowed", http.StatusMethodNotAllowed)
		return
	}
	var req struct {
		Username string `json:"username"`
		BookID   int    `json:"book_id"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request", http.StatusBadRequest)
		return
	}

	//önce silmeyi deniyoruz(eğer varsa beğeni kalkar)
	res, err := db.Exec("DELETE FROM likes WHERE username=? AND book_id=?", req.Username, req.BookID)
	if err != nil {
		http.Error(w, "Veritabanı hatası", http.StatusInternalServerError)
		return
	}
	rowsAffected, _ := res.RowsAffected()

	// eğer hiç silinmediyse yani beğeni yoktu, ekleyelim
	if rowsAffected == 0 {
		_, err := db.Exec("INSERT INTO likes(username, book_id) VALUES(?, ?)", req.Username, req.BookID)
		if err != nil {
			http.Error(w, "Beğeni eklenemedi", http.StatusInternalServerError)
			return
		}
		respondJSON(w, http.StatusOK, map[string]string{"message": "Kitap beğenildi"})
	} else {
		// zaten vardı, silindi
		respondJSON(w, http.StatusOK, map[string]string{"message": "Kitap beğenisi kaldırıldı"})
	}
}

//! LIKES END ----------------------------------------------------------------------------------

// ! COMMENTS
func addCommentHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "only POST allowd", http.StatusMethodNotAllowed)
		return
	}

	var c Comment //gelen json verisini karşılamak için bir değişken oluşturduk.

	err := json.NewDecoder(r.Body).Decode(&c) //json u oku ve c değişkenine doldur.
	if err != nil {
		http.Error(w, "invalid request", http.StatusBadRequest)
		return
	}

	if c.Content == "" { //boş yorum yapılmasını engelliyoruz.
		http.Error(w, "Yorum içeriği boş olamaz", http.StatusBadRequest)
		return
	}

	res, err := db.Exec("INSERT INTO comments(username, book_id, content) VALUES(?, ?, ?)", c.Username, c.BookID, c.Content)

	if err != nil {
		http.Error(w, "Veritabanı kayıt hatası", http.StatusInternalServerError)
		return
	}

	lastID, _ := res.LastInsertId() //eklenen yorumun id'sini alıyoruz.
	c.ID = int(lastID)

	respondJSON(w, http.StatusCreated, c) //oluşturulan yorumu json olarak döndürüyoruz.
}

func getCommentsHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "only GET allowed", http.StatusMethodNotAllowed)
		return
	}

	bookID := r.URL.Query().Get("book_id")
	if bookID == "" {
		http.Error(w, "book_id is required", http.StatusBadRequest)
		return
	}

	//veritabanından bu kitaba ait tüm yorumları çekelim
	rows, err := db.Query("SELECT id, username, book_id, content FROM comments WHERE book_id = ?", bookID)
	if err != nil {
		http.Error(w, "Veritabanı okuma hatası", http.StatusInternalServerError)
		return
	}

	defer rows.Close() //İş bitince bağlantıyı kapatıyoruz, yoksa bellek sızıntısı olur!

	//gelen satırları bir listeye dönüştürelim
	var comments []Comment
	for rows.Next() {
		var c Comment
		err := rows.Scan(&c.ID, &c.Username, &c.BookID, &c.Content)
		if err != nil {
			http.Error(w, "Veritabanı okuma hatası", http.StatusInternalServerError)
			return
		}
		comments = append(comments, c)
	}
	if comments == nil {
		comments = []Comment{} //eğer hiç yorum yoksa boş bir liste döndürelim, null değil.
	}
	respondJSON(w, http.StatusOK, comments) //yorum listesini json olarak döndürelim.

}

func deleteCommentHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		http.Error(w, "only DELETE allowed", http.StatusMethodNotAllowed)
		return
	}

	var req struct {
		CommentID int    `json:"comment_id"`
		Username  string `json:"username"`
	}

	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, "invalid request", http.StatusBadRequest)
		return
	}

	// Önce yorumun sahibini kontrol edeceğiz - sadece kendi yorumlarını silebilsin
	var commentOwner string
	err = db.QueryRow("SELECT username FROM comments WHERE id = ?", req.CommentID).Scan(&commentOwner)
	if err != nil {
		if err == sql.ErrNoRows {
			http.Error(w, "Yorum bulunamadı", http.StatusNotFound)
			return
		}
		http.Error(w, "Veritabanı hatası", http.StatusInternalServerError)
		return
	}

	// Sahibi kontrolü
	if commentOwner != req.Username {
		http.Error(w, "Bu yorumu silme yetkisine sahip değilsiniz", http.StatusForbidden)
		return
	}

	// Yorumu sil
	_, err = db.Exec("DELETE FROM comments WHERE id = ?", req.CommentID)
	if err != nil {
		http.Error(w, "Yorum silinirken hata oluştu", http.StatusInternalServerError)
		return
	}

	respondJSON(w, http.StatusOK, map[string]string{"message": "Yorum başarıyla silindi"})
}

//! COMMENTS END ---------------------------------------------------------------------------------------

func main() {
	initDB()

	http.Handle("/api/register", enableCORS(http.HandlerFunc(registerHandler)))
	http.Handle("/api/login", enableCORS(http.HandlerFunc(loginHandler)))
	http.Handle("/api/profile", enableCORS(http.HandlerFunc(profileHandler)))
	http.Handle("/api/updateProfile", enableCORS(http.HandlerFunc(updateProfileHandler)))
	http.Handle("/api/changeAccountSettings", enableCORS(http.HandlerFunc(changeAccountSettingsHandler)))

	// profile photo endpoints
	http.Handle("/api/uploadProfilePhoto", enableCORS(http.HandlerFunc(uploadProfilePhotoHandler)))
	http.Handle("/api/removeProfilePhoto", enableCORS(http.HandlerFunc(removeProfilePhotoHandler)))

	http.Handle("/api/addBook", enableCORS(http.HandlerFunc(addBookHandler)))
	http.Handle("/api/getBooks", enableCORS(http.HandlerFunc(getBooksHandler)))
	http.Handle("/api/deleteBook", enableCORS(http.HandlerFunc(deleteBookHandler)))
	http.Handle("/api/updateBook", enableCORS(http.HandlerFunc(updateBookHandler)))

	http.Handle("/api/findUsers", enableCORS(http.HandlerFunc(findUsersHandler)))
	http.Handle("/api/followUser", enableCORS(http.HandlerFunc(followUserHandler)))
	http.Handle("/api/getFollowers", enableCORS(http.HandlerFunc(getFollowersHandler)))
	http.Handle("/api/getFollowersList", enableCORS(http.HandlerFunc(getFollowersListHandler)))
	http.Handle("/api/getFollowingList", enableCORS(http.HandlerFunc(getFollowingListHandler)))
	http.Handle("/api/checkFollowStatus", enableCORS(http.HandlerFunc(checkFollowStatusHandler)))
	http.Handle("/api/getNotifications", enableCORS(http.HandlerFunc(getNotificationsHandler)))

	http.Handle("/api/getAllBooks", enableCORS(http.HandlerFunc(getAllBooksHandler)))

	http.Handle("/api/debugBooksColumns", enableCORS(http.HandlerFunc(debugBooksColumnsHandler)))

	http.Handle("/api/toggleLike", enableCORS(http.HandlerFunc(toggleLikeHandler)))

	http.Handle("/api/addComment", enableCORS(http.HandlerFunc(addCommentHandler)))
	http.Handle("/api/getComments", enableCORS(http.HandlerFunc(getCommentsHandler)))
	http.Handle("/api/deleteComment", enableCORS(http.HandlerFunc(deleteCommentHandler)))

	// Messaging endpoints
	http.Handle("/api/messages/startDirect", enableCORS(http.HandlerFunc(startDirectHandler)))
	http.Handle("/api/messages/createGroup", enableCORS(http.HandlerFunc(createGroupHandler)))
	http.Handle("/api/messages/conversations", enableCORS(http.HandlerFunc(getConversationsHandler)))
	http.Handle("/api/messages/history", enableCORS(http.HandlerFunc(getMessageHistoryHandler)))
	http.Handle("/api/messages/send", enableCORS(http.HandlerFunc(sendMessageHandler)))
	http.Handle("/api/messages/unreadCount", enableCORS(http.HandlerFunc(getUnreadCountHandler)))
	http.Handle("/api/messages/markRead", enableCORS(http.HandlerFunc(markReadHandler)))
	http.Handle("/api/messages/groupMembers", enableCORS(http.HandlerFunc(getGroupMembersHandler)))
	http.Handle("/api/messages/addGroupMember", enableCORS(http.HandlerFunc(addGroupMemberHandler)))
	http.Handle("/api/messages/removeGroupMember", enableCORS(http.HandlerFunc(removeGroupMemberHandler)))

	port := ":8000"
	log.Printf("Go API sunucusu http://localhost%s adresinde dinlemede...", port)
	log.Fatal(http.ListenAndServe(port, nil))

}

// ! MESSAGING HANDLERS
func startDirectHandler(w http.ResponseWriter, r *http.Request) {
	var req struct {
		CurrentUser string `json:"currentUser"`
		TargetUser  string `json:"targetUser"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Check if direct conversation already exists
	var convID int
	query := `
		SELECT c.id FROM conversations c
		JOIN conversation_members m1 ON c.id = m1.conversation_id
		JOIN conversation_members m2 ON c.id = m2.conversation_id
		WHERE c.type = 'direct' AND m1.username = ? AND m2.username = ?`
	err := db.QueryRow(query, req.CurrentUser, req.TargetUser).Scan(&convID)

	if err == sql.ErrNoRows {
		// Create new conversation
		res, err := db.Exec("INSERT INTO conversations(type) VALUES('direct')")
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		id, _ := res.LastInsertId()
		convID = int(id)
		
		db.Exec("INSERT INTO conversation_members(conversation_id, username) VALUES(?, ?)", convID, req.CurrentUser)
		db.Exec("INSERT INTO conversation_members(conversation_id, username) VALUES(?, ?)", convID, req.TargetUser)
	} else if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	respondJSON(w, http.StatusOK, map[string]int{"conversationId": convID})
}

func createGroupHandler(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Name        string   `json:"name"`
		CurrentUser string   `json:"currentUser"`
		Members     []string `json:"members"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	res, err := db.Exec("INSERT INTO conversations(type, name) VALUES('group', ?)", req.Name)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	id, _ := res.LastInsertId()
	convID := int(id)

	// Add members including the creator
	allMembers := append(req.Members, req.CurrentUser)
	for _, m := range allMembers {
		db.Exec("INSERT OR IGNORE INTO conversation_members(conversation_id, username) VALUES(?, ?)", convID, m)
	}

	respondJSON(w, http.StatusOK, map[string]int{"conversationId": convID})
}

type ConversationDTO struct {
	ID            int    `json:"id"`
	Type          string `json:"type"`
	Name          string `json:"name"`
	OtherUsername string `json:"otherUsername"` // for direct
	OtherPhoto    string `json:"otherPhoto"`
	LastMessage   string `json:"lastMessage"`
	LastMessageAt string `json:"lastMessageAt"`
	UnreadCount   int    `json:"unreadCount"`
}

func getConversationsHandler(w http.ResponseWriter, r *http.Request) {
	currentUser := r.URL.Query().Get("currentUser")
	if currentUser == "" {
		http.Error(w, "currentUser needed", http.StatusBadRequest)
		return
	}

	// Complex query to get conversations and the last message
	query := `
		SELECT 
			c.id, c.type, COALESCE(c.name, ''),
			(SELECT username FROM conversation_members WHERE conversation_id = c.id AND username != ? LIMIT 1) as other_user,
			(SELECT COALESCE(photo, '') FROM users WHERE username = (SELECT username FROM conversation_members WHERE conversation_id = c.id AND username != ? LIMIT 1)) as other_photo,
			(SELECT content FROM messages WHERE conversation_id = c.id ORDER BY created_at DESC LIMIT 1) as last_msg,
			(SELECT created_at FROM messages WHERE conversation_id = c.id ORDER BY created_at DESC LIMIT 1) as last_time,
			(SELECT COUNT(*) FROM messages m WHERE m.conversation_id = c.id AND m.created_at > cm.last_read_time) as unread_count
		FROM conversations c
		JOIN conversation_members cm ON c.id = cm.conversation_id
		WHERE cm.username = ?
		ORDER BY last_time DESC
	`
	rows, err := db.Query(query, currentUser, currentUser, currentUser)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var convs []ConversationDTO
	for rows.Next() {
		var c ConversationDTO
		var otherUser, otherPhoto, lastMsg, lastTime sql.NullString
		if err := rows.Scan(&c.ID, &c.Type, &c.Name, &otherUser, &otherPhoto, &lastMsg, &lastTime, &c.UnreadCount); err != nil {
			continue
		}
		c.OtherUsername = otherUser.String
		c.OtherPhoto = otherPhoto.String
		c.LastMessage = lastMsg.String
		c.LastMessageAt = lastTime.String

		if c.Type == "group" {
			c.OtherUsername = c.Name
		}
		convs = append(convs, c)
	}

	respondJSON(w, http.StatusOK, convs)
}

type MessageDTO struct {
	ID             int    `json:"id"`
	ConversationID int    `json:"conversationId"`
	SenderUsername string `json:"senderUsername"`
	Content        string `json:"content"`
	CreatedAt      string `json:"createdAt"`
}

func getMessageHistoryHandler(w http.ResponseWriter, r *http.Request) {
	convID := r.URL.Query().Get("conversationId")
	rows, err := db.Query("SELECT id, conversation_id, sender_username, content, created_at FROM messages WHERE conversation_id = ? ORDER BY created_at ASC", convID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var msgs []MessageDTO
	for rows.Next() {
		var m MessageDTO
		if err := rows.Scan(&m.ID, &m.ConversationID, &m.SenderUsername, &m.Content, &m.CreatedAt); err != nil {
			continue
		}
		msgs = append(msgs, m)
	}
	respondJSON(w, http.StatusOK, msgs)
}

func sendMessageHandler(w http.ResponseWriter, r *http.Request) {
	var req struct {
		ConversationID int    `json:"conversationId"`
		SenderUsername string `json:"senderUsername"`
		Content        string `json:"content"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	res, err := db.Exec("INSERT INTO messages(conversation_id, sender_username, content) VALUES(?, ?, ?)", req.ConversationID, req.SenderUsername, req.Content)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	id, _ := res.LastInsertId()

	// Update sender's read time so their own message isn't unread
	db.Exec("UPDATE conversation_members SET last_read_time = CURRENT_TIMESTAMP WHERE conversation_id = ? AND username = ?", req.ConversationID, req.SenderUsername)

	respondJSON(w, http.StatusOK, map[string]int{"messageId": int(id)})
}

func getUnreadCountHandler(w http.ResponseWriter, r *http.Request) {
	username := r.URL.Query().Get("currentUser")
	if username == "" {
		respondJSON(w, http.StatusOK, map[string]int{"count": 0})
		return
	}

	var count int
	query := `
		SELECT COUNT(*) FROM messages m
		JOIN conversation_members cm ON m.conversation_id = cm.conversation_id
		WHERE cm.username = ? AND m.created_at > cm.last_read_time
	`
	err := db.QueryRow(query, username).Scan(&count)
	if err != nil {
		count = 0
	}
	respondJSON(w, http.StatusOK, map[string]int{"count": count})
}

func markReadHandler(w http.ResponseWriter, r *http.Request) {
	var req struct {
		ConversationID int    `json:"conversationId"`
		CurrentUser    string `json:"currentUser"`
	}
	json.NewDecoder(r.Body).Decode(&req)

	db.Exec("UPDATE conversation_members SET last_read_time = CURRENT_TIMESTAMP WHERE conversation_id = ? AND username = ?", req.ConversationID, req.CurrentUser)
	respondJSON(w, http.StatusOK, map[string]string{"status": "ok"})
}

type GroupMemberDTO struct {
	Username string `json:"username"`
	Photo    string `json:"photo"`
}

func getGroupMembersHandler(w http.ResponseWriter, r *http.Request) {
	convID := r.URL.Query().Get("conversationId")
	if convID == "" {
		http.Error(w, "conversationId required", http.StatusBadRequest)
		return
	}

	query := `
		SELECT cm.username, COALESCE(u.photo, '') 
		FROM conversation_members cm 
		JOIN users u ON cm.username = u.username 
		WHERE cm.conversation_id = ?`

	rows, err := db.Query(query, convID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var members []GroupMemberDTO
	for rows.Next() {
		var m GroupMemberDTO
		if err := rows.Scan(&m.Username, &m.Photo); err != nil {
			continue
		}
		members = append(members, m)
	}

	respondJSON(w, http.StatusOK, members)
}

func addGroupMemberHandler(w http.ResponseWriter, r *http.Request) {
	var req struct {
		ConversationID int    `json:"conversationId"`
		Username       string `json:"username"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Verify user exists first
	var exists int
	db.QueryRow("SELECT COUNT(*) FROM users WHERE username = ?", req.Username).Scan(&exists)
	if exists == 0 {
		http.Error(w, "Kullanıcı bulunamadı", http.StatusNotFound)
		return
	}

	_, err := db.Exec("INSERT INTO conversation_members(conversation_id, username) VALUES(?, ?)", req.ConversationID, req.Username)
	if err != nil {
		http.Error(w, "Bu kullanıcı zaten grupta veya eklenemedi", http.StatusInternalServerError)
		return
	}

	respondJSON(w, http.StatusOK, map[string]string{"status": "ok"})
}

func removeGroupMemberHandler(w http.ResponseWriter, r *http.Request) {
	var req struct {
		ConversationID int    `json:"conversationId"`
		Username       string `json:"username"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	_, err := db.Exec("DELETE FROM conversation_members WHERE conversation_id = ? AND username = ?", req.ConversationID, req.Username)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	respondJSON(w, http.StatusOK, map[string]string{"status": "ok"})
}

// bu kodu hata çıktığı için ekliyorum
func debugBooksColumnsHandler(w http.ResponseWriter, r *http.Request) {
	rows, err := db.Query("PRAGMA table_info(books)")
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	defer rows.Close()

	type Col struct {
		Cid     int    `json:"cid"`
		Name    string `json:"name"`
		Type    string `json:"type"`
		NotNull int    `json:"notnull"`
		Pk      int    `json:"pk"`
	}

	var cols []Col
	for rows.Next() {
		var c Col
		var dflt sql.NullString
		if err := rows.Scan(&c.Cid, &c.Name, &c.Type, &c.NotNull, &dflt, &c.Pk); err != nil {
			http.Error(w, err.Error(), 500)
			return
		}
		cols = append(cols, c)
	}

	respondJSON(w, 200, cols)
}
