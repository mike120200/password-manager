package password

import (
	"errors"
	"password_manager/service/aes"
	dbfilekit "password_manager/service/dbfile_Kit"
	"strconv"

	"github.com/gookit/color"
	"go.etcd.io/bbolt"
	"go.uber.org/zap"
)

const nonceHexLen = 12

type PasswordService struct {
	db     *bbolt.DB
	logger *zap.Logger
	aesSrv aes.AesInterface
}
type PasswordData struct {
	Key      string `json:"key"`
	Platform string `json:"platform"`
	Password string `json:"password"`
}

func NewPasswordData(key, platform, password string) *PasswordData {
	return &PasswordData{
		Key:      key,
		Password: password,
		Platform: platform}
}

// NewPasswordService 创建一个新的 PasswordService 并打开 BoltDB
// NewPasswordService 创建一个新的PasswordService实例
func NewPasswordService(aesSrv aes.AesInterface, db *bbolt.DB) *PasswordService {
	// 返回一个PasswordService实例，包含db、logger和aesSrv
	return &PasswordService{
		db:     db,
		logger: zap.L(),
		aesSrv: aesSrv,
	}
}

// SavePassword 存储密码
func (srv *PasswordService) SavePassword(key, password, platform string) error {
	if key == "" {
		srv.logger.Error("key is empty")
		return errors.New("key is empty")
	}
	if password == "" {
		srv.logger.Error("password is empty")
		return errors.New("password is empty")
	}
	// 先检查 key 是否存在
	var exists bool
	err := srv.db.View(func(tx *bbolt.Tx) error {
		bucket := tx.Bucket([]byte(dbfilekit.PasswordBucketName))
		if bucket == nil {
			return errors.New("bucket not found")
		}
		exists = bucket.Get([]byte(key)) != nil
		return nil
	})
	if err != nil {
		srv.logger.Error("check key failed:", zap.Error(err))
		return err
	}
	//存在就返回错误
	if exists {
		srv.logger.Debug("save password failed, key has been set ", zap.String("key", key))
		return errors.New("key has been set")
	}
	//第一次存因此newKey参赛可以为空
	err = srv.updateDb(key, password, platform, "")
	if err != nil {
		srv.logger.Error("save password failed:", zap.Error(err))
		return err
	}
	return nil
}

// GetPasswordWithKey 获取指定 key 的密码
func (srv *PasswordService) GetPasswordWithKey(key string) (string, string, error) {
	if key == "" {
		srv.logger.Debug("key is empty")
		return "", "", errors.New("key is empty")
	}

	var encryptedValue []byte
	var platform []byte
	var platformLen int
	// 从 BoltDB 获取密码
	err := srv.db.View(func(tx *bbolt.Tx) error {
		bucket := tx.Bucket([]byte(dbfilekit.PasswordBucketName))
		if bucket == nil {
			srv.logger.Error("password bucket not found")
			return errors.New("password bucket not found")
		}
		encryptedValue = bucket.Get([]byte(key))
		if encryptedValue == nil {
			return errors.New("key:" + key + " not found")
		}
		bucket = tx.Bucket([]byte(dbfilekit.PlatformLenBucketName))
		if bucket == nil {
			srv.logger.Error("platform bucket not found")
			return errors.New("platform bucket not found")
		}
		//获取平台信息长度
		platformLenByte := bucket.Get([]byte(key))
		if platformLenByte == nil {
			platformLen = 0
		} else {
			var err error
			platformLen, err = strconv.Atoi(string(platformLenByte))
			if err != nil {
				srv.logger.Error("convert platformLen failed:", zap.Error(err))
				return err
			}
		}

		return nil
	})
	if err != nil {
		// srv.logger.Error("get password failed:", zap.Error(err))
		return "", "", err
	}

	if platformLen > 0 {
		platform = encryptedValue[:platformLen]
	}
	// 解密密码
	password, err := srv.decryptValue(encryptedValue[platformLen:])
	if err != nil {
		srv.logger.Error("decrypt password failed:", zap.Error(err))
		return "", "", err
	}

	return string(password), string(platform), nil
}

// GetAllPasswords 获取所有存储的密码
func (srv *PasswordService) GetAllPasswords() (map[string]PasswordData, error) {
	passwordsData := make(map[string]PasswordData)

	err := srv.db.View(func(tx *bbolt.Tx) error {
		//获取密码
		passwordBucket := tx.Bucket([]byte(dbfilekit.PasswordBucketName))
		if passwordBucket == nil {
			srv.logger.Error("password bucket not found")
			return errors.New("password bucket not found")
		}
		platformBucket := tx.Bucket([]byte(dbfilekit.PlatformLenBucketName))
		if platformBucket == nil {
			srv.logger.Error("platform bucket not found")
			return errors.New("platform bucket not found")
		}
		err := passwordBucket.ForEach(func(k, v []byte) error {
			//截取平台信息
			platformLenByte := platformBucket.Get(k)
			if platformLenByte == nil {
				password, err := srv.decryptValue(v)
				if err != nil {
					srv.logger.Error("decrypt password failed:", zap.Error(err))
					return err
				}
				passwordsData[string(k)] = PasswordData{
					Key:      string(k),
					Password: string(password),
				}
			} else {
				platformLen, err := strconv.Atoi(string(platformLenByte))
				if err != nil {
					srv.logger.Error("convert platformLen failed:", zap.Error(err))
					return err
				}
				platform := v[:platformLen]
				// 解密密码
				password, err := srv.decryptValue(v[platformLen:])
				if err != nil {
					srv.logger.Error("decrypt password failed:", zap.Error(err))
					return err
				}
				passwordsData[string(k)] = PasswordData{
					Key:      string(k),
					Password: string(password),
					Platform: string(platform),
				}
			}
			return nil
		})
		if err != nil {
			srv.logger.Error("get all passwords failed:", zap.Error(err))
			return err
		}
		return nil
	})
	if err != nil {
		srv.logger.Error("get all passwords failed:", zap.Error(err))
		return nil, err
	}

	return passwordsData, nil
}

// UpdatePassword 更新密码
func (srv *PasswordService) UpdatePassword(key, newPassword, newPlatform, newKey string) error {
	var (
		password    string
		platform    string
		oldPassword string
		oldPlatform string
		err         error
	)
	if key == "" {
		srv.logger.Error("key is empty")
		return errors.New("key is empty")
	}
	if newPassword != "" {
		password = newPassword
	}
	if newPlatform != "" {
		platform = newPlatform
	}

	// 先检查 key 是否存在
	var exists bool
	var oldValue []byte
	err = srv.db.View(func(tx *bbolt.Tx) error {
		bucket := tx.Bucket([]byte(dbfilekit.PasswordBucketName))
		if bucket == nil {
			return errors.New("password bucket not found")
		}
		value := bucket.Get([]byte(key))
		exists = value != nil
		oldValue = value
		//获取平台信息长度
		platformLen, err := srv.getPlatformLen(key, tx)
		if err != nil {
			srv.logger.Error("get key "+key+" platform len failed:", zap.Error(err))
			return err
		}
		//获取旧的密码和平台信息
		oldPasswordByte, err := srv.decryptValue(oldValue[platformLen:])
		if err != nil {
			srv.logger.Error("decrypt password failed:", zap.Error(err))
			return err
		}
		oldPassword = string(oldPasswordByte)
		srv.logger.Sugar().Debugf("oldPassword:" + oldPassword)
		oldPlatform = string(oldValue[:platformLen])
		//判断是否需要更新平台信息
		if platform == "" {
			platform = oldPlatform
		}
		//判断是否需要更新密码
		if password == "" {
			password = oldPassword
		}
		srv.logger.Sugar().Debugf("key:" + key + " password:" + password + " platform:" + platform)
		return nil
	})
	if err != nil {
		srv.logger.Error("get key "+key+" failed:", zap.Error(err))
		return errors.New("key:" + key + " not found")
	}
	if !exists {
		srv.logger.Sugar().Debugf("key:" + key + " not found")
		return errors.New("key:" + key + " not found")
	}

	srv.logger.Sugar().Debugf("password:" + password)
	err = srv.updateDb(key, password, platform, newKey)
	if err != nil {
		srv.logger.Error("update db failed:", zap.Error(err))
		return err
	}
	if newKey != "" {
		color.Green.Println("key updated successfully:" + key + " -> " + newKey)
	}
	if newPassword != "" {
		color.Green.Println("password updated successfully:" + oldPassword + " -> " + newPassword)
	}
	if newPlatform != "" {
		color.Green.Println("platform updated successfully:" + oldPlatform + " -> " + newPlatform)
	}

	return nil
}

// getPlatformLen 获取平台信息的长度
func (srv *PasswordService) getPlatformLen(key string, tx *bbolt.Tx) (int, error) {
	if key == "" {
		srv.logger.Error("key is empty")
		return 0, errors.New("key is empty")
	}
	bucket := tx.Bucket([]byte(dbfilekit.PlatformLenBucketName))
	if bucket == nil {
		srv.logger.Error("platform bucket not found")
		return 0, errors.New("platform bucket not found")
	}
	platformLenByte := bucket.Get([]byte(key))
	if platformLenByte == nil {
		srv.logger.Error("key:" + key + " not found")
		return 0, errors.New("key:" + key + " not found")
	}
	platformLen, err := strconv.Atoi(string(platformLenByte))
	if err != nil {
		srv.logger.Error("convert platformLen failed:", zap.Error(err))
		return 0, err
	}
	return platformLen, nil
}

// updateDb 更新数据库
func (srv *PasswordService) updateDb(key string, password string, platform string, newKey string) error {
	if key == "" {
		srv.logger.Error("key is empty")
		return errors.New("key is empty")
	}
	if password == "" {
		srv.logger.Error("password is empty")
		return errors.New("password is empty")
	}
	// 将平台信息转化为byte数组
	platformByte := []byte(platform)
	// 获取数组长度
	platformLen := len(platformByte)
	//将数组长度转化为字符串
	platformLenStr := strconv.Itoa(platformLen)
	// 将平台信息长度字符串转化为byte数组
	platformLenByte := []byte(platformLenStr)

	// 加密密码
	cipherPassword, nonce, err := srv.aesSrv.Encrypt(password)
	if err != nil {
		srv.logger.Error("encrypt password failed:", zap.Error(err))
		return err
	}
	// 拼接平台信息和nonce
	encryptedValue := append(platformByte, nonce...)
	encryptedValue = append(encryptedValue, cipherPassword...)

	// 将密码存入 BoltDB
	err = srv.db.Update(func(tx *bbolt.Tx) error {
		if newKey == "" {
			//没有修改key就直接更新
			err := srv.updateWithTx(key, encryptedValue, platformLenByte, tx)
			if err != nil {
				srv.logger.Error("updateWithTx failed:", zap.Error(err))
				return err
			}
		} else {
			//更新
			err := srv.updateWithTx(newKey, encryptedValue, platformLenByte, tx)
			if err != nil {
				srv.logger.Error("updateWithTx failed:", zap.Error(err))
				return err
			}
			//删除之前的
			err = srv.deleteWithTx(key, tx)
			if err != nil {
				srv.logger.Error("deleteWithTx failed:", zap.Error(err))
				return err
			}
		}

		return nil
	})

	if err != nil {
		srv.logger.Error("save password failed:", zap.Error(err))
		return err
	}
	return nil
}

// DeletePassword 删除密码
func (srv *PasswordService) DeletePassword(key string) error {
	if key == "" {
		srv.logger.Error("key is empty")
		return errors.New("key is empty")
	}
	// 先检查 key 是否存在
	var exists bool
	err := srv.db.View(func(tx *bbolt.Tx) error {
		bucket := tx.Bucket([]byte(dbfilekit.PasswordBucketName))
		if bucket == nil {
			return errors.New("bucket not found")
		}
		exists = bucket.Get([]byte(key)) != nil
		return nil
	})
	if err != nil || !exists {
		return errors.New("key:" + key + " not found")
	}
	//执行删除操作
	err = srv.db.Update(func(tx *bbolt.Tx) error {
		return srv.deleteWithTx(key, tx)
	})
	if err != nil {
		srv.logger.Error("delete password failed:", zap.Error(err))
		return err
	}
	return nil
}

// deleteWithTx在数据库中执行操作
func (srv *PasswordService) deleteWithTx(key string, tx *bbolt.Tx) error {
	srv.logger.Info("deleteInDb key:", zap.String("key", key))
	bucket := tx.Bucket([]byte(dbfilekit.PasswordBucketName))
	if bucket == nil {
		srv.logger.Error("password bucket not found")
		return errors.New("password bucket not found")
	}
	err := bucket.Delete([]byte(key))
	if err != nil {
		return err
	}
	bucket = tx.Bucket([]byte(dbfilekit.PlatformLenBucketName))
	if bucket == nil {
		srv.logger.Error("platformLen bucket not found")
		return errors.New("platform bucket not found")
	}
	err = bucket.Delete([]byte(key))
	if err != nil {
		return err
	}
	return nil
}

// updateWithTx 使用tx操作
func (srv *PasswordService) updateWithTx(key string, encryptedValue []byte, platformLenByte []byte, tx *bbolt.Tx) error {
	srv.logger.Info("updateWithTx key:", zap.String("key", key))
	// 存入密码
	bucket := tx.Bucket([]byte(dbfilekit.PasswordBucketName))
	if bucket == nil {
		srv.logger.Error("password bucket not found")
		return errors.New("password bucket not found")
	}
	if err := bucket.Put([]byte(key), encryptedValue); err != nil {
		srv.logger.Error("save password failed:", zap.Error(err))
		return err
	}
	// 存入平台长度
	bucket = tx.Bucket([]byte(dbfilekit.PlatformLenBucketName))
	if bucket == nil {
		srv.logger.Error("platformLen bucket not found")
		return errors.New("platformLen bucket not found")
	}
	if err := bucket.Put([]byte(key), platformLenByte); err != nil {
		srv.logger.Error("save platformLen failed:", zap.Error(err))
		return err
	}
	return nil
}

// decryptValue 解密存储的密码
func (srv *PasswordService) decryptValue(encryptedValue []byte) ([]byte, error) {
	// 分割值获取 nonce
	nonce := encryptedValue[:nonceHexLen]
	cipherData := encryptedValue[nonceHexLen:]

	// 解密密码
	password, err := srv.aesSrv.Decrypt(cipherData, nonce)
	if err != nil {
		srv.logger.Error("decrypt password failed:", zap.Error(err))
		return nil, err
	}
	return password, nil
}
