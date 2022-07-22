package handlers

import (
	"ToDo/Claims"
	"ToDo/database/helper"
	"ToDo/models"
	"encoding/json"
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/golang-jwt/jwt"
	"github.com/lib/pq"
	"log"
	"net/http"
	"strconv"
	"time"
)

var JwtKey = []byte("secureSecretText")

func Login(w http.ResponseWriter, r *http.Request) {
	var credentials models.Credential
	err := json.NewDecoder(r.Body).Decode(&credentials)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if credentials.Username == "" || credentials.Password == "" {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	passkey, err := helper.LoginUser(credentials.Username)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
	}
	if passkey.Password != credentials.Password {
		w.WriteHeader(http.StatusUnauthorized)
	}

	expirationTime := time.Now().Add(time.Hour * 7)
	sessionID, err := helper.CreateSession(passkey.ID, expirationTime)
	if err != nil {
		log.Printf("LogIn : Error in Creating the session")
		w.WriteHeader(http.StatusUnauthorized)
	}
	claims := &Claims.Claims{
		SessionID: sessionID,
		ID:        passkey.ID,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedToken, err := token.SignedString([]byte("secureSecretText"))
	tokenByte, err := json.Marshal(signedToken)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}
	fmt.Println(signedToken)
	_, _ = w.Write(tokenByte)

}
func Signup(writer http.ResponseWriter, request *http.Request) {
	var user models.CreateUser
	decoder := json.NewDecoder(request.Body).Decode(&user)
	if decoder != nil {
		writer.WriteHeader(http.StatusBadRequest)
		return
	}

	userId, err := helper.CreateUser(user)
	if err != nil {
		writer.WriteHeader(http.StatusInternalServerError)
		return
	}
	errs := json.NewEncoder(writer).Encode(userId)
	if errs != nil {
		writer.WriteHeader(http.StatusInternalServerError)
	}

}
func AddTask(writer http.ResponseWriter, request *http.Request) {
	var user models.AddTask
	err := json.NewDecoder(request.Body).Decode(&user)
	if err != nil {
		writer.WriteHeader(http.StatusBadRequest)
		return
	}
	uc := request.Context().Value("claims").(*Claims.Claims)
	_, err = helper.AddTask(user.Task, uc.ID, user.Status)
	if err != nil {
		writer.WriteHeader(http.StatusInternalServerError)
		return
	}

}
func SignOut(writer http.ResponseWriter, request *http.Request) {
	uc := request.Context().Value("claims").(*Claims.Claims)
	deleteSession := helper.DeleteSession(uc.SessionID)
	if deleteSession != nil {
		log.Printf("Delete Error : Session can't be deleted")
		writer.WriteHeader(http.StatusInternalServerError)
		return
	}

}
func FetchAllTask(writer http.ResponseWriter, request *http.Request) {
	uc := request.Context().Value("claims").(*Claims.Claims)
	searchText := request.URL.Query().Get("searchText")
	taskStatus := request.URL.Query().Get("status")
	pageNO := request.URL.Query().Get("page")
	Page, _ := strconv.Atoi(pageNO)
	Limit, err := strconv.Atoi(request.URL.Query().Get("limit"))
	if Limit == 0 {
		Limit = 5
	}
	status := make(pq.StringArray, 0)
	var isStatus bool
	if taskStatus == "active" || taskStatus == "draft" {
		isStatus = true
		status = append(status, taskStatus)
	} else {
		status = append(status, "active")
		status = append(status, "draft")
	}
	user, err := helper.GetAllTask(uc.ID, searchText, status, isStatus, Page, Limit)
	if err != nil {
		log.Printf("Fetch Error : Todos can't be fetched")
		writer.WriteHeader(http.StatusBadRequest)
		return
	}
	errs := json.NewEncoder(writer).Encode(user)
	if errs != nil {
		writer.WriteHeader(http.StatusInternalServerError)
	}
}

func DeleteTask(writer http.ResponseWriter, request *http.Request) {
	TaskID := chi.URLParam(request, "id")

	deleteError := helper.DeleteTask(TaskID)
	if deleteError != nil {
		log.Printf("DeleteTask : Error in archiving this task")
		writer.WriteHeader(http.StatusInternalServerError)
		return
	}

}

func UpdateTask(writer http.ResponseWriter, request *http.Request) {
	var user models.Task

	err := json.NewDecoder(request.Body).Decode(&user)
	if err != nil {
		writer.WriteHeader(http.StatusBadRequest)
		return
	}
	TaskID := chi.URLParam(request, "id")
	updateError := helper.UpdateTasks(TaskID, user.Task)
	if updateError != nil {
		log.Printf("UpdateTask : Error in updating this task")
		writer.WriteHeader(http.StatusInternalServerError)
		return
	}
}
func CompleteTask(writer http.ResponseWriter, request *http.Request) {

	TaskID := chi.URLParam(request, "id")
	markComplete := helper.MarkComplete(TaskID)
	if markComplete != nil {
		log.Printf("Incomplete: Error in marking this task")
		writer.WriteHeader(http.StatusInternalServerError)
		return
	}

}
