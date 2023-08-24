package utils

import (
	"fmt"
	"image/jpeg"
	"io"
	"log"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"strings"
	"time"

	"slices"

	"github.com/anaskhan96/soup"
	"github.com/karmdip-mi/go-fitz"
)

func SetMedia() {

}

func GetColor() {

}

func GetPdf() (result_url_name string, result_url string) {
	start := time.Now()
	const BASE_URL = "http://lmk-lipetsk.ru/"
	const URL = "http://lmk-lipetsk.ru/main_razdel/shedule/index.php"
	resp, _ := soup.Get(URL)
	doc := soup.HTMLParse(resp).FindAll("a")
	for _, value := range doc {
		if strings.Contains(value.FullText(), "Расписание занятий на") {
			result_url = value.Attrs()["href"]
			result_url_name = value.FullText()
		}

	}
	pdf, err := os.Create("shedule.pdf")
	if err != nil {
		log.Fatal(err, "log pdf")
	}
	defer pdf.Close()

	res, err := http.Get(BASE_URL + result_url)
	if err != nil {
		log.Fatal(err, "log response")
	}
	defer res.Body.Close()

	io.Copy(pdf, res.Body)
	fmt.Println(time.Since(start), " time GetPdf")
	return result_url_name, result_url

}

func Convert() []string {
	start := time.Now()

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
				if strings.Contains(path, "page-") && !slices.Contains(fileList, path) {
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
	fmt.Println(time.Since(start))
	return fileList
}
