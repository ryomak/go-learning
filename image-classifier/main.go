package main

import (
	"fmt"
	brain "github.com/ryomak/deep-learning-go"
	"github.com/ryomak/deep-learning-go/image/image_classifier"
)

func main() {
	c := brain.Config{
		ModelFile:    "model.json",
		EpochNum:     50,
		LearningRate: 0.01,
		MFactor:      0.2,
		Debug:        false,
		HiddenNum:    50,
	}
	i := iclassifier.Init(
		[]string{
			"たんぽぽ",
			"バラ",
			"ラベンダー",
		},
		"dataset",
		30,
		30,
	)
	b := brain.Init(&c, i)
	b.LoadModel()
	patterns, err := i.MakePattern()
	if err != nil {
		panic(err)
	}
	b.Train(patterns)
	gazou, err := i.Decode("input.jpg")
	if err != nil {
		panic(err)
	}
	output, err := i.Encode(b.Output(gazou))
	fmt.Println(b.Output(gazou))
	fmt.Println("this picture may be ", output)
	sum := float64(len(patterns))
	correct := 0.0
	for _, p := range patterns {
		actual, _ := i.Encode(b.Model.Update(p[0]))
		except, _ := i.Encode(p[1])
		if actual == except {
			correct++
		}
	}
	fmt.Printf("correct:%v, sum:%v  %0.1f％", correct, sum, 100*correct/sum)
}
