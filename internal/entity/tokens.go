package entity

// Tokens is a pair of access and refresh token which uses for auth operations
//
//	@Description	Pair of access and refresh token which uses for auth operations
type Tokens struct {
	AccessToken  string `json:"accessToken" example:"header.payload.signature"`
	RefreshToken string `json:"refreshToken" example:"header.payload.signature"`
}

// TokensUIDs is additional layer of abstraction which we need for security reasons and for performing auth operations
type TokensUIDs struct {
	AccessTokenUID  string `json:"accessTokenUID"`
	RefreshTokenUID string `json:"refreshTokenUID"`
}
