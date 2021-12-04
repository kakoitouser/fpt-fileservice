package models

import "github.com/kakoitouser/ftp-fileservice/internal/utils"

type File struct {
	UID  utils.UID
	Name string
	Path string
	Data []byte
	Size int64
}
