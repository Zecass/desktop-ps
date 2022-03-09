package settings

import (
	"errors"
	"fmt"
	"image/color"
	"os"
	"strconv"
	"strings"

	"gopkg.in/yaml.v2"
)

type Settings struct {
	Resolution  string `yaml:"resolution"`
	ResolutionX float64
	ResolutionY float64

	BackgroundColorHex  string `yaml:"backgroundColorHex"`
	BackgroundColorRGBA color.RGBA

	Dot struct {
		Radius          float64 `yaml:"radius"`
		ColorHex        string  `yaml:"colorHex"`
		ColorRGBA       color.RGBA
		BorderWidth     float64 `yaml:"borderWidth"`
		BorderColorHex  string  `yaml:"borderColorHex"`
		BorderColorRGBA color.RGBA
	} `yaml:"dot"`

	Font struct {
		Size      float64 `yaml:"size"`
		ColorHex  string  `yaml:"colorHex"`
		ColorRGBA color.RGBA
	} `yaml:"font"`

	Link struct {
		Width     float64 `yaml:"width"`
		ColorHex  string  `yaml:"colorHex"`
		ColorRGBA color.RGBA
	} `yaml:"link"`
}

func GetSettings() (s *Settings, err error) {
	f, err := os.Open("settings.yml")
	if err != nil {
		return nil, err
	}

	defer f.Close()

	decoder := yaml.NewDecoder(f)
	err = decoder.Decode(&s)
	if err != nil {
		return nil, err
	}

	err = s.parseSettigns()
	if err != nil {
		return nil, err
	}

	return
}

func (s *Settings) parseSettigns() (err error) {
	err = s.parseResolution()
	if err != nil {
		return err
	}

	if s.BackgroundColorHex == "" {
		s.BackgroundColorHex = "#323232"
		s.BackgroundColorRGBA = color.RGBA{50, 50, 50, 255}
	} else {
		s.BackgroundColorRGBA = hexToRGBA(s.BackgroundColorHex)
	}

	if s.Dot.Radius == 0 {
		s.Dot.Radius = (s.ResolutionY * 6) / 1800
	}

	if s.Dot.BorderWidth == 0 {
		s.Dot.BorderWidth = (s.ResolutionY * 1) / 1800
	}

	if s.Dot.ColorHex == "" {
		s.Dot.ColorHex = "#006464"
		s.Dot.ColorRGBA = color.RGBA{0, 100, 100, 255}
	} else {
		s.Dot.ColorRGBA = hexToRGBA(s.Dot.ColorHex)
	}

	if s.Dot.BorderColorHex == "" {
		s.Dot.BorderColorHex = "#000000"
		s.Dot.BorderColorRGBA = color.RGBA{0, 0, 0, 255}
	} else {
		s.Dot.BorderColorRGBA = hexToRGBA(s.Dot.BorderColorHex)
	}

	if s.Font.Size == 0 {
		s.Font.Size = (s.ResolutionY * 30) / 1800
	}

	if s.Font.ColorHex == "" {
		s.Font.ColorHex = "#FFFFFF"
		s.Font.ColorRGBA = color.RGBA{255, 255, 255, 255}
	} else {
		s.Font.ColorRGBA = hexToRGBA(s.Font.ColorHex)
	}

	if s.Link.Width == 0 {
		s.Link.Width = (s.ResolutionY * 1) / 1800
	}

	if s.Link.ColorHex == "" {
		s.Link.ColorHex = "#646464"
		s.Link.ColorRGBA = color.RGBA{100, 100, 100, 255}
	} else {
		s.Link.ColorRGBA = hexToRGBA(s.Link.ColorHex)
	}

	return
}

func (s *Settings) parseResolution() (err error) {
	var x, y float64

	r := strings.Split(s.Resolution, "x")
	if len(r) != 2 {
		err = errors.Unwrap(fmt.Errorf("invalid resolution: %s", s.Resolution))
		return err
	}

	x, err = strconv.ParseFloat(r[0], 64)
	if err != nil {
		return err
	}

	y, err = strconv.ParseFloat(r[1], 64)
	if err != nil {
		return err
	}

	s.ResolutionX = x
	s.ResolutionY = y

	return
}

func hexToRGBA(hex string) color.RGBA {
	var r, g, b, a int64 = 0, 0, 0, 255

	if len(hex) == 7 {
		r, _ = strconv.ParseInt(hex[1:3], 16, 0)
		g, _ = strconv.ParseInt(hex[3:5], 16, 0)
		b, _ = strconv.ParseInt(hex[5:7], 16, 0)
	}

	return color.RGBA{uint8(r), uint8(g), uint8(b), uint8(a)}
}
