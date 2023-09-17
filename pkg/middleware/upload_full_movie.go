package middleware

import (
	"io"
	"io/ioutil"
	"net/http"

	"github.com/labstack/echo/v4"
)

func UploadFullMovie(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		file, err := c.FormFile("full_movie")

		if file != nil {
			if err != nil {
				return c.JSON(http.StatusBadRequest, err)
			}

			src, err := file.Open()
			if err != nil {
				return c.JSON(http.StatusBadRequest, err)
			}
			defer src.Close()

			tempFile, err := ioutil.TempFile("uploads/full_movie", "video-full-movie-*.mp4")
			if err != nil {
				return c.JSON(http.StatusBadRequest, err)
			}
			defer tempFile.Close()

			if _, err = io.Copy(tempFile, src); err != nil {
				return c.JSON(http.StatusBadRequest, err)
			}

			data := tempFile.Name()
			filename := data[19:] // split uploads/

			c.Set("dataFullMovie", filename)
			return next(c)
		}

		c.Set("dataFullMovie", "")
		return next(c)
	}
}
