package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"

	_ "github.com/lib/pq"
)

const CONFIG_FILE = "config.json"

type Album struct {
	ID     int    `json:"id"`
	Title  string `json:"title"`
	Artist string `json:"artist"`
	Genre  string `json:"genre"`
	Year   string `json:"year"`
}

type Config struct {
	User     string `json:"user"`
	Password string `json:"password"`
	Host     string `json:"host"`
	Port     int    `json:"port"`
	Database string `json:"db"`
}

var album *Album
var db *sql.DB
var err error

func main() {
	album = &Album{}
	// Read the config information
	config := getConfigInfo(CONFIG_FILE)
	connectionString := generateConnectionString(&config)
	db, err = sql.Open("postgres", connectionString)
	if err != nil {
		log.Fatal(err)
	}

	defer db.Close()

	http.HandleFunc("/album/{id}", album.albumHandler)

	if err = http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal(err)
	}

}

func getConfigInfo(configFilename string) Config {
	// Read the contents of the config file
	var data []byte
	if data, err = os.ReadFile(configFilename); err != nil {
		log.Fatal(err)
	}

	// Unmarshal the JSON data into a Config object
	config := Config{}
	if err := json.Unmarshal(data, &config); err != nil {
		log.Fatal(err)
	}

	return config
}

func generateConnectionString(config *Config) string {
	connectionString := fmt.Sprintf(
		"postgresql://%s:%s@%s:%d/%s?sslmode=disable",
		config.User,
		config.Password,
		config.Host,
		config.Port,
		config.Database,
	)

	return connectionString
}

func (album *Album) getAlbumInfo(id int) {
	// Query the data
	query := fmt.Sprintf("SELECT * FROM albums WHERE id = %d", id)
	rows, err := db.Query(query)
	if err != nil {
		album.ID = 0
		album.Title = "N/A"
		album.Artist = "N/A"
		album.Genre = "N/A"
		album.Year = "N/A"

		return
	}

	defer rows.Close()

	// Load the values of the query into the `Album` struct
	rows.Next()
	if err = rows.Scan(
		&album.ID,
		&album.Title,
		&album.Artist,
		&album.Genre,
		&album.Year,
	); err != nil {
		album.ID = 0
		album.Title = "N/A"
		album.Artist = "N/A"
		album.Genre = "N/A"
		album.Year = "N/A"
	}

	return
}

func (album *Album) getJsonData() []byte {
	// Marshal the `Album` struct into a JSON object
	jsonData, err := json.Marshal(album)
	if err != nil {
		log.Fatal(err)
	}

	// Return the JSON
	return jsonData
}

func (album *Album) albumHandler(w http.ResponseWriter, r *http.Request) {
	var id int
	// Parse the album ID from the URL path
	idString := r.PathValue("id")
	log.Printf("%s GET %s\n",
		r.RemoteAddr,
		r.RequestURI,
	)
	if id, err = strconv.Atoi(idString); err != nil {
		fmt.Fprint(w, "Invalid ID")

		return
	}
	// Parse the album information from the album
	album.getAlbumInfo(id)
	// Generate the JSON data and convert it into a string
	jsonData := album.getJsonData()
	jsonDataString := string(jsonData)
	// Send the JSON data in a response
	fmt.Fprintf(w, jsonDataString)

	return
}
