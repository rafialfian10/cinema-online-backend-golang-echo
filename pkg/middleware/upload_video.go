package middleware

import (
	"io"
	"io/ioutil"
	"net/http"

	"github.com/labstack/echo/v4"
)

func UploadVideo(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		file, err := c.FormFile("trailer")

		if file != nil {
			if err != nil {
				return c.JSON(http.StatusBadRequest, err)
			}

			src, err := file.Open()
			if err != nil {
				return c.JSON(http.StatusBadRequest, err)
			}
			defer src.Close()

			tempFile, err := ioutil.TempFile("uploads/trailer", "video-*.mp4")
			if err != nil {
				return c.JSON(http.StatusBadRequest, err)
			}
			defer tempFile.Close()

			if _, err = io.Copy(tempFile, src); err != nil {
				return c.JSON(http.StatusBadRequest, err)
			}

			data := tempFile.Name()
			filename := data[16:] // split uploads/

			c.Set("dataVideo", filename)
			return next(c)
		}

		c.Set("dataVideo", "")
		return next(c)
	}
}
