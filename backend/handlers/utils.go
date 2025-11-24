package handlers

import (
	"strings"
)

// cleanImagePath deja solo la ruta relativa del asset
// Ej: "../assets/images/celeste_images/img.png" -> "celeste_images/img.png"
func cleanImagePath(path string) string {
	// Normalizamos separadores a slash normal
	path = strings.ReplaceAll(path, "\\", "/")

	// Lista de prefijos a eliminar
	prefixes := []string{
		"../assets/images/",
		"assets/images/",
		"images/",
	}

	for _, prefix := range prefixes {
		if idx := strings.Index(path, prefix); idx != -1 {
			return path[idx+len(prefix):]
		}
	}

	// Si no coincide con ningun prefijo, devolvemos el path tal cual
	// Esto permite que rutas como "celeste_images/foto.jpg" funcionen correctamente
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
