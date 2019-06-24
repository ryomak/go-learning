package util

import (
	"errors"
	"image"
	"image/color"
	"image/jpeg"
	_ "image/png"
	"log"
	"math"
	"os"
	"path/filepath"
	"strings"

	"encoding/json"
	"github.com/goml/gobrain"
)

func Bin(n int) []float64 {
	f := make([]float64, 8)
	for i := uint(0); i < 8; i++ {
		f[i] = float64((n >> i) & 1)
	}
	return f
}

func Dec(d []float64) int {
	n := 0
	for i, v := range d {
		if v > 0.9 {
			n += 1 << uint(i)
		}
	}
	return n
}

func DecodeImage(fname string) ([]float64, error) {
	f, err := os.Open(fname)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	src, _, err := image.Decode(f)
	if err != nil {
		return nil, err
	}

	bounds := src.Bounds()
	w, h := bounds.Dx(), bounds.Dy()
	if w < h {
		w = h
	} else {
		h = w
	}
	//image to float
	bb := make([]float64, w*h*3)
	for y := 0; y < h; y++ {
		for x := 0; x < w; x++ {
			r, g, b, _ := src.At(x, y).RGBA()
			bb[y*w*3+x*3] = float64(r>>8) / 255.0
			bb[y*w*3+x*3+1] = float64(g>>8) / 255.0
			bb[y*w*3+x*3+2] = float64(b>>8) / 255.0
		}
	}
	return bb, nil
}

func EncodeImage(b []float64) error {
	sideLen := int(math.Sqrt(float64(len(b) / 3)))
	log.Println("sidelen:", sideLen)
	img := image.NewRGBA(image.Rect(0, 0, sideLen, sideLen))
	for i := 0; i < len(b); i = i + 3 {
		img.Set(i/3%sideLen, (i/3)/sideLen, color.RGBA{uint8(b[i] * 255), uint8(b[i+1] * 255), uint8(b[i+2] * 255), 255})
	}
	f, err := os.OpenFile("out.jpeg", os.O_WRONLY|os.O_CREATE, 0600)
	if err != nil {
		return err
	}
	defer f.Close()
	return jpeg.Encode(f, img, nil)
}

func LoadImageSet(path string) (map[string][]float64, error) {
	result := map[string][]float64{}
	f, err := os.Open(filepath.Join("dataset", path))
	if err != nil {
		return nil, err
	}
	defer f.Close()
	names, err := f.Readdirnames(-1)
	if err != nil {
		return nil, err
	}
	for _, name := range names {
		if strings.Index(name, ".DS_Store") != -1 {
			continue
		}
		fname := filepath.Join("dataset", path, name)
		ff, err := DecodeImage(fname)
		if err != nil {
			return nil, err
		}
		result[name] = ff
	}
	return result, nil
}

func LoadModel() (*gobrain.FeedForward, error) {
	f, err := os.Open("model.json")
	if err != nil {
		return nil, nil
	}
	defer f.Close()

	ff := &gobrain.FeedForward{}
	err = json.NewDecoder(f).Decode(ff)
	if err != nil {
		return nil, err
	}
	return ff, nil
}

func MakeModel() (*gobrain.FeedForward, error) {
	ff := &gobrain.FeedForward{}
	patterns := [][][]float64{}
	learnSet, err := LoadImageSet("learn")
	if err != nil {
		log.Println("learnSet failed")
		return nil, err
	}
	answerSet, err := LoadImageSet("answer")
	if err != nil {
		log.Println("answerSet failed")
		return nil, err
	}
	for key, d := range learnSet {
		patterns = append(patterns, [][]float64{d, answerSet[key]})
	}
	if len(patterns) == 0 || len(patterns[0][0]) == 0 {
		return nil, errors.New("No images found")
	}
	ff.Init(len(patterns[0][0]), 40, len(patterns[0][1]))
	log.Println("train start")
	ff.Train(patterns, 1000, 0.6, 0.4, false)
	log.Println("train finished")
	return ff, nil
}

func SaveModel(ff *gobrain.FeedForward) error {
	f, err := os.Create("model.json")
	if err != nil {
		return err
	}
	defer f.Close()
	return json.NewEncoder(f).Encode(ff)
}
