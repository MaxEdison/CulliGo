package captchasolver

import (
	"fmt"

	ocr "github.com/ranghetto/go_ocr_space"
)

func Solver(path string) (string, error) {
	api_key := "YOUR API KEY"

	config := ocr.InitConfig(api_key, "eng", ocr.OCREngine2)

	result, err := config.ParseFromLocal(path)

	if err != nil {
		fmt.Println(err)
	}

	value := result.JustText()

	return value, err

}
