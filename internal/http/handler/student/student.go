package student

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"

	"github.com/anshbadoni30/students-api/internal/types"
	"github.com/anshbadoni30/students-api/internal/utils/response"
	"github.com/go-playground/validator/v10"
)

func New() http.HandlerFunc{
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
			validatorErrs:= err.(validator.ValidationErrors)
			response.WriteJson(w,http.StatusBadRequest,response.ValidateError(validatorErrs))
			return
		}

		response.WriteJson(w,http.StatusCreated,map[string]string{"success":"OK"})
	}
}