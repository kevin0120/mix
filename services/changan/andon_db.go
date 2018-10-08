package changan

import (
	"fmt"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mssql"
)

type AndonDB struct {
	cfg *ConfigAndonDB
	eng *gorm.DB
}

func (adb *AndonDB) StartService() error {
	sConn := fmt.Sprintf("sqlserver://%s:%s@%s?database=%s",
		adb.cfg.User,
		adb.cfg.Password,
		adb.cfg.Url,
		adb.cfg.DBName)

	_db, err := gorm.Open("mssql", sConn)

	_db.AutoMigrate(&TighteningResults{})

	adb.eng = _db

	return err
}

func (adb *AndonDB) StopService() error {
	adb.eng.Close()

	return nil
}

func (adb *AndonDB) InsertResult(result *TighteningResults) bool {
	adb.eng.Create(result)

	return true
}
