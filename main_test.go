package main

import (
	"reflect"
	"testing"
)

func Test_applesToBag(t *testing.T) {
	type args struct {
		appleNumber uint32
	}
	tests := []struct {
		name    string
		args    args
		want    []Bag
		wantErr bool
	}{
		{
			name:    "input number - 5",
			args:    struct{ appleNumber uint32 }{appleNumber: 5},
			want:    []Bag{{1}, {2}, {3}, {4}, {5}},
			wantErr: false,
		},
		{
			name:    "invalid input - 0",
			args:    struct{ appleNumber uint32 }{appleNumber: 0},
			want:    nil,
			wantErr: true,
		},
		{
			name: "input number - 10",
			args: struct{ appleNumber uint32 }{appleNumber: 10},
			want: []Bag{{1}, {2}, {3},
				{4}, {5}, {6},
				{7}, {8}, {9},
				{10}},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := applesToBag(tt.args.appleNumber)
			if (err != nil) != tt.wantErr {
				t.Errorf("applesToBag() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("applesToBag() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_boxesToCars(t *testing.T) {
	type args struct {
		boxes        []Box
		numberOfCars uint32
	}

	boxes := []Box{
		{Bags: []Bag{
			{Apple: Apple(1)},
			{Apple: Apple(2)},
			{Apple: Apple(3)},
		}},
		{Bags: []Bag{
			{Apple: Apple(3)},
			{Apple: Apple(4)},
			{Apple: Apple(5)},
		}},
	}

	cars := []Car{
		{Boxes: []Box{
			{Bags: []Bag{
				{Apple: Apple(1)},
				{Apple: Apple(2)},
				{Apple: Apple(3)},
			}},
		}},
		{Boxes: []Box{
			{Bags: []Bag{
				{Apple: Apple(3)},
				{Apple: Apple(4)},
				{Apple: Apple(5)},
			}},
		}},
	}

	boxesWithDiffrentNumberOfBags := []Box{
		{Bags: []Bag{
			{Apple: Apple(1)},
			{Apple: Apple(2)},
			{Apple: Apple(3)},
		}},
		{Bags: []Bag{
			{Apple: Apple(4)},
			{Apple: Apple(5)},
			{Apple: Apple(6)},
		}},
		{Bags: []Bag{
			{Apple: Apple(7)},
			{Apple: Apple(8)},
			{Apple: Apple(9)},
		}},
	}

	carsWithExtraBoxes := []Car{
		{Boxes: []Box{
			{Bags: []Bag{
				{Apple: Apple(1)},
				{Apple: Apple(2)},
				{Apple: Apple(3)},
			}},
			{Bags: []Bag{
				{Apple: Apple(7)},
				{Apple: Apple(8)},
				{Apple: Apple(9)},
			}},
		}},
		{Boxes: []Box{
			{Bags: []Bag{
				{Apple: Apple(4)},
				{Apple: Apple(5)},
				{Apple: Apple(6)},
			}},
		}},
	}

	tests := []struct {
		name    string
		args    args
		want    []Car
		wantErr bool
	}{
		{
			name: "valid input packed to 2 cars",
			args: struct {
				boxes        []Box
				numberOfCars uint32
			}{boxes: boxes, numberOfCars: 2},
			want:    cars,
			wantErr: false,
		},
		{
			name: "invalid input - number of cars is zero",
			args: struct {
				boxes        []Box
				numberOfCars uint32
			}{boxes: nil, numberOfCars: 0},
			want:    nil,
			wantErr: true,
		},
		{
			name: "empty input packed to 2 cars",
			args: struct {
				boxes        []Box
				numberOfCars uint32
			}{boxes: boxesWithDiffrentNumberOfBags, numberOfCars: 2},
			want:    carsWithExtraBoxes,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := boxesToCars(tt.args.boxes, tt.args.numberOfCars)
			if (err != nil) != tt.wantErr {
				t.Errorf("boxesToCars() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("boxesToCars() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_packageBagsToBoxes(t *testing.T) {
	type args struct {
		bags                []Bag
		numberOfApplesInBox uint32
	}

	bags := []Bag{
		{Apple: Apple(1)},
		{Apple: Apple(3)},
		{Apple: Apple(2)},
		{Apple: Apple(4)},
	}
	boxesTwoApples := []Box{
		{Bags: []Bag{
			{Apple: Apple(4)},
			{Apple: Apple(3)},
		}},
		{Bags: []Bag{
			{Apple: Apple(2)},
			{Apple: Apple(1)},
		}},
	}
	boxesThreeApples := []Box{
		{Bags: []Bag{
			{Apple: Apple(4)},
			{Apple: Apple(3)},
			{Apple: Apple(2)},
		}},
		{Bags: []Bag{
			{Apple: Apple(1)},
		}},
	}

	tests := []struct {
		name    string
		args    args
		want    []Box
		wantErr bool
	}{
		{
			name: "valid input - 2 apples in box",
			args: struct {
				bags                []Bag
				numberOfApplesInBox uint32
			}{bags: bags, numberOfApplesInBox: 2},
			want:    boxesTwoApples,
			wantErr: false,
		},
		{
			name: "invalid input - number of apples in bag is zero",
			args: struct {
				bags                []Bag
				numberOfApplesInBox uint32
			}{bags: bags, numberOfApplesInBox: 0},
			want:    nil,
			wantErr: true,
		},
		{
			name: "valid input - 3 apples in box",
			args: struct {
				bags                []Bag
				numberOfApplesInBox uint32
			}{bags: bags, numberOfApplesInBox: 3},
			want:    boxesThreeApples,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := packageBagsToBoxes(tt.args.bags, tt.args.numberOfApplesInBox)
			if (err != nil) != tt.wantErr {
				t.Errorf("packageBagsToBoxes() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("packageBagsToBoxes() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_carsFormatString(t *testing.T) {
	type args struct {
		cars []Car
	}
	cars := []Car{
		{Boxes: []Box{
			{Bags: []Bag{
				{Apple: Apple(1)},
				{Apple: Apple(2)},
				{Apple: Apple(3)},
			}},
		}},
		{Boxes: []Box{
			{Bags: []Bag{
				{Apple: Apple(3)},
				{Apple: Apple(4)},
				{Apple: Apple(5)},
			}},
		}},
	}

	carsWithDiffrentNumbersOfBoxes := []Car{
		{Boxes: []Box{
			{Bags: []Bag{
				{Apple: Apple(1)},
				{Apple: Apple(2)},
				{Apple: Apple(3)},
			}},
			{Bags: []Bag{
				{Apple: Apple(6)},
				{Apple: Apple(7)},
			}},
		}},
		{Boxes: []Box{
			{Bags: []Bag{
				{Apple: Apple(3)},
				{Apple: Apple(4)},
				{Apple: Apple(5)},
			}},
		}},
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "2 cars with 1 box with 3 bags in each",
			args: struct{ cars []Car }{cars: cars},
			want: "Машина: 1, Ящик: 1, Яблоки:  [1, 2, 3]\n" +
				"Машина: 2, Ящик: 2, Яблоки:  [3, 4, 5]\n",
		},
		{
			name: "empty input",
			args: struct{ cars []Car }{cars: nil},
			want: "",
		},
		{
			name: "cars with different numbers of boxes in each",
			args: struct{ cars []Car }{cars: carsWithDiffrentNumbersOfBoxes},
			want: "Машина: 1, Ящик: 1, Яблоки:  [1, 2, 3]\n" +
				"Машина: 1, Ящик: 3, Яблоки:  [6, 7]\n" +
				"Машина: 2, Ящик: 2, Яблоки:  [3, 4, 5]\n",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := carsFormatString(tt.args.cars); got != tt.want {
				t.Errorf("carsFormatString() = %v, want %v", got, tt.want)
			}
		})
	}
}
