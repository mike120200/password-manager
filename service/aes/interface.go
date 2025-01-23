package aes

type AesInterface interface {
	// 加密
	Encrypt(plainText string) ([]byte, []byte, error)
	// 解密
	Decrypt(cipherData, nonce []byte) ([]byte, error)
}
