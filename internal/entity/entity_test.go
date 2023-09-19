package entity

import (
	"testing"
	"time"
)

func TestCustomTimeUnmarshalJson(t *testing.T) {

	testCases := []struct {
		input  []byte
		expect CustomTime
		err    error
	}{
		{
			[]byte(`"2023-09-15T21:09"`),
			CustomTime(time.Date(2023, 9, 15, 21, 9, 0, 0, time.UTC)),
			nil,
		},
		{
			[]byte(`"invalid-date"`),
			CustomTime{},
			&time.ParseError{},
		},
	}

	for _, test := range testCases {
		date := CustomTime{}
		err := date.UnmarshalJSON(test.input)

		if err != nil {
			if _, ok := err.(*time.ParseError); !ok {
				t.Errorf("Expected a time.ParseError, but got: %T", err)
			}
		} else if !time.Time(date).Equal(time.Time(test.expect)) {
			t.Errorf("Expected %v, but got %v", time.Time(test.expect), time.Time(date))
		}
	}

}

type NewRouteParams struct {
	id       string
	name     string
	distance float64
}

func TestNewRoute(t *testing.T) {

	testCases := []NewRouteParams{
		{
			id:       "1",
			name:     "test",
			distance: 10,
		},
		{
			id:       "15",
			name:     "new route",
			distance: 15,
		},
	}

	for _, test := range testCases {

		route := NewRoute(test.id, test.name, test.distance)

		if route.ID != test.id || route.Name != test.name || route.Distance != test.distance {
			t.Errorf("Expected %v, %v, %v, but got %v, %v, %v", test.id, test.name, test.distance, route.ID, route.Name, route.Distance)
		}
	}
}

func TestRouteStart(t *testing.T) {

	testCases := []Route{
		{
			ID:           "1",
			Name:         "test",
			Distance:     10,
			Status:       "pending",
			FreightPrice: 15,
			StartedAt:    time.Now(),
		},
		{
			ID:           "15",
			Name:         "new route",
			Distance:     15,
			Status:       "pending",
			FreightPrice: 20,
			StartedAt:    time.Now(),
		},
	}

	for _, test := range testCases {
		newDate := time.Now().Add(time.Hour * 24 * 30)
		test.Start(newDate)

		if test.Status != "started" || !test.StartedAt.Equal(newDate) {
			t.Errorf("Expected status: %v, startedAt: %v, but got %v, %v", "started", newDate, test.Status, test.StartedAt)
		}
	}
}

type RouteFinishCases struct {
	route Route
	date  time.Time
}

func TestRouteFinish(t *testing.T) {
	testCases := []RouteFinishCases{
		{
			Route{
				ID:           "1",
				Name:         "test",
				Distance:     10,
				Status:       "started",
				FreightPrice: 15,
				StartedAt:    time.Now(),
			},
			time.Now(),
		},
		{
			Route{
				ID:           "15",
				Name:         "New route",
				Distance:     15,
				Status:       "started",
				FreightPrice: 20,
				StartedAt:    time.Date(2023, 12, 24, 21, 9, 0, 0, time.UTC),
			},
			time.Date(2023, 12, 25, 8, 0, 0, 0, time.UTC),
		},
	}

	for _, test := range testCases {
		test.route.Finish(test.date)
		if test.route.Status != "finished" || !test.route.FinishedAt.Equal(test.date) {
			t.Errorf("Expected status: %v, finishedAt: %v, but got %v, %v", "finished", test.date, test.route.Status, test.route.FinishedAt)
		}
	}

}

func TestNewFreight(t *testing.T) {
	testCases := []float64{1.0, 8.5}

	for _, test := range testCases {
		freight := NewFreight(test)
		if freight.PricePerKm != test {
			t.Errorf("Expected %v, but got %v", test, freight.PricePerKm)
		}
	}
}

func TestFreightCalculate(t *testing.T) {
	testCases := []struct {
		Freight Freight
		Route   Route
		expect  float64
	}{
		{
			Freight{PricePerKm: 2.0},
			Route{Distance: 10.0},
			20.0,
		},
		{
			Freight{PricePerKm: 3.0},
			Route{Distance: 15.0},
			45.0,
		},
	}

	for _, test := range testCases {
		freight := test.Freight
		route := test.Route

		freight.Calculate(&route)

		if route.FreightPrice != test.expect {
			t.Errorf("Expected %v, but got %v", test.expect, route.FreightPrice)
		}
	}
}
