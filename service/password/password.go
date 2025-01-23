package password

import (
	"errors"
	"password_manager/service/aes"
	dbfilekit "password_manager/service/dbfile_Kit"

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
	Key   string
	Value string
}

// NewPasswordService 创建一个新的 PasswordService 并打开 BoltDB
func NewPasswordService(aesSrv aes.AesInterface, db *bbolt.DB) *PasswordService {
	return &PasswordService{
		db:     db,
		logger: zap.L(),
		aesSrv: aesSrv,
	}
}

// SavePassword 存储密码
func (srv *PasswordService) SavePassword(key, password string) error {
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
		bucket := tx.Bucket([]byte(dbfilekit.BucketName))
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

	// 加密密码
	cipherPassword, nonce, err := srv.aesSrv.Encrypt(password)
	if err != nil {
		srv.logger.Error("encrypt password failed:", zap.Error(err))
		return err
	}
	encryptedValue := append(nonce, cipherPassword...)

	// 将密码存入 BoltDB
	err = srv.db.Update(func(tx *bbolt.Tx) error {
		bucket := tx.Bucket([]byte(dbfilekit.BucketName))
		if bucket == nil {
			srv.logger.Error("bucket not found")
			return errors.New("bucket not found")
		}
		return bucket.Put([]byte(key), []byte(encryptedValue))
	})

	if err != nil {
		srv.logger.Error("save password failed:", zap.Error(err))
		return err
	}
	return nil
}

// GetPasswordWithKey 获取指定 key 的密码
func (srv *PasswordService) GetPasswordWithKey(key string) (string, error) {
	if key == "" {
		srv.logger.Debug("key is empty")
		return "", errors.New("key is empty")
	}

	var encryptedValue []byte

	// 从 BoltDB 获取密码
	err := srv.db.View(func(tx *bbolt.Tx) error {
		bucket := tx.Bucket([]byte(dbfilekit.BucketName))
		if bucket == nil {
			return errors.New("bucket not found")
		}
		encryptedValue = bucket.Get([]byte(key))
		if encryptedValue == nil {
			return errors.New("key:" + key + " not found")
		}
		return nil
	})
	if err != nil {
		// srv.logger.Error("get password failed:", zap.Error(err))
		return "", err
	}

	// 解密密码
	password, err := srv.decryptValue(encryptedValue)
	if err != nil {
		srv.logger.Error("decrypt password failed:", zap.Error(err))
		return "", err
	}

	return string(password), nil
}

// GetAllPasswords 获取所有存储的密码
func (srv *PasswordService) GetAllPasswords() (map[string]string, error) {
	passwords := make(map[string]string)

	err := srv.db.View(func(tx *bbolt.Tx) error {
		bucket := tx.Bucket([]byte(dbfilekit.BucketName))
		if bucket == nil {
			return errors.New("bucket not found")
		}
		return bucket.ForEach(func(k, v []byte) error {
			password, err := srv.decryptValue(v)
			if err != nil {
				srv.logger.Error("decrypt password failed:", zap.Error(err))
				return err
			}
			passwords[string(k)] = string(password)
			return nil
		})
	})
	if err != nil {
		srv.logger.Error("get all passwords failed:", zap.Error(err))
		return nil, err
	}

	return passwords, nil
}

// UpdatePassword 更新密码
func (srv *PasswordService) UpdatePassword(key, newPassword string) error {
	if key == "" {
		srv.logger.Error("key is empty")
		return errors.New("key is empty")
	}
	if newPassword == "" {
		srv.logger.Error("newPassword is empty")
		return errors.New("newPassword is empty")
	}

	// 先检查 key 是否存在
	var exists bool
	var oldPassword string
	err := srv.db.View(func(tx *bbolt.Tx) error {
		bucket := tx.Bucket([]byte(dbfilekit.BucketName))
		if bucket == nil {
			return errors.New("bucket not found")
		}
		value := bucket.Get([]byte(key))
		exists = value != nil
		var err error
		value, err = srv.decryptValue(value)
		if err != nil {
			return err
		}
		oldPassword = string(value)
		return nil
	})
	if err != nil || !exists {
		return errors.New("key:" + key + " not found")
	}

	// 加密新密码
	cipherPassword, nonce, err := srv.aesSrv.Encrypt(newPassword)
	if err != nil {
		srv.logger.Error("encrypt new password failed", zap.Error(err))
		return err
	}
	encryptedValue := append(nonce, cipherPassword...)

	// 更新 BoltDB 中的密码
	err = srv.db.Update(func(tx *bbolt.Tx) error {
		bucket := tx.Bucket([]byte(dbfilekit.BucketName))
		if bucket == nil {
			return errors.New("bucket not found")
		}
		return bucket.Put([]byte(key), []byte(encryptedValue))
	})

	if err != nil {
		srv.logger.Error("update password failed:", zap.Error(err))
		return err
	}

	color.Green.Println("password updated successfully:" + oldPassword + " -> " + newPassword)
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
		bucket := tx.Bucket([]byte(dbfilekit.BucketName))
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
		bucket := tx.Bucket([]byte(dbfilekit.BucketName))
		if bucket == nil {
			return errors.New("bucket not found")
		}
		err := bucket.Delete([]byte(key))
		if err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		srv.logger.Error("delete password failed:", zap.Error(err))
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
