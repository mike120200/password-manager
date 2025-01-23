package dbfilekit_test

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

func TestInit(t *testing.T) {
	secretKeyInstance := secretkey.NewSecretKeyWithFilePath("./test.gob")
	instance := dbfilekit.NewDBKitWithFilePath("./", secretKeyInstance)
	if err := instance.Init(); err != nil {
		t.Log(err.Error())
		return
	}
}

func TestBackup(t *testing.T) {
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
	t.Logf("main db:%s", values)

	//备份
	if err := dbfileKitInstance.BackupDB(); err != nil {
		t.Log(err.Error())
		return
	}
	//获取备份的数据库实例
	if err := dbfileKitInstance.InitFromBackupFile(); err != nil {
		t.Log(err.Error())
		return
	}
	backupDB, err := dbfileKitInstance.GetDB()
	if err != nil {
		t.Log(err.Error())
		return
	}
	defer backupDB.Close()

	//获取备份的密码键值对
	passwordBackupInstance := password.NewPasswordService(aesInstance, backupDB)
	backupValues, err := passwordBackupInstance.GetAllPasswords()
	if err != nil {
		t.Log(err.Error())
		return
	}
	t.Logf("backup db:%s", backupValues)
	if !assert.Equal(backupValues, values) {
		t.Error("backup db is not equal to main db")
		return
	}
}

func init() {
	zaplog.LoggerInit()
}
