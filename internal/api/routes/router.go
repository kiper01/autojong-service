package routes

import (
	"main/internal/api/middlewares"
	"net/http"
	"strings"

	"github.com/google/uuid"
	httpSwagger "github.com/swaggo/http-swagger"
)

type IHandler interface {
	PostRequest(w http.ResponseWriter, r *http.Request)
	GetRequest(w http.ResponseWriter, r *http.Request, id uuid.UUID)
	ListOfRequests(w http.ResponseWriter, r *http.Request)
}

func NewRouter(
	auth middlewares.Auth,
	handler IHandler,
) http.Handler {
	mux := http.NewServeMux()

	mux.Handle("/swagger.json", http.FileServer(http.Dir("docs")))
	mux.Handle("/swagger/", httpSwagger.Handler(
		httpSwagger.URL("/swagger.json"),
	))

	mux.HandleFunc("/v3/api-docs/autojong-request-service", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "docs/swagger.json")
	})

	mux.HandleFunc("/api/v1/request/post", middlewares.CorsMiddleware(
		func(w http.ResponseWriter, r *http.Request) {
			switch r.Method {
			case http.MethodPost:
				http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					handler.PostRequest(w, r)
				}).ServeHTTP(w, r)
			case http.MethodOptions:
				return
			default:
				http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			}
		}))

	mux.HandleFunc("/api/v1/request/list", middlewares.CorsMiddleware(
		func(w http.ResponseWriter, r *http.Request) {
			switch r.Method {
			case http.MethodGet:
				middlewares.AuthMiddleware(
					middlewares.AdminMiddleware(
						http.HandlerFunc(
							func(w http.ResponseWriter, r *http.Request) {
								handler.ListOfRequests(w, r)
							})), auth).ServeHTTP(w, r)
			case http.MethodOptions:
				return
			default:
				http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			}
		}))

	mux.HandleFunc("/api/v1/request/get/", middlewares.CorsMiddleware(
		func(w http.ResponseWriter, r *http.Request) {
			parts := strings.Split(r.URL.Path, "/")
			id, err := uuid.Parse(parts[len(parts)-1])
			if err != nil || id == uuid.Nil {
				http.Error(w, "Bad request", http.StatusBadRequest)
				return
			}

			switch r.Method {
			case http.MethodGet:
				middlewares.AuthMiddleware(
					middlewares.AdminMiddleware(
						http.HandlerFunc(
							func(w http.ResponseWriter, r *http.Request) {
								handler.GetRequest(w, r, id)
							})), auth).ServeHTTP(w, r)
			case http.MethodOptions:
				return
			default:
				http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			}
		}))

	return mux
}
