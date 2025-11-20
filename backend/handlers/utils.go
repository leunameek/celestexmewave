package handlers

import (
	"strings"
)

// cleanImagePath extracts the relative path from the stored path
// e.g., "../assets/images/celeste_images/img.png" -> "celeste_images/img.png"
func cleanImagePath(path string) string {
	// Normalize separators to forward slashes
	path = strings.ReplaceAll(path, "\\", "/")

	// Remove ../assets/images/ prefix if present
	if idx := strings.Index(path, "assets/images/"); idx != -1 {
		return path[idx+len("assets/images/"):]
	}

	// Fallback: try to find just "images/"
	if idx := strings.Index(path, "images/"); idx != -1 {
		return path[idx+len("images/"):]
	}

	// Safety fallback: just return the filename to avoid breaking the URL structure
	// This might break the image link if it's in a subdir, but it prevents 404s on the route itself
	parts := strings.Split(path, "/")
	if len(parts) > 0 {
		return parts[len(parts)-1]
	}

	return path
}
