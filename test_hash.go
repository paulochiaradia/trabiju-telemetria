package main

import (
	"fmt"

	"golang.org/x/crypto/bcrypt"
)

func main() {
	// Hash que está no banco
	hashFromDB := "$2a$10$N9qo8uLOickgx2ZMRZoMye.mQKyVZV8Cgp1.JEE8SL8BI7.FjAy6."

	// Senha que estamos testando
	password := "admin123"

	// Testar se o hash confere
	err := bcrypt.CompareHashAndPassword([]byte(hashFromDB), []byte(password))
	if err != nil {
		fmt.Printf("❌ Hash NÃO confere: %v\n", err)

		// Vamos gerar um novo hash para ver como deve ser
		newHash, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
		fmt.Printf("Hash que deveria ser: %s\n", string(newHash))
	} else {
		fmt.Println("✅ Hash confere perfeitamente!")
	}
}
