package main

import (
	"flag"
	"log"
	"os"

	"github.com/ryomak/go-learning/image-resolution/util"
)

func main() {
	flag.Parse()

	ff, err := util.LoadModel()
	if err != nil {
		log.Fatal(err)
	}
	if ff == nil {
		log.Println("making model file since not found")
		ff, err = util.MakeModel()
		if err != nil {
			log.Fatal(err)
		}
		log.Println("makeModel complete")
		err = util.SaveModel(ff)
		if err != nil {
			log.Fatal(err)
		}
		log.Println("saveModel complete")
	}
	if flag.NArg() == 0 {
		flag.Usage()
		os.Exit(2)
	}
	for _, arg := range flag.Args() {
		input, err := util.DecodeImage(arg)
		if err != nil {
			log.Fatal(err)
		}
		util.EncodeImage(ff.Update(input))
	}
}
