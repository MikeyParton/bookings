package handlers

import (
	"net/http"

	"github.com/MikeyParton/bookings/pkg/config"
	"github.com/MikeyParton/bookings/pkg/models"
	"github.com/MikeyParton/bookings/pkg/render"
)

var Repo *Repository

type Repository struct {
	App *config.AppConfig
}

func NewRepo(a *config.AppConfig) *Repository {
	return &Repository{
		App: a,
	}
}

func NewHandlers(r *Repository) {
	Repo = r
}

func (m *Repository) Home(w http.ResponseWriter, r *http.Request) {
	remoteIP := r.RemoteAddr

	m.App.Session.Put(r.Context(), "remote_ip", remoteIP)
	render.RenderTemplate(w, "home.page.html", &models.TemplateData{})
}

// About is the about page handler
func (m *Repository) About(w http.ResponseWriter, r *http.Request) {
	remoteIP := m.App.Session.GetString(r.Context(), "remote_ip")
	render.RenderTemplate(w, "about.page.html", &models.TemplateData{
		StringMap: map[string]string{
			"test":      "greeting",
			"remote_ip": remoteIP,
		},
	})
}
