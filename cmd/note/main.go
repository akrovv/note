package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"note/config"
	"note/internal/adapter/repository"
	"note/internal/controllers/http/middleware"
	v1 "note/internal/controllers/http/v1"
	_ "note/internal/logger"
	"note/internal/service"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
	"go.uber.org/zap"
)

const openCons = 10

func main() {
	cfg, err := config.LoadConfig(".")

	if err != nil {
		log.Printf("error: %s\n", err)
		return
	}

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?", cfg.DBUser, cfg.DBDriver, cfg.DBHost, cfg.DBPort, cfg.DBName)
	dsn += "&charset=utf8"
	dsn += "&interpolateParams=true"
	dsn += "&parseTime=true"

	db, err := sql.Open("mysql", dsn)

	if err != nil {
		log.Printf("error: %s\n", err)
		return
	}
	defer db.Close()

	db.SetMaxOpenConns(openCons)

	err = db.Ping()

	if err != nil {
		log.Printf("error: %s\n", err)
		return
	}

	var (
		logger    = zap.L()
		notesRepo = repository.NewStorage(db)
	)

	handlerIndex := v1.NewIndexHandler(logger)
	noteService := service.NewService(notesRepo)
	handlersNotes := v1.NewNoteHandler(noteService, logger)

	router := mux.NewRouter()

	router.HandleFunc("/", handlerIndex.Index)
	router.HandleFunc("/note/{id}", handlersNotes.GetByID).Methods("GET")
	router.HandleFunc("/note", handlersNotes.Create).Methods("POST")
	router.HandleFunc("/note/{id}", handlersNotes.UpdateByID).Methods("PUT")
	router.HandleFunc("/note/{id}", handlersNotes.DeleteByID).Methods("DELETE")
	router.HandleFunc("/note", handlersNotes.GetAll).Methods("GET")

	siteMux := middleware.Logger(router, logger)

	logger.Info("Listennig on :8080")
	err = http.ListenAndServe(":8080", siteMux)

	if err != nil {
		log.Printf("error: %s\n", err)
	}
}
