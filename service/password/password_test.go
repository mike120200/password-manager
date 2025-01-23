package password_test

import (
	"math/rand"
	"os"
	zaplog "password_manager/common/log"
	"password_manager/service/aes"
	dbfilekit "password_manager/service/dbfile_Kit"
	"password_manager/service/password"
	secretkey "password_manager/service/secret_key"
	"strconv"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestSavePasswordAndGetPassword(t *testing.T) {
	os.Remove("./test.gob")
	os.Remove("./data.db")
	os.Remove("./data.backup.db")
	defer func() {
		os.Remove("./test.gob")
		os.Remove("./data.db")
		os.Remove("./data.backup.db")
	}()
	//初始化密钥实例
	secretKeyInstance := secretkey.NewSecretKeyWithFilePath("./test.gob")
	//初始化数据库实例
	dbfileKitInstance := dbfilekit.NewDBKitWithFilePath("./", secretKeyInstance)
	if err := dbfileKitInstance.Init(); err != nil {
		t.Log(err.Error())
		return
	}
	//获取数据库实例
	db, err := dbfileKitInstance.GetDB()
	if err != nil {
		t.Log(err.Error())
		return
	}
	defer db.Close()
	//获取密钥
	key, err := secretKeyInstance.GetSecretKey()
	if err != nil {
		t.Log(err.Error())
		return
	}
	aesInstance := aes.NewAesService(key)

	passwordInstance := password.NewPasswordService(aesInstance, db)
	// 设置随机种子
	src := rand.NewSource(time.Now().UnixNano())
	randObj := rand.New(src)
	num1 := strconv.Itoa(randObj.Intn(1000))
	num2 := strconv.Itoa(randObj.Intn(1000))

	//保存密码
	if err := passwordInstance.SavePassword("test"+num1, "test"+num2); err != nil {
		t.Log(err.Error())
		return
	}
	//查看密码
	value, err := passwordInstance.GetPasswordWithKey("test" + num1)
	if err != nil {
		t.Log(err.Error())
		return
	}
	t.Log(value)
}

func TestGetAllPassword(t *testing.T) {
	os.Remove("./test.gob")
	os.Remove("./data.db")
	os.Remove("./data.backup.db")
	defer func() {
		os.Remove("./test.gob")
		os.Remove("./data.db")
		os.Remove("./data.backup.db")
	}()
	//初始化密钥实例
	secretKeyInstance := secretkey.NewSecretKeyWithFilePath("./test.gob")
	//初始化数据库实例
	dbfileKitInstance := dbfilekit.NewDBKitWithFilePath("./", secretKeyInstance)
	if err := dbfileKitInstance.Init(); err != nil {
		t.Log(err.Error())
		return
	}
	//获取数据库实例
	db, err := dbfileKitInstance.GetDB()
	if err != nil {
		t.Log(err.Error())
		return
	}
	defer db.Close()
	//获取密钥
	key, err := secretKeyInstance.GetSecretKey()
	if err != nil {
		t.Log(err.Error())
		return
	}

	//初始化aes实例
	aesInstance := aes.NewAesService(key)
	passwordInstance := password.NewPasswordService(aesInstance, db)
	// 设置随机种子
	src := rand.NewSource(time.Now().UnixNano())
	randObj := rand.New(src)
	for i := 0; i < 10; i++ {
		num1 := strconv.Itoa(randObj.Intn(1000))
		num2 := strconv.Itoa(randObj.Intn(1000))

		if err := passwordInstance.SavePassword("test"+num1, "test"+num2); err != nil {
			t.Log(err.Error())
			return
		}
	}
	//获取所有的密码键值对
	values, err := passwordInstance.GetAllPasswords()
	if err != nil {
		t.Log(err.Error())
		return
	}
	t.Log(values)
}

func TestUpdatePassword(t *testing.T) {
	os.Remove("./test.gob")
	os.Remove("./data.db")
	os.Remove("./data.backup.db")
	defer func() {
		os.Remove("./test.gob")
		os.Remove("./data.db")
		os.Remove("./data.backup.db")
	}()
	//初始化密钥实例
	secretKeyInstance := secretkey.NewSecretKeyWithFilePath("./test.gob")
	//初始化数据库实例
	dbfileKitInstance := dbfilekit.NewDBKitWithFilePath("./", secretKeyInstance)
	if err := dbfileKitInstance.Init(); err != nil {
		t.Log(err.Error())
		return
	}
	//获取数据库实例
	db, err := dbfileKitInstance.GetDB()
	if err != nil {
		t.Log(err.Error())
		return
	}
	defer db.Close()
	//获取密钥
	key, err := secretKeyInstance.GetSecretKey()
	if err != nil {
		t.Log(err.Error())
		return
	}
	aesInstance := aes.NewAesService(key)

	passwordInstance := password.NewPasswordService(aesInstance, db)
	// 设置随机种子
	src := rand.NewSource(time.Now().UnixNano())
	randObj := rand.New(src)
	num1 := strconv.Itoa(randObj.Intn(1000))
	num2 := strconv.Itoa(randObj.Intn(1000))

	//保存密码
	if err := passwordInstance.SavePassword("test"+num1, "test"+num2); err != nil {
		t.Log(err.Error())
		return
	}
	//查看密码
	value, err := passwordInstance.GetPasswordWithKey("test" + num1)
	if err != nil {
		t.Log(err.Error())
		return
	}
	t.Logf("before value:%s", value)

	num3 := strconv.Itoa(randObj.Intn(1000))
	if err := passwordInstance.UpdatePassword("test"+num1, "test"+num3); err != nil {
		t.Log(err.Error())
		return
	}
	value, err = passwordInstance.GetPasswordWithKey("test" + num1)
	if err != nil {
		t.Log(err.Error())
		return
	}
	t.Logf("after value:%s", value)
}

func TestDeletePassword(t *testing.T) {
	assert_instance := assert.New(t)
	os.Remove("./test.gob")
	os.Remove("./data.db")
	os.Remove("./data.backup.db")
	defer func() {
		os.Remove("./test.gob")
		os.Remove("./data.db")
		os.Remove("./data.backup.db")
	}()
	//初始化密钥实例
	secretKeyInstance := secretkey.NewSecretKeyWithFilePath("./test.gob")
	//初始化数据库实例
	dbfileKitInstance := dbfilekit.NewDBKitWithFilePath("./", secretKeyInstance)
	if err := dbfileKitInstance.Init(); err != nil {
		t.Log(err.Error())
		return
	}
	//获取数据库实例
	db, err := dbfileKitInstance.GetDB()
	if err != nil {
		t.Log(err.Error())
		return
	}
	defer db.Close()
	//获取密钥
	key, err := secretKeyInstance.GetSecretKey()
	if err != nil {
		t.Log(err.Error())
		return
	}
	aesInstance := aes.NewAesService(key)

	passwordInstance := password.NewPasswordService(aesInstance, db)
	// 设置随机种子
	src := rand.NewSource(time.Now().UnixNano())
	randObj := rand.New(src)
	num1 := strconv.Itoa(randObj.Intn(1000))
	num2 := strconv.Itoa(randObj.Intn(1000))

	//保存密码
	if err := passwordInstance.SavePassword("test"+num1, "test"+num2); err != nil {
		t.Log(err.Error())
		return
	}
	//查看密码
	value, err := passwordInstance.GetPasswordWithKey("test" + num1)
	if err != nil {
		t.Log(err.Error())
		return
	}
	t.Logf("before delete check:%s", value)

	if err := passwordInstance.DeletePassword("test" + num1); err != nil {
		t.Log(err.Error())
		return
	}

	_, err = passwordInstance.GetPasswordWithKey("test" + num1)
	if !assert_instance.Error(err) {
		t.Error("error should not be nil")
	}

}

func init() {
	zaplog.LoggerInit()
}
