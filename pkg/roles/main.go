package roles

var roles = map[string]bool{
	"hyperspace": true,
	"admin":      true,
	"superuser":  true,
}

func IsValidRole(role string) bool {
	_, ok := roles[role]
	return ok
}
