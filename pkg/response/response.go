package response

import (
	"encoding/json"
	"net/http"
)

// Response is the standard structure for all API responses
type Response struct {
	Success bool        `json:"success"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
	Errors  interface{} `json:"errors,omitempty"`
	Code    int         `json:"code"`
	Meta    interface{} `json:"meta,omitempty"`
}

// sendJSON is the base function to send a JSON response
func sendJSON(w http.ResponseWriter, status int, res Response) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	res.Code = status
	if err := json.NewEncoder(w).Encode(res); err != nil {
		// Fallback to http.Error if encoding fails
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}
}

// Success returns a successful response with data
func Success(w http.ResponseWriter, status int, message string, data interface{}) {
	sendJSON(w, status, Response{
		Success: true,
		Message: message,
		Data:    data,
	})
}

// Error returns a simple error response
func Error(w http.ResponseWriter, status int, message string) {
	sendJSON(w, status, Response{
		Success: false,
		Message: message,
	})
}

// ErrorWithDetails returns an error response with additional details (e.g., validation errors)
func ErrorWithDetails(w http.ResponseWriter, status int, message string, errors interface{}) {
	sendJSON(w, status, Response{
		Success: false,
		Message: message,
		Errors:  errors,
	})
}

// PaginationMeta holds pagination metadata
type PaginationMeta struct {
	CurrentPage int   `json:"current_page"`
	TotalPages  int   `json:"total_pages"`
	PageSize    int   `json:"page_size"`
	TotalItems  int64 `json:"total_items"`
}

// Pagination returns a successful response with pagination metadata
func Pagination(w http.ResponseWriter, status int, message string, data interface{}, meta PaginationMeta) {
	sendJSON(w, status, Response{
		Success: true,
		Message: message,
		Data:    data,
		Meta:    meta,
	})
}

// JSON is maintained for backward compatibility but wrapped in the new standard format
func JSON(w http.ResponseWriter, status int, data interface{}) {
	if status >= 400 {
		// If it's an error, pass data to details
		ErrorWithDetails(w, status, "Request failed", data)
		return
	}
	Success(w, status, "Success", data)
}
