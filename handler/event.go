package handler

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/labstack/echo"
	"github.com/maps90/librarian"
	"github.com/mataharimall/micro/api"
	"github.com/mataharimall/micro/helper"
)

type EventsList struct {
	Request  interface{}
	Response interface{}
}

type SearchEventPayload struct {
	Data struct {
		Search struct {
			EventName    string `json:"event_name" query:"event_name"`
			LocationName string `json:"location_name" query:"location_name"`
			City         string `json:"city" query:"city"`
			MinStartDate string `json:"min_start_date" query:"min_start_date"`
			MaxEndDate   string `json:"max_end_date" query:"max_end_date"`
			MinPrice     string `json:"min_price" query:"min_price"`
			MaxPrice     string `json:"max_price" query:"max_price"`
		} `json:"search"`
		Limit  string `json:"limit" query:"limit"`
		Offset string `json:"offset" query:"offset"`
	} `json:"data"`
}

func GetEventList(c echo.Context) error {
	loket, ok := librarian.Get("loket").(*api.Loket)
	if !ok {
		return echo.NewHTTPError(http.StatusInternalServerError, "Internal Server Error")

	}

	loket.CacheOn().Post("/v3/event", "form", "")
	return helper.BuildJSON(c, loket.Response.Data, loket.Error)
}

func GetEventDetail(c echo.Context) error {
	loket, ok := librarian.Get("loket").(*api.Loket)
	if !ok {
		return echo.NewHTTPError(http.StatusInternalServerError, "Internal Server Error")

	}

	url := fmt.Sprintf(`/v3/event/%s`, c.Param("event_id"))
	loket.CacheOn().Post(url, "form", "")
	return helper.BuildJSON(c, loket.Response.Data, loket.Error)
}

func SearchEvent(c echo.Context) error {
	loket, ok := librarian.Get("loket").(*api.Loket)
	if !ok {
		return echo.NewHTTPError(http.StatusInternalServerError, "Internal Server Error")
	}

	s := SearchEventPayload{}

	if err := c.Bind(&s); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	search_json, err := json.Marshal(s)
	if err != nil {
		loket.Error = err
	}

	loket.CacheOn().Post("/v3/event_search", "json", string(search_json))
	return helper.BuildJSON(c, loket.Response.Data, loket.Error)
}

func GetEventCities(c echo.Context) error {
	loket, ok := librarian.Get("loket").(*api.Loket)
	if !ok {
		return echo.NewHTTPError(http.StatusInternalServerError, "Internal Server Error")
	}

	loket.CacheOn().Post("/v3/event_list/cities", "form", "")
	return helper.BuildJSON(c, loket.Response.Data, loket.Error)
}
