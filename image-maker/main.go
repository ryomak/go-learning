package main

import (
	brain "github.com/ryomak/deep-learning-go"
	"github.com/ryomak/deep-learning-go/image/image_maker"
)

func main() {
	c := brain.Config{
		ModelFile:    "model.json",
		EpochNum:     500,
		LearningRate: 0.6,
		MFactor:      0.4,
		Debug:        true,
		HiddenNum:    40,
	}
	i := imaker.Init("dataset/learn", "dataset/answer", "out1.jpeg")

	b := brain.Init(&c, i)
	err := b.LoadModel()
	if err != nil {
		panic(err)
	}
	aa, _ := i.MakePattern()
	b.Train(aa)
	gazou, err := i.Decode("nissy.jpeg")
	if err != nil {
		panic(err)
	}
	i.Encode(gazou)
	i.OutputFile = "out2.jpg"
	i.Encode(b.Output(gazou))
}
