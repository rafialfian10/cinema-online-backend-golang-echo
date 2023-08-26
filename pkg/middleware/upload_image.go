package middleware

import (
	"fmt"
	"io"
	"io/ioutil"
	"net/http"

	"github.com/labstack/echo/v4"
)

func UploadImage(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		file, err := c.FormFile("thumbnail")

		if file != nil {
			if err != nil {
				return c.JSON(http.StatusBadRequest, err)
			}

			src, err := file.Open()
			if err != nil {
				return c.JSON(http.StatusBadRequest, err)
			}
			defer src.Close()

			tempFile, err := ioutil.TempFile("uploads/thumbnail", "image-*.png")
			if err != nil {
				return c.JSON(http.StatusBadRequest, err)
			}
			defer tempFile.Close()

			if _, err = io.Copy(tempFile, src); err != nil {
				return c.JSON(http.StatusBadRequest, err)
			}

			data := tempFile.Name()
			filename := data[18:] // split uploads/
			fmt.Println("thumbnail", filename)

			c.Set("dataImage", filename)
			return next(c)
		}

		c.Set("dataImage", "")
		return next(c)
	}
}
