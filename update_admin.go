package main

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/go-sql-driver/mysql"
	"golang.org/x/crypto/bcrypt"
)

func main() {
	// Conectar ao banco
	db, err := sql.Open("mysql", "root:senha123@tcp(localhost:3306)/trabiju_telemetria")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// Gerar hash para admin123
	hash, err := bcrypt.GenerateFromPassword([]byte("admin123"), bcrypt.DefaultCost)
	if err != nil {
		log.Fatal(err)
	}

	// Atualizar no banco
	_, err = db.Exec("UPDATE usuarios SET senha = ? WHERE email = ?", string(hash), "admin@gestaotelemetria.com")
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("✅ Senha do admin atualizada com sucesso!")
	fmt.Printf("Hash gerado: %s\n", string(hash))
}
