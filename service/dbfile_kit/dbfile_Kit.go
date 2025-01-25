package dbfilekit

import (
	"bufio"
	"crypto/sha256"
	"errors"
	"fmt"
	"io"
	"os"
	secretkey "password_manager/service/secret_key"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/gookit/color"
	"go.etcd.io/bbolt"
	"go.uber.org/zap"
)

const (
	PasswordBucketName    = "passwords"
	PlatformLenBucketName = "platformsLen"
	FileDBName            = "data.db"
	BackupDBName          = "data.backup.db"
	timeStampKey          = "time-stamp"
	TimeStampBucket       = "time_stamp_bucket"
	BackupIntervalSeconds = 500
)

// 验证DBFileKit接口是否实现
var _ DBFileKit = (*DBKitImpl)(nil)

type DBKitImpl struct {
	logger       *zap.Logger
	db           *bbolt.DB
	dirPath      string
	secretKeySrv secretkey.SecretKeyInterface
}

// NewDBKit 初始化数据库管理器
func NewDBKit(secretKeySrv secretkey.SecretKeyInterface) *DBKitImpl {
	return &DBKitImpl{
		logger:       zap.L(),
		secretKeySrv: secretKeySrv,
	}
}

func NewDBKitWithFilePath(dirPath string, secretKeySrv secretkey.SecretKeyInterface) *DBKitImpl {
	return &DBKitImpl{
		logger:       zap.L(),
		dirPath:      dirPath,
		secretKeySrv: secretKeySrv,
	}
}

// Init初始化数据库
func (srv *DBKitImpl) Init() error {
	var dir string
	if srv.dirPath != "" {
		dir = srv.dirPath
	} else {
		// 获取当前可执行文件的路径
		exePath, err := os.Executable()
		if err != nil {
			srv.logger.Error("can't get executable path", zap.Error(err))
			return err
		}
		dir = filepath.Dir(exePath)
	}

	dbFile := filepath.Join(dir, FileDBName)
	// srv.logger.Sugar().Debugf("dbFile: %s", dbFile)
	// 检查数据库文件是否存在
	if srv.isDbFileExist(dbFile) {
		db, err := srv.initDB(dbFile)
		if err != nil {
			srv.logger.Error("init db fail:", zap.Error(err))
			return err
		}
		// 设置数据库实例
		srv.db = db
		defer func() {
			timeStampStr := strconv.FormatInt(time.Now().Unix(), 10)
			if err := srv.saveTimeStamp(timeStampStr); err != nil {
				fmt.Println("save time stamp fail:" + err.Error())
			}
		}()
		//开始进行备份流程
		//检查时间戳
		result, err := srv.checkTimeStamp()
		if err != nil {
			srv.logger.Error("check time stamp fail:", zap.Error(err))
			return err
		}
		//时间超过了备份间隔，则进行备份
		if result {
			if err := srv.BackupDB(); err != nil {
				return err
			}
		}
		return nil
	} else {
		backupFile := filepath.Join(dir, BackupDBName)
		if srv.isDbFileExist(backupFile) {
			color.Yellow.Printf("db file not found.")
			//开始恢复流程
			result := srv.restoreDataProcess(backupFile, dbFile)
			if !result {
				return errors.New("restore db fail")
			}
			db, err := srv.initDB(dbFile)
			if err != nil {
				srv.logger.Error("init db fail:", zap.Error(err))
				return err
			}
			defer func() {
				timeStampStr := strconv.FormatInt(time.Now().Unix(), 10)
				if err := srv.saveTimeStamp(timeStampStr); err != nil {
					fmt.Println("save time stamp fail:" + err.Error())
					return
				}
			}()
			// 设置数据库实例
			srv.db = db
			return nil
		} else {
			//主数据库文件和备份数据库文件都不存在的情况
			//创建数据库文件
			file, err := createFile(dbFile)
			if err != nil {
				srv.logger.Error("create db file fail:", zap.Error(err))
				return err
			}
			file.Close()
			//创建备份数据库文件
			file, err = createFile(backupFile)
			if err != nil {
				srv.logger.Error("create backup db file fail:", zap.Error(err))
				return err
			}
			file.Close()
			//重新设置密钥
			if err := srv.secretKeySrv.SetSecretKey(); err != nil {
				return err
			}

			fmt.Println("init db file and backup file success")
			//初始化db
			db, err := srv.initDB(dbFile)
			if err != nil {
				srv.logger.Error("init db fail:", zap.Error(err))
				return err
			}
			// 设置数据库实例
			srv.db = db
			//保存时间戳
			timeStampStr := strconv.FormatInt(time.Now().Unix(), 10)
			if err := srv.saveTimeStamp(timeStampStr); err != nil {
				srv.logger.Error("save time stamp fail:", zap.Error(err))
				return err
			}
			return nil

		}
	}

}
func (srv *DBKitImpl) InitFromBackupFile() error {
	var dir string
	if srv.dirPath != "" {
		dir = srv.dirPath
	} else {
		// 获取当前可执行文件的路径
		exePath, err := os.Executable()
		if err != nil {
			srv.logger.Error("can't get executable path", zap.Error(err))
			return err
		}
		dir = filepath.Dir(exePath)
	}
	backupFile := filepath.Join(dir, BackupDBName)
	// 打开或创建数据库文件
	db, err := bbolt.Open(backupFile, 0600, nil)
	if err != nil {
		return err
	}
	// 确保数据库打开后创建bucket
	err = db.Update(func(tx *bbolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists([]byte(PasswordBucketName))
		return err
	})
	if err != nil {
		db.Close()
		srv.logger.Error("create password bucket fail:", zap.Error(err))
		return err
	}
	// 确保数据库打开后创建bucket
	err = db.Update(func(tx *bbolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists([]byte(PlatformLenBucketName))
		return err
	})
	if err != nil {
		db.Close()
		srv.logger.Error("create platformLen bucket fail:", zap.Error(err))
		return err
	}
	srv.db = db
	return nil

}

// Close 关闭数据库
func (srv *DBKitImpl) Close() {
	srv.db.Close()
}

// IsDbFileExist 检查数据库文件是否存在
func (srv *DBKitImpl) isDbFileExist(dbPath string) bool {
	_, err := os.Stat(dbPath)
	return !os.IsNotExist(err)
}

// initDB 初始化数据库实例
func (srv *DBKitImpl) initDB(dbFile string) (*bbolt.DB, error) {
	// 打开或创建数据库文件
	db, err := bbolt.Open(dbFile, 0664, nil)
	if err != nil {
		return nil, err
	}
	// 确保数据库打开后创建bucket
	err = db.Update(func(tx *bbolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists([]byte(PasswordBucketName))
		return err
	})
	if err != nil {
		db.Close()
		srv.logger.Error("create password bucket fail:", zap.Error(err))
		return nil, err
	}
	// 确保数据库打开后创建bucket
	err = db.Update(func(tx *bbolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists([]byte(PlatformLenBucketName))
		return err
	})
	if err != nil {
		db.Close()
		srv.logger.Error("create platformLen bucket fail:", zap.Error(err))
		return nil, err
	}
	return db, nil
}

// RestoreDB 恢复数据库
func (srv *DBKitImpl) RestoreDB() error {
	var dir string
	if srv.dirPath != "" {
		dir = srv.dirPath
	} else {
		// 获取当前可执行文件的路径
		exePath, err := os.Executable()
		if err != nil {
			srv.logger.Error("can't get executable path", zap.Error(err))
			return err
		}
		dir = filepath.Dir(exePath)
	}

	dbFile := filepath.Join(dir, FileDBName)
	backupFile := filepath.Join(dir, BackupDBName)
	result := srv.restoreDataProcess(backupFile, dbFile)
	if !result {
		return errors.New("restore db fail")
	}
	return nil
}

// restoreDataProcess 恢复的流程
func (srv *DBKitImpl) restoreDataProcess(backupPath, mainDBPath string) bool {
	reader := bufio.NewReader(os.Stdin)
	color.Yellow.Println("Do you want to restore the data? (y/n): ")
	for {
		input, _ := reader.ReadString('\n')               // 读取输入
		input = strings.TrimSpace(strings.ToLower(input)) // 去除换行符并转换为小写

		if input == "y" {
			fmt.Println("Restoring data...")
			err := srv.restoreDB(backupPath, mainDBPath)
			if err != nil {
				srv.logger.Error("restore db fail:", zap.Error(err))
				return false
			}
			color.Green.Println("Data restored successfully.")
			return true
		} else if input == "n" {
			fmt.Println("Data will not be restored.")
			return false
		} else {
			fmt.Println("Invalid input. Please enter 'y' or 'n'.")
		}
	}
}

// restoreDB 恢复数据
func (srv *DBKitImpl) restoreDB(backupPath, mainDBPath string) error {
	// 确保备份文件存在
	if _, err := os.Stat(backupPath); os.IsNotExist(err) {
		srv.logger.Error("backup file does not exist")
		return errors.New("backup file does not exist")
	}
	//删除旧的数据库
	if err := os.Remove(mainDBPath); err != nil && !os.IsNotExist(err) {
		srv.logger.Error("failed to remove main db file:" + err.Error())
		return err
	}
	file, err := createFile(mainDBPath)
	if err != nil {
		srv.logger.Error("failed to create new db file:" + err.Error())
		return err
	}
	file.Close()
	// 复制备份到主数据库
	err = copyFile(backupPath, mainDBPath)
	if err != nil {
		srv.logger.Error("failed to copy backup file:" + err.Error())
		return err
	}
	return nil
}

// copyFile 复制文件
func copyFile(src, dst string) error {
	in, err := os.Open(src)
	if err != nil {
		return err
	}
	defer in.Close()

	out, err := createFile(dst)
	if err != nil {
		return err
	}
	defer out.Close()

	_, err = io.Copy(out, in)
	if err != nil {
		return err
	}
	return nil
}

// dDB获取数据库
func (srv *DBKitImpl) GetDB() (*bbolt.DB, error) {
	if srv.db == nil {
		srv.logger.Error("db is nil")
		return nil, errors.New("db is nil")
	}
	return srv.db, nil
}

// BackupDB 备份数据
func (srv *DBKitImpl) BackupDB() error {
	srv.logger.Debug("BackupDB begin")
	var dir string
	if srv.dirPath != "" {
		dir = srv.dirPath
	} else {
		// 获取当前可执行文件的路径
		exePath, err := os.Executable()
		if err != nil {
			srv.logger.Error("can't get executable path", zap.Error(err))
			return err
		}
		dir = filepath.Dir(exePath)
	}

	// 数据库文件路径
	dbFile := filepath.Join(dir, FileDBName)
	backupFile := filepath.Join(dir, BackupDBName)

	//删除旧的备份文件
	if err := os.Remove(backupFile); err != nil && !os.IsNotExist(err) {
		srv.logger.Error("failed to remove backup file", zap.Error(err))
		return err
	}
	//重新创建
	file, err := createFile(backupFile)
	if err != nil {
		srv.logger.Error("failed to create backup file:" + err.Error())
		return err
	}
	file.Close()
	// 复制数据库文件到备份文件
	err = copyFile(dbFile, backupFile)
	if err != nil {
		srv.logger.Error("failed to copy db file:" + err.Error())
		return err
	}
	return nil
}

// createFile 创建文件,并给予0600权限，只让当前用户对文件进行读写操作
func createFile(name string) (*os.File, error) {
	return os.OpenFile(name, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0600)
}

// saveTimeStamp 保存时间戳
func (srv *DBKitImpl) saveTimeStamp(timeStampStr string) error {
	data := []byte(timeStampKey)
	hash := sha256.Sum256(data)
	err := srv.db.Update(func(tx *bbolt.Tx) error {
		bucket, err := tx.CreateBucketIfNotExists([]byte(TimeStampBucket))
		if err != nil {
			return err
		}

		err = bucket.Put(hash[:], []byte(timeStampStr))
		if err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		srv.logger.Error("failed to save timestamp", zap.Error(err))
		return err
	}
	srv.logger.Debug("timestamp saved successfully")
	return nil
}

// checkTimeStamp 检查时间戳
func (srv *DBKitImpl) checkTimeStamp() (bool, error) {
	data := []byte(timeStampKey)
	hash := sha256.Sum256(data)
	var lastTimestamp int64
	timeStamp := time.Now().Unix()
	err := srv.db.View(func(tx *bbolt.Tx) error {
		bucket := tx.Bucket([]byte(TimeStampBucket))
		if bucket == nil {
			return errors.New("bucket not found")
		}
		var err error
		value := bucket.Get(hash[:])
		srv.logger.Sugar().Debugf("value: %s", string(value))
		lastTimestamp, err = strconv.ParseInt(string(value), 10, 64)
		if err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		srv.logger.Error("failed to check timestamp", zap.Error(err))
		return false, err
	}
	srv.logger.Debug("timestamp checked successfully")
	result := timeStamp >= lastTimestamp+BackupIntervalSeconds
	if result {
		// 将时间戳转换为 time.Time
		t := time.Unix(timeStamp, 0)

		// 格式化时间为字符串
		formattedTime := t.Format("2006-01-02 15:04:05")
		color.Gray.Println("Last backup time:", formattedTime)
	} else {
		t := time.Unix(lastTimestamp, 0)
		formattedTime := t.Format("2006-01-02 15:04:05")
		color.Gray.Println("Last backup time:", formattedTime)
	}
	return result, nil
}
