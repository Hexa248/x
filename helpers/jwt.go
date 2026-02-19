package helpers

func GenerateJWT(email, role, secret string) (string, error) {
	return email + ":" + role + ":" + secret, nil
}
