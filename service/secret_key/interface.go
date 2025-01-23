package secretkey

type SecretKeyInterface interface {

	// GetSecretKey 获取密钥
	GetSecretKey() (string, error)

	// SetSecretKey 设置密钥
	SetSecretKey() error
}
