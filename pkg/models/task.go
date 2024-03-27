// Copyright (c) 2024. LubyRuffy. All rights reserved.

package models

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Task struct {
	gorm.Model
	UUID        string `json:"taskid" gorm:"type:varchar(36);not null;uniqueIndex"`
	Title       string // 任务标题
	Description string // 任务描述
	QueryA      string // 参赛者A提交的查询语句
	QueryB      string // 参赛者B提交的查询语句
	AScore      int    // 参赛者A的得分
	BScore      int    // 参赛者B的得分
}

func NewTask(q1, q2 string) *Task {
	uid, err := uuid.NewUUID()
	if err != nil {
		panic(err)
	}
	return &Task{
		UUID:   uid.String(),
		QueryA: q1,
		QueryB: q2,
	}
}
