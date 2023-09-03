package imageprocess

import (
	"errors"
	"fmt"
	"image/jpeg"
	"io"
	"log"
	"miniproj/entity"
	"net/http"
	"os"
	"path/filepath"

	compression "github.com/nurlantulemisov/imagecompression"
)

type ProductImages struct {
	FolderName string
}

func (pi *ProductImages) ImageProcessing(product *entity.Product) ([]string, error) {
	var compressedImgUrls []string

	for i, value := range product.ProductImages {
		fileName := fmt.Sprintf("%s-%v.jpg", product.ProductName, i)
		_, err := pi.DownloadImageFile(value, fileName)
		if err != nil {
			errSmt := fmt.Sprintf("Error while downloading images err:%v", err)
			log.Println(errSmt)
			return nil, errors.New(errSmt)
		}
		compressedFile, err := pi.ImageCompression(fileName)
		if err != nil {
			errSmt := fmt.Sprintf("Error while compressing images err:%v", err)
			log.Println(errSmt)
			return nil, errors.New(errSmt)
		}
		compressedImgUrls = append(compressedImgUrls, compressedFile)
	}
	return compressedImgUrls, nil
}

func (pi *ProductImages) DownloadImageFile(URL, fileName string) (string, error) {
	//Get the response bytes from the url
	response, err := http.Get(URL)
	if err != nil {
		return "", err
	}
	defer response.Body.Close()

	if response.StatusCode != 200 {
		return "", errors.New("Received 200 response code")
	}

	//Create a empty file
	file, err := os.Create(fileName)
	if err != nil {
		return "", err
	}
	defer file.Close()

	//Write the bytes to the fiel
	_, err = io.Copy(file, response.Body)
	if err != nil {
		return fileName, err
	}

	return fileName, nil
}

func (pi *ProductImages) ImageCompression(fileName string) (string, error) {
	file, err := os.Open(fileName)
	compressedFileName := fmt.Sprintf("Compressed-%v", fileName)
	relativePath := filepath.Join(pi.FolderName, compressedFileName)
	compressedFileName, err = filepath.Abs(relativePath)
	if err != nil {
		log.Fatalf(err.Error())
		return compressedFileName, err
	}
	if err != nil {
		log.Fatalf(err.Error())
		return compressedFileName, err
	}
	fmt.Println("compressed", fileName)
	img, err := jpeg.Decode(file)

	if err != nil {
		log.Fatalf(err.Error())
		return compressedFileName, err
	}

	compressing, _ := compression.New(96)
	compressingImage := compressing.Compress(img)

	f, err := os.Create(compressedFileName)
	if err != nil {
		log.Fatalf("error creating file: %s", err)
		return compressedFileName, err
	}
	defer func(f *os.File) {
		err := f.Close()
		if err != nil {
			log.Fatalf(err.Error())

		}
	}(f)

	err = jpeg.Encode(f, compressingImage, nil)
	if err != nil {
		log.Fatalf(err.Error())
		return compressedFileName, err
	}

	fmt.Println("compressedFileName", compressedFileName)
	fmt.Println(err)
	return compressedFileName, nil
}
