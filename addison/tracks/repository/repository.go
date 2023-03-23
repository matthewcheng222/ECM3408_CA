package repository

import (
	"database/sql"
	"log"

	_ "github.com/mattn/go-sqlite3"
)

type Repository struct {
	DB *sql.DB
}

var repo Repository

func Init() {
	db, err := sql.Open("sqlite3", "AddisonTracks.db")

	if err != nil {
		log.Fatal("Database Initialisation Failed")
		return
	}
	repo = Repository{DB: db}
}

func Create() int {
	const sql = "CREATE TABLE IF NOT EXISTS Tracks(Id TEXT PRIMARY KEY, Audio TEXT)"
	_, err := repo.DB.Exec(sql)

	if err != nil {
		return -1
	}
	return 0
}

func Clear() int {
	const sql = "DELETE FROM Tracks"
	_, err := repo.DB.Exec(sql)

	if err != nil {
		return -1
	}
	return 0
}

func Update(t Track) int64 {
	const sql = "UPDATE Tracks SET Audio = ? WHERE Id = ?"
	stmt, err := repo.DB.Prepare(sql)
	if err != nil {
		return -1
	}
	defer stmt.Close()

	res, err := stmt.Exec(t.Audio, t.Id)
	if err != nil {
		return -1
	}

	n, err := res.RowsAffected()
	if err != nil {
		return -1
	}

	return n
}

func Insert(t Track) int64 {
	const sql = "INSERT INTO Tracks (Id, Audio) VALUES (?, ?)"
	stmt, err := repo.DB.Prepare(sql)
	if err != nil {
		return -1
	}
	defer stmt.Close()

	res, err := stmt.Exec(t.Id, t.Audio)
	if err != nil {
		return -1
	}

	n, err := res.RowsAffected()
	if err != nil {
		return -1
	}

	return n
}

func List() ([]Track, int) {
	const sql = "SELECT * FROM Tracks"
	stmt, err := repo.DB.Prepare(sql)
	if err != nil {
		return []Track{}, -1
	}
	defer stmt.Close()

	tracks := make([]Track, 0)
	rows, err := stmt.Query()
	if err != nil {
		return []Track{}, -1
	}

	for rows.Next() {
		tempTrack := Track{}
		err := rows.Scan(&tempTrack.Id, &tempTrack.Audio)
		if err != nil {
			return []Track{}, 0
		}
		tracks = append(tracks, tempTrack)
	}

	return tracks, len(tracks)
}

func Read(id string) (Track, int64) {
	const sql = "SELECT * FROM Tracks WHERE Id = ?"
	stmt, err := repo.DB.Prepare(sql)
	if err != nil {
		return Track{}, -1
	}
	defer stmt.Close()

	var t Track
	row := stmt.QueryRow(id)
	err = row.Scan(&t.Id, &t.Audio)
	if err != nil {
		return Track{}, 0
	}
	return t, 1
}

func Delete(id string) int64 {
	const sql = "DELETE FROM Tracks WHERE Id = ?"
	stmt, err := repo.DB.Prepare(sql)
	if err != nil {
		return -1
	}
	defer stmt.Close()

	executed, err := stmt.Exec(id)
	if err != nil {
		return -1
	}

	n, err := executed.RowsAffected()
	if err != nil {
		return -1
	}

	return n
}
