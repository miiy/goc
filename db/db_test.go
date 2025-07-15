package db

import "testing"

func TestNewDB(t *testing.T) {
	db, err := NewDB(Config{
		Driver:   "mysql",
		Host:     "127.0.0.1",
		Port:     "3306",
		Username: "root",
		Password: "123456",
		Database: "goctest",
	})
	if err != nil {
		t.Fatal(err)
	}

	ret := make(map[string]interface{})
	err = db.Gorm().Raw("select * from users limit 1").Scan(&ret).Error
	if err != nil {
		t.Error(err)
	}
	t.Log(ret)
}
