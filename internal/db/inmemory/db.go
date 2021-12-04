package inmemory

import (
	"errors"

	"github.com/kakoitouser/ftp-fileservice/internal/models"
	"github.com/kakoitouser/ftp-fileservice/internal/utils"
)

var Files []*models.File = []*models.File{
	{utils.UID(1), "file.txt", "/files/", []byte("hello world"), 128},
}

func GetFileByUID(uid string) (*models.File, error) {
	UID := utils.UID(uid)
	for i := range Files {
		if UID == Files[i].UID {
			return Files[i], nil
		}
	}
	return nil, errors.New("file not found")
}
