package main

import (
	"fmt"

	brain "github.com/ryomak/deep-learning-go"
	iclassifier "github.com/ryomak/deep-learning-go/image/image_classifier"
)

func main() {
	c := brain.Config{
		ModelFile:  "model.json",
		EpochNum:   50,
		Bias:       true,
		Hiddens:    []int{1000, 100},
		Activation: 5,
		Mode:       0,
	}
	i := iclassifier.Init(
		[]string{
			"tida",
			"yuna",
			"lulu",
		},
		"dataset",
		30,
		30,
	)
  b := brain.Init(&c, i)
  patterns, err := i.MakePattern()
  if err != nil {
    panic(err)
  }
  b.NewAdamTrainer(0.01, len(patterns))
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
		actual, _ := i.Encode(b.Output(p.Input))
		except, _ := i.Encode(p.Response)
		if actual == except {
			correct++
		}
	}
	fmt.Printf("correct:%v, sum:%v  %0.1f％", correct, sum, 100*correct/sum)
}
