package aes

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"errors"
	"io"

	"go.uber.org/zap"
)

var _ AesInterface = (*AesService)(nil)

var AesSrv *AesService

type AesService struct {
	block  cipher.Block
	logger *zap.Logger
}

func NewAesService(secretKey string) *AesService {
	block, err := aes.NewCipher([]byte(secretKey))
	if err != nil {
		panic(err)
	}
	return &AesService{
		block:  block,
		logger: zap.L(),
	}
}

// Encrypt加密
func (srv *AesService) Encrypt(plainText string) ([]byte, []byte, error) {
	if plainText == "" {
		srv.logger.Error("plainText is empty")
		return nil, nil, errors.New("plainText is empty")
	}
	gcm, err := cipher.NewGCM(srv.block)
	if err != nil {
		return nil, nil, err
	}
	//创建随机的nonce
	nonce := make([]byte, gcm.NonceSize())
	if _, err = io.ReadFull(rand.Reader, nonce); err != nil {
		return nil, nil, err
	}
	// 加密并附加认证标签
	cipherData := gcm.Seal(nil, nonce, []byte(plainText), nil)
	return cipherData, nonce, nil
}

// Decrypt解密
func (srv *AesService) Decrypt(cipherData, nonce []byte) ([]byte, error) {
	//参数检测
	if cipherData == nil {
		srv.logger.Error("cipherData is nil")
		return nil, errors.New("cipherData is nil")
	}
	if nonce == nil {
		srv.logger.Error("nonceHex is nil")
		return nil, errors.New("nonceHex is nil")
	}
	gcm, err := cipher.NewGCM(srv.block)
	if err != nil {
		return nil, err
	}

	// 解密并验证
	plainData, err := gcm.Open(nil, nonce, cipherData, nil)
	if err != nil {
		return nil, err
	}
	return plainData, nil
}
