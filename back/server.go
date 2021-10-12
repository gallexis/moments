package back

import (
	"fmt"
	"net/http"
	"runtime/debug"
	"time"

	"steelseries/back/controller"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
)

type ApiServer struct {
	router *mux.Router
	server *http.Server
}

func (a ApiServer) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	defer func() {
		if r := recover(); r != nil {
			fmt.Printf("recovered from panic : [%v] - stack trace : \n [%s]\n", r, debug.Stack())

			writer.WriteHeader(http.StatusInternalServerError)
			writer.Write([]byte("Server error\n"))
		}
	}()

	// handle CORS
	writer.Header().Set("Access-Control-Allow-Origin", "*")
	writer.Header().Set("Access-Control-Allow-Headers", "x-api-key, Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, access-control-allow-origin, access-control-allow-headers")
	writer.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
	if request.Method == "OPTIONS" {
		return
	}

	a.router.ServeHTTP(writer, request)
}

func (s *ApiServer) Start() error {
	fmt.Println("Server listening on http://localhost:8899")
	return s.server.ListenAndServe()
}
func (s *ApiServer) Close() error {
	fmt.Println("Closing connection...")
	return s.server.Close()
}

func NewApiServer() ApiServer {
	server := ApiServer{
		router: mux.NewRouter().Schemes("https", "http", "").Subrouter(),
	}

	server.server = &http.Server{
		Addr:         ":8899",
		Handler:      handlers.CompressHandler(server),
		ReadTimeout:  time.Second * 30,
		WriteTimeout: time.Second * 120,
		IdleTimeout:  time.Second * 60,
	}

	server.router.HandleFunc("/moments", controller.MomentsController).
		Methods("POST").
		Name("sendToken")

	return server
}
