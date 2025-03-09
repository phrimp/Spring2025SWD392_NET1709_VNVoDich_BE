package handlers

import (
	"authservice/internal/repository"
	"authservice/internal/services"
	"authservice/utils"
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v4"
)

type Claims struct {
	jwt.RegisteredClaims
	UserID   uint   `json:"userId"`
	Username string `json:"username"`
	Email    string `json:"email"`
	Role     string `json:"role"`
}

var (
	tokens  map[string]*jwt.Token  = make(map[string]*jwt.Token)
	blocked map[string]bool        = make(map[string]bool)
	claims  map[*jwt.Token]*Claims = make(map[*jwt.Token]*Claims)
)

func BlockUser(username string) {
	blocked_status := blocked[username]
	if blocked_status {
		fmt.Println("User is already blocked")
		return
	}
	blocked[username] = true
}

func UnBlockUser(username string) {
	blocked[username] = false
}

func HandleBlockToken() fiber.Handler {
	return func(c *fiber.Ctx) error {
		username := c.Query("username")
		block_time_str := c.Query("time")
		block_time_int, err := strconv.Atoi(block_time_str)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "invalid block time",
			})
		}
		block_time := time.Now().Unix() + int64(block_time_int)
		BlockUser(username)
		fmt.Println("User JWT Token blocked:", username, "and will be unblocked at", time.Unix(block_time, 0))

		go func(block_time_int int, username string) {
			time.Sleep(time.Duration(block_time_int) * time.Second)
			UnBlockUser(username)
		}(block_time_int, username)

		return c.Status(fiber.StatusOK).JSON(fiber.Map{
			"message": "user is blocked",
		})
	}
}

func HandleUnblockToken() fiber.Handler {
	return func(c *fiber.Ctx) error {
		username := c.Query("username")
		if username == "" {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "username is required",
			})
		}
		UnBlockUser(username)
		fmt.Println("Unblock user jwt attempt successfully")
		return c.Status(fiber.StatusOK).JSON(fiber.Map{
			"message": "user jwt is unblocked",
		})
	}
}

func generateNewToken(user *repository.User) *jwt.Token {
	claim := Claims{
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * 24)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Issuer:    "auth-service",
		},
		UserID:   user.ID,
		Username: user.Username,
		Email:    user.Email,
		Role:     user.Role,
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claim)
	claims[token] = &claim
	return token
}

func HandleLogin() fiber.Handler {
	return func(c *fiber.Ctx) error {
		var req LoginRequest
		if err := c.BodyParser(&req); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Invalid request",
			})
		}

		// Forward to user service
		user, err := services.GetUserFromUserService(utils.SERVICES_ROUTES.UserService, req.Username, req.Password)
		if err != nil {
			fmt.Println("Login error:", err)
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "Invalid credentials",
			})
		}

		is_blocked := blocked[user.Username]

		if is_blocked {
			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
				"error": "user jwt is blocked",
			})
		}

		if user.Status != "Active" {
			fmt.Println("Serious error attempt: user status is not active but is not black list, username:", user.Username)
			BlockUser(user.Username)
			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
				"error": "invalid login attempt",
			})
		}

		var token *jwt.Token
		if existingToken, ok := tokens[user.Username]; ok {
			claim, ok := claims[existingToken]
			if !ok || time.Now().After(claim.ExpiresAt.Time) {
				fmt.Println("token is expired, generate a new one")
				token = generateNewToken(user)
				tokens[user.Username] = token
			} else {
				token = existingToken
			}
		} else {
			token = generateNewToken(user)
			tokens[user.Username] = token
		}

		tokenString, err := token.SignedString([]byte(os.Getenv("JWT_SECRET")))
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Could not generate token",
			})
		}

		return c.JSON(LoginResponse{
			Token: tokenString,
			User:  *user,
		})
	}
}

func HandleRegister() fiber.Handler {
	return func(c *fiber.Ctx) error {
		var req RegisterRequest
		if err := c.BodyParser(&req); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid Request"})
		}

		err := services.AddUserUserService(utils.SERVICES_ROUTES.UserService, req.Username, req.Password, req.Email, req.Role)
		if err != nil {
			fmt.Println("Register Error:", err)
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
		}
		return c.Status(fiber.StatusCreated).JSON(fiber.Map{"message": "User Created"})
	}
}
