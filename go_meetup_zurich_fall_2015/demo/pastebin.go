package main

import (
	"log"
	"math/rand"
	"net/http"
	"strconv"
	"time"

	"github.com/StephanDollberg/go-json-rest-middleware-jwt"
	"github.com/ant0ine/go-json-rest/rest"
	"golang.org/x/crypto/bcrypt"
)

var PasteStore = make(map[uint32]*Paste)
var UserStore = map[string]*User{}

func StoreGet(id uint32) (*Paste, bool) {
	paste, ok := PasteStore[id]
	return paste, ok
}

func StoreCreate(paste *Paste) {
	paste.Id = rand.Uint32()
	PasteStore[paste.Id] = paste
}

func StoreCreateWithUser(*PasteWithUser) {

}

func StoreGetAll() []*Paste {
	res := make([]*Paste, 0)
	for _, v := range PasteStore {
		res = append(res, v)
	}

	return res
}

func StoreGetUser(id string) (*User, bool) {
	user, ok := UserStore[id]
	return user, ok
}

type User struct {
	Id       string
	Password []byte
}

type Paste struct {
	Id    uint32 `json:"id"`
	Title string `json:"title"`
	Body  string `json:"body"`
}

type PasteWithUser struct {
	Id    uint32 `json:"id"`
	Title string `json:"title"`
	Body  string `json:"body"`
	User  string `json:"user"`
}

func GetPaste(w rest.ResponseWriter, r *rest.Request) { // HL
	parsedId, err := strconv.ParseUint(r.PathParam("id"), 10, 32) // HL
	if err != nil {
		rest.NotFound(w, r)
		return
	}

	if paste, ok := StoreGet(uint32(parsedId)); ok {
		w.WriteJson(paste) // HL
		return
	}

	rest.Error(w, "not found", http.StatusNotFound)
}

func PostPaste(w rest.ResponseWriter, r *rest.Request) {
	paste := &Paste{}
	err := r.DecodeJsonPayload(&paste) // HL
	if err != nil {
		rest.Error(w, "Invalid Paste", http.StatusBadRequest)
		return
	}

	StoreCreate(paste)

	w.WriteJson(paste)
}

func PostPasteWithUser(w rest.ResponseWriter, r *rest.Request) {
	paste := &PasteWithUser{}
	err := r.DecodeJsonPayload(&paste)
	if err != nil {
		rest.Error(w, "Invalid Paste", http.StatusBadRequest)
		return
	}

	paste.User = r.Env["REMOTE_USER"].(string) // HL

	StoreCreateWithUser(paste)

	w.WriteJson(paste)
}

func GetAll(w rest.ResponseWriter, r *rest.Request) {
	all := StoreGetAll()
	w.WriteJson(all)
}

func authUser(userId string, passwordClaim string) bool {
	if user, ok := StoreGetUser(userId); ok {
		return bcrypt.CompareHashAndPassword(user.Password, []byte(passwordClaim)) == nil
	}
	return false
}

func MakeApiSimple() *rest.Api {
	api := rest.NewApi()             // HL
	api.Use(rest.DefaultDevStack...) // HL

	router, err := rest.MakeRouter( // HL
		&rest.Route{"GET", "/pastes/:id", GetPaste},
		&rest.Route{"POST", "/pastes", PostPaste},
		&rest.Route{"GET", "/pastes", GetAll},
	)

	if err != nil {
		log.Fatal(err)
	}

	api.SetApp(router) // HL
	return api
}

func MakeApi() *rest.Api {
	api := rest.NewApi()
	api.Use(rest.DefaultDevStack...)

	jwtMiddleware := &jwt.JWTMiddleware{
		Key:           []byte("super secret key"),
		Realm:         "Spacebook",
		Timeout:       time.Hour,
		MaxRefresh:    time.Hour * 24,
		Authenticator: authUser,
	}

	api.Use(&rest.IfMiddleware{
		Condition: func(request *rest.Request) bool {
			return request.URL.Path == "/pasteswithuser"
		},
		IfTrue: jwtMiddleware,
	})

	router, err := rest.MakeRouter(
		&rest.Route{"GET", "/pastes/:id", GetPaste},
		&rest.Route{"POST", "/pastes", PostPaste},
		&rest.Route{"POST", "/pasteswithuser", PostPasteWithUser},
		&rest.Route{"GET", "/pastes", GetAll},
		&rest.Route{"POST", "/login", jwtMiddleware.LoginHandler},
	)

	if err != nil {
		log.Fatal(err)
	}

	api.SetApp(router)
	return api
}

func main() {

	api := MakeApi()
	http.Handle("/", http.FileServer(http.Dir("client/app/")))
	http.Handle("/api/", http.StripPrefix("/api", api.MakeHandler()))
	http.ListenAndServe("127.0.0.1:20001", nil)
}

func MyRestHandler(w rest.ResponseWriter, r *rest.Request) {
	w.WriteJson(map[string]string{"hello": "world"})
}

func SetupGoJSONRest() {
	api := rest.NewApi()

	router, _ := rest.MakeRouter(
		&rest.Route{"GET", "/", MyRestHandler},
	)

	api.SetApp(router)
	http.Handle("/", api.MakeHandler())
	http.ListenAndServe("localhost:8080", nil)
}
