package handlers

import (
	"strings"
)

// cleanImagePath deja solo la ruta relativa del asset
// Ej: "../assets/images/celeste_images/img.png" -> "celeste_images/img.png"
func cleanImagePath(path string) string {
	// Normalizamos separadores a slash normal
	path = strings.ReplaceAll(path, "\\", "/")

	// Quitamos el prefijo ../assets/images/ si esta
	if idx := strings.Index(path, "assets/images/"); idx != -1 {
		return path[idx+len("assets/images/"):]
	}

	// Plan B: buscamos "images/"
	if idx := strings.Index(path, "images/"); idx != -1 {
		return path[idx+len("images/"):]
	}

	// Ultimo recurso: devolvemos solo el filename para no romper la URL (puede quedar en blanco)
	parts := strings.Split(path, "/")
	if len(parts) > 0 {
		return parts[len(parts)-1]
	}

	return path
}

// imageURL arma la ruta publica para servir la imagen y salta el warning de ngrok
func imageURL(path string) string {
	cleaned := cleanImagePath(path)
	if cleaned == "" {
		return ""
	}
	return "/api/products/images/" + cleaned + "?ngrok-skip-browser-warning=true"
}
