package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"time"
)

type LoginRequest struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}

// общая функция для HTTP-клиента
func httpClient() *http.Client {
	return &http.Client{
		Timeout: 10 * time.Second,
	}
}

// doJSONRequest выполняет JSON-запрос с необязательным токеном и парсит ответ в out.
func doJSONRequest(method, url, token string, payload any, out any) error {
	var body io.Reader
	if payload != nil {
		var buf bytes.Buffer
		if err := json.NewEncoder(&buf).Encode(payload); err != nil {
			return fmt.Errorf("encode payload: %w", err)
		}
		body = &buf
	}

	req, err := http.NewRequest(method, url, body)
	if err != nil {
		return fmt.Errorf("new request: %w", err)
	}
	if payload != nil {
		req.Header.Set("Content-Type", "application/json")
	}
	if token != "" {
		req.Header.Set("Authorization", "Bearer "+token)
	}

	resp, err := httpClient().Do(req)
	if err != nil {
		return fmt.Errorf("do request: %w", err)
	}
	defer resp.Body.Close()

	respBody, _ := io.ReadAll(resp.Body)
	if resp.StatusCode >= 400 {
		return fmt.Errorf("HTTP %d: %s", resp.StatusCode, strings.TrimSpace(string(respBody)))
	}

	if out != nil {
		if err := json.Unmarshal(respBody, out); err != nil {
			return fmt.Errorf("decode response: %w", err)
		}
	}

	return nil
}

// loginAndGetToken логинится и возвращает JWT-токен из заголовка Authorization.
func loginAndGetToken(baseURL, login, password string) (string, error) {
	url := strings.TrimRight(baseURL, "/") + "/api/user/login"
	reqBody := LoginRequest{
		Login:    login,
		Password: password,
	}

	var buf bytes.Buffer
	if err := json.NewEncoder(&buf).Encode(reqBody); err != nil {
		return "", fmt.Errorf("encode login: %w", err)
	}

	req, err := http.NewRequest(http.MethodPost, url, &buf)
	if err != nil {
		return "", fmt.Errorf("new request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := httpClient().Do(req)
	if err != nil {
		return "", fmt.Errorf("do login request: %w", err)
	}
	defer resp.Body.Close()

	bodyBytes, _ := io.ReadAll(resp.Body)
	if resp.StatusCode >= 400 {
		return "", fmt.Errorf("login failed: HTTP %d: %s", resp.StatusCode, strings.TrimSpace(string(bodyBytes)))
	}

	authHeader := resp.Header.Get("Authorization")
	if authHeader == "" {
		return "", fmt.Errorf("no Authorization header in login response")
	}
	const prefix = "Bearer "
	if !strings.HasPrefix(authHeader, prefix) {
		return "", fmt.Errorf("unexpected Authorization header: %s", authHeader)
	}
	return strings.TrimPrefix(authHeader, prefix), nil
}

// -------- Меню cli --------

func printUsage() {
	fmt.Println("GopherKeeper client")
	fmt.Println("Usage:")
	fmt.Println("  gopherkeeper-client <command> [flags]")
	fmt.Println()
	fmt.Println("Commands:")
	fmt.Println("  version")
	fmt.Println("  register")
	fmt.Println("  login")
	fmt.Println("  add-password, get-password, list-passwords")
	fmt.Println("  add-text, get-text, list-texts")
	fmt.Println("  add-binary, get-binary, list-binaries")
	fmt.Println("  add-card, get-card, list-cards")
}

func main() {
	if len(os.Args) < 2 {
		printUsage()
		os.Exit(1)
	}

	cmd := os.Args[1]
	args := os.Args[2:]

	switch cmd {
	case "version":
		cmdVersion()
	case "login":
		cmdLogin(args)
	case "register":
		cmdRegister(args)

	case "add-text":
		cmdAddText(args)
	case "get-text":
		cmdGetText(args)
	case "list-texts":
		cmdListTexts(args)

	case "add-card":
		cmdAddCard(args)
	case "get-card":
		cmdGetCard(args)
	case "list-cards":
		cmdListCards(args)

	case "add-password":
		cmdAddPassword(args)
	case "get-password":
		cmdGetPassword(args)
	case "list-passwords":
		cmdListPasswords(args)

	case "add-binary":
		cmdAddBinary(args)
	case "get-binary":
		cmdGetBinary(args)
	case "list-binaries":
		cmdListBinaries(args)

	default:
		fmt.Println("Unknown command:", cmd)
		printUsage()
		os.Exit(1)
	}
}
