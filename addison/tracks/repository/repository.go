package repository

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/mattn/go-sqlite3"
)

type Repository struct {
	DB *sql.DB
}

var repo Repository

func Init() {
	if db, err := sql.Open("sqlite3", "AddisonTracks.db"); err == nil {
		repo = Repository{DB: db}
	} else {
		log.Fatal("Database initialisation")
	}
}

func Create() int {
	const sql = "CREATE TABLE IF NOT EXISTS Tracks(Id TEXT PRIMARY KEY, Audio TEXT)"
	if _, err := repo.DB.Exec(sql); err == nil {
		return 0
	} else {
		return -1
	}
}

func Clear() int {
	const sql = "DELETE FROM Tracks"
	if _, err := repo.DB.Exec(sql); err == nil {
		return 0
	} else {
		return -1
	}
}

func Update(t Track) int64 {
	const sql = "UPDATE Tracks SET Audio = ? WHERE Id = ?"
	if stmt, err := repo.DB.Prepare(sql); err == nil {
		defer stmt.Close()
		if res, err := stmt.Exec(t.Audio, t.Id); err == nil {
			if n, err := res.RowsAffected(); err == nil {
				fmt.Println("Update")
				return n
			} else {
				return -1
			}
		} else {
			return -1
		}
	} else {
		return -1
	}
}

func Insert(t Track) int64 {
	const sql = "INSERT INTO Tracks (Id, Audio) VALUES (?, ?)"
	if stmt, err := repo.DB.Prepare(sql); err == nil {
		defer stmt.Close()
		if res, err := stmt.Exec(t.Id, t.Audio); err == nil {
			if n, err := res.RowsAffected(); err == nil {
				fmt.Println("Insert")
				return n
			} else {
				return -1
			}
		} else {
			return -1
		}
	} else {
		return -1
	}
}

func List() ([]Track, int) {
	const sql = "SELECT * FROM Tracks"
	if stmt, err := repo.DB.Prepare(sql); err == nil {
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
	} else {
		return []Track{}, -1
	}
}

func Read(id string) (Track, int64) {
	const sql = "SELECT * FROM Tracks WHERE Id = ?"
	if stmt, err := repo.DB.Prepare(sql); err == nil {
		defer stmt.Close()
		var t Track
		row := stmt.QueryRow(id)
		if err := row.Scan(&t.Id, &t.Audio); err == nil {
			return t, 1
		} else {
			return Track{}, 0
		}
	}
	return Track{}, -1
}

func Delete(id string) int64 {
	const sql = "DELETE FROM Tracks WHERE Id = ?"
	if stmt, err := repo.DB.Prepare(sql); err == nil {
		defer stmt.Close()
		if executed, err := stmt.Exec(id); err == nil {
			if n, err := executed.RowsAffected(); err == nil {
				return n
			} else {
				return -1
			}
		} else {
			return -1
		}
	} else {
		return -1
	}
}
