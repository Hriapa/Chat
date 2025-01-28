package apiserver

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
	"webServer/chat"
	"webServer/model"
	"webServer/protocol"
)

const (
	sessionName        = `authorizing`
	ctxKeyUser  ctxKey = iota
)

type ctxKey int

// Регистрация нового пользователя
func (s *ApiServer) createUser() http.HandlerFunc {
	type request struct {
		UserName string `json:"username"`
		Password string `json:"password"`
	}
	return func(w http.ResponseWriter, r *http.Request) {
		setupCORS(&w)
		if r.Method != http.MethodPost {
			return
		}
		if !s.store.Connect {
			s.error(w, http.StatusInternalServerError, errConnectToDatabase)
			return
		}
		req := &request{}

		if err := json.NewDecoder(r.Body).Decode(req); err != nil {
			s.error(w, http.StatusBadRequest, err)
			return
		}
		user := &model.User{
			Login:    req.UserName,
			Password: req.Password,
		}

		if err := s.store.User().AddUser(user); err != nil {
			s.error(w, http.StatusUnprocessableEntity, err)
			return
		}

		//user.Private()
		//s.result(w, r, http.StatusCreated, user)

		session, err := s.sessionStore.Get(r, sessionName)
		if err != nil {
			s.error(w, http.StatusInternalServerError, err)
		}
		session.Values["user_id"] = user.Id

		err = s.sessionStore.Save(r, w, session)
		if err != nil {
			s.error(w, http.StatusInternalServerError, err)
		}

		//Добавляем пользователя в список

		s.wsChat.NewUser <- &model.UserName{Id: user.Id, Name: user.Login}

		s.result(w, http.StatusOK, nil)
	}
}

// Авторизация
func (s *ApiServer) authorization() http.HandlerFunc {
	type request struct {
		UserName string `json:"username"`
		Password string `json:"password"`
	}
	return func(w http.ResponseWriter, r *http.Request) {
		setupCORS(&w)
		if r.Method != http.MethodPost {
			return
		}
		if !s.store.Connect {
			s.error(w, http.StatusInternalServerError, errConnectToDatabase)
			return
		}
		req := &request{}

		if err := json.NewDecoder(r.Body).Decode(req); err != nil {
			s.error(w, http.StatusBadRequest, err)
			return
		}
		user, err := s.store.User().FindByLogin(req.UserName)
		if err != nil || !user.CheckPassword(req.Password) {
			s.error(w, http.StatusUnauthorized, errIncorrectLoginOrPassword)
			return
		}

		session, err := s.sessionStore.Get(r, sessionName)
		if err != nil {
			s.error(w, http.StatusInternalServerError, err)
		}
		session.Values["user_id"] = user.Id

		err = s.sessionStore.Save(r, w, session)
		if err != nil {
			s.error(w, http.StatusInternalServerError, err)
		}

		s.result(w, http.StatusOK, nil)
	}
}

// Аунтификация
func (s *ApiServer) auntification(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		setupCORS(&w)
		session, err := s.sessionStore.Get(r, sessionName)
		if err != nil {
			s.error(w, http.StatusInternalServerError, err)
			return
		}
		id, ok := session.Values["user_id"]
		if !ok {
			s.error(w, http.StatusUnauthorized, errNotAuthenticated)
			return
		}

		_, err = s.store.User().FindById(id.(int))
		if err != nil {
			s.error(w, http.StatusUnauthorized, err)
			return
		}
		next.ServeHTTP(w, r.WithContext(context.WithValue(r.Context(), ctxKeyUser, id.(int))))
	})
}

// Сброс пароля
func (s *ApiServer) restorePassword() http.HandlerFunc {
	type request struct {
		UserName string `json:"username"`
	}
	type response struct {
		Password string `json:"password"`
	}
	return func(w http.ResponseWriter, r *http.Request) {
		setupCORS(&w)
		if r.Method != http.MethodPost {
			return
		}
		req := &request{}
		if err := json.NewDecoder(r.Body).Decode(req); err != nil {
			s.error(w, http.StatusBadRequest, err)
			return
		}
		user, err := s.store.User().FindByLogin(req.UserName)
		if err != nil {
			s.error(w, http.StatusInternalServerError, errUserNotExist)
			return
		}
		user.CreateNewPassword()
		user.EncryptPassword()
		if err = s.store.User().UpdatePassword(user); err != nil {
			log.Println("error update user password:", err.Error())
			s.error(w, http.StatusInternalServerError, err)
			return
		}
		resp := &response{}
		resp.Password = user.Password
		s.result(w, http.StatusOK, resp)
	}
}

// Выход из авторизации
func (s *ApiServer) logout() http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {
		setupCORS(&w)
		if r.Method != http.MethodPost {
			return
		}
		session, err := s.sessionStore.Get(r, sessionName)
		if err != nil {
			return
		}
		session.Options.MaxAge = -1
		session.Save(r, w)
		http.Redirect(w, r, `/login`, http.StatusAccepted)
	}
}

// Отправка имени пользователя

func (s *ApiServer) userName() http.HandlerFunc {
	type response struct {
		Name string `json:"name"`
	}
	return func(w http.ResponseWriter, r *http.Request) {
		setupCORS(&w)
		if r.Method != http.MethodGet {
			return
		}
		var (
			err  error
			name string
		)
		user := &model.UserInfo{}
		user.Id, err = strconv.Atoi(r.URL.Query().Get("id"))
		if err != nil {
			log.Println(err)
			s.error(w, http.StatusInternalServerError, err)
			return
		}
		err = s.store.User().GetUserInfo(user)
		if err != nil {
			log.Println(err)
			s.error(w, http.StatusInternalServerError, err)
			return
		}
		if user.Name == "" {
			name = user.Login
		} else {
			name = user.CreateUserName()
		}
		resp := &response{
			Name: name,
		}
		s.result(w, http.StatusOK, resp)
	}
}

// Смена пароля

func (s *ApiServer) changePassword() http.HandlerFunc {
	type request struct {
		Id  int    `json:"id"`
		Old string `json:"old"`
		New string `json:"new"`
	}
	return func(w http.ResponseWriter, r *http.Request) {
		setupCORS(&w)
		if r.Method != http.MethodPost {
			return
		}
		password := &request{}
		if err := json.NewDecoder(r.Body).Decode(password); err != nil {
			s.error(w, http.StatusInternalServerError, err)
			return
		}

		user, err := s.store.User().FindById(password.Id)
		if err != nil {
			s.error(w, http.StatusUnauthorized, errUserNotFound)
			return
		}
		// Проверяем пароль
		if !user.CheckPassword(password.Old) {
			s.error(w, http.StatusUnauthorized, errIncorrectPassword)
			return
		}
		//Обновляем пароль
		user.Password = password.New
		user.EncryptPassword()
		if err = s.store.User().UpdatePassword(user); err != nil {
			log.Println("error update user password:", err.Error())
			s.error(w, http.StatusInternalServerError, err)
			return
		}
		s.result(w, http.StatusOK, nil)
	}
}

// Запуск чата
func (s *ApiServer) chat() http.HandlerFunc {
	type response struct {
		User int `json:"user"`
	}
	return func(w http.ResponseWriter, r *http.Request) {
		setupCORS(&w)
		if r.Method != http.MethodGet {
			return
		}

		//формируем список пользователей если он пустой
		if len(s.wsChat.Store.User().UsersList) == 0 {
			if err := s.store.User().GetUsersList(); err != nil {
				log.Println("errror, select users from database", err.Error())
				s.result(w, http.StatusInternalServerError, nil)
				return
			}
		}
		user := r.Context().Value(ctxKeyUser).(int)

		resp := &response{
			User: user,
		}

		//отправляем на клиент
		s.result(w, http.StatusOK, resp)
	}
}

// Работа с Websocket
func (s *ApiServer) wsServe() http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {
		setupCORS(&w)

		id, err := strconv.Atoi(r.URL.Query().Get("id"))
		if err != nil {
			log.Println(err)
		}

		s.wsChat.Upgrader.CheckOrigin = func(r *http.Request) bool { return true }

		con, err := s.wsChat.Upgrader.Upgrade(w, r, nil)
		if err != nil {
			log.Println(err)
			return
		}

		log.Println("Client connected:", con.RemoteAddr().String())

		if len(s.wsChat.Store.User().UsersList) == 0 {
			err = s.store.User().GetUsersList()
			if err != nil {
				log.Println("error in select user list function", err.Error())
			}
		}

		client := chat.NewClient(int(id), con, s.wsChat)

		if val, ok := s.wsChat.Store.User().UsersList[int(id)]; ok {
			val.Online = true
		}

		go client.ReadMessages()

		go client.WriteMessages()

		//проверка хранилища собираемых сообщений
		go client.AssemblerStoreCleaner()

		// Обработка ошибок, когда из хранилища удалено недособранное сообщение по timeout
		go client.AssemblingStoreErrorProcessing()

		// Отправляем список пользователей
		client.SendUsersList()
		// Отправляем остальным клиентам изменение статуса на online
		s.wsChat.Register <- client
	}
}

func (s *ApiServer) error(w http.ResponseWriter, code int, err error) {
	s.result(w, code, err.Error())
}

func (s *ApiServer) result(w http.ResponseWriter, code int, data interface{}) {
	w.WriteHeader(code)
	if data != nil {
		json.NewEncoder(w).Encode(data)
	}
}

func setupCORS(w *http.ResponseWriter) {
	(*w).Header().Set("Access-Control-Allow-Origin", "http://localhost:8080")
	(*w).Header().Set("Access-Control-Allow-Methods", "POST, PUT, GET, OPTIONS, PATCH, DELETE")
	(*w).Header().Set("Content-Type", "application/json")
	(*w).Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
	(*w).Header().Set("Access-Control-Allow-Credentials", "true")
}

// SPA

type hookedResponseWriter struct {
	http.ResponseWriter
	got404 bool
}

func (hrw *hookedResponseWriter) WriteHeader(status int) {
	if status == http.StatusNotFound {
		// Don't actually write the 404 header, just set a flag.
		hrw.got404 = true
	} else {
		hrw.ResponseWriter.WriteHeader(status)
	}
}

func (hrw *hookedResponseWriter) Write(p []byte) (int, error) {
	if hrw.got404 {
		// No-op, but pretend that we wrote len(p) bytes to the writer.
		return len(p), nil
	}

	return hrw.ResponseWriter.Write(p)
}

func intercept404(handler, on404 http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		hookedWriter := &hookedResponseWriter{ResponseWriter: w}
		handler.ServeHTTP(hookedWriter, r)

		if hookedWriter.got404 {
			on404.ServeHTTP(w, r)
		}
	})
}

func serveFileContents(file string, files http.FileSystem) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Restrict only to instances where the browser is looking for an HTML file
		if !strings.Contains(r.Header.Get("Accept"), "text/html") {
			w.WriteHeader(http.StatusNotFound)
			fmt.Fprint(w, "404 not found")

			return
		}

		// Open the file and return its contents using http.ServeContent
		index, err := files.Open(file)
		if err != nil {
			w.WriteHeader(http.StatusNotFound)
			fmt.Fprintf(w, "%s not found", file)

			return
		}

		fi, err := index.Stat()
		if err != nil {
			w.WriteHeader(http.StatusNotFound)
			fmt.Fprintf(w, "%s not found", file)

			return
		}

		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		http.ServeContent(w, r, fi.Name(), fi.ModTime(), index)
	}
}

//Тестовый обработчик для отладки протокола

func (s *ApiServer) test() http.HandlerFunc {
	var request = make([]byte, 2)
	var raw []byte
	return func(w http.ResponseWriter, r *http.Request) {
		setupCORS(&w)
		if r.Method != http.MethodPost {
			return
		}

		r.Body.Read(request)

		control := protocol.NewControlMessage()
		data := protocol.NewDataMessage()
		ack := protocol.NewAckMessage()

		switch request[0] {
		case uint8(protocol.ControlMessageTitle):
			raw = control.ControlCommandTestCoder(request[1])
		case uint8(protocol.DataMessageTitle):
			raw = data.DataMessageTestCoder(request[1])
		case uint8(protocol.AckMessageTitle):
			raw = ack.AckMessageTestCoder(request[1])
		}

		w.Header().Set("Content-Type", "data")
		_, err := w.Write(raw)
		if err != nil {
			log.Println(err)
		}
	}
}
