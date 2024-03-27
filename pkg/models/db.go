// Copyright (c) 2024. LubyRuffy. All rights reserved.

package models

import (
	"github.com/glebarez/sqlite"
	"gorm.io/gorm"
)

var (
	db *gorm.DB
)

func Get() *gorm.DB {
	if db != nil {
		return db
	}

	var err error
	db, err = gorm.Open(sqlite.Open("fofapk.db"), &gorm.Config{})
	if err != nil {
		panic(err)
	}
	err = db.AutoMigrate(&Result{}, &Task{})
	if err != nil {
		panic(err)
	}
	return db
}
