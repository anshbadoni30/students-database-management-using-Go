package sqlite

import (
	"database/sql"
	"fmt"

	"github.com/anshbadoni30/students-api/internal/config"
	"github.com/anshbadoni30/students-api/internal/types"
	_ "github.com/mattn/go-sqlite3"
)

type Sqlite struct {
	Db *sql.DB
}

func CreateTable(cfg *config.Config) (*Sqlite,error){
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

func (s *Sqlite)GetStudentById(id int64) (types.Student,error){
	stmt,err:=s.Db.Prepare("SELECT * FROM students WHERE ID = ? LIMIT 1")
	if err!=nil{
		return types.Student{}, err
	}
	defer stmt.Close()

	var student types.Student

	err = stmt.QueryRow(id).Scan(&student.Id,&student.Name,&student.Email,&student.Age) //when we get the output from query then it will stored in student object by using scan function
	if err!=nil{
		if err==sql.ErrNoRows{
			return types.Student{},fmt.Errorf("no student found with id %s",fmt.Sprint(id))
		}
		return types.Student{},fmt.Errorf("query error: %w",err)
	}
	return student,nil
}

func (s *Sqlite)GetStudents()([]types.Student,error){
	stmt,err:= s.Db.Prepare("Select * from students")
	if err!=nil{
		return nil,err
	}
	defer stmt.Close()

	rows,er:=stmt.Query()
	if er!=nil{
		return nil,er
	}
	defer rows.Close()
	var students []types.Student

	for rows.Next(){
		var student types.Student
		er:=rows.Scan(&student.Id,&student.Name,&student.Email,&student.Age)
		if er!=nil{
			return nil,er
		}
		students=append(students, student)
	}
	return students,nil
}