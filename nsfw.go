package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"github.com/koyachi/go-nude"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
)

func upload(c echo.Context) error {

	file, err := c.FormFile("file")
	if err != nil {
		return err
	}
	src, err := file.Open()
	if err != nil {
		return err
	}
	defer src.Close()

	dst, err := os.Create(file.Filename)
	if err != nil {
		return err
	}
	defer dst.Close()

	if _, err = io.Copy(dst, src); err != nil {
		return err
	}
	imagePath := file.Filename

	isNude, err := nude.IsNude(imagePath)
	if err != nil {
		return err
	}
	
	os.Remove(file.Filename)
	if isNude == true {
		return c.HTML(http.StatusOK, fmt.Sprintf("<p>File %s uploaded successfully, the image does contain nudity [NSFW]. </p>", file.Filename))
	} else {
		return c.HTML(http.StatusOK, fmt.Sprintf("<p>File %s uploaded successfully, the image does not contain nudity [SAFE]. </p>", file.Filename))
	}
}

func main() {
	e := echo.New()

	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	e.Static("/", "public")
	e.POST("/upload", upload)

	e.Logger.Fatal(e.Start(":1323"))
}
