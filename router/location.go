package router

import (
	"context"
	"net/http"

	"github.com/oscto/ky3k"

	"gateway.oscto.icu/handler/location"
	"github.com/gin-gonic/gin"
	"go-micro.dev/v4"
)

var LocationServiceName = "location.oscto.icu"

type LocationBody struct {
	City struct {
		GeoNameID int `json:"GeoNameID"`
		Names     struct {
			De   string `json:"de"`
			En   string `json:"en"`
			Es   string `json:"es"`
			Fr   string `json:"fr"`
			Ja   string `json:"ja"`
			PtBR string `json:"pt-BR"`
			Ru   string `json:"ru"`
			ZhCN string `json:"zh-CN"`
		} `json:"Names"`
	} `json:"City"`
	Continent struct {
		Code      string `json:"Code"`
		GeoNameID int    `json:"GeoNameID"`
		Names     struct {
			De   string `json:"de"`
			En   string `json:"en"`
			Es   string `json:"es"`
			Fr   string `json:"fr"`
			Ja   string `json:"ja"`
			PtBR string `json:"pt-BR"`
			Ru   string `json:"ru"`
			ZhCN string `json:"zh-CN"`
		} `json:"Names"`
	} `json:"Continent"`
	Country struct {
		GeoNameID         int    `json:"GeoNameID"`
		IsInEuropeanUnion bool   `json:"IsInEuropeanUnion"`
		IsoCode           string `json:"IsoCode"`
		Names             struct {
			De   string `json:"de"`
			En   string `json:"en"`
			Es   string `json:"es"`
			Fr   string `json:"fr"`
			Ja   string `json:"ja"`
			PtBR string `json:"pt-BR"`
			Ru   string `json:"ru"`
			ZhCN string `json:"zh-CN"`
		} `json:"Names"`
	} `json:"Country"`
	Location struct {
		AccuracyRadius int     `json:"AccuracyRadius"`
		Latitude       float64 `json:"Latitude"`
		Longitude      float64 `json:"Longitude"`
		MetroCode      int     `json:"MetroCode"`
		TimeZone       string  `json:"TimeZone"`
	} `json:"Location"`
	Postal struct {
		Code string `json:"Code"`
	} `json:"Postal"`
	RegisteredCountry struct {
		GeoNameID         int    `json:"GeoNameID"`
		IsInEuropeanUnion bool   `json:"IsInEuropeanUnion"`
		IsoCode           string `json:"IsoCode"`
		Names             struct {
			De   string `json:"de"`
			En   string `json:"en"`
			Es   string `json:"es"`
			Fr   string `json:"fr"`
			Ja   string `json:"ja"`
			PtBR string `json:"pt-BR"`
			Ru   string `json:"ru"`
			ZhCN string `json:"zh-CN"`
		} `json:"Names"`
	} `json:"RegisteredCountry"`
	RepresentedCountry struct {
		GeoNameID         int         `json:"GeoNameID"`
		IsInEuropeanUnion bool        `json:"IsInEuropeanUnion"`
		IsoCode           string      `json:"IsoCode"`
		Names             interface{} `json:"Names"`
		Type              string      `json:"Type"`
	} `json:"RepresentedCountry"`
	Subdivisions []struct {
		GeoNameID int    `json:"GeoNameID"`
		IsoCode   string `json:"IsoCode"`
		Names     struct {
			En   string `json:"en"`
			Fr   string `json:"fr"`
			ZhCN string `json:"zh-CN"`
		} `json:"Names"`
	} `json:"Subdivisions"`
	Traits struct {
		IsAnonymousProxy    bool `json:"IsAnonymousProxy"`
		IsSatelliteProvider bool `json:"IsSatelliteProvider"`
	} `json:"Traits"`
}

func Location(r *gin.Engine) {
	r.GET("/ip", func(ctx *gin.Context) {
		clientIP := ctx.ClientIP()
		client := NewLocationClient()
		rsp, err := client.Call(context.Background(), &location.CallRequest{ClientIp: clientIP})

		if err != nil {

		}

		var location LocationBody
		ky3k.StringToJson(rsp.Location, &location)
		ctx.JSON(http.StatusOK, location)
	})

}
func NewLocationClient() location.LocationService {

	service := micro.NewService()
	return location.NewLocationService(LocationServiceName, service.Client())
}
