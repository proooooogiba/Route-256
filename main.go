package main

import (
	"errors"
	"fmt"
	"sort"
)

type Apple uint32

type Bag struct {
	Apple Apple
}

type Box struct {
	Bags []Bag
}

type Car struct {
	Boxes []Box
}

func applesToBag(appleNumber uint32) ([]Bag, error) {
	if appleNumber == 0 {
		return nil, errors.New("number of apples can't be zero")
	}
	var i uint32
	var bags = make([]Bag, appleNumber)

	for i = 0; i < appleNumber; i++ {
		bags[i] = Bag{Apple: Apple(i + 1)}
	}
	return bags, nil
}

func packageBagsToBoxes(bags []Bag, numberOfApplesInBox uint32) ([]Box, error) {
	if numberOfApplesInBox == 0 {
		return nil, errors.New("number of apples in box can't be zero")
	}

	sort.Slice(bags, func(i, j int) bool {
		return bags[i].Apple > bags[j].Apple
	})

	var numberOfBoxes, i uint32
	bagsLength := uint32(len(bags))
	numberOfBoxes = bagsLength / numberOfApplesInBox
	var boxes = make([]Box, numberOfBoxes)

	for i = 0; i < numberOfBoxes; i++ {
		boxes[i] = Box{Bags: bags[i*numberOfApplesInBox : i*numberOfApplesInBox+numberOfApplesInBox]}
	}

	if bagsLength%numberOfApplesInBox == 0 {
		return boxes, nil
	}
	numberOfApplesInLastBox := bagsLength % numberOfApplesInBox
	boxes = append(boxes, Box{})
	bagsLength++
	boxes[numberOfBoxes] = Box{Bags: bags[bagsLength-1-numberOfApplesInLastBox : bagsLength-1]}

	return boxes, nil
}

func boxesToCars(boxes []Box, numberOfCars uint32) ([]Car, error) {
	if numberOfCars == 0 {
		return nil, errors.New("number of cars can't be zero")
	}

	var i uint32
	numberOfBoxes := uint32(len(boxes))
	boxesInCar := numberOfBoxes / numberOfCars

	var cars = make([]Car, numberOfCars)

	for i = 0; i < numberOfCars; i++ {
		cars[i] = Car{Boxes: make([]Box, boxesInCar)}
	}

	numberOfCarsWithExtraBox := numberOfBoxes % numberOfCars
	for i = 0; i < numberOfCarsWithExtraBox; i++ {
		cars[i] = Car{Boxes: make([]Box, boxesInCar+1)}
	}

	for i = 0; i < numberOfBoxes; i++ {
		cars[i%numberOfCars].Boxes[i/numberOfCars] = boxes[i]
	}

	return cars, nil
}

func carsFormatString(cars []Car) string {
	var formatString string
	carNumber := len(cars)
	for i := 0; i < carNumber; i++ {
		boxNumber := len(cars[i].Boxes)
		for j := 0; j < boxNumber; j++ {
			formatString += fmt.Sprint("Машина: ", i+1, ", Ящик: ", j*carNumber+i+1, ", Яблоки:  [")
			bagsNumber := len(cars[i].Boxes[j].Bags)
			for k := 0; k < bagsNumber-1; k++ {
				formatString += fmt.Sprint(cars[i].Boxes[j].Bags[k].Apple, ", ")
			}
			if bagsNumber > 0 {
				formatString += fmt.Sprint(cars[i].Boxes[j].Bags[bagsNumber-1].Apple)
			}
			formatString += fmt.Sprint("]\n")
		}
	}
	return formatString
}

func main() {
	bags, err := applesToBag(100)
	if err != nil {
		fmt.Println(fmt.Sprintf("Error: %s", err))
		return
	}
	boxes, err := packageBagsToBoxes(bags, 10)
	if err != nil {
		fmt.Println(fmt.Sprintf("Error: %s", err))
		return
	}
	cars, err := boxesToCars(boxes, 2)
	if err != nil {
		fmt.Println(fmt.Sprintf("Error: %s", err))
		return
	}
	formatCarString := carsFormatString(cars)
	fmt.Println(formatCarString)
}
