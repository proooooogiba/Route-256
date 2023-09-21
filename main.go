package main

import (
	"errors"
	"fmt"
	"log"
	"sort"
	"strconv"
	"strings"
)

type Apple uint32

type Box struct {
	Bags []*Apple
}

type Car struct {
	Boxes []*Box
}

func applesToBag(appleNumber uint32) ([]*Apple, error) {
	var bags = make([]*Apple, appleNumber)

	for i := range bags {
		apple := Apple(i + 1)
		bags[i] = &apple
	}
	return bags, nil
}

func packageBagsToBoxes(bags []*Apple, numberOfApplesInBox int) ([]*Box, error) {
	if numberOfApplesInBox <= 0 {
		return nil, errors.New("number of apples in box can't be zero or negative")
	}

	sort.Slice(bags, func(i, j int) bool {
		return *bags[i] > *bags[j]
	})

	bagsLength := len(bags)

	numberOfBoxes := bagsLength / numberOfApplesInBox
	if bagsLength%numberOfApplesInBox != 0 {
		numberOfBoxes++
	}

	var boxes = make([]*Box, numberOfBoxes)

	for i := range boxes {
		start := i * numberOfApplesInBox
		end := (i + 1) * numberOfApplesInBox
		if end > bagsLength {
			end = bagsLength
		}
		boxes[i] = &Box{Bags: bags[start:end]}
	}

	return boxes, nil
}

func boxesToCars(boxes []*Box, numberOfCars int) ([]*Car, error) {
	if numberOfCars <= 0 {
		return nil, errors.New("number of cars can't be zero or negative")
	}

	var cars = make([]*Car, numberOfCars)

	for i := range boxes {
		if cars[i%numberOfCars] == nil {
			cars[i%numberOfCars] = &Car{}
		}
		cars[i%numberOfCars].Boxes = append(cars[i%numberOfCars].Boxes, boxes[i])
	}

	return cars, nil
}

func carsFormatString(cars []*Car) string {
	builder := strings.Builder{}
	carNumber := len(cars)
	for i, car := range cars {
		for j, box := range car.Boxes {
			builder.WriteString("Машина: ")
			builder.WriteString(strconv.Itoa(i + 1))
			builder.WriteString(", Ящик: ")
			builder.WriteString(strconv.Itoa(j*carNumber + i + 1))
			builder.WriteString(", Яблоки:  [")
			for k, apple := range box.Bags {
				builder.WriteString(strconv.Itoa(int(*apple)))
				if k != len(box.Bags)-1 {
					builder.WriteString(", ")
				}
			}
			builder.WriteString("]\n")
		}
	}
	formatString := builder.String()
	return formatString
}

func main() {
	bags, err := applesToBag(100)
	if err != nil {
		log.Fatal(err)
	}
	boxes, err := packageBagsToBoxes(bags, 7)
	if err != nil {
		log.Fatal(err)
	}
	cars, err := boxesToCars(boxes, 3)
	if err != nil {
		log.Fatal(err)
	}
	formatCarString := carsFormatString(cars)
	fmt.Println(formatCarString)
}
