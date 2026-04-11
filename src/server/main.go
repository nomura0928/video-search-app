package main

import (
	"database/sql"  // データベースの操作
	"encoding/json" // JSONのエンコードとデコード
	"fmt"
	"log"
	"net/http" // HTTPでWebサーバーを立てる
	// "strings"

	_ "modernc.org/sqlite"
)

type Rating struct {
	Source string `json:"Source"`
	Value  string `json:"Value"`
}

type Info struct {
	ImdbID   string   `json:"imdbID"`
	Title    string   `json:"Title"`
	Year     string   `json:"Year"`
	Genre    string   `json:"Genre"`
	Type     string   `json:"Type"`
	Director string   `json:"Director"`
	Writer   string   `json:"Writer"`
	Actor    string   `json:"Actor"`
	Ratings  []Rating `json:"Ratings"`
	Plot     string   `json:"Plot"`
	Poster   string   `json:"Poster"`
	Response string   `json:"Response"`
}

type DBInfo struct {
	ImdbID string `json:"imdbID"`
	Title  string `json:"Title"`
	Year   string `json:"Year"`
	Poster string `json:"Poster"`
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

	http.HandleFunc("/search", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost {
			HandleAddToList(w,r,db)
		} else if r.Method == http.MethodDelete {
			HandleDeleteFromList(w,r,db)
		} else {
			http.Error(w, "Unsupported method", http.StatusMethodNotAllowed)
		}
	})

	http.HandleFunc("/favorites", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet {

		} else {
			http.Error(w, "Unsupported method", http.StatusMethodNotAllowed)
		}
	})

	log.Println("Server started on: 8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
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