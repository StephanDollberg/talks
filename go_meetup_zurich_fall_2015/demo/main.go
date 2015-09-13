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

type User struct {
	Id       string
	Password []byte
}

type Pic struct {
	Id          uint32   `json:"id"`
	Title       string   `json:"title"`
	Description string   `json:"description"`
	ImageUrl    string   `json:"image_url"`
	Comments    []uint32 `json:"comments"`
	Liked       []string `json:"liked"`
}

type Comment struct {
	Id     uint32 `json:"id"`
	Body   string `json:"body"`
	UserId string `json:"userid"`
	PicId  uint32 `json:"picid"`
}

var commentStore = map[uint32]*Comment{}
var picStore = map[uint32]*Pic{}
var userStore = map[string]*User{}

func hashPassword(password string) []byte {
	pw, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)

	if err != nil {
		log.Fatal(err)
	}

	return pw
}

func initState() {
	userStore["Alice"] = &User{"Alice", hashPassword("alice")}
	userStore["Bob"] = &User{"Bob", hashPassword("bob")}
	userStore["Eve"] = &User{"Eve", hashPassword("eve")}

	picStore[1] = &Pic{1, "Moon", "Moon from the ISS",
		"img/pics/moon.jpg", []uint32{1, 2}, []string{"Alice", "Bob", "Eve"}}
	picStore[2] = &Pic{2, "Solar Eclipse", "Solar eclipse from ISS ",
		"img/pics/solareclipse.jpg", []uint32{3, 4}, []string{"Alice", "Eve"}}
	picStore[3] = &Pic{3, "In the Cloud", "Clouds over the tropical waters of the West Pacific Ocean.",
		"img/pics/typhoon.jpg", []uint32{}, []string{}}
	picStore[4] = &Pic{4, "Aurora", "Aurora with solar pannels",
		"img/pics/aurora.jpg", []uint32{}, []string{}}

	commentStore[1] = &Comment{1, "So cool!", "Alice", 1}
	commentStore[2] = &Comment{2, "Indeed", "Bob", 1}
	commentStore[3] = &Comment{3, "Love it", "Eve", 2}
	commentStore[4] = &Comment{4, "+1", "Alice", 2}
}

func singleComment(w rest.ResponseWriter, r *rest.Request) {
	parsedId, err := strconv.ParseUint(r.PathParam("id"), 10, 32)
	if err != nil {
		rest.NotFound(w, r)
		return
	}

	if pic, ok := commentStore[uint32(parsedId)]; ok {
		w.WriteJson(pic)
	} else {
		rest.NotFound(w, r)
	}
}

func singlePic(w rest.ResponseWriter, r *rest.Request) {
	parsedId, err := strconv.ParseUint(r.PathParam("id"), 10, 32)
	if err != nil {
		rest.NotFound(w, r)
		return
	}

	if pic, ok := picStore[uint32(parsedId)]; ok {
		w.WriteJson(pic)
	} else {
		rest.NotFound(w, r)
	}
}

func likePic(w rest.ResponseWriter, r *rest.Request) {
	parsedId, err := strconv.ParseUint(r.PathParam("id"), 10, 32)
	if err != nil {
		rest.NotFound(w, r)
		return
	}

	if pic, ok := picStore[uint32(parsedId)]; ok {
		for _, v := range pic.Liked {
			if v == r.Env["REMOTE_USER"].(string) {
				w.WriteHeader(http.StatusOK)
				return
			}
		}

		pic.Liked = append(pic.Liked, r.Env["REMOTE_USER"].(string))
		w.WriteHeader(http.StatusCreated)
	}
}

func allPics(w rest.ResponseWriter, r *rest.Request) {
	pics := make([]*Pic, 0)

	for _, pic := range picStore {
		pics = append(pics, pic)
	}

	w.WriteJson(pics)
}

func createComment(w rest.ResponseWriter, r *rest.Request) {
	c := &Comment{}
	err := r.DecodeJsonPayload(c)
	if err != nil {
		rest.Error(w, "invalid comment", http.StatusBadRequest)
		return
	}

	if pic, ok := picStore[c.PicId]; ok {
		c.UserId = r.Env["REMOTE_USER"].(string)
		c.Id = rand.Uint32()
		commentStore[c.Id] = c
		pic.Comments = append(pic.Comments, c.Id)
		w.WriteJson(c)
	}
}

func authUser(userId string, passwordClaim string) bool {
	if user, ok := userStore[userId]; ok {
		return bcrypt.CompareHashAndPassword(user.Password, []byte(passwordClaim)) == nil
	}
	return false
}

func main() {
	initState()

	jwtMiddleware := &jwt.JWTMiddleware{
		Key:           []byte("super secret key"),
		Realm:         "Spacebook",
		Timeout:       time.Hour,
		MaxRefresh:    time.Hour * 24,
		Authenticator: authUser,
	}

	api := rest.NewApi()
	api.Use(rest.DefaultDevStack...)
	api.Use(&rest.IfMiddleware{
		Condition: func(request *rest.Request) bool {
			return request.URL.Path != "/login" && request.Method == "POST"
		},
		IfTrue: jwtMiddleware,
	})

	router, err := rest.MakeRouter(
		&rest.Route{"GET", "/pics/:id", singlePic},
		&rest.Route{"POST", "/pics/:id/like", likePic},
		&rest.Route{"GET", "/pics", allPics},
		&rest.Route{"GET", "/comments/:id", singleComment},
		&rest.Route{"POST", "/comments", createComment},
		&rest.Route{"POST", "/login", jwtMiddleware.LoginHandler},
	)

	if err != nil {
		log.Fatal(err)
	}

	api.SetApp(router)

	http.Handle("/", http.FileServer(http.Dir("client/app/")))
	http.Handle("/api/", http.StripPrefix("/api", api.MakeHandler()))
	http.ListenAndServe("localhost:3001", nil)
}

func MyHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Hello World"))
}

func MyRestHandler(w rest.ResponseWriter, r *rest.Request) {
	w.WriteJson(map[string]string{"hello": "world"})
}
