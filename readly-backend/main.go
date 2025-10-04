package main

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"golang.org/x/crypto/bcrypt"
	_ "modernc.org/sqlite"
)

type User struct {
	FirstName string `json:"firstname"`
	LastName  string `json:"lastname"`
	Email     string `json:"email"`
	Username  string `json:"username"`
	Password  string `json:"password"`
	DOB       string `json:"dob"`
	Gender    string `json:"gender"`
	Bio       string `json:"bio"`
}
type Book struct {
	ID          int    `json:"id"`
	Title       string `json:"title"`
	Description string `json:"description"`
	Image       string `json:"image"`
	Username    string `json:"username"`
}

var db *sql.DB

func initDB() {
	var err error
	db, err = sql.Open("sqlite", "./users.db")
	if err != nil {
		log.Fatal("Veritabanı bağlantı hatası:", err)
	}
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
		"bio" TEXT
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

	createFollowTable:= `CREATE TABLE IF NOT EXISTS follow (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		follower TEXT NOT NULL,
		following TEXT NOT NULL,
		UNIQUE(follower, following)
	);`
	_, err = db.Exec(createFollowTable)
	if err != nil {
		log.Fatal("Follow tablosunu oluşturma hatası:", err)
	}


	log.Println("Veritabanı başarıyla hazırlandı.")

}

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

func profileHandler(w http.ResponseWriter, r *http.Request) {
	username := r.URL.Query().Get("username")
	log.Println("Requested username:", username)

	var user User
	err := db.QueryRow(`SELECT firstname, lastname, username, dob, gender, bio FROM users WHERE username=?`, username).
		Scan(&user.FirstName, &user.LastName, &user.Username, &user.DOB, &user.Gender, &user.Bio)
	if err != nil {
		log.Println("DB fetch error:", err)
		respondJSON(w, http.StatusNotFound, map[string]string{"error": "Kullanıcı bulunamadı"})
		return
	}

	if user.Bio == "" {
		user.Bio = ""
	}

	log.Println("Fetched user:", user)
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
		"UPDATE users SET firstname = ?, lastname = ?, bio = ? WHERE username = ?",
		user.FirstName, user.LastName, user.Bio, user.Username,
	)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	respondJSON(w, http.StatusOK, map[string]string{"message": "Profil başarıyla güncellendi!"})
}

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
	if username == "" {
		http.Error(w, "username is required", http.StatusBadRequest)
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

func findUsersHandler(w http.ResponseWriter, r *http.Request) {
    query := r.URL.Query().Get("q")
    currentUser := r.URL.Query().Get("current") 

    rows, err := db.Query(`SELECT username, firstname FROM users WHERE username LIKE ?`, "%"+query+"%")
    if err != nil {
        respondJSON(w, http.StatusInternalServerError, map[string]string{"error":"DB error"})
        return
    }
    defer rows.Close()

    var users []map[string]interface{}
    for rows.Next() {
        var username, firstname string
        rows.Scan(&username, &firstname)
        var count int
        db.QueryRow(`SELECT COUNT(*) FROM follow WHERE follower=? AND following=?`, currentUser, username).Scan(&count)

        users = append(users, map[string]interface{}{
            "username": username,
            "firstname": firstname,
            "isFollowing": count > 0,
        })
    }
    respondJSON(w, http.StatusOK, users)
}


func followUserHandler(w http.ResponseWriter, r *http.Request) {
    var req struct{
        Follower string `json:"follower"`
        Following string `json:"following"`
        Action string `json:"action"`
    }
    json.NewDecoder(r.Body).Decode(&req)

    if req.Action == "follow" {
        _, err := db.Exec(`INSERT OR IGNORE INTO follow(follower, following) VALUES(?,?)`, req.Follower, req.Following)
        if err != nil {
            respondJSON(w, 500, map[string]string{"error":"DB error"})
            return
        }
    } else {
        _, err := db.Exec(`DELETE FROM follow WHERE follower=? AND following=?`, req.Follower, req.Following)
        if err != nil {
            respondJSON(w, 500, map[string]string{"error":"DB error"})
            return
        }
    }
    respondJSON(w, 200, map[string]string{"message":"success"})
}


func getFollowersHandler(w http.ResponseWriter, r *http.Request) {
    username := r.URL.Query().Get("username")
    var followersCount, followingCount int

	db.QueryRow("SELECT COUNT(*) FROM follow WHERE following = ? AND follower != ?", username, username).Scan(&followersCount)
	db.QueryRow("SELECT COUNT(*) FROM follow WHERE follower = ? AND following != ?", username, username).Scan(&followingCount)


    respondJSON(w, http.StatusOK, map[string]int{
        "followers": followersCount,
        "following": followingCount,
    })
}


func getAllBooksHandler(w http.ResponseWriter, r *http.Request) {
    rows, err := db.Query("SELECT id, title, description, image, username FROM books")
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }
    defer rows.Close()

    var books []Book
    for rows.Next() {
        var b Book
        if err := rows.Scan(&b.ID, &b.Title, &b.Description, &b.Image, &b.Username); err != nil {
            http.Error(w, err.Error(), http.StatusInternalServerError)
            return
        }
        books = append(books, b)
    }
    respondJSON(w, http.StatusOK, books)
}


func getFollowersListHandler(w http.ResponseWriter, r *http.Request) {
    username := r.URL.Query().Get("username")

    rows, err := db.Query(`
        SELECT f.follower 
        FROM follow f
        WHERE f.following = ?`, username)
    if err != nil {
        http.Error(w, err.Error(), 500)
        return
    }
    defer rows.Close()

    var followers []string
    for rows.Next() {
        var uname string
        if err := rows.Scan(&uname); err != nil {
            http.Error(w, err.Error(), 500)
            return
        }
        followers = append(followers, uname)
    }
    respondJSON(w, 200, followers)
}

func getFollowingListHandler(w http.ResponseWriter, r *http.Request) {
    username := r.URL.Query().Get("username")

    rows, err := db.Query(`
        SELECT f.following 
        FROM follow f
        WHERE f.follower = ?`, username)
    if err != nil {
        http.Error(w, err.Error(), 500)
        return
    }
    defer rows.Close()

    var following []string
    for rows.Next() {
        var uname string
        if err := rows.Scan(&uname); err != nil {
            http.Error(w, err.Error(), 500)
            return
        }
        following = append(following, uname)
    }
    respondJSON(w, 200, following)
}


func main() {
	initDB()

	http.Handle("/api/register", enableCORS(http.HandlerFunc(registerHandler)))
	http.Handle("/api/login", enableCORS(http.HandlerFunc(loginHandler)))
	http.Handle("/api/profile", enableCORS(http.HandlerFunc(profileHandler)))
	http.Handle("/api/updateProfile", enableCORS(http.HandlerFunc(updateProfileHandler)))

	http.Handle("/api/addBook", enableCORS(http.HandlerFunc(addBookHandler)))
	http.Handle("/api/getBooks", enableCORS(http.HandlerFunc(getBooksHandler)))
	http.Handle("/api/deleteBook", enableCORS(http.HandlerFunc(deleteBookHandler)))
	http.Handle("/api/updateBook", enableCORS(http.HandlerFunc(updateBookHandler)))

	http.Handle("/api/findUsers", enableCORS(http.HandlerFunc(findUsersHandler)))
	http.Handle("/api/followUser", enableCORS(http.HandlerFunc(followUserHandler)))
	http.Handle("/api/getFollowers", enableCORS(http.HandlerFunc(getFollowersHandler)))
	http.Handle("/api/getFollowersList", enableCORS(http.HandlerFunc(getFollowersListHandler)))
	http.Handle("/api/getFollowingList", enableCORS(http.HandlerFunc(getFollowingListHandler)))


	http.Handle("/api/getAllBooks", enableCORS(http.HandlerFunc(getAllBooksHandler)))

	http.Handle("/api/debugBooksColumns", enableCORS(http.HandlerFunc(debugBooksColumnsHandler)))


	port := ":8000"
	log.Printf("Go API sunucusu http://localhost%s adresinde dinlemede...", port)
	log.Fatal(http.ListenAndServe(port, nil))

}


//bu kodu hata çıktığı için ekliyorum
func debugBooksColumnsHandler(w http.ResponseWriter, r *http.Request) {
	rows, err := db.Query("PRAGMA table_info(books)")
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	defer rows.Close()

	type Col struct {
		Cid   int    `json:"cid"`
		Name  string `json:"name"`
		Type  string `json:"type"`
		NotNull int `json:"notnull"`
		Pk    int    `json:"pk"`
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
