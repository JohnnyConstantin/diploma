package main

import (
	"diploma/internal/models"
	"diploma/pkg/config"
	"diploma/pkg/version"
	"flag"
	"fmt"
	"net/http"
	"os"
	"strings"
)

func cmdVersion() {
	fmt.Printf("GopherKeeper client\nVersion: %s\nBuild date: %s\n", version.Version, version.BuildDate)
}

func cmdRegister(args []string) {
	fs := flag.NewFlagSet("register", flag.ExitOnError)
	baseURL := fs.String("base-url", config.DefaultBaseURL, "Base URL of server")
	login := fs.String("login", "", "User login")
	password := fs.String("password", "", "User password")
	_ = fs.Parse(args)

	if *login == "" || *password == "" {
		fmt.Println("login and password are required")
		fs.Usage()
		os.Exit(1)
	}

	url := strings.TrimRight(*baseURL, "/") + "/api/user/register"
	reqBody := LoginRequest{Login: *login, Password: *password}
	var resp map[string]any

	if err := doJSONRequest(http.MethodPost, url, "", reqBody, &resp); err != nil {
		fmt.Println("register error:", err)
		os.Exit(1)
	}
	fmt.Println("Registered successfully:", resp)
}

func cmdLogin(args []string) {
	fs := flag.NewFlagSet("login", flag.ExitOnError)
	baseURL := fs.String("base-url", config.DefaultBaseURL, "Base URL of server")
	login := fs.String("login", "", "User login")
	password := fs.String("password", "", "User password")
	_ = fs.Parse(args)

	if *login == "" || *password == "" {
		fmt.Println("login and password are required")
		fs.Usage()
		os.Exit(1)
	}

	token, err := loginAndGetToken(*baseURL, *login, *password)
	if err != nil {
		fmt.Println("login error:", err)
		os.Exit(1)
	}
	fmt.Println("Login OK. Token (JWT):")
	fmt.Println(token)
}

func cmdAddPassword(args []string) {
	fs := flag.NewFlagSet("add-password", flag.ExitOnError)
	baseURL := fs.String("base-url", config.DefaultBaseURL, "Base URL of server")
	userLogin := fs.String("login", "", "User login")
	userPassword := fs.String("password", "", "User password")
	id := fs.String("id", "", "Password record ID")
	serviceLogin := fs.String("plogin", "", "Service login for this password")
	servicePassword := fs.String("ppass", "", "Service password")
	meta := fs.String("meta", "", "Optional meta information")
	_ = fs.Parse(args)

	if *userLogin == "" || *userPassword == "" || *id == "" || *serviceLogin == "" || *servicePassword == "" {
		fmt.Println("login, password, id, plogin and ppass are required")
		fs.Usage()
		os.Exit(1)
	}

	token, err := loginAndGetToken(*baseURL, *userLogin, *userPassword)
	if err != nil {
		fmt.Println("login error:", err)
		os.Exit(1)
	}

	var metaPtr *string
	if *meta != "" {
		metaPtr = meta
	}

	reqBody := models.PasswordRequest{
		ID:       *id,
		Login:    *serviceLogin,
		Password: *servicePassword,
		Meta:     metaPtr,
	}
	url := strings.TrimRight(*baseURL, "/") + "/api/v1/passwords"
	var resp models.PasswordResponse

	if err := doJSONRequest(http.MethodPost, url, token, reqBody, &resp); err != nil {
		fmt.Println("add-password error:", err)
		os.Exit(1)
	}
	fmt.Printf("Password record saved: %+v\n", resp)
}

func cmdGetPassword(args []string) {
	fs := flag.NewFlagSet("get-password", flag.ExitOnError)
	baseURL := fs.String("base-url", config.DefaultBaseURL, "Base URL of server")
	userLogin := fs.String("login", "", "User login")
	userPassword := fs.String("password", "", "User password")
	id := fs.String("id", "", "Password record ID")
	_ = fs.Parse(args)

	if *userLogin == "" || *userPassword == "" || *id == "" {
		fmt.Println("login, password and id are required")
		fs.Usage()
		os.Exit(1)
	}

	token, err := loginAndGetToken(*baseURL, *userLogin, *userPassword)
	if err != nil {
		fmt.Println("login error:", err)
		os.Exit(1)
	}

	url := strings.TrimRight(*baseURL, "/") + "/api/v1/passwords/" + *id
	var resp models.PasswordResponse
	if err := doJSONRequest(http.MethodGet, url, token, nil, &resp); err != nil {
		fmt.Println("get-password error:", err)
		os.Exit(1)
	}
	fmt.Printf("Password record: %+v\n", resp)
}

func cmdListPasswords(args []string) {
	fs := flag.NewFlagSet("list-passwords", flag.ExitOnError)
	baseURL := fs.String("base-url", config.DefaultBaseURL, "Base URL of server")
	userLogin := fs.String("login", "", "User login")
	userPassword := fs.String("password", "", "User password")
	_ = fs.Parse(args)

	if *userLogin == "" || *userPassword == "" {
		fmt.Println("login and password are required")
		fs.Usage()
		os.Exit(1)
	}

	token, err := loginAndGetToken(*baseURL, *userLogin, *userPassword)
	if err != nil {
		fmt.Println("login error:", err)
		os.Exit(1)
	}

	url := strings.TrimRight(*baseURL, "/") + "/api/v1/passwords"
	var resp []models.PasswordResponse
	if err := doJSONRequest(http.MethodGet, url, token, nil, &resp); err != nil {
		fmt.Println("list-passwords error:", err)
		os.Exit(1)
	}
	fmt.Printf("Password records: %+v\n", resp)
}

func cmdAddText(args []string) {
	fs := flag.NewFlagSet("add-text", flag.ExitOnError)
	baseURL := fs.String("base-url", config.DefaultBaseURL, "Base URL of server")
	userLogin := fs.String("login", "", "User login")
	userPassword := fs.String("password", "", "User password")
	id := fs.String("id", "", "Text ID (title)")
	text := fs.String("text", "", "Text content")
	meta := fs.String("meta", "", "Optional meta")
	_ = fs.Parse(args)

	if *userLogin == "" || *userPassword == "" || *id == "" || *text == "" {
		fmt.Println("login, password, id and text are required")
		fs.Usage()
		os.Exit(1)
	}

	token, err := loginAndGetToken(*baseURL, *userLogin, *userPassword)
	if err != nil {
		fmt.Println("login error:", err)
		os.Exit(1)
	}

	var metaPtr *string
	if *meta != "" {
		metaPtr = meta
	}

	reqBody := models.TextRequest{
		ID:   *id,
		Text: *text,
		Meta: metaPtr,
	}
	url := strings.TrimRight(*baseURL, "/") + "/api/v1/texts"
	var resp models.TextResponse
	if err := doJSONRequest(http.MethodPost, url, token, reqBody, &resp); err != nil {
		fmt.Println("add-text error:", err)
		os.Exit(1)
	}
	fmt.Printf("Text record saved: %+v\n", resp)
}

func cmdGetText(args []string) {
	fs := flag.NewFlagSet("get-text", flag.ExitOnError)
	baseURL := fs.String("base-url", config.DefaultBaseURL, "Base URL of server")
	userLogin := fs.String("login", "", "User login")
	userPassword := fs.String("password", "", "User password")
	id := fs.String("id", "", "Text ID")
	_ = fs.Parse(args)

	if *userLogin == "" || *userPassword == "" || *id == "" {
		fmt.Println("login, password and id are required")
		fs.Usage()
		os.Exit(1)
	}

	token, err := loginAndGetToken(*baseURL, *userLogin, *userPassword)
	if err != nil {
		fmt.Println("login error:", err)
		os.Exit(1)
	}

	url := strings.TrimRight(*baseURL, "/") + "/api/v1/texts/" + *id
	var resp models.TextResponse
	if err := doJSONRequest(http.MethodGet, url, token, nil, &resp); err != nil {
		fmt.Println("get-text error:", err)
		os.Exit(1)
	}
	fmt.Printf("Text record: %+v\n", resp)
}

func cmdListTexts(args []string) {
	fs := flag.NewFlagSet("list-texts", flag.ExitOnError)
	baseURL := fs.String("base-url", config.DefaultBaseURL, "Base URL of server")
	userLogin := fs.String("login", "", "User login")
	userPassword := fs.String("password", "", "User password")
	_ = fs.Parse(args)

	if *userLogin == "" || *userPassword == "" {
		fmt.Println("login and password are required")
		fs.Usage()
		os.Exit(1)
	}

	token, err := loginAndGetToken(*baseURL, *userLogin, *userPassword)
	if err != nil {
		fmt.Println("login error:", err)
		os.Exit(1)
	}

	url := strings.TrimRight(*baseURL, "/") + "/api/v1/texts"
	var resp []models.TextResponse
	if err := doJSONRequest(http.MethodGet, url, token, nil, &resp); err != nil {
		fmt.Println("list-texts error:", err)
		os.Exit(1)
	}
	fmt.Printf("Text records: %+v\n", resp)
}

func cmdAddCard(args []string) {
	fs := flag.NewFlagSet("add-card", flag.ExitOnError)
	baseURL := fs.String("base-url", config.DefaultBaseURL, "Base URL of server")
	userLogin := fs.String("login", "", "User login")
	userPassword := fs.String("password", "", "User password")
	number := fs.String("number", "", "Card number")
	holder := fs.String("holder", "", "Card holder name")
	expire := fs.String("expire", "", "Expire date (e.g. 12/30)")
	cvv := fs.String("cvv", "", "CVV")
	meta := fs.String("meta", "", "Optional meta")
	_ = fs.Parse(args)

	if *userLogin == "" || *userPassword == "" || *number == "" || *holder == "" || *expire == "" || *cvv == "" {
		fmt.Println("login, password, number, holder, expire and cvv are required")
		fs.Usage()
		os.Exit(1)
	}

	token, err := loginAndGetToken(*baseURL, *userLogin, *userPassword)
	if err != nil {
		fmt.Println("login error:", err)
		os.Exit(1)
	}

	var metaPtr *string
	if *meta != "" {
		metaPtr = meta
	}

	reqBody := models.CardRequest{
		Number: *number,
		Holder: *holder,
		Expire: *expire,
		CVV:    *cvv,
		Meta:   metaPtr,
	}
	url := strings.TrimRight(*baseURL, "/") + "/api/v1/cards"
	var resp models.CardResponse
	if err := doJSONRequest(http.MethodPost, url, token, reqBody, &resp); err != nil {
		fmt.Println("add-card error:", err)
		os.Exit(1)
	}
	fmt.Printf("Card saved: %+v\n", resp)
}

func cmdGetCard(args []string) {
	fs := flag.NewFlagSet("get-card", flag.ExitOnError)
	baseURL := fs.String("base-url", config.DefaultBaseURL, "Base URL of server")
	userLogin := fs.String("login", "", "User login")
	userPassword := fs.String("password", "", "User password")
	cardPAN := fs.String("card-pan", "", "Masked card PAN (e.g. 4600********5363)")
	_ = fs.Parse(args)

	if *userLogin == "" || *userPassword == "" || *cardPAN == "" {
		fmt.Println("login, password and card-pan are required")
		fs.Usage()
		os.Exit(1)
	}

	token, err := loginAndGetToken(*baseURL, *userLogin, *userPassword)
	if err != nil {
		fmt.Println("login error:", err)
		os.Exit(1)
	}

	url := strings.TrimRight(*baseURL, "/") + "/api/v1/cards/" + *cardPAN
	var resp models.CardResponse
	if err := doJSONRequest(http.MethodGet, url, token, nil, &resp); err != nil {
		fmt.Println("get-card error:", err)
		os.Exit(1)
	}
	fmt.Printf("Card: %+v\n", resp)
}

func cmdListCards(args []string) {
	fs := flag.NewFlagSet("list-cards", flag.ExitOnError)
	baseURL := fs.String("base-url", config.DefaultBaseURL, "Base URL of server")
	userLogin := fs.String("login", "", "User login")
	userPassword := fs.String("password", "", "User password")
	_ = fs.Parse(args)

	if *userLogin == "" || *userPassword == "" {
		fmt.Println("login and password are required")
		fs.Usage()
		os.Exit(1)
	}

	token, err := loginAndGetToken(*baseURL, *userLogin, *userPassword)
	if err != nil {
		fmt.Println("login error:", err)
		os.Exit(1)
	}

	url := strings.TrimRight(*baseURL, "/") + "/api/v1/cards"
	var resp []models.CardResponse
	if err := doJSONRequest(http.MethodGet, url, token, nil, &resp); err != nil {
		fmt.Println("list-cards error:", err)
		os.Exit(1)
	}
	fmt.Printf("Cards: %+v\n", resp)
}

func cmdAddBinary(args []string) {
	fs := flag.NewFlagSet("add-binary", flag.ExitOnError)
	baseURL := fs.String("base-url", config.DefaultBaseURL, "Base URL of server")
	userLogin := fs.String("login", "", "User login")
	userPassword := fs.String("password", "", "User password")
	id := fs.String("id", "", "Binary record ID")
	data := fs.String("data", "", "Binary data as string (you can pre-encode base64)")
	meta := fs.String("meta", "", "Optional meta")
	_ = fs.Parse(args)

	if *userLogin == "" || *userPassword == "" || *id == "" || *data == "" {
		fmt.Println("login, password, id and data are required")
		fs.Usage()
		os.Exit(1)
	}

	token, err := loginAndGetToken(*baseURL, *userLogin, *userPassword)
	if err != nil {
		fmt.Println("login error:", err)
		os.Exit(1)
	}

	var metaPtr *string
	if *meta != "" {
		metaPtr = meta
	}

	reqBody := models.BinaryRequest{
		ID:   *id,
		Data: *data,
		Meta: metaPtr,
	}
	url := strings.TrimRight(*baseURL, "/") + "/api/v1/binaries"
	var resp models.BinaryResponse
	if err := doJSONRequest(http.MethodPost, url, token, reqBody, &resp); err != nil {
		fmt.Println("add-binary error:", err)
		os.Exit(1)
	}
	fmt.Printf("Binary record saved: %+v\n", resp)
}

func cmdGetBinary(args []string) {
	fs := flag.NewFlagSet("get-binary", flag.ExitOnError)
	baseURL := fs.String("base-url", config.DefaultBaseURL, "Base URL of server")
	userLogin := fs.String("login", "", "User login")
	userPassword := fs.String("password", "", "User password")
	id := fs.String("id", "", "Binary record ID")
	_ = fs.Parse(args)

	if *userLogin == "" || *userPassword == "" || *id == "" {
		fmt.Println("login, password and id are required")
		fs.Usage()
		os.Exit(1)
	}

	token, err := loginAndGetToken(*baseURL, *userLogin, *userPassword)
	if err != nil {
		fmt.Println("login error:", err)
		os.Exit(1)
	}

	url := strings.TrimRight(*baseURL, "/") + "/api/v1/binaries/" + *id
	var resp models.BinaryResponse
	if err := doJSONRequest(http.MethodGet, url, token, nil, &resp); err != nil {
		fmt.Println("get-binary error:", err)
		os.Exit(1)
	}
	fmt.Printf("Binary record: %+v\n", resp)
}

func cmdListBinaries(args []string) {
	fs := flag.NewFlagSet("list-binaries", flag.ExitOnError)
	baseURL := fs.String("base-url", config.DefaultBaseURL, "Base URL of server")
	userLogin := fs.String("login", "", "User login")
	userPassword := fs.String("password", "", "User password")
	_ = fs.Parse(args)

	if *userLogin == "" || *userPassword == "" {
		fmt.Println("login and password are required")
		fs.Usage()
		os.Exit(1)
	}

	token, err := loginAndGetToken(*baseURL, *userLogin, *userPassword)
	if err != nil {
		fmt.Println("login error:", err)
		os.Exit(1)
	}

	url := strings.TrimRight(*baseURL, "/") + "/api/v1/binaries"
	var resp []models.BinaryResponse
	if err := doJSONRequest(http.MethodGet, url, token, nil, &resp); err != nil {
		fmt.Println("list-binaries error:", err)
		os.Exit(1)
	}
	fmt.Printf("Binary records: %+v\n", resp)
}
