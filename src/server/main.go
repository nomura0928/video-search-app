package main

import (
	"database/sql"  // データベースの操作
	"encoding/json" // JSONのエンコードとデコード
	"fmt"
	"log"
	"net/http" // HTTPでWebサーバーを立てる
	"os"
	"strings"

	_ "github.com/go-sql-driver/mysql"
)

// 構造体を定義
type Item struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

func main() {
	// 環境変数からDB接続情報を取得
	dbHost := os.Getenv("DB_HOST")
	dbUser := os.Getenv("DB_USER")
	dbPassword := os.Getenv("DB_PASSWORD")
	dbName := os.Getenv("DB_NAME")

	dsn := fmt.Sprintf("%s:%s@tcp(%s:3306)/%s", dbUser, dbPassword, dbHost, dbName)
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	// itemsにGETリクエストがあった場合の処理
	http.HandleFunc("/items", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet {
			handleGetItems(w, r, db)
		} else if r.Method == http.MethodPost {
			handlePostItem(w, r, db)
		} else if r.Method == http.MethodDelete {
			handleDelete(w, r, db)
		} else {
			http.Error(w, "Unsupported method", http.StatusMethodNotAllowed)
		}
	})

	http.HandleFunc("/items/", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodDelete {
			handleDelete(w, r, db)
		} else if r.Method == http.MethodPut {
            handlePut(w, r, db)
        }else {
			http.Error(w, "Unsupported method", http.StatusMethodNotAllowed)
		}
	})

	// サーバーを起動
	log.Println("Server started on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func handleGetItems(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	// データベースからデータを取得
	rows, err := db.Query("SELECT id, name FROM items")
	// エラーハンドリング
	if err != nil {
		http.Error(w, fmt.Sprintf("Database query failed: %v", err), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	// 取得したデータを格納するスライスを定義
	var items []Item
	for rows.Next() {
		var item Item
		if err := rows.Scan(&item.ID, &item.Name); err != nil {
			http.Error(w, fmt.Sprintf("Row scan failed: %v", err), http.StatusInternalServerError)
			return
		}
		// 取得したレコードをitemsに追加
		items = append(items, item)
	}

	// JSON形式でレスポンスを返す
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(items); err != nil {
		http.Error(w, fmt.Sprintf("JSON encoding failed: %v", err), http.StatusInternalServerError)
	}
}

func handlePostItem(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	var item Item
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&item); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	if item.Name == "" {
		http.Error(w, "Name must be provided", http.StatusBadRequest)
		return
	}

	stmt, err := db.Prepare("INSERT INTO items(name) VALUES(?)")
	if err != nil {
		http.Error(w, fmt.Sprintf("Database query failed: %v", err), http.StatusInternalServerError)
		return
	}
	defer stmt.Close()

	_, err = stmt.Exec(item.Name)
	if err != nil {
		http.Error(w, fmt.Sprintf("Database query failed: %v", err), http.StatusInternalServerError)
		return
	}

	// w.WriteHeader(http.StatusCreated)
	// w.Header().Set("Content-Type", "application/json")
	// json.NewEncoder(w).Encode(item)
}

func handleDelete(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	id := strings.TrimPrefix(r.URL.Path, "/items/")

	stmt, err := db.Prepare("DELETE FROM items WHERE id = ?")
	if err != nil {
        http.Error(w, fmt.Sprintf("Failed to retrieve affected rows: %v", err), http.StatusInternalServerError)
        return
	}
	defer stmt.Close()

	result, err := stmt.Exec(id)
	if err != nil {
        http.Error(w, fmt.Sprintf("Failed to execute SQL statement: %v", err), http.StatusInternalServerError)
        return
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
        http.Error(w, fmt.Sprintf("Failed to retrieve affected rows: %v", err), http.StatusInternalServerError)
        return
	}

    var msg string
    if rowsAffected == 0 {
        msg = "No item found with the provided ID"
    } else {
        msg = "Successfully delete a item"
    }

    response := map[string]string {"message": msg}
    w.WriteHeader(http.StatusOK)
    w.Header().Set("Content-Type", "application/json")
    if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, fmt.Sprintf("JSON encoding failed: %v", err), http.StatusInternalServerError)
        return
	}
}

func handlePut(w http.ResponseWriter, r *http.Request, db *sql.DB) {
    var newname Item
    decoder := json.NewDecoder(r.Body)
    err := decoder.Decode(&newname)
    if err != nil {
        http.Error(w, "Invalid JSON", http.StatusBadRequest)
        return
    }

    if newname.Name == "" {
        http.Error(w, fmt.Sprintf("Name must be provided: %v", err), http.StatusBadRequest)
        return
    }

    id := strings.TrimPrefix(r.URL.Path, "/items/")

    log.Print(newname)

    stmt, err := db.Prepare("UPDATE items SET name = ? WHERE id = ?")
    if err != nil {
        http.Error(w, fmt.Sprintf("Failed to retrieve affected rows: %v", err), http.StatusBadRequest)
        return
    }
    defer stmt.Close()

    result, err := stmt.Exec(newname.Name, id)
	if err != nil {
        http.Error(w, fmt.Sprintf("Failed to execute SQL statement: %v", err), http.StatusBadRequest)
        return
	}

    rowsAffected, err := result.RowsAffected()
	if err != nil {
        http.Error(w, fmt.Sprintf("Failed to retrieve affected rows: %v", err), http.StatusInternalServerError)
        return
	}

    var msg string
    if rowsAffected == 0 {
        msg = "No item found with the provided ID"
    } else {
        msg = "Successfully upadate a name"
    }

    response := map[string]string {"message": msg}
    w.WriteHeader(http.StatusOK)
    w.Header().Set("Content-Type", "application/json")
    if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, fmt.Sprintf("JSON encoding failed: %v", err), http.StatusInternalServerError)
        return
	}
}