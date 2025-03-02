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

// AdminService in adminservice/internal/services/admin_service.go
func (s *AdminService) GetAllUsers(page, limit int, role, status, search string, isVerified *bool, dateFrom, dateTo, sort, sortDir string) (*models.PaginatedResponse, error) {
	req := fasthttp.AcquireRequest()
	resp := fasthttp.AcquireResponse()
	defer fasthttp.ReleaseRequest(req)
	defer fasthttp.ReleaseResponse(resp)

	// Build URL with query parameters
	url := fmt.Sprintf("%s/user/get-all-user?page=%d&limit=%d",
		s.config.ExternalServices.UserService, page, limit)

	// Add filters to the query string
	if role != "" {
		url += fmt.Sprintf("&role=%s", role)
	}

	if status != "" {
		url += fmt.Sprintf("&status=%s", status)
	}

	if search != "" {
		url += fmt.Sprintf("&search=%s", search)
	}

	if isVerified != nil {
		url += fmt.Sprintf("&is_verified=%t", *isVerified)
	}

	if dateFrom != "" {
		url += fmt.Sprintf("&created_from=%s", dateFrom)
	}

	if dateTo != "" {
		url += fmt.Sprintf("&created_to=%s", dateTo)
	}

	if sort != "" {
		url += fmt.Sprintf("&sort=%s", sort)
	}

	if sortDir != "" {
		url += fmt.Sprintf("&sort_dir=%s", sortDir)
	}

	utils.BuildRequest(req, "GET", nil, s.config.APIKey, url)

	if err := fasthttp.Do(req, resp); err != nil {
		return nil, fmt.Errorf("failed to connect to user service: %v", err)
	}

	if resp.StatusCode() != fasthttp.StatusOK {
		return nil, fmt.Errorf("user service returned error: %s", resp.Body())
	}

	// Parse response
	var response models.PaginatedResponse

	if err := json.Unmarshal(resp.Body(), &response); err != nil {
		return nil, fmt.Errorf("failed to parse user service response: %v", err)
	}

	return &response, nil
}

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
