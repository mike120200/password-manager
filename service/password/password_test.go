package password_test

import (
	"errors"
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
	assert := assert.New(t)
	os.Remove("./test.gob")
	os.Remove("./data.db")
	os.Remove("./data.backup.db")

	defer func() {
		os.Remove("./test.gob")

		os.Remove("./data.db")
		os.Remove("./data.backup.db")
	}()

	var testCases = []struct {
		test_name string
		key       string
		password  string
		platform  string
	}{
		{"common_test", "test-1", "test-password-1", "google"},
		{"common_test", "test-2", "test-password-2", "edge"},
		{"empty_key_test", "", "test-password-3", "google"},
		{"empty_password_test", "test-4", "", "google"},
		{"empty_platform_test", "test-5", "test-password-5", ""},
	}
	var expectedValue = []struct {
		err      error
		password string
		platform string
	}{
		{nil, "test-password-1", "google"},
		{nil, "test-password-2", "edge"},
		{errors.New("key is empty"), "", ""},
		{errors.New("password is empty"), "", ""},
		{nil, "test-password-5", ""},
	}
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
	for i, testCase := range testCases {
		t.Run(testCase.test_name, func(t *testing.T) {

			//保存密码
			err := passwordInstance.SavePassword(testCase.key, testCase.password, testCase.platform)
			if !assert.Equal(expectedValue[i].err, err) {
				t.Errorf("expect error %v, but got %v", expectedValue[i].err, err)
				return
			}
			if err != nil {
				return
			}
			//查看密码
			value, platform, err := passwordInstance.GetPasswordWithKey(testCase.key)
			if err != nil {
				t.Error(err.Error())
				return
			}
			t.Log(value)
			t.Log(platform)
			if !assert.Equal(expectedValue[i].password, value) {
				t.Errorf("expect password %v, but got %v", expectedValue[i].password, value)
				return
			}
			if !assert.Equal(expectedValue[i].platform, platform) {
				t.Errorf("expect platform %v, but got %v", expectedValue[i].platform, platform)
				return
			}
		})
	}
}

func TestGetAllPassword(t *testing.T) {
	assert := assert.New(t)
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

	values, err := passwordInstance.GetAllPasswords()
	if err != nil {
		t.Log(err.Error())
		return
	}

	if !assert.Equal(0, len(values)) {
		t.Errorf("expect len %v, but got %v", 0, len(values))
		return
	}
	// 设置随机种子
	src := rand.NewSource(time.Now().UnixNano())
	randObj := rand.New(src)
	for i := 0; i < 10; i++ {
		num1 := strconv.Itoa(randObj.Intn(1000))
		num2 := strconv.Itoa(randObj.Intn(1000))

		if err := passwordInstance.SavePassword("test"+num1, "password-test"+num2, "google"); err != nil {
			t.Log(err.Error())
			return
		}
	}
	//获取所有的密码键值对
	values, err = passwordInstance.GetAllPasswords()
	if err != nil {
		t.Log(err.Error())
		return
	}
	if !assert.Equal(10, len(values)) {
		t.Errorf("expect len %v, but got %v", 10, len(values))
		return
	}
	t.Log(values)

}

func TestUpdatePassword(t *testing.T) {
	assert := assert.New(t)
	os.Remove("./test.gob")
	os.Remove("./data.db")
	os.Remove("./data.backup.db")
	defer func() {
		os.Remove("./test.gob")
		os.Remove("./data.db")
		os.Remove("./data.backup.db")
	}()
	var testCases = []struct {
		test_name string
		key       string
		password  string
		platform  string
		newKey    string
	}{
		{"common_test_changeAll", "test_init", "test-password-1", "google", ""},
		{"common_test_changeAll", "test_init", "test-password-2", "edge", ""},
		{"empty_key_test", "", "test-password-3", "google", ""},
		{"empty_password_test", "test_init", "", "google", ""},
		{"empty_platform_test", "test_init", "test-password-5", "", ""},
		{"common_test_changeKey", "test_init", "", "", "test_new"},
	}
	var expectedValue = []struct {
		err      error
		key      string
		password string
		platform string
	}{
		{nil, "test_init", "test-password-1", "google"},
		{nil, "test_init", "test-password-2", "edge"},
		{errors.New("key is empty"), "test_init", "", ""},
		{nil, "test_init", "test-password-2", "google"},
		{nil, "test_init", "test-password-5", "google"},
		{nil, "test_new", "test-password-5", "google"},
	}
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

	//保存密码
	if err := passwordInstance.SavePassword("test_init", "test_init_password", "wechat"); err != nil {
		t.Log(err.Error())
		return
	}
	//查看密码
	value, platform, err := passwordInstance.GetPasswordWithKey("test_init")
	if err != nil {
		t.Log(err.Error())
		return
	}
	t.Logf("before value:%s", value)
	t.Logf("before platform:%s", platform)

	for i, testCase := range testCases {
		t.Run(testCase.test_name, func(t *testing.T) {
			err := passwordInstance.UpdatePassword(testCase.key, testCase.password, testCase.platform, testCase.newKey)
			if !assert.Equal(expectedValue[i].err, err) {
				t.Errorf("expect err %v, but got %v", expectedValue[i].err, err)
				return
			}
			if err != nil {
				return
			}
			value, platform, err = passwordInstance.GetPasswordWithKey(expectedValue[i].key)
			if err != nil {
				t.Log(err.Error())
				return
			}
			t.Logf("after value:%s", value)
			t.Logf("after platform:%s", platform)
			if !assert.Equal(expectedValue[i].password, value) {
				t.Errorf("expect password %v, but got %v", expectedValue[i].password, value)
				return
			}
			if !assert.Equal(expectedValue[i].platform, platform) {
				t.Errorf("expect platform %v, but got %v", expectedValue[i].platform, platform)
			}
		})
	}
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
	if err := passwordInstance.SavePassword("test"+num1, "test"+num2, "google"); err != nil {
		t.Log(err.Error())
		return
	}
	//查看密码
	value, platform, err := passwordInstance.GetPasswordWithKey("test" + num1)
	if err != nil {
		t.Log(err.Error())
		return
	}
	t.Logf("before delete check:%s", value)
	t.Logf("before delete platform:%s", platform)

	if err := passwordInstance.DeletePassword("test" + num1); err != nil {
		t.Log(err.Error())
		return
	}

	err = passwordInstance.DeletePassword("test" + num1)
	if !assert_instance.Error(err) {
		t.Error("error should not be nil")
	}

}

func init() {
	zaplog.LoggerInit()
}
