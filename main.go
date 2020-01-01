package main

import (
	"log"
	"net/http"
	"os"
	"time"

	"github.com/faiface/beep"
	"github.com/faiface/beep/mp3"
	"github.com/faiface/beep/speaker"
	"github.com/kelseyhightower/envconfig"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"gopkg.in/go-playground/validator.v9"
)

func main() {
	var opts Options

	if err := envconfig.Process("doorbell", &opts); err != nil {
		log.Fatal(err)
	}

	validate := validator.New()
	if err := validate.Struct(opts); err != nil {
		log.Fatal(err)
	}

	e := echo.New()

	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: opts.CORSAllowOrigins,
		AllowMethods: opts.CORSAllowMethods,
	}))

	e.GET("/play", handlePlay)

	e.Logger.Fatal(e.Start(":8080"))
}

func handlePlay(c echo.Context) error {
	f, err := os.Open("/usr/local/share/doorbell/doorbell.mp3")
	if err != nil {
		return err
	}

	streamer, format, err := mp3.Decode(f)
	if err != nil {
		return err
	}
	defer streamer.Close()

	speaker.Init(format.SampleRate, format.SampleRate.N(time.Second / 10))

	done := make(chan bool)
	speaker.Play(beep.Seq(streamer, beep.Callback(func() {
		done <- true
	})))

	<-done

	return c.NoContent(http.StatusOK)
}
