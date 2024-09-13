package handler

import (
	"github.com/go-chi/chi"

	"github.com/go-chi/cors"

	"github.com/madeindra/mock-interview/server/internal/config"
	"github.com/madeindra/mock-interview/server/internal/data"
	"github.com/madeindra/mock-interview/server/internal/elevenlab"
	"github.com/madeindra/mock-interview/server/internal/middleware"
	"github.com/madeindra/mock-interview/server/internal/openai"
)

type handler struct {
	ai openai.Client
	el elevenlab.Client
	db *data.Database
}

func NewHandler(cfg config.AppConfig) *chi.Mux {
	h := &handler{
		ai: openai.NewOpenAI(cfg.APIKey),
		el: elevenlab.NewElevenLab(cfg.TTSAPIKey),
		db: data.New(cfg.DBPath),
	}

	r := chi.NewRouter()

	r.Use(cors.Handler(cors.Options{
		AllowedOrigins: cfg.CORSOrigins,
		AllowedMethods: cfg.CORSMethods,
		AllowedHeaders: cfg.CORSHeaders,
	}))

	r.Get("/chat/status", h.Status)
	r.Post("/chat/start", h.StartChat)

	r.Group(func(r chi.Router) {
		r.Use(middleware.BasicAuth)
		r.Post("/chat/answer", h.AnswerChat)
		r.Get("/chat/end", h.EndChat)
	})

	return r
}
