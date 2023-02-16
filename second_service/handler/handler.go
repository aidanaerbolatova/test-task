package handler

import (
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"log"
	"net/http"
	"net/rpc"
	"regexp"
	"serv2/models"

	"github.com/go-chi/chi"
)

type Handler struct {
	RpcClient *rpc.Client
}

func (h *Handler) CreateUser(w http.ResponseWriter, r *http.Request) {
	var user *models.User
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		log.Fatal(err)
	}

	req, err := http.NewRequest("POST", "http://localhost:8082/generate-salt", nil)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	type response struct {
		Salt string `json:"salt"`
	}

	var generatedSalt *response
	if err := json.NewDecoder(res.Body).Decode(&generatedSalt); err != nil {
		log.Fatal(err)
	}

	user.Password = createHash(user.Password, generatedSalt.Salt)

	var userReply *models.User
	if err = h.RpcClient.Call("App.GetUser", user.Email, &userReply); err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	if userReply.Email != "" || !checkUser(user.Email) {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	var reply string
	if err = h.RpcClient.Call("App.CreateUser", user, &reply); err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
}

func (h *Handler) GetUser(w http.ResponseWriter, r *http.Request) {
	email := chi.URLParam(r, "email")
	var user *models.User

	if err := h.RpcClient.Call("App.GetUser", email, &user); err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	if user.Email == "" {
		http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
		return
	}

	if err := json.NewEncoder(w).Encode(user); err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
	}
}

func checkUser(email string) bool {
	return regexp.MustCompile(`^[a-z0-9._%+\-]+@[a-z0-9.\-]+\.[a-z]{2,4}$`).MatchString(email)
}

func createHash(password, salt string) string {
	hasher := md5.New()
	hasher.Write([]byte(salt + password))
	return hex.EncodeToString(hasher.Sum(nil))
}
