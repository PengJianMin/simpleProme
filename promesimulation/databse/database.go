package database

import (
	"database/sql"
	"fmt"
	"github.com/go-sql-driver/mysql"
	"github.com/spf13/viper"
	"log"
)
type Album struct {
	ID     int64
	Title  string
	Artist string
	Price  float64
}

var db *sql.DB
var cfg mysql.Config

func init() {
	cfg , _ = getDbConfig()
	db = GetDBHandle(cfg)
}

//取得mysql登录的配置信息
func getDbConfig() (mysql.Config,error) {
	config := viper.New()
	config.AddConfigPath("./database/")
	config.SetConfigName("init")
	config.SetConfigType("yaml")
	if err := config.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			fmt.Println("Configuration File Not Found!")
		} else {
			fmt.Println("Errors in Configuration Content!")
		}
		return mysql.Config{},err
	}
	var result mysql.Config
	result.User = config.GetString("database.user")
	result.Passwd = config.GetString("database.pwd")
	result.Net = "tcp"
	result.Addr = config.GetString("database.host")
	result.DBName = config.GetString("database.dbname")

	return result,nil
}

//获得数据库处理器
func GetDBHandle(cfg mysql.Config) *sql.DB{
	db,err := sql.Open("mysql",cfg.FormatDSN())
	if err != nil {
		log.Fatal(err)
	}
	pingErr := db.Ping()
	if pingErr != nil {
		log.Fatal(pingErr)
	}
	fmt.Println("Mysql Connected!")
	return  db
}

func GetConnectedDbHandle() *sql.DB{
	return db
}


//请求全部数据
func GetALLAlbumsInfo() ([]Album, error) {
	// An albums slice to hold data from returned rows.
	var albums []Album

	rows, err := db.Query("SELECT * FROM album")
	if err != nil {
		return nil, fmt.Errorf("albumsByArtist %q: %v",err)
	}
	defer rows.Close()
	// Loop through rows, using Scan to assign column data to struct fields.
	for rows.Next() {
		var alb Album
		if err := rows.Scan(&alb.ID, &alb.Title, &alb.Artist, &alb.Price); err != nil {
			return nil, fmt.Errorf("albumsByArtist  %v",  err)
		}
		albums = append(albums, alb)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("albumsByArtist  %v", err)
	}
	return albums, nil
}

//根据ID请求数据
// albumByID queries for the album with the specified ID.
func GetAlbumInfoByID(id int64) (Album, error) {
	// An album to hold data from the returned row.
	var alb Album

	row := db.QueryRow("SELECT * FROM album WHERE id = ?", id)
	if err := row.Scan(&alb.ID, &alb.Title, &alb.Artist, &alb.Price); err != nil {
		if err == sql.ErrNoRows {
			return alb, fmt.Errorf("albumsById %d: no such album", id)
		}
		return alb, fmt.Errorf("albumsById %d: %v", id, err)
	}
	return alb, nil
}

//新增数据
// addAlbum adds the specified album to the database,
// returning the album ID of the new entry
func AddAlbum(alb Album) (int64, error) {
	result, err := db.Exec("INSERT INTO album (title, artist, price) VALUES (?, ?, ?)", alb.Title, alb.Artist, alb.Price)
	if err != nil {
		return 0, fmt.Errorf("addAlbum: %v", err)
	}
	id, err := result.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("addAlbum: %v", err)
	}
	return id, nil
}

func DelAlbumInfoByID(id int64)  {
	result, err :=  db.Exec("DELETE FROM album WHERE id = ?", id)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(result.RowsAffected())
}
