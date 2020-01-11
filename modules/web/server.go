package web

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"time"

	"galched-bot/modules/settings"
	"galched-bot/modules/youtube"
)

type (
	WebServer struct {
		server http.Server
		r      *youtube.Requester

		users  map[string]string
		authed map[string]struct{}
	}
)

func New(s *settings.Settings, r *youtube.Requester) *WebServer {
	srv := http.Server{
		Addr: s.QueueAddress,
	}

	webServer := &WebServer{
		server: srv,
		r:      r,

		users:  s.LoginUsers,
		authed: make(map[string]struct{}, 10),
	}

	http.HandleFunc("/", func(writer http.ResponseWriter, request *http.Request) {
		if request.Method != http.MethodGet {
			return
		}

		if webServer.IsAuthorized(request) {
			http.ServeFile(writer, request, "web/index.html")
		} else {
			http.Redirect(writer, request, "/login", 301)
		}
	})

	http.HandleFunc("/login", func(writer http.ResponseWriter, request *http.Request) {
		if request.Method == http.MethodGet {
			http.ServeFile(writer, request, "web/login.html")
			return
		} else if request.Method != http.MethodPost {
			return
		}

		login := request.FormValue("login")
		pwd := request.FormValue("password")

		log.Print("web: trying to log in with user: ", login)

		if webServer.IsRegistered(login, pwd) {
			webServer.Authorize(writer)
		} else {
			log.Print("web: incorrect password attempt")
		}
		http.Redirect(writer, request, "/", http.StatusSeeOther)
	})

	http.HandleFunc("/scripts.js", func(writer http.ResponseWriter, request *http.Request) {
		if request.Method != http.MethodGet {
			return
		}

		http.ServeFile(writer, request, "web/scripts.js")
	})

	http.HandleFunc("/style.css", func(writer http.ResponseWriter, request *http.Request) {
		if request.Method != http.MethodGet {
			return
		}

		http.ServeFile(writer, request, "web/style.css")
	})

	http.HandleFunc("/queue", func(writer http.ResponseWriter, request *http.Request) {
		if !webServer.IsAuthorized(request) {
			http.Error(writer, "not authorized", http.StatusUnauthorized)
			return
		}

		switch request.Method {
		case http.MethodGet:
		case http.MethodPost:
			body, err := ioutil.ReadAll(request.Body)
			if err != nil {
				log.Print("web: cannot read body msg, %v", err)
				return
			}
			id := string(body)
			if len(id) != youtube.YoutubeIDLength && len(id) > 0 {
				log.Printf("web: incorrect data in body, <%s>", id)
				return
			}
			r.Remove(id)
		default:
			return
		}
		writer.Header().Set("Content-Type", "application/json")
		writer.Header().Set("Access-Control-Allow-Origin", "*")
		json.NewEncoder(writer).Encode(webServer.r.List())
	})

	return webServer
}

func (s WebServer) Start() error {
	go func() {
		s.server.ListenAndServe()
	}()
	return nil
}

func (s WebServer) Stop(ctx context.Context) error {
	return s.server.Shutdown(ctx)
}

func (s WebServer) IsAuthorized(request *http.Request) bool {
	if cookie, err := request.Cookie("session"); err == nil {
		if _, ok := s.authed[cookie.Value]; ok {
			return true
		}
	}
	return false
}

func (s WebServer) IsRegistered(login, pwd string) bool {
	return s.users[login] == pwd
}

func (s WebServer) Authorize(response http.ResponseWriter) {
	var byteKey = make([]byte, 16)
	rand.Read(byteKey[:])
	stringKey := hex.EncodeToString(byteKey)

	s.authed[stringKey] = struct{}{}
	log.Print("web: authenticated new user")

	expires := time.Now().AddDate(0, 1, 0)
	http.SetCookie(response, &http.Cookie{
		Name:    "session",
		Value:   stringKey,
		Path:    "/",
		Expires: expires,
	})
}
