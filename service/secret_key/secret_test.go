package secretkey_test

import (
	zaplog "password_manager/common/log"
	secretkey "password_manager/service/secret_key"
	"testing"
)

// 测试保存密钥
func TestSaveSecretKey(t *testing.T) {
	instance := secretkey.NewSecretKeyWithFilePath("./test.gob")
	if err := instance.SetSecretKey(); err != nil {
		t.Error(err)
		return
	}
}

// 测试获取密钥
func TestGetSecretKey(t *testing.T) {
	instance := secretkey.NewSecretKeyWithFilePath("./test.gob")
	key, err := instance.GetSecretKey()
	if err != nil {
		t.Error(err)
		return
	}
	t.Log(key)
}

func init() {
	zaplog.LoggerInit()
}
