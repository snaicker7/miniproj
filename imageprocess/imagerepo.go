package imageprocess

import "miniproj/entity"

type ImageInterface interface {
	ImageProcessing(product *entity.Product) ([]string, error)
	DownloadImageFile(URL, fileName string) (string, error)
	ImageCompression(fileName string) (string, error)
}
