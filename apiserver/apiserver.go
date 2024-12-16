package apiserver

import (
	"io/fs"
	"log"
	"net/http"
	"os"
	"webServer/chat"
	"webServer/store"

	"github.com/gorilla/sessions"
)

type ApiServer struct {
	Config       *Config
	store        *store.Store
	router       *http.ServeMux
	sessionStore *sessions.CookieStore
	wsChat       *chat.ChatManager
}

func New(store *store.Store, sessions *sessions.CookieStore) *ApiServer {
	return &ApiServer{
		Config:       NewConfig(),
		store:        store,
		router:       http.NewServeMux(),
		sessionStore: sessions,
		wsChat:       chat.NewChatManager(store),
	}
}

func (api *ApiServer) Start() error {
	server := http.Server{
		Addr:    api.Config.ServerAddr,
		Handler: api.router,
	}

	//Create user list
	if err := api.store.User().GetUsersList(); err != nil {
		log.Println("errror, select users from database", err.Error())
	}

	api.configureRouter()

	go api.wsChat.Run()

	log.Println("Server is started on address: " + server.Addr)

	return server.ListenAndServe()
}

func (s *ApiServer) configureRouter() {

	var frontend fs.FS = os.DirFS("frontend/dist")
	httpFS := http.FS(frontend)
	fileServer := http.FileServer(httpFS)
	serveIndex := serveFileContents("index.html", httpFS)

	s.router.Handle(`/`, intercept404(fileServer, serveIndex))
	s.router.HandleFunc(`/register`, s.createUser())
	s.router.HandleFunc(`/login`, s.authorization())
	s.router.HandleFunc(`/logout`, s.logout())
	s.router.HandleFunc(`/restore`, s.restorePassword())
	s.router.HandleFunc(`/userupdate`, s.updateUser())
	s.router.HandleFunc(`/changepass`, s.changePassword())
	s.router.HandleFunc(`/chat`, s.auntification(s.chat()))
	s.router.HandleFunc(`/username`, s.userName())
	s.router.HandleFunc(`/ws`, s.wsServe())

	//Тестовый обработчик
	s.router.HandleFunc(`/test`, s.test())
}
