package handlers

import (
	"fmt"
	"gateway/internal/middleware"
	"gateway/internal/routes"
	"math"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/valyala/fasthttp"
)

type UserServiceHandler struct {
	userServiceURL string
}

func NewUserService(userServiceURL string) *UserServiceHandler {
	return &UserServiceHandler{
		userServiceURL: userServiceURL,
	}
}

func (h *UserServiceHandler) HandleGetPublicUserProfile() fiber.Handler {
	return func(c *fiber.Ctx) error {
		req := fasthttp.AcquireRequest()
		resp := fasthttp.AcquireResponse()
		defer fasthttp.ReleaseRequest(req)
		defer fasthttp.ReleaseResponse(resp)

		username := c.Params("username")
		if username == "" {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Username is required",
			})
		}

		transformer := func(originalData interface{}) (interface{}, error) {
			userData, ok := originalData.(map[string]interface{})
			if !ok {
				return nil, fmt.Errorf("unexpected response format")
			}

			// Create a properly structured public profile response
			result := make(map[string]interface{})

			// Include basic information
			result["username"] = userData["username"]
			result["full_name"] = userData["full_name"]
			result["picture"] = userData["picture"]
			result["role"] = userData["role"]
			result["is_verified"] = userData["is_verified"]

			// Convert specific timestamp to more private format
			if createdAt, ok := userData["created_at"].(string); ok {
				if t, err := time.Parse(time.RFC3339, createdAt); err == nil {
					result["join_date"] = t.Format("January 2006") // Only month and year
				}
			}

			// Handle role-specific public information
			roleInfo := make(map[string]interface{})

			switch userData["role"] {
			case "Tutor":
				if tutor, ok := userData["tutor"].(map[string]interface{}); ok {
					// Include tutor information that's appropriate for public viewing
					roleInfo["bio"] = tutor["bio"]
					roleInfo["qualifications"] = tutor["qualifications"]
					roleInfo["teaching_style"] = tutor["teaching_style"]
					roleInfo["is_available"] = tutor["is_available"]

					// Include specialties but limit detailed information
					if specialties, ok := tutor["specialties"].([]interface{}); ok {
						simplifiedSpecialties := make([]map[string]interface{}, 0)
						for _, s := range specialties {
							if specialty, ok := s.(map[string]interface{}); ok {
								simplifiedSpecialties = append(simplifiedSpecialties, map[string]interface{}{
									"subject": specialty["subject"],
									"level":   specialty["level"],
								})
							}
						}
						roleInfo["specialties"] = simplifiedSpecialties
					}

					// Calculate and include average rating if available
					if reviews, ok := tutor["tutor_reviews"].([]interface{}); ok && len(reviews) > 0 {
						var totalRating float64
						for _, r := range reviews {
							if review, ok := r.(map[string]interface{}); ok {
								if rating, ok := review["rating"].(float64); ok {
									totalRating += rating
								}
							}
						}
						avgRating := totalRating / float64(len(reviews))
						roleInfo["average_rating"] = math.Round(avgRating*10) / 10 // Round to 1 decimal place
						roleInfo["review_count"] = len(reviews)
					}
				}
			case "Parent":
				// For parents, include minimal public information
				if _, ok := userData["created_at"].(string); ok {
					if t, err := time.Parse(time.RFC3339, userData["created_at"].(string)); err == nil {
						roleInfo["join_year"] = t.Year()
					}
				}
			}

			result["role_info"] = roleInfo

			return result, nil
		}

		return routes.CustomForwardRequest(
			req,
			resp,
			c,
			h.userServiceURL+"/user/get-public-user?username="+username,
			"GET",
			nil,
			transformer,
		)
	}
}

func (h *UserServiceHandler) HandleAllGetUser() fiber.Handler {
	return func(c *fiber.Ctx) error {
		req := fasthttp.AcquireRequest()
		resp := fasthttp.AcquireResponse()
		defer fasthttp.ReleaseRequest(req)
		defer fasthttp.ReleaseResponse(resp)
		query_url := fmt.Sprintf("?page=%s&limit=%s", c.Query("page"), c.Query("limit"))
		return routes.GetAllUser(req, resp, c, h.userServiceURL+"/user/get-all-user"+query_url)
	}
}

func (h *UserServiceHandler) HandleGetMe() fiber.Handler {
	return func(c *fiber.Ctx) error {
		req := fasthttp.AcquireRequest()
		resp := fasthttp.AcquireResponse()
		defer fasthttp.ReleaseRequest(req)
		defer fasthttp.ReleaseResponse(resp)
		claims, ok := c.Locals("user").(*middleware.Claims)
		if !ok {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "cannot find username in token claim"})
		}
		current_username := claims.Username
		query_url := fmt.Sprintf("?username=%s", current_username)

		transformer := func(originalData interface{}) (interface{}, error) {
			userData, ok := originalData.(map[string]interface{})
			if !ok {
				return nil, fmt.Errorf("unexpected response format")
			}

			// Create a structured user profile response
			result := make(map[string]interface{})

			// Basic user information (appropriate for self-view)
			result["id"] = userData["ID"]
			result["username"] = userData["username"]
			result["email"] = userData["email"]
			result["full_name"] = userData["full_name"]
			result["role"] = userData["role"]
			result["picture"] = userData["picture"]
			result["phone"] = userData["phone"]
			result["is_verified"] = userData["is_verified"]
			result["status"] = userData["status"]
			result["created_at"] = userData["created_at"]
			result["updated_at"] = userData["updated_at"]

			if userData["last_login_at"] != nil {
				result["last_login_at"] = userData["last_login_at"]
			}

			// Account settings and preferences
			accountSettings := make(map[string]interface{})

			// Add role-specific data
			switch userData["role"] {
			case "Tutor":
				if tutor, ok := userData["tutor"].(map[string]interface{}); ok {
					tutorProfile := make(map[string]interface{})

					// Basic tutor profile information
					tutorProfile["bio"] = tutor["bio"]
					tutorProfile["qualifications"] = tutor["qualifications"]
					tutorProfile["teaching_style"] = tutor["teaching_style"]
					tutorProfile["is_available"] = tutor["is_available"]
					tutorProfile["demo_video_url"] = tutor["demo_video_url"]
					tutorProfile["image"] = tutor["image"]

					// Include specialties
					if specialties, ok := tutor["specialties"].([]interface{}); ok {
						tutorProfile["specialties"] = specialties
					}

					// Include availability
					if availability, ok := tutor["availability"].(map[string]interface{}); ok {
						tutorProfile["availability"] = availability
					}

					// Include summary of courses
					if courses, ok := tutor["courses"].([]interface{}); ok {
						coursesSummary := make(map[string]interface{})
						coursesSummary["count"] = len(courses)
						coursesSummary["subjects"] = getUniqueSubjects(courses)
						tutorProfile["courses_summary"] = coursesSummary
					}

					// Include summary of reviews
					if reviews, ok := tutor["reviews"].([]interface{}); ok {
						reviewsSummary := make(map[string]interface{})
						reviewsSummary["count"] = len(reviews)
						if len(reviews) > 0 {
							reviewsSummary["average_rating"] = calculateAverageRating(reviews)
						}
						tutorProfile["reviews_summary"] = reviewsSummary
					}

					result["tutor_profile"] = tutorProfile
				}

			case "Parent":
				if parent, ok := userData["parent"].(map[string]interface{}); ok {
					parentProfile := make(map[string]interface{})

					// Basic parent profile information
					parentProfile["preferred_language"] = parent["preferred_language"]
					parentProfile["notifications_enabled"] = parent["notifications_enabled"]

					// Include children summary
					if children, ok := parent["children"].([]interface{}); ok {
						childrenData := make([]map[string]interface{}, 0)

						for _, child := range children {
							if childMap, ok := child.(map[string]interface{}); ok {
								childData := make(map[string]interface{})
								childData["id"] = childMap["id"]
								childData["full_name"] = childMap["full_name"]
								childData["age"] = childMap["age"]
								childData["grade_level"] = childMap["grade_level"]
								childData["learning_goals"] = childMap["learning_goals"]

								childrenData = append(childrenData, childData)
							}
						}

						parentProfile["children"] = childrenData
					}

					result["parent_profile"] = parentProfile
				}
			}

			// Add account settings
			accountSettings["google_linked"] = userData["google_token"] != nil && userData["google_token"] != ""
			result["account_settings"] = accountSettings

			return result, nil
		}
		return routes.GetMe(req, resp, c, h.userServiceURL+"/user"+query_url, transformer)
	}
}

func (h *UserServiceHandler) HandleDeleteMe() fiber.Handler {
	return func(c *fiber.Ctx) error {
		req := fasthttp.AcquireRequest()
		resp := fasthttp.AcquireResponse()
		defer fasthttp.ReleaseRequest(req)
		defer fasthttp.ReleaseResponse(resp)
		claims, ok := c.Locals("user").(*middleware.Claims)
		if !ok {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "cannot find username in token claim"})
		}
		current_username := claims.Username
		query_url := fmt.Sprintf("?username=%s", current_username)
		return routes.DeleteMe(req, resp, c, h.userServiceURL+"/delete"+query_url)
	}
}

func (h *UserServiceHandler) HandleCancelDeleteMe() fiber.Handler {
	return func(c *fiber.Ctx) error {
		req := fasthttp.AcquireRequest()
		resp := fasthttp.AcquireResponse()
		defer fasthttp.ReleaseRequest(req)
		defer fasthttp.ReleaseResponse(resp)
		claims, ok := c.Locals("user").(*middleware.Claims)
		if !ok {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "cannot find username in token claim"})
		}
		current_username := claims.Username
		query_url := fmt.Sprintf("?username=%s", current_username)
		return routes.CancelDeleteMe(req, resp, c, h.userServiceURL+"/delete/cancel"+query_url)
	}
}

func (h *UserServiceHandler) HandleUpdateMe() fiber.Handler {
	return func(c *fiber.Ctx) error {
		req := fasthttp.AcquireRequest()
		resp := fasthttp.AcquireResponse()
		defer fasthttp.ReleaseRequest(req)
		defer fasthttp.ReleaseResponse(resp)
		claims, ok := c.Locals("user").(*middleware.Claims)
		if !ok {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "cannot find username in token claim"})
		}
		current_username := claims.Username
		query_url := fmt.Sprintf("?username=%s", current_username)
		return routes.UpdateMe(req, resp, c, h.userServiceURL+"/user/update"+query_url)
	}
}

func (h *UserServiceHandler) HandleUpdateMePassword() fiber.Handler {
	return func(c *fiber.Ctx) error {
		req := fasthttp.AcquireRequest()
		resp := fasthttp.AcquireResponse()
		defer fasthttp.ReleaseRequest(req)
		defer fasthttp.ReleaseResponse(resp)
		claims, ok := c.Locals("user").(*middleware.Claims)
		if !ok {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "cannot find username in token claim"})
		}
		current_username := claims.Username
		query_url := fmt.Sprintf("?username=%s&new_password=%s&cur_password=%s", current_username, c.Query("new_password"), c.Query("cur_password"))
		return routes.UpdateMePassword(req, resp, c, h.userServiceURL+"/user/update/password"+query_url)
	}
}

func calculateAverageRating(reviews []interface{}) float64 {
	var totalRating float64
	var count int

	for _, r := range reviews {
		if review, ok := r.(map[string]interface{}); ok {
			if rating, ok := review["rating"].(float64); ok {
				totalRating += rating
				count++
			}
		}
	}

	if count == 0 {
		return 0
	}

	return math.Round((totalRating/float64(count))*10) / 10
}

func getUniqueSubjects(courses []interface{}) []string {
	subjectSet := make(map[string]bool)

	for _, c := range courses {
		if course, ok := c.(map[string]interface{}); ok {
			if subject, ok := course["subject"].(string); ok {
				subjectSet[subject] = true
			}
		}
	}

	subjects := make([]string, 0, len(subjectSet))
	for subject := range subjectSet {
		subjects = append(subjects, subject)
	}

	return subjects
}
