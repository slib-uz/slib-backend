package context

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/labstack/echo/v4"
	"slib.uz/src/core/application/response"
)

func GetBody[T any](c echo.Context) (*T, error) {
	data := new(T)
	if err := c.Bind(data); err != nil {
		message := fmt.Sprintf("Failed to bind request body: %v", err)
		return nil, response.NewFailResponse(400, message)
	}
	if err := c.Validate(data); err != nil {
		return nil, response.NewFailResponse(400, err.Error())
	}
	return data, nil
}

func GetStringListQueryParam(c echo.Context, param string) []string {
	raw := c.QueryParam(param)
	if raw == "" {
		return []string{}
	}
	parts := strings.Split(raw, ",")
	result := make([]string, 0, len(parts))

	for _, p := range parts {
		p = strings.TrimSpace(p)
		if p != "" {
			result = append(result, p)
		}
	}
	return result
}

func GetStringQueryParam(c echo.Context, param string) *string {
	value := c.QueryParam(param)
	if value == "" {
		return nil
	}
	return &value
}

func GetUintListQueryParam(c echo.Context, param string) ([]uint, error) {
	raw := c.QueryParam(param)
	if raw == "" {
		return []uint{}, nil
	}
	parts := strings.Split(raw, ",")
	ids := make([]uint, 0, len(parts))

	for _, p := range parts {
		p = strings.TrimSpace(p)
		if p == "" {
			continue
		}
		val, err := strconv.ParseUint(p, 10, 64)
		if err != nil {
			return nil, err
		}
		ids = append(ids, uint(val))
	}
	return ids, nil
}

func GetIntPathParam(c echo.Context, param string) (int, error) {
	value := c.Param(param)
	if value == "" {
		message := fmt.Sprintf("%s is required", param)
		return 0, response.NewFailResponse(400, message)
	}

	id, err := strconv.Atoi(value)
	if err != nil {
		message := fmt.Sprintf("Invalid %s: %v", param, err)
		return 0, response.NewFailResponse(400, message)
	}

	return id, nil
}

func GetIntQueryParam(c echo.Context, param string, defaultValue int) int {
	value := c.QueryParam(param)
	if value == "" {
		return defaultValue
	}

	id, err := strconv.Atoi(value)
	if err != nil {
		message := fmt.Sprintf("Invalid %s: %v", param, err)
		panic(response.NewFailResponse(400, message))
	}
	return id
}

func GetBoolQueryParam(c echo.Context, param string) *bool {
	value := c.QueryParam(param)
	if value == "" {
		return nil
	}

	parsed, err := strconv.ParseBool(value)
	if err != nil {
		message := fmt.Sprintf("Invalid %s: %v", param, err)
		panic(response.NewFailResponse(400, message))
	}
	return &parsed
}

func GetPagingParams(c echo.Context) (_page int, _size int) {
	page := c.QueryParam("page")
	size := c.QueryParam("page_size")

	if page == "" {
		_page = 1
	} else {
		_page, _ = strconv.Atoi(page)
	}

	if size == "" {
		_size = 10
	} else {
		_size, _ = strconv.Atoi(size)
	}

	return _page, _size
}

const dateLayout = "2006-01-02"

// GetDateRangeQueryParams parses ?start=YYYY-MM-DD&end=YYYY-MM-DD
// and returns (startOfDay, endOfDay) in Asia/Tashkent, converted to UTC.
func GetDateRangeQueryParams(c echo.Context, startParam, endParam string) (time.Time, time.Time) {
	loc := time.FixedZone("UTC+5", 5*60*60)

	// --- START DATE ---
	var startDate time.Time
	startStr := c.QueryParam(startParam)
	if startStr != "" {
		// Parse in +05 zone, so “2025-07-01” → 2025-07-01 00:00 in +05
		t, err := time.ParseInLocation(dateLayout, startStr, loc)
		if err != nil {
			panic(response.NewFailResponse(400,
				fmt.Sprintf("Invalid %s: %v", startParam, err)))
		}
		startDate = time.Date(t.Year(), t.Month(), t.Day(),
			0, 0, 0, 0, loc)
		// DB odatda UTC saqlaydi → konvert
		startDate = startDate.UTC()
	}

	// --- END DATE ---
	var endDate time.Time
	endStr := c.QueryParam(endParam)
	if endStr != "" {
		t, err := time.ParseInLocation(dateLayout, endStr, loc)
		if err != nil {
			panic(response.NewFailResponse(400,
				fmt.Sprintf("Invalid %s: %v", endParam, err)))
		}
		// 23:59:59.999 999 999 (kun oxiri)
		endDate = time.Date(t.Year(), t.Month(), t.Day(),
			23, 59, 59, int(time.Second-time.Nanosecond), loc)
	} else {
		now := time.Now().In(loc)
		endDate = time.Date(now.Year(), now.Month(), now.Day(),
			23, 59, 59, int(time.Second-time.Nanosecond), loc)
	}
	endDate = endDate.UTC()

	return startDate, endDate
}
