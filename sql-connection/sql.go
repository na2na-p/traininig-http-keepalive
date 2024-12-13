package main

import (
	"database/sql"
	"fmt"
	"log"
	"time"

	_ "github.com/mattn/go-sqlite3" // SQLite driver
)

const (
	CONN_MAX_LIFETIME  = 60 * time.Second
	CONN_MAX_IDLETIME  = 30 * time.Second
	TEST_QUERY_LATENCY = 5 * time.Second
)

func main() {
	db, err := sql.Open("sqlite3", ":memory:")
	if err != nil {
		log.Fatalf("Failed to open database: %v", err)
	}
	defer db.Close()

	// テーブル作成
	_, err = db.Exec("CREATE TABLE users (id INTEGER PRIMARY KEY, name TEXT)")
	if err != nil {
		log.Fatalf("Failed to create table: %v", err)
	}

	// 接続プールを1つに制限
	db.SetMaxOpenConns(1)
	db.SetMaxIdleConns(1)
	// 最大接続寿命は設定しない
	db.SetConnMaxLifetime(CONN_MAX_LIFETIME)
	// 最大アイドル時間を設定
	db.SetConnMaxIdleTime(CONN_MAX_IDLETIME)

	// データを挿入
	_, err = db.Exec("INSERT INTO users (name) VALUES (?)", "Alice")
	if err != nil {
		log.Fatalf("Failed to insert: %v", err)
	}

	// 待機して寿命を超えさせる
	fmt.Println("Waiting before query...")
	time.Sleep(6 * time.Second)

	// クエリを実行
	var name string
	err = db.QueryRow(fmt.Sprintf(`
		SELECT name 
		FROM users 
		WHERE id = 1 
	`)).Scan(&name)
	if err != nil {
		log.Fatalf("Query failed: %v", err)
	}

	fmt.Printf("Retrieved user: %s\n", name)

	// 確認のためのログを出力
	stats := db.Stats()
	fmt.Printf("Stats: %+v\n", stats)
}
