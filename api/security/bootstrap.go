package security

func Bootstrap() *Service {
	return newDefaultService()
}
