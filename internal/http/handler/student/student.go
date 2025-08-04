package student

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/anshbadoni30/students-api/internal/storage"
	"github.com/anshbadoni30/students-api/internal/types"
	"github.com/anshbadoni30/students-api/internal/utils/response"
	"github.com/go-playground/validator/v10"
)

func New( storage storage.Storage) http.HandlerFunc{
	return func(w http.ResponseWriter, r *http.Request){
		var student types.Student

		err:= json.NewDecoder(r.Body).Decode(&student)
		//Checking of errors first
		if errors.Is(err,io.EOF){
			response.WriteJson(w,http.StatusBadRequest,response.GeneralError(fmt.Errorf("empty body")))
			return
		}

		if err!=nil{
			response.WriteJson(w,http.StatusBadRequest,response.GeneralError(err))
		}

		//validating request
		errr:=validator.New().Struct(student)
		if errr!=nil{
			validatorErrs:= errr.(validator.ValidationErrors)
			response.WriteJson(w,http.StatusBadRequest,response.ValidateError(validatorErrs))
			return
		}

		
		lastid,err:=storage.CreateStudent(student.Name,student.Email,student.Age)

		slog.Info("User Created Successfully",slog.String("User Id=",fmt.Sprint(lastid)))
		if err!=nil{
			response.WriteJson(w,http.StatusInternalServerError,err)
		}

		response.WriteJson(w,http.StatusCreated,map[string]int64{"id":lastid})
	}
}

func Getbyid( storage storage.Storage) http.HandlerFunc{
	return func(w http.ResponseWriter, r *http.Request){

	id:=r.PathValue("id")
	slog.Info("getting a student",slog.String("id",id))

	//converting string id in int id
	intid,er:=strconv.ParseInt(id,10,64)
	if er!=nil{
		response.WriteJson(w,http.StatusBadRequest,response.GeneralError(er))
		return
	}
	
	student,er:=storage.GetStudentById(intid)
	if er!=nil{
		response.WriteJson(w,http.StatusInternalServerError,response.GeneralError(er))
		return
	}

	response.WriteJson(w,http.StatusOK,student)

	} 
}

func GetList( storage storage.Storage) http.HandlerFunc{
	return func(w http.ResponseWriter, r *http.Request){
		slog.Info("Getting all students")
		students,err:=storage.GetStudents()
		if err!=nil{
			response.WriteJson(w,http.StatusInternalServerError,err)
			return
		}

		response.WriteJson(w,http.StatusOK,students)
	}
}

func ReplaceName( storage storage.Storage) http.HandlerFunc{
	return func(w http.ResponseWriter, r *http.Request){

		username:=r.PathValue("username")
		id:=r.PathValue("id")

		//converting string id in int id
		intid,er:=strconv.ParseInt(id,10,64)
		if er!=nil{
		response.WriteJson(w,http.StatusBadRequest,response.GeneralError(er))
		return
		}
		rows,err:=storage.NameChange(username,intid)
		if err!=nil{
			response.WriteJson(w,http.StatusInternalServerError,response.GeneralError(err))
			return 
		}
		response.WriteJson(w,http.StatusOK,map[string]int64 {"rows Affected":rows})

	}
}

func Delete(storage storage.Storage) http.HandlerFunc{
	return func(w http.ResponseWriter,r *http.Request){
		id:=r.PathValue("id")

		//converting string id in int id
		intid,er:=strconv.ParseInt(id,10,64)
		if er!=nil{
		response.WriteJson(w,http.StatusBadRequest,response.GeneralError(er))
		return
		}
		rows,err:= storage.DeleteRecord(intid)
		if err!=nil{
			response.WriteJson(w,http.StatusInternalServerError,response.GeneralError(err))
			return
		}
		response.WriteJson(w,http.StatusOK,map[string]int64 {"rows Affected":rows})
	}
}