package entity

type Tokens struct {
	AccessToken  string `json:"accessToken"`
	RefreshToken string `json:"refreshToken"`
}

type TokensUIDs struct {
	AccessTokenUID  string `json:"accessTokenUID"`
	RefreshTokenUID string `json:"refreshTokenUID"`
}
