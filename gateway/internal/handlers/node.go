package handlers

import (
	"gateway/internal/routes"

	"github.com/gofiber/fiber/v2"
	"github.com/valyala/fasthttp"
)

type NodeServiceHandler struct {
	nodeServiceUrl string
}

func NewNodeServiceHandler(nodeServiceURL string) *NodeServiceHandler {
	return &NodeServiceHandler{
		nodeServiceUrl: nodeServiceURL,
	}
}

func (h *NodeServiceHandler) HandleGetAllCourse() fiber.Handler {
	return func(c *fiber.Ctx) error {
		req := fasthttp.AcquireRequest()
		resp := fasthttp.AcquireResponse()
		defer fasthttp.ReleaseRequest(req)
		defer fasthttp.ReleaseResponse(resp)
		return routes.GetAllCourse(req, resp, c, h.nodeServiceUrl+"/courses/")
	}
}

func (h *NodeServiceHandler) HandleGetACourse() fiber.Handler {
	return func(c *fiber.Ctx) error {
		req := fasthttp.AcquireRequest()
		resp := fasthttp.AcquireResponse()
		defer fasthttp.ReleaseRequest(req)
		defer fasthttp.ReleaseResponse(resp)
		return routes.GetACourse(req, resp, c, h.nodeServiceUrl+"/courses/"+c.Params("id"))
	}
}
