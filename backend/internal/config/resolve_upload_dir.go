package config

import "os"

// resolveUploadDir intenta encontrar el directorio de assets
func resolveUploadDir(defaultPath string) string {
	// Si viene por env y existe, usamos ese
	if _, err := os.Stat(defaultPath); err == nil {
		return defaultPath
	}

	// Si no, probamos rutas comunes
	candidates := []string{
		"../assets/images",
		"assets/images",
		"../../assets/images",
		"/app/assets/images", // Docker fallback
	}

	for _, path := range candidates {
		if _, err := os.Stat(path); err == nil {
			return path
		}
	}

	// Si nada funciona, devolvemos el default y que explote con error claro despues
	return defaultPath
}
