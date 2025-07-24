package main

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"strings"
	"time"
)

func main() {
	// Token recebido
	token := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyX2lkIjoxLCJlbWFpbCI6ImFkbWluQGdlc3Rhb3RlbGVtZXRyaWEuY29tIiwicm9sZV9pZCI6MywiZW1wcmVzYV9pZCI6MSwiaXNzIjoidHJhYmlqdS10ZWxlbWV0cmlhIiwic3ViIjoiYWRtaW5AZ2VzdGFvdGVsZW1ldHJpYS5jb20iLCJleHAiOjE3NTM0NDY3MzMsIm5iZiI6MTc1MzM2MDMzMywiaWF0IjoxNzUzMzYwMzMzfQ.8KvhBIMrFwUuonGUD-uIJ_HQR700oufrUTv7zRvB2Hg"

	// Dividir o token em partes
	parts := strings.Split(token, ".")
	if len(parts) != 3 {
		fmt.Println("Token inválido")
		return
	}

	// Decodificar header
	header, _ := base64.RawURLEncoding.DecodeString(parts[0])
	fmt.Println("🔐 HEADER:")
	fmt.Println(string(header))
	fmt.Println()

	// Decodificar payload
	payload, _ := base64.RawURLEncoding.DecodeString(parts[1])
	fmt.Println("📋 PAYLOAD:")

	// Parse do payload como JSON
	var claims map[string]interface{}
	json.Unmarshal(payload, &claims)

	// Exibir claims formatados
	for key, value := range claims {
		switch key {
		case "exp", "nbf", "iat":
			// Converter timestamp para data legível
			if timestamp, ok := value.(float64); ok {
				date := time.Unix(int64(timestamp), 0)
				fmt.Printf("  %s: %v (%s)\n", key, value, date.Format("02/01/2006 15:04:05"))
			}
		default:
			fmt.Printf("  %s: %v\n", key, value)
		}
	}

	fmt.Println()
	fmt.Println("🔑 ASSINATURA:")
	fmt.Println(parts[2])
}
