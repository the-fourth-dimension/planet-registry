package roles

var roles = []string{"superuser", "admin", "hyperspace"}

func IsValidRole(roleValue int) bool {
	return roleValue >= 0 && roleValue < len(roles)
}
