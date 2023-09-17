package middleware

import (
	"io"
	"io/ioutil"
	"net/http"

	"github.com/labstack/echo/v4"
)

func UploadPhoto(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		file, err := c.FormFile("photo")

		if file != nil {
			if err != nil {
				return c.JSON(http.StatusBadRequest, err)
			}

			src, err := file.Open()
			if err != nil {
				return c.JSON(http.StatusBadRequest, err)
			}
			defer src.Close()

			tempFile, err := ioutil.TempFile("uploads/photo", "image-photo-*.png")
			if err != nil {
				return c.JSON(http.StatusBadRequest, err)
			}
			defer tempFile.Close()

			if _, err = io.Copy(tempFile, src); err != nil {
				return c.JSON(http.StatusBadRequest, err)
			}

			data := tempFile.Name()
			filename := data[14:] // split uploads/
			// fmt.Println("photo", filename)

			c.Set("dataPhoto", filename)
			return next(c)
		}

		c.Set("dataPhoto", "")
		return next(c)
	}
}
