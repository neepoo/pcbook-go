package service

import (
	"bytes"
	"fmt"
	"github.com/google/uuid"
	"log"
	"os"
	"path/filepath"
	"sync"
)

type ImageStore interface {
	Save(laptopID string, imageType string, imageData bytes.Buffer) (string, error)
}

type DiskImageStore struct {
	mutex       sync.RWMutex
	imageFolder string
	images      map[string]*ImageInfo
}
func (d *DiskImageStore) Save(laptopID string, imageType string, imageData bytes.Buffer) (string, error) {
	imageID, err := uuid.NewRandom()
	if err != nil {
		return "", fmt.Errorf("cannot genearate uuid %v", err)
	}
	imagePath := filepath.Join(d.imageFolder, imageID.String()+imageType)
	//imagePath := fmt.Sprintf("%s/%s%s", d.imageFolder, imageID, imageType)
	log.Println(d.imageFolder, imageID, imageType, imagePath)
	file, err := os.Create(imagePath)
	if err != nil {
		return "", fmt.Errorf("cannot create image file")
	}
	defer file.Close()
	_, err = imageData.WriteTo(file)
	if err != nil{
		return "", fmt.Errorf("cannot write image to file :%v", err)
	}

	d.mutex.Lock()
	defer d.mutex.Unlock()

	d.images[imageID.String()] = &ImageInfo{
		LaptopID: laptopID,
		Type:     imageType,
		Path:     imagePath,
	}
	return imageID.String(), nil
}

func NewDiskImageStore(imageFolder string) *DiskImageStore {
	return &DiskImageStore{
		imageFolder: imageFolder,
		images:      map[string]*ImageInfo{},
	}
}

type ImageInfo struct {
	LaptopID string
	Type     string
	Path     string
}
