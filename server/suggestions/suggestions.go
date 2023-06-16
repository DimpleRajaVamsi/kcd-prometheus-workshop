package suggestions

import "net/http"

// Suggestions interface defines the API exposed by Suggestions APP
type Suggestions interface {
	Suggest(w http.ResponseWriter, req *http.Request)
}
