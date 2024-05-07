package lib

func MakeAuthHeader(token string) string {
	return "Bearer " + token
}
