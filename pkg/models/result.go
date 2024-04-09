// Copyright (c) 2024. LubyRuffy. All rights reserved.

package models

import "gorm.io/gorm"

type From uint8

const (
	FromBoth From = iota
	FromA
	FromB
)

type Result struct {
	gorm.Model     `json:"-"`
	UUID           string `json:"taskid" gorm:"type:varchar(36);not null;index"` // 对应task
	Host           string `json:"host" gorm:"index"`
	IP             string `json:"ip" gorm:"index"`
	Port           string `json:"port" `
	ASOrganization string `json:"as_organization" `
	Protocol       string `json:"protocol" `
	Domain         string `json:"domain" `
	Title          string `json:"title" `
	FID            string `json:"fid" `
	CertsSubjectCN string `json:"certs_subject_cn" `
	From           From   `json:"from" `                           // 0 表示相同，1 表示A多，2表示B多
	Score          int    `json:"score" gorm:"not null;default:0"` // 加减分
	Done           bool   // 是否完成
}

func (r Result) Fields() []string {
	return []string{
		"host",
		"ip",
		"port",
		"as_organization",
		"protocol",
		"domain",
		"title",
		"fid",
		"certs_subject_cn",
	}
}
