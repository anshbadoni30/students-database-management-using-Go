package sqlite

import (
	"database/sql"

	"github.com/anshbadoni30/students-api/internal/config"
	_ "github.com/mattn/go-sqlite3"
)

type Sqlite struct {
	Db *sql.DB
}

func New(cfg *config.Config) (*Sqlite,error){
	db,err:= sql.Open("sqlite3", cfg.StoragePath)
	if err!=nil{
		return nil,err
	}

	_, errr:=db.Exec(`CREATE TABLE IF NOT EXISTS students (
	id INTEGER PRIMARY KEY AUTOINCREMENT,
	name TEXT,
	email TEXT,
	age INTEGER) 
	`)

	if errr!=nil{
		return nil,errr
	}

	return &Sqlite{
		Db: db,
	},nil
}

func (s *Sqlite)CreateStudent(name string, email string, age int) (int64,error){
	stmt,err:=s.Db.Prepare("INSERT INTO STUDENTS (name,email,age) VALUES (?,?,?)")
	if err!=nil{
		return 0,err
	}
	defer stmt.Close()

	result,er:=stmt.Exec(name,email,age)
	if er!=nil{
		return 0,er
	}
	lastid,er:=result.LastInsertId()
	if er!=nil{
		return 0,er
	}
	
	return lastid,nil
}