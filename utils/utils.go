package utils

import (
	"fmt"
	"image/jpeg"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/anaskhan96/soup"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/karmdip-mi/go-fitz"
)

func SetMedia() {

}

func GetColor() {

}

func GetPdf() string {
	start := time.Now()
	const BASE_URL = "http://lmk-lipetsk.ru/"
	const URL = "http://lmk-lipetsk.ru/main_razdel/shedule/index.php"
	resp, _ := soup.Get(URL)
	result_url := ""
	result_url_name := ""

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

	fmt.Println(time.Since(start))
	return result_url_name

}

func Convert(c chan tgbotapi.InputMediaPhoto, cap string) {

	start := time.Now()

	doc, err := fitz.New("../lmkbot/shedule.pdf")
	if err != nil {
		panic(err)
	}

	for n := 0; n < doc.NumPage(); n++ {
		img, err := doc.Image(n)
		if err != nil {
			panic(err)
		}
		f, err := os.Create(filepath.Join("../lmkbot/img/shedule/", fmt.Sprintf("page-%d.png", n)))
		fmt.Println(f.Name())
		if err != nil {
			panic(err)
		}
		err = jpeg.Encode(f, img, &jpeg.Options{Quality: jpeg.DefaultQuality})
		if err != nil {
			panic(err)
		}
		if n == 1 {

			added := tgbotapi.NewInputMediaPhoto(tgbotapi.FilePath(f.Name()))
			added.Caption = cap

			c <- added
		} else {
			added := tgbotapi.NewInputMediaPhoto(tgbotapi.FilePath(f.Name()))
			c <- added
		}
		if n == doc.NumPage()-1 {
			close(c)
		}

		f.Close()

	}

	fmt.Println(time.Since(start), " convert")
}
