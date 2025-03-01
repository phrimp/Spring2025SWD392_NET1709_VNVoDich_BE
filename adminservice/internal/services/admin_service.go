package services

import (
	"adminservice/internal/config"
	"adminservice/internal/models"
	"adminservice/utils"
	"encoding/json"
	"fmt"

	"github.com/valyala/fasthttp"
)

type AdminService struct {
	config *config.Config
}

func NewAdminService(cfg *config.Config) *AdminService {
	return &AdminService{
		config: cfg,
	}
}

func (s *AdminService) GetAllUsers(page, limit int, role, status, search string) ([]models.User, int, error) {
	req := fasthttp.AcquireRequest()
	resp := fasthttp.AcquireResponse()
	defer fasthttp.ReleaseRequest(req)
	defer fasthttp.ReleaseResponse(resp)

	// Build URL with query parameters
	url := fmt.Sprintf("%s/user/get-all-user?page=%d&limit=%d",
		s.config.ExternalServices.UserService, page, limit)

	if role != "" {
		url += fmt.Sprintf("&role=%s", role)
	}

	if status != "" {
		url += fmt.Sprintf("&status=%s", status)
	}

	if search != "" {
		url += fmt.Sprintf("&search=%s", search)
	}

	utils.BuildRequest(req, "GET", nil, s.config.APIKey, url)

	if err := fasthttp.Do(req, resp); err != nil {
		return nil, 0, fmt.Errorf("failed to connect to user service: %v", err)
	}

	if resp.StatusCode() != fasthttp.StatusOK {
		return nil, 0, fmt.Errorf("user service returned error: %s", resp.Body())
	}

	// Parse response
	var response struct {
		Data       []models.User `json:"data"`
		Pagination struct {
			Total int `json:"total"`
		} `json:"pagination"`
	}

	if err := json.Unmarshal(resp.Body(), &response); err != nil {
		// If the structure doesn't match, try just unmarshaling the users array
		var users []models.User
		if innerErr := json.Unmarshal(resp.Body(), &users); innerErr != nil {
			return nil, 0, fmt.Errorf("failed to parse user data: %v", err)
		}

		// If successful, use the array length as an estimate for the total
		total := len(users) + (page-1)*limit
		return users, total, nil
	}

	return response.Data, response.Pagination.Total, nil
}

// DeleteUser deletes a user
func (s *AdminService) DeleteUser(id uint) error {
	req := fasthttp.AcquireRequest()
	resp := fasthttp.AcquireResponse()
	defer fasthttp.ReleaseRequest(req)
	defer fasthttp.ReleaseResponse(resp)

	url := fmt.Sprintf("%s/user/delete/%d", s.config.ExternalServices.UserService, id)
	utils.BuildRequest(req, "DELETE", nil, s.config.APIKey, url)

	if err := fasthttp.Do(req, resp); err != nil {
		return fmt.Errorf("failed to connect to user service: %v", err)
	}

	if resp.StatusCode() != fasthttp.StatusOK {
		return fmt.Errorf("user service returned error: %s", resp.Body())
	}

	return nil
}

// UpdateUserStatus updates a user's status
func (s *AdminService) UpdateUserStatus(id uint, status string) error {
	req := fasthttp.AcquireRequest()
	resp := fasthttp.AcquireResponse()
	defer fasthttp.ReleaseRequest(req)
	defer fasthttp.ReleaseResponse(resp)

	// Create request body
	bodyJSON, err := json.Marshal(map[string]string{"status": status})
	if err != nil {
		return fmt.Errorf("failed to create request body: %v", err)
	}

	url := fmt.Sprintf("%s/user/status/%d", s.config.ExternalServices.UserService, id)
	utils.BuildRequest(req, "PUT", bodyJSON, s.config.APIKey, url)

	if err := fasthttp.Do(req, resp); err != nil {
		return fmt.Errorf("failed to connect to user service: %v", err)
	}

	if resp.StatusCode() != fasthttp.StatusOK {
		return fmt.Errorf("user service returned error: %s", resp.Body())
	}

	return nil
}

// GetAllCourses retrieves all courses from the node service with optional filtering
func (s *AdminService) GetAllCourses(page, limit int, subject, status, grade, search string) ([]models.Course, int, error) {
	req := fasthttp.AcquireRequest()
	resp := fasthttp.AcquireResponse()
	defer fasthttp.ReleaseRequest(req)
	defer fasthttp.ReleaseResponse(resp)

	// Build URL with query parameters
	url := fmt.Sprintf("%s/courses?page=%d&pageSize=%d",
		s.config.ExternalServices.NodeService, page, limit)

	if subject != "" {
		url += fmt.Sprintf("&subject=%s", subject)
	}

	if status != "" {
		url += fmt.Sprintf("&status=%s", status)
	}

	if grade != "" {
		url += fmt.Sprintf("&grade=%s", grade)
	}

	if search != "" {
		url += fmt.Sprintf("&title=%s", search)
	}

	utils.BuildRequest(req, "GET", nil, s.config.APIKey, url)

	if err := fasthttp.Do(req, resp); err != nil {
		return nil, 0, fmt.Errorf("failed to connect to node service: %v", err)
	}

	if resp.StatusCode() != fasthttp.StatusOK {
		return nil, 0, fmt.Errorf("node service returned error: %s", resp.Body())
	}

	// Parse response
	var coursesResponse struct {
		Data       []models.Course `json:"data"`
		Pagination struct {
			Total int `json:"total"`
		} `json:"pagination"`
	}

	if err := json.Unmarshal(resp.Body(), &coursesResponse); err != nil {
		return nil, 0, fmt.Errorf("failed to parse course data: %v", err)
	}

	return coursesResponse.Data, coursesResponse.Pagination.Total, nil
}
