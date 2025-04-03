package main

import "net/http"

func (cfg *apiConfig) handleAdminDashboard(w http.ResponseWriter, r *http.Request) {
	cfg.Render(w, r, "templates/pages/admin.html", nil)
}
