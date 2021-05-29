package service

import (
	"archive/tar"
	"compress/gzip"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"sort"
	"strconv"
	"time"

	"github.com/akhilrex/hammond/db"
	"github.com/akhilrex/hammond/internal/sanitize"
	"github.com/akhilrex/hammond/models"
	uuid "github.com/satori/go.uuid"
)

func CreateAttachment(path, originalName string, size int64, contentType, userId string) (*db.Attachment, error) {
	model := &db.Attachment{
		Path:         path,
		OriginalName: originalName,
		Size:         size,
		ContentType:  contentType,
		UserID:       userId,
	}
	tx := db.DB.Create(&model)

	if tx.Error != nil {
		return nil, tx.Error
	}
	return model, nil
}

func CreateQuickEntry(model models.CreateQuickEntryModel, attachmentId, userId string) (*db.QuickEntry, error) {
	toCreate := &db.QuickEntry{
		AttachmentID: attachmentId,
		UserID:       userId,
		Comments:     model.Comments,
	}
	tx := db.DB.Create(&toCreate)

	if tx.Error != nil {
		return nil, tx.Error
	}
	return toCreate, nil
}
func GetAllQuickEntries(sorting string) (*[]db.QuickEntry, error) {
	return db.GetAllQuickEntries(sorting)
}
func GetQuickEntriesForUser(userId, sorting string) (*[]db.QuickEntry, error) {
	return db.GetQuickEntriesForUser(userId, sorting)
}
func GetQuickEntryById(id string) (*db.QuickEntry, error) {
	return db.GetQuickEntryById(id)
}
func SetQuickEntryAsProcessed(id string) error {
	return db.SetQuickEntryAsProcessed(id, time.Now())

}

func GetAttachmentById(id string) (*db.Attachment, error) {
	return db.GetAttachmentById(id)
}

func GetAllBackupFiles() ([]string, error) {
	var files []string
	folder := createConfigFolderIfNotExists("backups")
	err := filepath.Walk(folder, func(path string, info os.FileInfo, err error) error {
		if !info.IsDir() {
			files = append(files, path)
		}
		return nil
	})
	sort.Sort(sort.Reverse(sort.StringSlice(files)))
	return files, err
}

func GetFileSize(path string) (int64, error) {
	info, err := os.Stat(path)
	if err != nil {
		return 0, err
	}
	return info.Size(), nil
}

func deleteOldBackup() {
	files, err := GetAllBackupFiles()
	if err != nil {
		return
	}
	if len(files) <= 5 {
		return
	}

	toDelete := files[5:]
	for _, file := range toDelete {
		fmt.Println(file)
		DeleteFile(file)
	}
}

func GetFilePath(originalName string) string {
	dataPath := os.Getenv("DATA")
	return path.Join(dataPath, getFileName(originalName))
}
func getFileName(orig string) string {

	ext := filepath.Ext(orig)
	return uuid.NewV4().String() + ext

}

func CreateBackup() (string, error) {

	backupFileName := "hammond_backup_" + time.Now().Format("2006.01.02_150405") + ".tar.gz"
	folder := createConfigFolderIfNotExists("backups")
	configPath := os.Getenv("CONFIG")
	tarballFilePath := path.Join(folder, backupFileName)
	file, err := os.Create(tarballFilePath)
	if err != nil {
		return "", errors.New(fmt.Sprintf("Could not create tarball file '%s', got error '%s'", tarballFilePath, err.Error()))
	}
	defer file.Close()

	dbPath := path.Join(configPath, "hammond.db")
	_, err = os.Stat(dbPath)
	if err != nil {
		return "", errors.New(fmt.Sprintf("Could not find db file '%s', got error '%s'", dbPath, err.Error()))
	}
	gzipWriter := gzip.NewWriter(file)
	defer gzipWriter.Close()

	tarWriter := tar.NewWriter(gzipWriter)
	defer tarWriter.Close()

	err = addFileToTarWriter(dbPath, tarWriter)
	if err == nil {
		deleteOldBackup()
	}
	return backupFileName, err
}

func addFileToTarWriter(filePath string, tarWriter *tar.Writer) error {
	file, err := os.Open(filePath)
	if err != nil {
		return errors.New(fmt.Sprintf("Could not open file '%s', got error '%s'", filePath, err.Error()))
	}
	defer file.Close()

	stat, err := file.Stat()
	if err != nil {
		return errors.New(fmt.Sprintf("Could not get stat for file '%s', got error '%s'", filePath, err.Error()))
	}

	header := &tar.Header{
		Name:    filePath,
		Size:    stat.Size(),
		Mode:    int64(stat.Mode()),
		ModTime: stat.ModTime(),
	}

	err = tarWriter.WriteHeader(header)
	if err != nil {
		return errors.New(fmt.Sprintf("Could not write header for file '%s', got error '%s'", filePath, err.Error()))
	}

	_, err = io.Copy(tarWriter, file)
	if err != nil {
		return errors.New(fmt.Sprintf("Could not copy the file '%s' data to the tarball, got error '%s'", filePath, err.Error()))
	}

	return nil
}
func httpClient() *http.Client {
	client := http.Client{
		CheckRedirect: func(r *http.Request, via []*http.Request) error {
			//	r.URL.Opaque = r.URL.Path
			return nil
		},
	}

	return &client
}

func createFolder(folder string, parent string) string {
	folder = cleanFileName(folder)
	//str := stringy.New(folder)
	folderPath := path.Join(parent, folder)
	if _, err := os.Stat(folderPath); os.IsNotExist(err) {
		os.MkdirAll(folderPath, 0777)
		changeOwnership(folderPath)
	}
	return folderPath
}

func createDataFolderIfNotExists(folder string) string {
	dataPath := os.Getenv("DATA")
	return createFolder(folder, dataPath)
}
func createConfigFolderIfNotExists(folder string) string {
	dataPath := os.Getenv("CONFIG")
	return createFolder(folder, dataPath)
}

func cleanFileName(original string) string {
	return sanitize.Name(original)
}
func changeOwnership(path string) {
	uid, err1 := strconv.Atoi(os.Getenv("PUID"))
	gid, err2 := strconv.Atoi(os.Getenv("PGID"))
	fmt.Println(path)
	if err1 == nil && err2 == nil {
		fmt.Println(path + " : Attempting change")
		os.Chown(path, uid, gid)
	}

}
func DeleteFile(filePath string) error {
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		return err
	}
	if err := os.Remove(filePath); err != nil {
		return err
	}
	return nil
}
func checkError(err error) {
	if err != nil {
		panic(err)
	}
}
