package handlers

type TokenResponse struct {
	AccessToken string `json:"access" example:"eyJhbGciOiJIUzUxMiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3NTE5MDUzMDQsInN1YiI6ImZiNjJhYTgxLTExNzItNGM3My04ZmMzLWNkNWE0NDYzNDZiYSJ9.SZHR-VexEcSNwe1GbmiG0p8lQVMTLH9MOIWV2N3I4ZMXEtYWF4Zcm4SKeaGFND7JCZ858VmId1WgPXKxTzF_iA"`
}

type UserIdResponse struct {
	UserId string `json:"userId" example:"eyJhbGciOiJIUzUxMiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3NTE5MDUzMDQsInN1YiI6ImZiNjJhYTgxLTExNzItNGM3My04ZmMzLWNkNWE0NDYzNDZiYSJ9.SZHR-VexEcSNwe1GbmiG0p8lQVMTLH9MOIWV2N3I4ZMXEtYWF4Zcm4SKeaGFND7JCZ858VmId1WgPXKxTzF_iA"`
}

type ErrorResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
}

// Эта структура только для примера ответа в swagger
type errBadRequestResp struct {
	Success bool   `json:"success" example:"false"`
	Message string `json:"message" example:"signup error: uncorrect uuid"`
}

var _ = errBadRequestResp{}

// Эта структура только для примера ответа в swagger
type errSignupResp struct {
	Success bool   `json:"success" example:"false"`
	Message string `json:"message" example:"signup errror: error on server side or user already exist"`
}

var _ = errSignupResp{}

// Эта структура только для примера ответа в swagger
type errGetUserIdResp struct {
	Success bool   `json:"success" example:"false"`
	Message string `json:"message" example:"getUserId error: error on server side"`
}

var _ = errGetUserIdResp{}
