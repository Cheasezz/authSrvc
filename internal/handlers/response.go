package handlers

import (
	"fmt"
	"net/http"

	"github.com/Cheasezz/authSrvc/internal/core"
	"github.com/gin-gonic/gin"
)

type TokenResponse struct {
	AccessToken string `json:"access" example:"eyJhbGciOiJIUzUxMiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3NTE5MDUzMDQsInN1YiI6ImZiNjJhYTgxLTExNzItNGM3My04ZmMzLWNkNWE0NDYzNDZiYSJ9.SZHR-VexEcSNwe1GbmiG0p8lQVMTLH9MOIWV2N3I4ZMXEtYWF4Zcm4SKeaGFND7JCZ858VmId1WgPXKxTzF_iA"`
	// Это поле только для swagger и имеет значение только в dev моде.
	RefreshCookie string `json:"refresh_cookie,omitempty" example:"refreshToken=eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzZXNzaW9uX2lkIjoiZDdhOTk2YjctOWZlNi00YjRlLWI4NWItODM3YzAyN2RmNDU3Iiwic3ViIjoiZmI2MmFhODEtMTE3Mi00YzczLThmYzMtY2Q1YTQ0NjM0NmJhIiwiZXhwIjoxNzUzNjE4MzUxfQ.SdTNSloBSOnxHJeq6FWlN3UuiyZBOzL9P5OQVp23Wlg; Path=/; Max-Age=86400; HttpOnly; SameSite=Lax"`
}

func newTokensResponse(c *gin.Context, tkns *core.TokenPairResult) {
	c.SetSameSite(http.SameSiteLaxMode)
	c.SetCookie(refreshTknCookie, tkns.Refresh, int(tkns.RefreshTTL.Seconds()), "", "", false, true)

	if gin.Mode() == gin.ReleaseMode {
		c.JSON(http.StatusOK, TokenResponse{AccessToken: tkns.Access})
	} else {
		c.JSON(http.StatusOK, TokenResponse{
			AccessToken: tkns.Access,
			RefreshCookie: fmt.Sprintf(
				"refreshToken=%s; Path=/; Max-Age=86400; HttpOnly; SameSite=Lax",
				tkns.Refresh),
		})

	}
}

type MeResponse struct {
	UserId string `json:"userId" example:"eyJhbGciOiJIUzUxMiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3NTE5MDUzMDQsInN1YiI6ImZiNjJhYTgxLTExNzItNGM3My04ZmMzLWNkNWE0NDYzNDZiYSJ9.SZHR-VexEcSNwe1GbmiG0p8lQVMTLH9MOIWV2N3I4ZMXEtYWF4Zcm4SKeaGFND7JCZ858VmId1WgPXKxTzF_iA"`
}

type ErrorResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
}

// Эта структура только для примера ответа в swagger
type errBadRequestResp struct {
	Success bool   `json:"success" example:"false"`
	Message string `json:"message" example:"create session: error: uncorrect uuid"`
}

var _ = errBadRequestResp{}

// Эта структура только для примера ответа в swagger
type errTokenIssuanceResp struct {
	Success bool   `json:"success" example:"false"`
	Message string `json:"message" example:"create session error: error on server side or user already exist"`
}

var _ = errTokenIssuanceResp{}

// Эта структура только для примера ответа в swagger
type errMeResp struct {
	Success bool   `json:"success" example:"false"`
	Message string `json:"message" example:"me error: error on server side"`
}

var _ = errMeResp{}

// Эта структура только для примера ответа в swagger
type errAuthResp struct {
	Success bool   `json:"success" example:"false"`
	Message string `json:"message" example:"error authorization"`
}

var _ = errAuthResp{}
