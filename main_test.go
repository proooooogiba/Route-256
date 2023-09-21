package main

import (
	"reflect"
	"testing"
)

func Test_applesToBag(t *testing.T) {
	type args struct {
		appleNumber uint32
	}
	apple1 := Apple(1)
	apple2 := Apple(2)
	apple3 := Apple(3)
	apple4 := Apple(4)
	apple5 := Apple(5)
	apple6 := Apple(6)
	apple7 := Apple(7)
	apple8 := Apple(8)
	apple9 := Apple(9)
	apple10 := Apple(10)

	tests := []struct {
		name    string
		args    args
		want    []*Apple
		wantErr bool
	}{
		{
			name:    "input number - 5",
			args:    struct{ appleNumber uint32 }{appleNumber: 5},
			want:    []*Apple{&apple1, &apple2, &apple3, &apple4, &apple5},
			wantErr: false,
		},
		{
			name:    "invalid input - 0",
			args:    struct{ appleNumber uint32 }{appleNumber: 0},
			want:    []*Apple{},
			wantErr: false,
		},
		{
			name: "input number - 10",
			args: struct{ appleNumber uint32 }{appleNumber: 10},
			want: []*Apple{
				&apple1, &apple2, &apple3,
				&apple4, &apple5, &apple6,
				&apple7, &apple8, &apple9,
				&apple10,
			},
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
		boxes        []*Box
		numberOfCars int
	}
	apple1 := Apple(1)
	apple2 := Apple(2)
	apple3 := Apple(3)
	apple4 := Apple(4)
	apple5 := Apple(5)
	apple6 := Apple(6)
	apple7 := Apple(7)
	apple8 := Apple(8)
	apple9 := Apple(9)

	boxes := []*Box{
		{Bags: []*Apple{&apple1, &apple2, &apple3}},
		{Bags: []*Apple{&apple3, &apple4, &apple5}},
	}

	cars := []*Car{
		{Boxes: []*Box{
			{Bags: []*Apple{&apple1, &apple2, &apple3}},
		}},
		{Boxes: []*Box{
			{Bags: []*Apple{&apple3, &apple4, &apple5}},
		}},
	}

	boxesWithDifferentNumberOfBags := []*Box{
		{Bags: []*Apple{&apple1, &apple2, &apple3}},
		{Bags: []*Apple{&apple4, &apple5, &apple6}},
		{Bags: []*Apple{&apple7, &apple8, &apple9}},
	}

	carsWithExtraBoxes := []*Car{
		{Boxes: []*Box{
			{Bags: []*Apple{&apple1, &apple2, &apple3}},
			{Bags: []*Apple{&apple7, &apple8, &apple9}},
		}},
		{Boxes: []*Box{
			{Bags: []*Apple{&apple4, &apple5, &apple6}},
		}},
	}

	tests := []struct {
		name    string
		args    args
		want    []*Car
		wantErr bool
	}{
		{
			name: "valid input packed to 2 cars",
			args: struct {
				boxes        []*Box
				numberOfCars int
			}{boxes: boxes, numberOfCars: 2},
			want:    cars,
			wantErr: false,
		},
		{
			name: "invalid input - number of cars is zero",
			args: struct {
				boxes        []*Box
				numberOfCars int
			}{boxes: nil, numberOfCars: 0},
			want:    nil,
			wantErr: true,
		},
		{
			name: "empty input packed to 2 cars",
			args: struct {
				boxes        []*Box
				numberOfCars int
			}{boxes: boxesWithDifferentNumberOfBags, numberOfCars: 2},
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
		bags                []*Apple
		numberOfApplesInBox int
	}

	apple1 := Apple(1)
	apple2 := Apple(2)
	apple3 := Apple(3)
	apple4 := Apple(4)

	bags := []*Apple{&apple1, &apple2, &apple3, &apple4}
	boxesTwoApples := []*Box{
		{Bags: []*Apple{&apple4, &apple3}},
		{Bags: []*Apple{&apple2, &apple1}},
	}
	boxesThreeApples := []*Box{
		{Bags: []*Apple{
			&apple4,
			&apple3,
			&apple2,
		}},
		{Bags: []*Apple{
			&apple1,
		}},
	}

	tests := []struct {
		name    string
		args    args
		want    []*Box
		wantErr bool
	}{
		{
			name: "valid input - 2 apples in box",
			args: struct {
				bags                []*Apple
				numberOfApplesInBox int
			}{bags: bags, numberOfApplesInBox: 2},
			want:    boxesTwoApples,
			wantErr: false,
		},
		{
			name: "invalid input - number of apples in bag is zero",
			args: struct {
				bags                []*Apple
				numberOfApplesInBox int
			}{bags: bags, numberOfApplesInBox: 0},
			want:    nil,
			wantErr: true,
		},
		{
			name: "valid input - 3 apples in box",
			args: struct {
				bags                []*Apple
				numberOfApplesInBox int
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
		cars []*Car
	}

	apple1 := Apple(1)
	apple2 := Apple(2)
	apple3 := Apple(3)
	apple4 := Apple(4)
	apple5 := Apple(5)
	apple6 := Apple(6)
	apple7 := Apple(7)

	cars := []*Car{
		{Boxes: []*Box{
			{Bags: []*Apple{&apple1, &apple2, &apple3}},
		}},
		{Boxes: []*Box{
			{Bags: []*Apple{&apple3, &apple4, &apple5}},
		}},
	}

	carsWithDiffrentNumbersOfBoxes := []*Car{
		{Boxes: []*Box{
			{Bags: []*Apple{&apple1, &apple2, &apple3}},
			{Bags: []*Apple{&apple6, &apple7}},
		}},
		{Boxes: []*Box{
			{Bags: []*Apple{&apple3, &apple4, &apple5}},
		}},
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "2 cars with 1 box with 3 bags in each",
			args: struct{ cars []*Car }{cars: cars},
			want: "Машина: 1, Ящик: 1, Яблоки:  [1, 2, 3]\n" +
				"Машина: 2, Ящик: 2, Яблоки:  [3, 4, 5]\n",
		},
		{
			name: "empty input",
			args: struct{ cars []*Car }{cars: nil},
			want: "",
		},
		{
			name: "cars with different numbers of boxes in each",
			args: struct{ cars []*Car }{cars: carsWithDiffrentNumbersOfBoxes},
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
