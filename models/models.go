package models

import (
	"errors"
	"fmt"
	"log"
	"net/url"
	"os"
	"strings"

	_ "github.com/jackc/pgx/v4/stdlib"
	"github.com/jmoiron/sqlx"
	"github.com/joho/godotenv"
)

const (
	ShortUrlLength          = 5
	ShortUrlGenTriesIfExist = 3
)

var (
	ErrCantGenerateShortURL = errors.New("can't generate short url")
	ErrShortURLExist        = errors.New("short url is already exist")
	ErrShortUrlNotExist     = errors.New("short url doesn't exist")
)

type DB struct {
	*sqlx.DB
}

func InitDB() (*DB, error) {
	if err := godotenv.Load(".env"); err != nil {
		log.Fatal("Error when loading .env file")
	}

	conn, err := sqlx.Connect("pgx", os.Getenv("DATABASE"))
	if err != nil {
		return nil, err
	}
	db := &DB{DB: conn}
	return db, nil
}

func getFullURL(URL string) (string, error) {
	if !strings.Contains(URL, "http") {
		URL = fmt.Sprintf("http://%s", URL)
	}
	if _, err := url.ParseRequestURI(URL); err != nil {
		return "", fmt.Errorf("cant parse url: %v", err)
	}
	return URL, nil

}

func (db *DB) SaveURL(longURL string, customURL string) (string, error) {
	longURL, err := getFullURL(longURL)
	if err != nil {
		return "", fmt.Errorf("invalid long url: %v", err)
	}
	if customURL == "" {
		return db.SaveGeneratedURL(longURL)
	}
	return db.SaveCustomURL(longURL, customURL)
}

func (db *DB) SaveGeneratedURL(longURL string) (string, error) {
	var shortURL string
	urlGenerated := false
	for i := 0; i < ShortUrlGenTriesIfExist; i++ {
		shortURL = GenerateURL(ShortUrlLength)
		if exist, err := db.CheckShortURLExist(shortURL); err != nil {
			return "", err
		} else if exist {
			continue
		}
		urlGenerated = true
		break
	}
	if !urlGenerated {
		return "", ErrCantGenerateShortURL
	}
	return db.SaveShortURL(longURL, shortURL)
}

func (db *DB) SaveCustomURL(longURL, customURL string) (string, error) {
	if exist, err := db.CheckShortURLExist(customURL); err != nil {
		return "", err
	} else if exist {
		return "", ErrShortURLExist
	}
	return db.SaveShortURL(longURL, customURL)
}

func (db *DB) SaveShortURL(longURL string, shortURL string) (string, error) {
	_, err := db.Exec("INSERT INTO urls (short_url, long_url) VALUES ($1, $2)", shortURL, longURL)
	if err != nil {
		return "", fmt.Errorf("error save short url: %v", err)
	}
	return shortURL, nil
}

func (db *DB) CheckShortURLExist(shortURL string) (bool, error) {
	urls := []string{}
	err := db.Select(&urls, "SELECT short_url FROM urls WHERE short_url = $1", shortURL)
	if err != nil {
		return false, fmt.Errorf("error checking url exist: %v", err)
	}
	if len(urls) == 0 {
		return false, nil
	}
	return urls[0] == shortURL, nil
}

func (db *DB) GetLongURL(shortURL string) (string, error) {
	if exist, err := db.CheckShortURLExist(shortURL); err != nil {
		return "", err
	} else if !exist {
		return "", ErrShortUrlNotExist
	}
	longUrl := ""
	err := db.Get(&longUrl, "SELECT long_url FROM urls WHERE short_url = $1", shortURL)
	if err != nil {
		return "", fmt.Errorf("error select long url: %v", err)
	}
	return longUrl, nil
}
