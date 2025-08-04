package storage

import "github.com/anshbadoni30/students-api/internal/types"

type Storage interface {
	CreateStudent(name string, email string, age int) (int64, error)
	GetStudentById(id int64) (types.Student,error)
	GetStudents() ([]types.Student, error)
	NameChange(username string,id int64) (int64,error) 
	DeleteRecord(id int64) (int64,error)
}