package apiserver

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/Vladislav001/golang_study_http_rest_api/internal/app/model"
	"github.com/Vladislav001/golang_study_http_rest_api/internal/app/store"
	"github.com/google/uuid"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/gorilla/sessions"
	"github.com/sirupsen/logrus"
	"net/http"
	"time"
)

// Более простая версия сервера - не знает про запуск сервера, про http, а
// будет уметь только обрабатывать входящий запрос

const (
	sessionName        = "vladstudy"
	ctxKeyUser  ctxKey = iota
	ctxKeyRequestID
)

var (
	errIncorrectEmailOrPassword = errors.New("incorrect email or password")
	errNotAuthenticated         = errors.New("not authenticated")
)

type ctxKey int8

type server struct {
	router       *mux.Router
	logger       *logrus.Logger
	store        store.Store
	sessionStore sessions.Store
}

func newServer(store store.Store, sessionStore sessions.Store) *server {
	s := &server{
		router:       mux.NewRouter(),
		logger:       logrus.New(),
		store:        store,
		sessionStore: sessionStore,
	}

	s.configureRouter()

	return s
}

func (s *server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.router.ServeHTTP(w, r)
}

func (s *server) configureRouter() {
	// у каждого запроса будет уникальный X-Request-ID
	s.router.Use(s.setRequestID)
	s.router.Use(s.logRequest)
	// разрешить запросы с любых источников/доменов (если напр.с бразуера из за разных портов возникнет CORS)
	s.router.Use(handlers.CORS(handlers.AllowedOrigins([]string{"*"})))

	s.router.HandleFunc("/registration", s.handleUserRegistration()).Methods("POST")
	s.router.HandleFunc("/auth", s.handleUserAuth()).Methods("POST")

	// middleware будет работать для url-ов вида /private/...
	private := s.router.PathPrefix("/private").Subrouter()
	private.Use(s.authenticateUser)
	private.HandleFunc("/profile", s.handleGetProfile()).Methods("GET")

	// 404
	s.router.NotFoundHandler = s.handle404()
}

func (s *server) setRequestID(next http.Handler) http.Handler {
	return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		id := uuid.New().String()
		writer.Header().Set("X-Request-ID", id)
		next.ServeHTTP(writer, request.WithContext(context.WithValue(request.Context(), ctxKeyRequestID, id)))
	})
}

func (s *server) logRequest(next http.Handler) http.Handler {
	return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		logger := s.logger.WithFields(logrus.Fields{
			"remote_addr": request.RemoteAddr,
			"request_id":  request.Context().Value(ctxKeyRequestID),
		})

		// пример: started GET /endpoint
		logger.Infof("started %s %s", request.Method, request.RequestURI)

		start := time.Now()

		// т.к ResponseWriter - интерфейс, то определили свой responseWriter, чтобы http код залогировать
		rw := &responseWriter{writer, http.StatusOK}
		next.ServeHTTP(rw, request)

		logger.Infof(
			"completed with %d %s in %v",
			rw.code,
			http.StatusText(rw.code),
			time.Now().Sub(start),
		)
	})
}

func (s *server) authenticateUser(next http.Handler) http.Handler {
	return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		session, err := s.sessionStore.Get(request, sessionName)
		if err != nil {
			s.error(writer, request, http.StatusInternalServerError, err)
			return
		}

		id, ok := session.Values["user_id"]
		if !ok {
			s.error(writer, request, http.StatusUnauthorized, errNotAuthenticated)
			return
		}

		u, err := s.store.User().Find(id.(int))
		if err != nil {
			s.error(writer, request, http.StatusUnauthorized, errNotAuthenticated)
			return
		}

		// прикрепляем юзера к контексту запроса (чтобы в след. миддлевер и т.п не искать заного)
		next.ServeHTTP(writer, request.WithContext(context.WithValue(request.Context(), ctxKeyUser, u)))
	})
}

func (s *server) handleUserRegistration() http.HandlerFunc {
	// Для "отражения" входящих данных
	type request struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	return func(w http.ResponseWriter, r *http.Request) {
		req := &request{}
		if err := json.NewDecoder(r.Body).Decode(req); err != nil {
			s.error(w, r, http.StatusBadRequest, err)
			return
		}

		u := &model.User{
			Email:    req.Email,
			Password: req.Password,
		}

		if err := s.store.User().Create(u); err != nil {
			s.error(w, r, http.StatusUnprocessableEntity, err)
			return
		}

		u.Sanitize()
		s.respond(w, r, http.StatusCreated, u)
	}
}

func (s *server) handleUserAuth() http.HandlerFunc {
	// Для "отражения" входящих данных
	type request struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	return func(w http.ResponseWriter, r *http.Request) {
		req := &request{}
		if err := json.NewDecoder(r.Body).Decode(req); err != nil {
			s.error(w, r, http.StatusBadRequest, err)
			return
		}

		u, err := s.store.User().FindByEmail(req.Email)

		if err != nil || !u.ComparePassword(req.Password) {
			s.error(w, r, http.StatusUnauthorized, errIncorrectEmailOrPassword)
			return
		}

		session, err := s.sessionStore.Get(r, sessionName)
		if err != nil {
			s.error(w, r, http.StatusInternalServerError, err)
			return
		}

		session.Values["user_id"] = u.ID
		if err := s.sessionStore.Save(r, w, session); err != nil {
			s.error(w, r, http.StatusInternalServerError, err)
			return
		}

		s.respond(w, r, http.StatusOK, nil)
	}
}

func (s *server) handleGetProfile() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// на этом шаге уже подразумеваем, что юзер залогинен и записан в наш контекст
		s.respond(w, r, http.StatusOK, r.Context().Value(ctxKeyUser).(*model.User))
	}
}

func (s *server) handle404() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		s.respond(w, r, http.StatusNotFound, map[string]string{"error": `Page not found`})
	}
}

func (s *server) error(w http.ResponseWriter, r *http.Request, code int, err error) {
	s.respond(w, r, code, map[string]string{"error": err.Error()})
}

func (s *server) respond(w http.ResponseWriter, r *http.Request, code int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)

	if data != nil {
		json.NewEncoder(w).Encode(data)
	}
}
