package response
package response

import (
	"encoding/json"
	"net/http"

	"github.com/apariciocch/psicosstcloud/internal/dto"
)

// Success respuesta exitosa
func Success(w http.ResponseWriter, statusCode int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(data)
}

// SuccessPaginated respuesta exitosa paginada
func SuccessPaginated(w http.ResponseWriter, statusCode int, data interface{}, page, pageSize, total int) {
	totalPages := (total + pageSize - 1) / pageSize
	response := dto.PaginatedResponse{
		Data:      data,
		Page:      page,
		PageSize:  pageSize,
		Total:     total,
		TotalPages: totalPages,
		HasNext:   page < totalPages,
		HasPrev:   page > 1,
	}

	Success(w, statusCode, response)
}

// Error respuesta de error
func Error(w http.ResponseWriter, statusCode int, code, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)

	errResp := dto.ErrorResponse{
		Code:    code,
		Message: message,
	}

	json.NewEncoder(w).Encode(errResp)
}

// ErrorWithField respuesta de error con campo
func ErrorWithField(w http.ResponseWriter, statusCode int, code, message, field string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)

	errResp := dto.ErrorResponse{
		Code:    code,
		Message: message,
		Field:   field,
	}

	json.NewEncoder(w).Encode(errResp)
}

// ErrorWithDetails respuesta de error con detalles
func ErrorWithDetails(w http.ResponseWriter, statusCode int, code, message string, details interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)

	errResp := dto.ErrorResponse{
		Code:    code,
		Message: message,
		Details: details,
	}

	json.NewEncoder(w).Encode(errResp)
}
