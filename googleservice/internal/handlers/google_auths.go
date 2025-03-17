package handlers

import (
	"encoding/json"
	"fmt"
	"google-service/internal/config"
	"google-service/internal/middleware"
	"google-service/internal/services"
	"google-service/utils"
	"os"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v4"
	"github.com/valyala/fasthttp"
)

type GoogleHandler struct {
	oauthService *services.GoogleOAuthService
}

func NewGoogleHandler(config *config.GoogleOAuthConfig) *GoogleHandler {
	return &GoogleHandler{
		oauthService: services.NewGoogleOAuthService(config),
	}
}

func (h *GoogleHandler) HandleGoogleLogin(c *fiber.Ctx) error {
	state := c.Query("state")
	if state == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "State parameter required",
		})
	}

	url := h.oauthService.GetAuthURL(state)
	return c.Redirect(url)
}

func (h *GoogleHandler) HandleGoogleCallback(c *fiber.Ctx) error {
	state := c.Cookies("oauth_state")
	fmt.Println("COOKIE:", state, "\nQUERY:", c.Query("state"))
	if state != c.Query("state") {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Invalid state",
		})
	}

	code := c.Query("code")
	if code == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Authorization code is missing",
		})
	}

	token, err := h.oauthService.Exchange(c.Context(), code)
	if err != nil {

		fmt.Printf("Token exchange error: %v\n", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":   "Failed to exchange token",
			"details": err.Error(),
		})
	}

	userInfo, err := h.oauthService.GetUserInfo(token)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to get user info",
		})
	}

	err = AddUser(userInfo.Name, userInfo.Email, userInfo.Picture, token.AccessToken)
	if err != nil {
		fmt.Println("GOOGLE SERVICE: Save User to Database failed:", err)
	}

	id, err := GetUserID(userInfo.Email)
	if err != nil {
		fmt.Println("GOOGLE SERVICE: Get user id failed:", err)
	}

	claims := middleware.Claims{
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * 24)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Issuer:    "google-service",
		},
		UserID:   uint(id),
		Username: userInfo.Email,
		Email:    userInfo.Email,
		Role:     "Parent",
	}

	type User struct {
		Userid   uint   `json:"user_id"`
		Username string `json:"username"`
		Email    string `json:"email"`
		Role     string `json:"role"`
		Status   string `json:"status"`
	}

	jwt_token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	jwt_tokenString, err := jwt_token.SignedString([]byte(os.Getenv("JWT_SECRET")))
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "User is Signed in but could not generate jwt token",
		})
	}
	h.oauthService.StoreUserToken(userInfo.Email, token)

	return c.JSON(fiber.Map{
		"token": jwt_tokenString,
		"user": User{
			Userid:   uint(id),
			Username: userInfo.Email,
			Email:    userInfo.Email,
			Role:     "Parent",
			Status:   "Active",
		},
	})
}

func AddUser(name, email, picture, access_token string) error {
	req := &fasthttp.Request{}
	resp := &fasthttp.Response{}

	body := fmt.Sprintf(`{"username":"%s","password":"%s", "email":"%s", "role":"%s", "picture":"%s"}`, email, "", email, "Parent", picture)
	query := fmt.Sprintf("?google_token=%s", access_token)
	fmt.Println(config.Google_config.USER_SERVICE_URL + "/user/add" + query)
	utils.BuildRequest(req, "POST", []byte(body), config.Google_config.API_KEY, config.Google_config.USER_SERVICE_URL+"/user/add"+query)

	if err := fasthttp.Do(req, resp); err != nil {
		return fmt.Errorf("user service unavailable: %v", err)
	}

	if resp.StatusCode() >= 200 && resp.StatusCode() < 300 {
		return nil
	}
	return fmt.Errorf("add user failed: %s", string(resp.Body()))
}

func GetUserID(email string) (int, error) {
	req := &fasthttp.Request{}
	resp := &fasthttp.Response{}

	query := fmt.Sprintf("?email=%s", email)
	utils.BuildRequest(req, "GET", nil, config.Google_config.API_KEY, config.Google_config.USER_SERVICE_URL+"/user/get/id"+query)

	if err := fasthttp.Do(req, resp); err != nil {
		return 0, fmt.Errorf("user service unavailable: %v", err)
	}

	if resp.StatusCode() >= 200 && resp.StatusCode() < 300 {
		var originalData map[string]interface{}
		if err := json.Unmarshal(resp.Body(), &originalData); err != nil {
			fmt.Printf("Error parsing JSON response: %v\n", err)
			return 0, fmt.Errorf("error formatting response data: %s", err)
		}

		return originalData["user_id"].(int), nil
	}
	return 0, fmt.Errorf("get user id failed: %s", string(resp.Body()))
}
