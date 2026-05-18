package password
package password

import (
	"fmt"
	"regexp"
	"unicode"

	"golang.org/x/crypto/bcrypt"
)

// Manager gestiona operaciones de contraseña
type Manager struct {
	bcryptCost int
}

// NewManager crea una nueva instancia del manager de contraseña
func NewManager() *Manager {
	return &Manager{
		bcryptCost: bcrypt.DefaultCost,
	}
}

// Hash genera un hash seguro de la contraseña
func (m *Manager) Hash(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), m.bcryptCost)
	if err != nil {
		return "", fmt.Errorf("error hashing password: %w", err)
	}
	return string(hash), nil
}

// Compare verifica una contraseña contra su hash
func (m *Manager) Compare(hashedPassword, password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
	return err == nil
}

// ValidateStrength valida la fortaleza de una contraseña
func (m *Manager) ValidateStrength(password string) (bool, []string) {
	var errors []string

	if len(password) < 8 {
		errors = append(errors, "La contraseña debe tener al menos 8 caracteres")
	}

	if len(password) > 128 {
		errors = append(errors, "La contraseña no puede exceder 128 caracteres")
	}

	if !hasUpperCase(password) {
		errors = append(errors, "La contraseña debe contener al menos una mayúscula")
	}

	if !hasLowerCase(password) {
		errors = append(errors, "La contraseña debe contener al menos una minúscula")
	}

	if !hasNumber(password) {
		errors = append(errors, "La contraseña debe contener al menos un número")
	}

	if !hasSpecialChar(password) {
		errors = append(errors, "La contraseña debe contener al menos un carácter especial")
	}

	return len(errors) == 0, errors
}

// helper functions
func hasUpperCase(s string) bool {
	for _, r := range s {
		if unicode.IsUpper(r) {
			return true
		}
	}
	return false
}

func hasLowerCase(s string) bool {
	for _, r := range s {
		if unicode.IsLower(r) {
			return true
		}
	}
	return false
}

func hasNumber(s string) bool {
	for _, r := range s {
		if unicode.IsDigit(r) {
			return true
		}
	}
	return false
}

func hasSpecialChar(s string) bool {
	// Caracteres especiales permitidos
	specialChars := regexp.MustCompile(`[!@#$%^&*()_+\-=\[\]{};':"\\|,.<>\/?]`)
	return specialChars.MatchString(s)
}
