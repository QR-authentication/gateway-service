package main

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"image"
	"image/png"
	"os"

	"gocv.io/x/gocv"
)

func makeImage(base64String string) {
	imgData, err := base64.StdEncoding.DecodeString(base64String)
	if err != nil {
		fmt.Println("failed to decode Base64:", err)
		return
	}

	img, _, err := image.Decode(bytes.NewReader(imgData))
	if err != nil {
		fmt.Println("failed to decode image:", err)
		return
	}

	outputFile := "qr.png"

	outFile, err := os.Create(outputFile)
	if err != nil {
		fmt.Println("failed to create a file:", err)
		return
	}
	defer outFile.Close()

	err = png.Encode(outFile, img)
	if err != nil {
		fmt.Println("failed to write data into image:", err)
		return
	}
}

func main() {
	// Открываем камеру (по умолчанию используется камера с индексом 0)
	webcam, err := gocv.OpenVideoCapture(0)
	if err != nil {
		fmt.Println("Ошибка при открытии камеры:", err)
		return
	}
	defer webcam.Close()

	// Создаем окно для отображения видео
	window := gocv.NewWindow("Камера")
	defer window.Close()

	// Создаем Mat (матрицу) для хранения изображения
	img := gocv.NewMat()
	defer img.Close()

	// Главный цикл для захвата кадров с камеры и отображения их
	for {
		// Читаем очередной кадр с камеры
		if ok := webcam.Read(&img); !ok {
			fmt.Println("Не удалось захватить кадр с камеры")
			break
		}

		// Если кадр пустой, продолжаем
		if img.Empty() {
			continue
		}

		// Отображаем изображение в окне
		window.IMShow(img)

		// Выход по нажатию клавиши ESC
		if window.WaitKey(1) == 27 {
			break
		}
	}
}
