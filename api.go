package main

import (
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func fileNameWithoutExtension(fileName string) string {
	return strings.TrimSuffix(fileName, filepath.Ext(fileName))
}
func cleanFileName(fileName string) string {
	return strings.ReplaceAll(fileName, " ", "_")
}

func upload(c echo.Context) error {
	//-----------
	// Read file
	//-----------

	// Source
	file, err := c.FormFile("file")
	if err != nil {
		return err
	}
	src, err := file.Open()
	if err != nil {
		return err
	}
	defer src.Close()
	// log.Info("file: ", fileNameWithoutExtension(file.Filename))
	clean := cleanFileName(fileNameWithoutExtension(file.Filename))

	// Destination
	dst, err := ioutil.TempFile("uploads", "upload-*.cue")

	if err != nil {
		return err
	}
	defer os.Remove(dst.Name())

	ConvertFile(&src, dst)

	outputName := clean + ".lbl.txt"
	return c.Attachment(dst.Name(), outputName)

}

func main() {

	port := os.Getenv("PORT")
	if port == "" {
		log.Fatal("$PORT must be set")
	}

	e := echo.New()
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(middleware.CORS())
	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "Hello, World!")
	})

	e.POST("/upload", upload)
	e.Logger.Fatal(e.Start(":" + port))
}
