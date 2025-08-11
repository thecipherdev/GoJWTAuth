package mock

import "github.com/thecipherdev/goauth/model"

var Users = []model.User{
	{
		Username: "johndoe",
		Password: "$argon2id$v=19$m=65536,t=3,p=4$5+OVrxlJkVGWOzk+0ifqgA$/XB7Q6t+ihgUSDl+tBKwCWbuhHZQrexS3cIvejlReyI",
	},
}
