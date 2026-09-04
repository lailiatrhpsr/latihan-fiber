package helper

import (
	"context"
	"strconv"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"

	"latihan-fiber/pert4/app/model"
)

// memberi batas timeout untuk eksekusi ke database
func RequestContext(c *fiber.Ctx) (context.Context, context.CancelFunc) {
	return context.WithTimeout(c.UserContext(), 5*time.Second)
}

// membaca dan memvalidasi id integer dari URL /:id
func ParamID(c *fiber.Ctx) (int, bool) {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil || id < 1 {
		return 0, false
	}
	return id, true
}

var allowedSortHelper = map[string]bool{
	"id":    true,
	"name":  true,
	"grade": true,
}

func ParseListQuery(c *fiber.Ctx) model.ListQuery {
	q := model.ListQuery{
		Page:  1,
		Limit: 10,
		Sort:  "id",
		Order: "asc",
	}

	if p, err := strconv.Atoi(c.Query("page")); err == nil && p > 0 {
		q.Page = p
	}

	if l, err := strconv.Atoi(c.Query("limit")); err == nil && l > 0 {
		q.Limit = l
	}

	q.Search = strings.TrimSpace(c.Query("search"))

	sortParam := c.Query("sort")
	if allowedSortHelper[sortParam] {
		q.Sort = sortParam
	}

	orderParam := strings.ToLower(c.Query("order"))
	if orderParam == "desc" {
		q.Order = "desc"
	}

	if raw := c.Query("is_active"); raw != "" {
		val := strings.ToLower(raw) == "true"
		q.IsActive = &val
	}

	return q
}
