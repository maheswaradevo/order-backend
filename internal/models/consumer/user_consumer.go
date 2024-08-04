package consumer

type UserEvent struct {
	ID                  uint64 `json:"id"`
	Username            string `json:"username"`
	Email               string `json:"email"`
	AccessToken         string `json:"access_token"`
	RefreshToken        string `json:"refresh_token"`
	ExpiredToken        int64  `json:"expired_token"`
	ExpiredRefreshToken int64  `json:"expired_refresh_token"`
}
