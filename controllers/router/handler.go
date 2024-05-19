package router

import (
	"io"
	"net/http"
	"strings"

	"assignment/Ingress-Controller/controllers"
	ctrl "sigs.k8s.io/controller-runtime"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	ingressRouter controllers.IngressRouterEval
}

func NewHandler(ingressRouter controllers.IngressRouterEval) *Handler {
	return &Handler{
		ingressRouter: ingressRouter,
	}
}

func (h *Handler) Health(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{})
}

func (h *Handler) Route(c *gin.Context) {
	// Log the incoming HTTP request
	ctrl.Log.Info("Incoming HTTP request")

	// Extract the request host and get the first part
	reqHost := strings.Split(c.Request.Host, ".")[0]
	if reqHost == "" {
		c.String(http.StatusNotFound, "Host not found")
		return
	}
	ctrl.Log.Info("Searching route", "route", reqHost)

	// Attempt to route the request using the ingress router
	res, err := h.ingressRouter.RunNewRoute(reqHost)
	if err != nil {
		c.String(http.StatusNotFound, "Host not found")
		return
	}

	// Write the response from the routed request
	c.Status(res.StatusCode)
	for key, values := range res.Header {
		for _, value := range values {
			c.Writer.Header().Add(key, value)
		}
	}
	// Write the response body to the client
	_, err = io.Copy(c.Writer, res.Body)
	if err != nil {
		c.String(http.StatusInternalServerError, "Failed to return response")
	}
}
