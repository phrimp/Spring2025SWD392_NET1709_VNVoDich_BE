package services

import (
	"authservice/internal/repository"
	"authservice/utils"
	"encoding/json"
	"fmt"

	"github.com/gofiber/fiber/v2"
	"github.com/valyala/fasthttp"
)

func GetUserFromUserService(userServiceURL string, username, password string) (*repository.User, error) {
	req := &fasthttp.Request{}
	resp := &fasthttp.Response{}

	// Prepare request body
	body := fmt.Sprintf(`{"username":"%s","password":"%s"}`, username, password)
	utils.BuildRequest(req, "POST", []byte(body), utils.API_KEY, userServiceURL+"/user/get")

	if err := fasthttp.Do(req, resp); err != nil {
		return nil, fmt.Errorf("user service unavailable: %v", err)
	}

	if resp.StatusCode() != fiber.StatusOK {
		return nil, fmt.Errorf("invalid credentials")
	}

	var user repository.User
	if err := json.Unmarshal(resp.Body(), &user); err != nil {
		return nil, fmt.Errorf("failed to parse user data: %v", err)
	}

	return &user, nil
}

func AddUserUserService(userServiceURL, username, password, email, role string) error {
	req := &fasthttp.Request{}
	resp := &fasthttp.Response{}

	body := fmt.Sprintf(`{"username":"%s","password":"%s", "email":"%s", "role":"%s"}`, username, password, email, role)
	utils.BuildRequest(req, "POST", []byte(body), utils.API_KEY, userServiceURL+"/user/add")

	if err := fasthttp.Do(req, resp); err != nil {
		return fmt.Errorf("user service unavailable: %v", err)
	}

	if resp.StatusCode() != fiber.StatusOK {
		return fmt.Errorf("add user failed: %s", string(resp.Body()))
	}
	fmt.Println(string(resp.Body()))
	return nil
}
