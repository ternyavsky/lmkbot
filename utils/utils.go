package utils

import (
	"fmt"
	"image/jpeg"
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/karmdip-mi/go-fitz"
)

func Convert() {

	var files []string

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

			err = jpeg.Encode(f, img, &jpeg.Options{Quality: jpeg.DefaultQuality})
			if err != nil {
				panic(err)
			}

			f.Close()

		}
	}
}
