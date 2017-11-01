package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	_ "github.com/mattn/go-sqlite3"
)

type ImageInfo struct {
	Info Data
}

type Data struct {
	Id          int    `json:"id"`
	Title       string `json:"title"`
	Description string `json:"description"`
	Images      Links  `json:"images"`
}
type Links struct {
	Hidpi  string
	Normal string
	Teaser string
}

type Title map[string]string

func main() {

	db, _ := sql.Open("sqlite3", "cache/web.db")
	creatDB := "CREATE TABLE IF NOT EXISTS images (id INTEGER PRIMARY KEY UNIQUE NOT NULL,title STRING,description TEXT,link STRING,fileName STRING);"
	db.Exec(creatDB)
	imageNames := [10]string{"3872570-Floating-buttons", "3905149-Timer-App-Prototype", "3905639-Dream-of-Amsterdam", "3906066-The-City-of-Cats", "160781-Olimp-006", "618406-Cassette", "2685421-Tmnt-manhole", "1508295-Captain-America", "998546-Gibson-Les-Paul", "2380018-Explorer-Club-Laguna"}

	fmt.Println("Saving entries to database. Please wait...")
	for index := range imageNames {

		url := "https://api.dribbble.com/v1/shots/" + imageNames[index] + "?access_token=343a27a5a71a4ed9d29030354cde2b630030eeba9b73a8785b100ceef966bbe5"
		res, err := http.Get(url)

		if err != nil {
			panic(err)
		}
		defer res.Body.Close()

		body, err := ioutil.ReadAll(res.Body)

		if err != nil {
			panic(err)
		}

		var result Data

		err = json.Unmarshal(body, &result)

		if err != nil {
			panic(err)
		}
		//fmt.Println(result.Images.Normal)

		tx, _ := db.Begin()
		stmt, _ := tx.Prepare("Insert Into images (id,title,description,link,fileName) values(?,?,?,?,?)")

		_, err = stmt.Exec(result.Id, result.Title, result.Description, result.Images.Normal, result.Title+".jpg")
		tx.Commit()
		response, e := http.Get(result.Images.Normal)
		if e != nil {
			panic(e)
		}
		defer response.Body.Close()
		file, err := os.Create("cache/images/" + result.Title + ".jpg")
		if err != nil {
			panic(err)
		}
		_, err = io.Copy(file, response.Body)
		if err != nil {
			panic(err)
		}
		file.Close()
		fmt.Println(".")
	}
	fmt.Println("All entries saved")

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		//s := r.PostFormValue("ky") + "Hello"
		//fmt.Fprintf(w, s)
		var dt Data
		var fl string
		s := r.PostFormValue("val")
		if s != "" {
			q, err := db.Query("SELECT * FROM images WHERE title LIKE '%" + s + "%' OR description Like '%" + s + "%'")
			if err != nil {
				panic(err)
			}
			for q.Next() {
				q.Scan(&dt.Id, &dt.Title, &dt.Description, &dt.Images.Normal, &fl)

				res, err := json.MarshalIndent(dt, "", "   ")
				if err != nil {
					panic(err)
				}
				fmt.Fprintf(w, string(res))
			}
		}

	})

	log.Fatal(http.ListenAndServe(":8080", nil))
	fmt.Println("ServerRunning")
}
