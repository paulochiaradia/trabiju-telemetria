package validation

import (
	"regexp"
	"strings"
	"unicode"

	"github.com/go-playground/validator/v10"
)

var validate *validator.Validate

func init() {
	validate = validator.New()

	// Registrar validações customizadas
	validate.RegisterValidation("cpf", validateCPF)
	validate.RegisterValidation("strongpassword", validateStrongPassword)
	validate.RegisterValidation("phone", validatePhone)
}

// GetValidator retorna a instância do validador
func GetValidator() *validator.Validate {
	return validate
}

// ValidateStruct valida uma struct
func ValidateStruct(s interface{}) error {
	return validate.Struct(s)
}

// validateCPF valida CPF brasileiro
func validateCPF(fl validator.FieldLevel) bool {
	cpf := strings.ReplaceAll(fl.Field().String(), ".", "")
	cpf = strings.ReplaceAll(cpf, "-", "")

	// CPF deve ter 11 dígitos
	if len(cpf) != 11 {
		return false
	}

	// Verifica se todos os dígitos são iguais
	if strings.Count(cpf, string(cpf[0])) == 11 {
		return false
	}

	// Validação do primeiro dígito verificador
	sum := 0
	for i := 0; i < 9; i++ {
		digit := int(cpf[i] - '0')
		sum += digit * (10 - i)
	}
	remainder := sum % 11
	firstDigit := 0
	if remainder >= 2 {
		firstDigit = 11 - remainder
	}

	if int(cpf[9]-'0') != firstDigit {
		return false
	}

	// Validação do segundo dígito verificador
	sum = 0
	for i := 0; i < 10; i++ {
		digit := int(cpf[i] - '0')
		sum += digit * (11 - i)
	}
	remainder = sum % 11
	secondDigit := 0
	if remainder >= 2 {
		secondDigit = 11 - remainder
	}

	return int(cpf[10]-'0') == secondDigit
}

// validateStrongPassword valida senha forte
func validateStrongPassword(fl validator.FieldLevel) bool {
	password := fl.Field().String()

	// Mínimo 8 caracteres
	if len(password) < 8 {
		return false
	}

	var (
		hasUpper   = false
		hasLower   = false
		hasNumber  = false
		hasSpecial = false
	)

	for _, char := range password {
		switch {
		case unicode.IsUpper(char):
			hasUpper = true
		case unicode.IsLower(char):
			hasLower = true
		case unicode.IsNumber(char):
			hasNumber = true
		case unicode.IsPunct(char) || unicode.IsSymbol(char):
			hasSpecial = true
		}
	}

	// Deve ter pelo menos 3 dos 4 tipos de caracteres
	count := 0
	if hasUpper {
		count++
	}
	if hasLower {
		count++
	}
	if hasNumber {
		count++
	}
	if hasSpecial {
		count++
	}

	return count >= 3
}

// validatePhone valida telefone brasileiro
func validatePhone(fl validator.FieldLevel) bool {
	phone := fl.Field().String()

	// Remove caracteres não numéricos
	re := regexp.MustCompile(`[^\d]`)
	numbers := re.ReplaceAllString(phone, "")

	// Telefone brasileiro: 10 ou 11 dígitos (com ou sem DDD)
	// Formato: (11) 99999-9999 ou (11) 9999-9999
	if len(numbers) == 10 || len(numbers) == 11 {
		return true
	}

	// Com código do país: +55 11 99999-9999
	if len(numbers) == 13 && strings.HasPrefix(numbers, "55") {
		return true
	}

	return false
}

// GetValidationErrors converte erros de validação em mensagens amigáveis
func GetValidationErrors(err error) map[string]string {
	errors := make(map[string]string)

	if validationErrors, ok := err.(validator.ValidationErrors); ok {
		for _, validationError := range validationErrors {
			field := strings.ToLower(validationError.Field())
			tag := validationError.Tag()

			switch tag {
			case "required":
				errors[field] = "Este campo é obrigatório"
			case "email":
				errors[field] = "Formato de email inválido"
			case "min":
				if validationError.Kind().String() == "string" {
					errors[field] = "Deve ter pelo menos " + validationError.Param() + " caracteres"
				} else {
					errors[field] = "Valor mínimo é " + validationError.Param()
				}
			case "max":
				if validationError.Kind().String() == "string" {
					errors[field] = "Deve ter no máximo " + validationError.Param() + " caracteres"
				} else {
					errors[field] = "Valor máximo é " + validationError.Param()
				}
			case "cpf":
				errors[field] = "CPF inválido"
			case "strongpassword":
				errors[field] = "Senha deve ter pelo menos 8 caracteres com letras maiúsculas, minúsculas, números e símbolos"
			case "phone":
				errors[field] = "Formato de telefone inválido"
			case "oneof":
				errors[field] = "Valor deve ser um dos seguintes: " + validationError.Param()
			default:
				errors[field] = "Valor inválido"
			}
		}
	}

	return errors
}
