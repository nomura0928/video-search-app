package main

import (
	"database/sql"  // データベースの操作
	"encoding/json" // JSONのエンコードとデコード
	"fmt"
	"io"
	"log"
	"net/http" // HTTPでWebサーバーを立てる
	"os"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"

	_ "modernc.org/sqlite"
)

type DBInfo struct {
	ImdbID string `json:"imdbID"`
	Title  string `json:"Title"`
	Year   string `json:"Year"`
	Poster string `json:"Poster"`
}

type User struct {
    UserID   string `json:"user_id"`
    Password string `json:"password"`
}


func main() {
	db, err := sql.Open("sqlite", "./favoritelist.db")
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS favoritelist (
	imdbID PRIMARY KEY NOT NULL,
	Title TEXT NOT NULL,
	Year TEXT NOT NULL,
	Poster TEXT NOT NULL
	)`)
	if err != nil {
		log.Fatalf("Failed to create table: %v", err)
	}
	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS users (
	id INTEGER PRIMARY KEY AUTOINCREMENT,
	user_id TEXT NOT NULL UNIQUE,
	password TEXT NOT NULL
	)`)
	if err != nil {
		log.Fatalf("Failed to create table: %v", err)
	}

	http.HandleFunc("/favorites", func(w http.ResponseWriter, r *http.Request) {
		//すべてのリクエストを許可
		w.Header().Set("Access-Control-Allow-Origin","*")
		//許可するメソッドを指定
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, DELETE")
		//リクエストヘッダーにContent-Typeの使用を許可
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		//プリフライト
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}

		if !checkAuth(w,r){return}

		if r.Method == http.MethodPost {
			HandleAddToList(w,r,db)
		} else if r.Method == http.MethodDelete {
			HandleDeleteFromList(w,r,db)
		} else if r.Method == http.MethodGet {
			HandleGetList(w,r,db)
		} else {
			http.Error(w, "Unsupported method", http.StatusMethodNotAllowed)
		}
	})

	http.HandleFunc("/search", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin","*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, DELETE")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		} else if r.Method == http.MethodPost {
			HandleGetInfo(w,r,db)
		} else {
			http.Error(w, "Unsupported method", http.StatusMethodNotAllowed)
		}
	})

	http.HandleFunc("/register", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin","*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		} else if r.Method == http.MethodPost {
			HandleRegister(w,r,db);
		} else {
			http.Error(w, "Unsupported method", http.StatusMethodNotAllowed)
			return
		}
	})

	http.HandleFunc("/login", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin","*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		} else if r.Method == http.MethodPost {
			HandleLogin(w,r,db);
		} else {
			http.Error(w, "Unsupported method", http.StatusMethodNotAllowed)
			return
		}
	})

	http.HandleFunc("/verify", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin","*")
		w.Header().Set("Access-Control-Allow-Methods", "GET")
		w.Header().Set("Access-Control-Allow-Headers", "Authorization")

		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		} else if r.Method == http.MethodGet {
			checkAuth(w,r);
			return
		} else {
			http.Error(w, "Unsupported method", http.StatusMethodNotAllowed)
			return
		}
	})

	log.Println("Server started on: 8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func checkAuth(w http.ResponseWriter, r *http.Request) bool{
	tokenString := strings.TrimPrefix(r.Header.Get("Authorization"), "Bearer ")
	//検証
	token, err := jwt.Parse(tokenString, func(t *jwt.Token) (any, error) {
		return []byte(os.Getenv("JWT_SECRET")), nil
	})
	if err != nil || !token.Valid {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return false
	}
	return true
}

func HandleAddToList(w http.ResponseWriter, r *http.Request, db *sql.DB){
	var info DBInfo
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&info); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	stmt, err := db.Prepare("INSERT INTO favoritelist (imdbID,Title,Year,Poster) VALUES (?,?,?,?)")
	if err != nil {
		http.Error(w, fmt.Sprintf("Database query filed: %v", err), http.StatusInternalServerError)
		return
	}
	defer stmt.Close()

	_, err = stmt.Exec(info.ImdbID,info.Title,info.Year,info.Poster)
	if err != nil {
		http.Error(w, fmt.Sprintf("Database query failed: %v", err), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(info)
}

func HandleDeleteFromList(w http.ResponseWriter, r *http.Request, db *sql.DB){
	var info DBInfo
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&info); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	stmt, err := db.Prepare("DELETE FROM favoritelist WHERE imdbID = ?")
	if err != nil {
		http.Error(w, fmt.Sprintf("Database query failed: %v",err), http.StatusInternalServerError)
		return
	}
	defer stmt.Close()

	_, err = stmt.Exec(info.ImdbID)
	if err != nil {
		http.Error(w, fmt.Sprintf("Database query failed: %v",err), http.StatusInternalServerError)
		return
	}

	response := map[string]string{"message": "Successfully deleted"}
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func HandleGetList(w http.ResponseWriter, r *http.Request, db *sql.DB){
	rows,err := db.Query("SELECT * FROM favoritelist")
	if err != nil {
		http.Error(w, fmt.Sprintf("Database query failed: %v",err), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	list := []DBInfo{}
	for rows.Next(){
		var info DBInfo
		if err := rows.Scan(&info.ImdbID,&info.Title,&info.Year,&info.Poster); err != nil {
			http.Error(w, fmt.Sprintf("Row scan failed: %v",err), http.StatusInternalServerError)
			return
		}

		list = append(list,info)
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type","application/json")
	json.NewEncoder(w).Encode(list)
}

func HandleGetInfo(w http.ResponseWriter, r *http.Request, db *sql.DB){
	var req DBInfo
	req.ImdbID = ""
	req.Title = ""
	apikey := os.Getenv("OMDB_API_KEY")
	var url string
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w,"Invalid JSON", http.StatusBadRequest)
		return
	}

	if req.ImdbID != "" {
		url = fmt.Sprintf("https://www.omdbapi.com/?apikey=%s&i=%s",apikey,req.ImdbID)
	} else if req.Title != "" {
		url = fmt.Sprintf("https://www.omdbapi.com/?apikey=%s&t=%s&y=%s",apikey,req.Title,req.Year)
	} else {
		http.Error(w, "imdbID or Title is required", http.StatusBadRequest)
		return
	}

	resp, err := http.Get(url)
	if err != nil {
		http.Error(w, "Failed to fetch from OMDb", http.StatusBadGateway)
		return
	}
	defer resp.Body.Close()

	w.Header().Set("Content-Type","application/json")
	//resp.Bodyをwにそのままながす
	io.Copy(w, resp.Body)
}

func HandleRegister(w http.ResponseWriter, r *http.Request, db *sql.DB){
	var user User
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	if user.UserID == "" || user.Password == "" {
		http.Error(w, "user id and password is required", http.StatusBadRequest)
		return
	}

	//パスワードをバイトに変換してからハッシュ化、コストは重くするほど時間が掛かるが解読されづらくなる
	hashed, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		http.Error(w, "Failed to hash password", http.StatusInternalServerError)
		return
	}

	stmt, err := db.Prepare("INSERT INTO users (user_id,password) VALUES(?,?)")
	if err != nil {
		http.Error(w, fmt.Sprintf("Database query filed: %v", err), http.StatusInternalServerError)
		return
	}
	defer stmt.Close()

	_, err = stmt.Exec(user.UserID,string(hashed))
	if err != nil {
		http.Error(w, fmt.Sprintf("Database query failed:%v", err), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	w.Header().Set("Content-Type","application/json")
	json.NewEncoder(w).Encode(map[string]string{"message": "User registered successfully"})
}

func HandleLogin(w http.ResponseWriter, r *http.Request, db *sql.DB){
	var user User
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	//結果が一行しか返ってこない場合はQueryRow
	row := db.QueryRow("SELECT password FROM users WHERE user_id = ?", user.UserID)
	var hashedPassword string
	err = row.Scan(&hashedPassword)
	if err != nil {
		//idが違うのか、パスワードが違うのかを判別させないためエラー文は同じにする
		http.Error(w, "Invalid user_id or password", http.StatusUnauthorized)
		return
	}
	err = bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(user.Password))
	if err != nil {
		http.Error(w, "Invalid user_id or password", http.StatusUnauthorized)
		return
	}

	//ペイロード
	claims := jwt.MapClaims{
		"user_id": user.UserID,
		"exp": time.Now().Add(time.Hour).Unix(),
	}

	//トークン生成
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	//秘密鍵で書名
	secretKey := os.Getenv("JWT_SECRET")
	tokenString, err := token.SignedString([]byte(secretKey))
	if err != nil {
		http.Error(w, "Failed to generate token", http.StatusInternalServerError)
		return
	}

	//writeheaderを省略すると自動的にOKステータスが返る
	w.Header().Set("Content-Type","application/json")
	json.NewEncoder(w).Encode(map[string]string{"token": tokenString})
}

//参考用コード
// // 構造体を定義
// type Item struct {
// 	ID   int    `json:"id"` //jsonでのキー名を指定
// 	Name string `json:"name"`
// }

// func main() {
// 	// SQLiteのDBファイルを開く（なければ自動作成）
// 	db, err := sql.Open("sqlite", "./test.db")
// 	if err != nil {
// 		log.Fatalf("Failed to connect to database: %v", err) //Fatalで強制終了
// 	}
// 	defer db.Close() //deferで関数の終わりにコードを実行

// 	// テーブルが存在しない場合は作成
// 	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS items (
// 		id INTEGER PRIMARY KEY AUTOINCREMENT,
// 		name TEXT NOT NULL
// 	)`)
// 	if err != nil {
// 		log.Fatalf("Failed to create table: %v", err)
// 	}

// 	// itemsにGETリクエストがあった場合の処理
// 	http.HandleFunc("/items", func(w http.ResponseWriter, r *http.Request) {
// 		if r.Method == http.MethodGet {
// 			handleGetItems(w, r, db)
// 		} else if r.Method == http.MethodPost {
// 			handlePostItem(w, r, db)
// 		} else {
// 			http.Error(w, "Unsupported method", http.StatusMethodNotAllowed)
// 		}
// 	})

// 	http.HandleFunc("/items/", func(w http.ResponseWriter, r *http.Request) {
// 		if r.Method == http.MethodDelete {
// 			handleDelete(w, r, db)
// 		} else if r.Method == http.MethodPut {
//             handlePut(w, r, db)
//         }else {
// 			http.Error(w, "Unsupported method", http.StatusMethodNotAllowed)
// 		}
// 	})

// 	// サーバーを起動
// 	log.Println("Server started on :8080")
// 	log.Fatal(http.ListenAndServe(":8080", nil))
// }

// func handleGetItems(w http.ResponseWriter, r *http.Request, db *sql.DB) {
// 	// データベースからデータを取得
// 	// SELECTのときは.Query、それ以外の時は.Exec
// 	rows, err := db.Query("SELECT id, name FROM items")
// 	// エラーハンドリング
// 	if err != nil {
// 		http.Error(w, fmt.Sprintf("Database query failed: %v", err), http.StatusInternalServerError)
// 		return
// 	}
// 	defer rows.Close()

// 	// 取得したデータを格納するスライスを定義
// 	var items []Item
// 	//ひとつずつ取り出す
// 	for rows.Next() {
// 		var item Item
// 		//Scanで読み込む、変数のポインタを渡す
// 		if err := rows.Scan(&item.ID, &item.Name); err != nil {
// 			http.Error(w, fmt.Sprintf("Row scan failed: %v", err), http.StatusInternalServerError)
// 			return
// 		}
// 		// 取得したレコードをitemsに追加
// 		items = append(items, item)
// 	}

// 	// JSON形式でレスポンスを返す
// 	w.Header().Set("Content-Type", "application/json")
// 	w.WriteHeader(http.StatusOK)
// 	if err := json.NewEncoder(w).Encode(items); err != nil {
// 		http.Error(w, fmt.Sprintf("JSON encoding failed: %v", err), http.StatusInternalServerError)
// 	}
// }

// func handlePostItem(w http.ResponseWriter, r *http.Request, db *sql.DB) {
// 	var item Item
// 	decoder := json.NewDecoder(r.Body)
// 	if err := decoder.Decode(&item); err != nil {
// 		http.Error(w, "Invalid JSON", http.StatusBadRequest)
// 		return
// 	}

// 	if item.Name == "" {
// 		http.Error(w, "Name must be provided", http.StatusBadRequest)
// 		return
// 	}

// 	//prepareでSQLi対策
// 	stmt, err := db.Prepare("INSERT INTO items(name) VALUES(?)")
// 	if err != nil {
// 		http.Error(w, fmt.Sprintf("Database query failed: %v", err), http.StatusInternalServerError)
// 		return
// 	}
// 	defer stmt.Close()

// 	result, err := stmt.Exec(item.Name)
// 	if err != nil {
// 		http.Error(w, fmt.Sprintf("Database query failed: %v", err), http.StatusInternalServerError)
// 		return
// 	}

// 	id, err := result.LastInsertId()
// 	if err != nil {
// 		http.Error(w, fmt.Sprintf("Failed to get last insert id: %v", err), http.StatusInternalServerError)
// 		return
// 	}

// 	item.ID = int(id)
// 	w.WriteHeader(http.StatusCreated)
// 	w.Header().Set("Content-Type", "application/json")
// 	json.NewEncoder(w).Encode(item)
// }

// func handleDelete(w http.ResponseWriter, r *http.Request, db *sql.DB) {
// 	id := strings.TrimPrefix(r.URL.Path, "/items/")

// 	stmt, err := db.Prepare("DELETE FROM items WHERE id = ?")
// 	if err != nil {
//         http.Error(w, fmt.Sprintf("Failed to retrieve affected rows: %v", err), http.StatusInternalServerError)
//         return
// 	}
// 	defer stmt.Close()

// 	result, err := stmt.Exec(id)
// 	if err != nil {
//         http.Error(w, fmt.Sprintf("Failed to execute SQL statement: %v", err), http.StatusInternalServerError)
//         return
// 	}

// 	//削除された行数を取得
// 	rowsAffected, err := result.RowsAffected()
// 	if err != nil {
//         http.Error(w, fmt.Sprintf("Failed to retrieve affected rows: %v", err), http.StatusInternalServerError)
//         return
// 	}

//     var msg string
//     if rowsAffected == 0 {
//         msg = "No item found with the provided ID"
//     } else {
//         msg = "Successfully delete a item"
//     }

// 	//メッセージを書き込む
//     response := map[string]string {"message": msg}
//     w.WriteHeader(http.StatusOK)
//     w.Header().Set("Content-Type", "application/json")
//     if err := json.NewEncoder(w).Encode(response); err != nil {
// 		http.Error(w, fmt.Sprintf("JSON encoding failed: %v", err), http.StatusInternalServerError)
//         return
// 	}
// }

// func handlePut(w http.ResponseWriter, r *http.Request, db *sql.DB) {
//     var newname Item
//     decoder := json.NewDecoder(r.Body)
//     err := decoder.Decode(&newname)
//     if err != nil {
//         http.Error(w, "Invalid JSON", http.StatusBadRequest)
//         return
//     }

//     if newname.Name == "" {
//         http.Error(w, fmt.Sprintf("Name must be provided: %v", err), http.StatusBadRequest)
//         return
//     }

//     id := strings.TrimPrefix(r.URL.Path, "/items/")

//     log.Print(newname)

//     stmt, err := db.Prepare("UPDATE items SET name = ? WHERE id = ?")
//     if err != nil {
//         http.Error(w, fmt.Sprintf("Failed to retrieve affected rows: %v", err), http.StatusBadRequest)
//         return
//     }
//     defer stmt.Close()

//     result, err := stmt.Exec(newname.Name, id)
// 	if err != nil {
//         http.Error(w, fmt.Sprintf("Failed to execute SQL statement: %v", err), http.StatusBadRequest)
//         return
// 	}

//     rowsAffected, err := result.RowsAffected()
// 	if err != nil {
//         http.Error(w, fmt.Sprintf("Failed to retrieve affected rows: %v", err), http.StatusInternalServerError)
//         return
// 	}

//     var msg string
//     if rowsAffected == 0 {
//         msg = "No item found with the provided ID"
//     } else {
//         msg = "Successfully upadate a name"
//     }

//     response := map[string]string {"message": msg}
//     w.WriteHeader(http.StatusOK)
//     w.Header().Set("Content-Type", "application/json")
//     if err := json.NewEncoder(w).Encode(response); err != nil {
// 		http.Error(w, fmt.Sprintf("JSON encoding failed: %v", err), http.StatusInternalServerError)
//         return
// 	}
// }