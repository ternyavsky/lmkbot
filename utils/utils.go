package utils

import (
	"fmt"
	"image/jpeg"
	"log"
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/anaskhan96/soup"
	"github.com/karmdip-mi/go-fitz"
)

func SetMedia() {

}

func GetColor() {

}

func GetPdf() (result_url_name string, result_url string) {
	const URL = "http://lmk-lipetsk.ru/main_razdel/shedule/index.php"
	resp, _ := soup.Get(URL)
	doc := soup.HTMLParse(resp).FindAll("a")
	for _, value := range doc {
		if strings.Contains(value.FullText(), "Расписание занятий на") {
			result_url = value.Attrs()["href"]
			result_url_name = value.FullText()
		}

	}
	return result_url_name, result_url

}

func Convert() []string {

	var files []string
	var fileList []string
	root := "../"
	err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if path == "../lmkbot/shedule.pdf" {
			files = append(files, path)
		}
		return nil
	})
	if err != nil {
		panic(err)
	}
	for _, file := range files {
		doc, err := fitz.New(file)
		if err != nil {
			panic(err)
		}
		folder := strings.TrimSuffix(path.Base(file), filepath.Ext(path.Base(file)))

		// Extract pages as images
		for n := 0; n < doc.NumPage(); n++ {
			img, err := doc.Image(n)
			if err != nil {
				panic(err)
			}
			err = os.MkdirAll("img/"+folder, 0755)
			if err != nil {
				panic(err)
			}

			f, err := os.Create(filepath.Join("img/"+folder+"/", fmt.Sprintf("page-%d.png", n)))

			if err != nil {
				panic(err)
			}
			err = filepath.Walk("../lmkbot/", func(path string, info os.FileInfo, err error) error {
				if strings.Contains(path, "page-") {
					fileList = append(fileList, path)
				}
				return nil
			})

			if err != nil {
				log.Println(err)
			}

			err = jpeg.Encode(f, img, &jpeg.Options{Quality: jpeg.DefaultQuality})
			if err != nil {
				panic(err)
			}

			f.Close()

		}
	}
	return fileList
}
