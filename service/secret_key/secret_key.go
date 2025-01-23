package secretkey

import (
	"crypto/rand"
	"encoding/gob"
	"encoding/hex"
	"errors"
	"os"
	"path/filepath"

	"go.uber.org/zap"
)

const (
	keyFileName = "key.gob"
	keyLength   = 16
)

var _ SecretKeyInterface = (*SecretKey)(nil)

type SecretKey struct {
	logger   *zap.Logger
	filePath string
}

func NewSecretKey() *SecretKey {
	return &SecretKey{
		logger:   zap.L(),
		filePath: "",
	}
}

func NewSecretKeyWithFilePath(filePath string) *SecretKey {
	return &SecretKey{
		logger:   zap.L(),
		filePath: filePath,
	}
}

// GenerateRandomKey 生成指定长度的随机密钥
func (*SecretKey) generateRandomKey(length int) (string, error) {
	// 创建一个字节切片来存储随机数据
	key := make([]byte, length)

	// 从crypto/rand读取随机数据
	_, err := rand.Read(key)
	if err != nil {
		return "", err
	}

	// 将随机字节转换为十六进制字符串
	return hex.EncodeToString(key), nil
}

// SetSecretKey 将密钥保存到文件中
func (srv *SecretKey) SetSecretKey() error {
	key, err := srv.generateRandomKey(keyLength)
	if err != nil {
		srv.logger.Error("Failed to generate random key")
		return err
	}
	var path string
	// 如果指定了文件路径，则使用该路径
	if srv.filePath != "" {
		path = srv.filePath
	} else {

		exePath, err := os.Executable()
		if err != nil {
			srv.logger.Error("Failed to get executable path:", zap.Error(err))
			return err
		}

		path = filepath.Join(filepath.Dir(exePath), keyFileName)
	}
	// srv.logger.Sugar().Debugf("Saving key to file: %s", path)
	file, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0600)
	if err != nil {
		srv.logger.Error("Failed to open file:", zap.Error(err))
		return err
	}
	defer file.Close()
	// 创建gob编码器
	encoder := gob.NewEncoder(file)

	// 将Key结构体编码并写入文件
	err = encoder.Encode(key)
	if err != nil {
		srv.logger.Error("Error encoding key:", zap.Error(err))
		return err
	}
	srv.logger.Info("Key saved successfully")
	return nil
}

// GetSecretKey 从文件中读取密钥
func (srv *SecretKey) GetSecretKey() (string, error) {
	var path string
	// 如果指定了文件路径，则使用该路径
	if srv.filePath != "" {
		path = srv.filePath
	} else {

		exePath, err := os.Executable()
		if err != nil {
			srv.logger.Error("Failed to get executable path:", zap.Error(err))
			return "", err
		}

		path = filepath.Join(filepath.Dir(exePath), keyFileName)
	}
	// 打开gob文件
	file, err := os.Open(path)
	if err != nil {
		srv.logger.Error("Error opening file:", zap.Error(err))
		return "", err
	}
	defer file.Close()
	// 创建gob解码器
	decoder := gob.NewDecoder(file)

	// 创建一个变量来存储解码后的数据
	var key string

	// 解码文件内容到key变量中
	err = decoder.Decode(&key)
	if err != nil {
		srv.logger.Error("Error decoding file:", zap.Error(err))
		return "", err
	}
	if key == "" {
		return "", errors.New("secret key not found")
	}
	srv.logger.Info("Key loaded successfully")
	return key, nil
}
