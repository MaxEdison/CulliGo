package captchasolver

import (
	"fmt"

	ocr "github.com/ranghetto/go_ocr_space"
)

func Solver(path string, api_key string) (string, error) {

	config := ocr.InitConfig(api_key, "eng", ocr.OCREngine2)

	result, err := config.ParseFromLocal(path)

	if err != nil {
		fmt.Println(err)
	}

	value := result.JustText()

	return value, err

}
