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

//! MODELLER
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
}
type Book struct {
	ID          int    `json:"id"`
	Title       string `json:"title"`
	Description string `json:"description"`
	Image       string `json:"image"`
	Username    string `json:"username"`
}
type Like struct{
	ID int `json:"id"`
	Username string `json:"username"`
	BookID int `json:"book_id"`
}
type Comment struct{
	ID int `json:"id"`
	Username string `json:"username"`
	BookID int `json:"book_id"`
	Content string `json:"content"`
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
		UNIQUE(follower, following)
	);`
	_, err = db.Exec(createFollowTable)
	if err != nil {
		log.Fatal("Follow tablosunu oluşturma hatası:", err)
	}


	//photo sütunu yoksa ekle
	_, _ = db.Exec(`ALTER TABLE users ADD COLUMN photo TEXT`)
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



//! REGISTER - LOGIN 
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






//! PROFILE GET - UPDATE
func profileHandler(w http.ResponseWriter, r *http.Request) {
	username := r.URL.Query().Get("username")
	log.Println("Requested username:", username)

	var user User
	err := db.QueryRow(`SELECT firstname, lastname, username, dob, gender, bio, photo FROM users WHERE username=?`, username).
		Scan(&user.FirstName, &user.LastName, &user.Username, &user.DOB, &user.Gender, &user.Bio, &user.Photo)
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
		"UPDATE users SET firstname = ?, lastname = ?, bio = ?, photo = ? WHERE username = ?",
		user.FirstName, user.LastName, user.Bio, user.Photo, user.Username,
	)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	respondJSON(w, http.StatusOK, map[string]string{"message": "Profil başarıyla güncellendi!"})
}
//! PROFILE GET - UPDATE END ----------------------------------------------------------------------------------





//! BOOKS ADD - GET - DELETE - UPDATE
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
//! BOOKS ADD - GET - DELETE - UPDATE END ----------------------------------------------------------------------------------







func findUsersHandler(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query().Get("q")
	currentUser := r.URL.Query().Get("current")

	rows, err := db.Query(`SELECT username, firstname, photo FROM users WHERE username LIKE ?`, "%"+query+"%")
	if err != nil {
		respondJSON(w, http.StatusInternalServerError, map[string]string{"error": "DB error"})
		return
	}
	defer rows.Close()

	var users []map[string]interface{}
	for rows.Next() {
		var username, firstname, photo string
		rows.Scan(&username, &firstname, &photo)
		var count int
		db.QueryRow(`SELECT COUNT(*) FROM follow WHERE follower=? AND following=?`, currentUser, username).Scan(&count)

		users = append(users, map[string]interface{}{
			"username":    username,
			"firstname":   firstname,
			"photo":   photo,
			"isFollowing": count > 0,
		})
	}
	respondJSON(w, http.StatusOK, users)
}



//! FOLLOW - UNFOLLOW
func followUserHandler(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Follower  string `json:"follower"`
		Following string `json:"following"`
		Action    string `json:"action"`
	}
	json.NewDecoder(r.Body).Decode(&req)

	if req.Action == "follow" {
		_, err := db.Exec(`INSERT OR IGNORE INTO follow(follower, following) VALUES(?,?)`, req.Follower, req.Following)
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
	var followersCount, followingCount int

	db.QueryRow("SELECT COUNT(*) FROM follow WHERE following = ? AND follower != ?", username, username).Scan(&followersCount)
	db.QueryRow("SELECT COUNT(*) FROM follow WHERE follower = ? AND following != ?", username, username).Scan(&followingCount)

	respondJSON(w, http.StatusOK, map[string]int{
		"followers": followersCount,
		"following": followingCount,
	})
}
//! FOLLOW - UNFOLLOW END ----------------------------------------------------------------------------------






func getAllBooksHandler(w http.ResponseWriter, r *http.Request) {
	//frınetndden giriş yapmış kullanıcının adını alıyoruz.
	currentUser := r.URL.Query().Get("currentUser")

	query := `
        SELECT 
			b.id, b.title, b.description, b.image, b.username, 
			COALESCE(u.photo, ''),
			(SELECT COUNT(*) FROM likes WHERE book_id=b.id) as like_count,
			(SELECT COUNT(*) FROM likes WHERE book_id=b.id AND username=?) as is_liked,
			(SELECT COUNT(*) FROM comments WHERE book_id=b.id) as comment_count
        FROM books b 
        LEFT JOIN users u ON b.username = u.username`

	rows, err := db.Query(query, currentUser)

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
            "id":          id,
            "title":       title,
            "description": description,
            "image":       image,
            "username":    username,
            "user_photo":  userPhoto, 
			"like_count":  likeCount,
			"is_liked":    isLiked > 0,
			"comment_count": commentCount,
        })
    }
	respondJSON(w, http.StatusOK, books)
}





//! GET FOLLOWERS - FOLLOWING
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
//! GET FOLLOWERS - FOLLOWING END ----------------------------------------------------------------------------------






//! PROFILE PHOTO
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






//! LIKES
func toggleLikeHandler(w http.ResponseWriter, r *http.Request){
	if r.Method != http.MethodPost {
        http.Error(w, "only POST allowed", http.StatusMethodNotAllowed)
        return
    }
	var req struct{
		Username string `json:"username"`
		BookID int `json:"book_id"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil{
		http.Error(w, "invalid request", http.StatusBadRequest)
		return
	}

	//önce silmeyi deniyoruz(eğer varsa beğeni kalkar)
	res, err:=db.Exec("DELETE FROM likes WHERE username=? AND book_id=?", req.Username, req.BookID)
	if err!=nil{
		http.Error(w,"Veritabanı hatası", http.StatusInternalServerError)
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
	}else{
		// zaten vardı, silindi
		respondJSON(w, http.StatusOK, map[string]string{"message": "Kitap beğenisi kaldırıldı"})
	}
}
//! LIKES END ----------------------------------------------------------------------------------






//! COMMENTS 
func addCommentHandler(w http.ResponseWriter, r *http.Request){
	if r.Method != http.MethodPost {
		http.Error(w, "only POST allowd", http.StatusMethodNotAllowed)
		return
	}

	var c Comment  //gelen json verisini karşılamak için bir değişken oluşturduk.

	err := json.NewDecoder(r.Body).Decode(&c) //json u oku ve c değişkenine doldur.
	if err != nil {
		http.Error(w, "invalid request", http.StatusBadRequest)
		return
	}

	if c.Content == "" {         //boş yorum yapılmasını engelliyoruz.
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


func getCommentsHandler(w http.ResponseWriter, r *http.Request){
	if(r.Method != http.MethodGet){
		http.Error(w, "only GET allowed", http.StatusMethodNotAllowed)
		return
	}

	bookID := r.URL.Query().Get("book_id")
	if bookID == ""{
		http.Error(w, "book_id is required", http.StatusBadRequest)
		return
	}

	//veritabanından bu kitaba ait tüm yorumları çekelim
	rows, err := db.Query("SELECT id, username, book_id, content FROM comments WHERE book_id = ?", bookID)
	if err != nil {
		http.Error(w, "Veritabanı okuma hatası", http.StatusInternalServerError)
		return
	}

	defer rows.Close()  //İş bitince bağlantıyı kapatıyoruz, yoksa bellek sızıntısı olur!

	//gelen satırları bir listeye dönüştürelim
	var comments []Comment
	for rows.Next(){
		var c Comment
		err := rows.Scan(&c.ID, &c.Username, &c.BookID, &c.Content)
		if err!=nil{
			http.Error(w, "Veritabanı okuma hatası", http.StatusInternalServerError)
			return
		}
		comments = append(comments,c)
	}
	if comments == nil {
		comments = []Comment{}  //eğer hiç yorum yoksa boş bir liste döndürelim, null değil.
	}
	respondJSON(w, http.StatusOK, comments) //yorum listesini json olarak döndürelim.

}
//! COMMENTS END ---------------------------------------------------------------------------------------



func main() {
	initDB()

	http.Handle("/api/register", enableCORS(http.HandlerFunc(registerHandler)))
	http.Handle("/api/login", enableCORS(http.HandlerFunc(loginHandler)))
	http.Handle("/api/profile", enableCORS(http.HandlerFunc(profileHandler)))
	http.Handle("/api/updateProfile", enableCORS(http.HandlerFunc(updateProfileHandler)))

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

	http.Handle("/api/getAllBooks", enableCORS(http.HandlerFunc(getAllBooksHandler)))

	http.Handle("/api/debugBooksColumns", enableCORS(http.HandlerFunc(debugBooksColumnsHandler)))

	http.Handle("/api/toggleLike", enableCORS(http.HandlerFunc(toggleLikeHandler)))

	http.Handle("/api/addComment", enableCORS(http.HandlerFunc(addCommentHandler)))
	http.Handle("/api/getComments", enableCORS(http.HandlerFunc(getCommentsHandler)))


	port := ":8000"
	log.Printf("Go API sunucusu http://localhost%s adresinde dinlemede...", port)
	log.Fatal(http.ListenAndServe(port, nil))

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
