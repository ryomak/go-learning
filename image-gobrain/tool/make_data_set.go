package main

import (
	"image"
	"image/color"
	"image/jpeg"
	_ "image/png"
	"log"
	"os"
	"path/filepath"
)

func main() {
	f, err := os.Open(filepath.Join("dataset", "answer"))
	if err != nil {
		panic(err)
	}
	defer f.Close()
	names, err := f.Readdirnames(-1)
	if err != nil {
		panic(err)
	}
	for _, v := range names {
		log.Println(v, ":make mosaic")
		err := makeMosaic(v)
		if err != nil {
			log.Println(err)
		}
	}
}

func makeMosaic(fname string) error {
	f, err := os.Open(filepath.Join("dataset", "answer", fname))
	if err != nil {
		return err
	}
	defer f.Close()
	img, _, err := image.Decode(f)
	if err != nil {
		return err
	}
	//モザイク処理
	bounds := img.Bounds()
	dest := image.NewRGBA(bounds)
	block := 5 // モザイクの粒度
	for y := bounds.Min.Y + (block-1)/2; y < bounds.Max.Y; y = y + block {
		for x := bounds.Min.X + (block-1)/2; x < bounds.Max.X; x = x + block {
			var cr, cg, cb float32
			var alpha uint8
			for j := y - (block-1)/2; j <= y+(block-1)/2; j++ {
				for i := x - (block-1)/2; i <= x+(block-1)/2; i++ {
					if i >= 0 && j >= 0 && i < bounds.Max.X && j < bounds.Max.Y {
						c := color.RGBAModel.Convert(img.At(i, j))
						col := c.(color.RGBA)
						cr += float32(col.R)
						cg += float32(col.G)
						cb += float32(col.B)
						alpha = col.A
					}
				}
			}
			cr = cr / float32(block*block)
			cg = cg / float32(block*block)
			cb = cb / float32(block*block)
			for j := y - (block-1)/2; j <= y+(block-1)/2; j++ {
				for i := x - (block-1)/2; i <= x+(block-1)/2; i++ {
					if i >= 0 && j >= 0 && i < bounds.Max.X && j < bounds.Max.Y {
						dest.Set(i, j, color.RGBA{uint8(cr), uint8(cg), uint8(cb), alpha})
					}
				}
			}
		}
	}
	outputFile, err := os.Create(filepath.Join("dataset", "learn", fname))
	if err != nil {
		return err
	}
	defer outputFile.Close()
	err = jpeg.Encode(outputFile, dest, nil)
	if err != nil {
		return err
	}
	return nil
}
