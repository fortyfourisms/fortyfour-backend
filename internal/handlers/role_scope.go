package handlers

func isCompanyScopedRole(role string) bool {
	return role == "user" || role == "user_pic"
}
