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
		// TODO: Add test cases.
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
	tests := []struct {
		name    string
		args    args
		want    []Car
		wantErr bool
	}{
		// TODO: Add test cases.
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

func Test_carsFormatString(t *testing.T) {
	type args struct {
		cars []Car
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := carsFormatString(tt.args.cars); got != tt.want {
				t.Errorf("carsFormatString() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_packageBagsToBoxes(t *testing.T) {
	type args struct {
		bags                []Bag
		numberOfApplesInBox uint32
	}
	tests := []struct {
		name    string
		args    args
		want    []Box
		wantErr bool
	}{
		// TODO: Add test cases.
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
