package aes_test

import (
	zaplog "password_manager/common/log"
	"password_manager/service/aes"
	"testing"
)

func TestAes(t *testing.T) {

	key := "1234567890123456"
	plainText := "hello world"
	instance := aes.NewAesService(key)
	cipherData, nonce, err := instance.Encrypt(plainText)
	if err != nil {
		t.Error(err)
		return
	}
	t.Log(cipherData, nonce)
	t.Log(len(nonce))
	DecryptData, err := instance.Decrypt(cipherData, nonce)
	if err != nil {
		t.Error(err)
		return
	}
	t.Log(DecryptData)
}

func init() {
	if err := zaplog.LoggerInit(); err != nil {
		panic(err)
	}
}
