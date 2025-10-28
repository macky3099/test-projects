package respository

type Seat struct {
	ID       string
	Type     string // A, B, C
	Price    float64
	Reserved bool
}

type Sales struct {
	Revenue          float64
	ServiceTax       float64 `json:"Service Tax"`
	SwachBharatCess  float64 `json:"Swachh Bharat Cess"`
	KrishiKalyanCess float64 `json:"Krishi Kalyan Cess"`
}

var (
	AvailableSeats = map[int][]Seat{
		1: {
			{ID: "A1", Type: "A", Price: 320.0, Reserved: false},
			{ID: "A2", Type: "A", Price: 320.0, Reserved: false},
			{ID: "A3", Type: "A", Price: 320.0, Reserved: false},
			{ID: "A4", Type: "A", Price: 320.0, Reserved: false},
			{ID: "A5", Type: "A", Price: 320.0, Reserved: false},
			{ID: "A6", Type: "A", Price: 320.0, Reserved: false},
			{ID: "A7", Type: "A", Price: 320.0, Reserved: false},
			{ID: "A8", Type: "A", Price: 320.0, Reserved: false},
			{ID: "A9", Type: "A", Price: 320.0, Reserved: false},
			{ID: "B1", Type: "B", Price: 280.0, Reserved: false},
			{ID: "B2", Type: "B", Price: 280.0, Reserved: false},
			{ID: "B3", Type: "B", Price: 280.0, Reserved: false},
			{ID: "B4", Type: "B", Price: 280.0, Reserved: false},
			{ID: "B5", Type: "B", Price: 280.0, Reserved: false},
			{ID: "B6", Type: "B", Price: 280.0, Reserved: false},
			{ID: "C1", Type: "C", Price: 0, Reserved: true},
			{ID: "C2", Type: "C", Price: 240.0, Reserved: false},
			{ID: "C3", Type: "C", Price: 240.0, Reserved: false},
			{ID: "C4", Type: "C", Price: 240.0, Reserved: false},
			{ID: "C5", Type: "C", Price: 240.0, Reserved: false},
			{ID: "C6", Type: "C", Price: 240.0, Reserved: false},
			{ID: "C7", Type: "C", Price: 240.0, Reserved: false},
		},
		2: {
			{ID: "A1", Type: "A", Price: 320.0, Reserved: false},
			{ID: "A2", Type: "A", Price: 320.0, Reserved: false},
			{ID: "A3", Type: "A", Price: 320.0, Reserved: false},
			{ID: "A4", Type: "A", Price: 320.0, Reserved: false},
			{ID: "A5", Type: "A", Price: 320.0, Reserved: false},
			{ID: "A6", Type: "A", Price: 320.0, Reserved: false},
			{ID: "A7", Type: "A", Price: 320.0, Reserved: false},
			{ID: "B1", Type: "B", Price: 0, Reserved: true},
			{ID: "B2", Type: "B", Price: 280.0, Reserved: false},
			{ID: "B3", Type: "B", Price: 280.0, Reserved: false},
			{ID: "B4", Type: "B", Price: 280.0, Reserved: false},
			{ID: "B5", Type: "B", Price: 280.0, Reserved: false},
			{ID: "B6", Type: "B", Price: 280.0, Reserved: false},
			{ID: "C1", Type: "C", Price: 240.0, Reserved: false},
			{ID: "C2", Type: "C", Price: 240.0, Reserved: false},
			{ID: "C3", Type: "C", Price: 240.0, Reserved: false},
			{ID: "C4", Type: "C", Price: 240.0, Reserved: false},
			{ID: "C5", Type: "C", Price: 240.0, Reserved: false},
			{ID: "C6", Type: "C", Price: 240.0, Reserved: false},
			{ID: "C7", Type: "C", Price: 240.0, Reserved: false},
		},
		3: {
			{ID: "A1", Type: "A", Price: 320.0, Reserved: false},
			{ID: "A2", Type: "A", Price: 320.0, Reserved: false},
			{ID: "A3", Type: "A", Price: 320.0, Reserved: false},
			{ID: "A4", Type: "A", Price: 320.0, Reserved: false},
			{ID: "A5", Type: "A", Price: 320.0, Reserved: false},
			{ID: "A6", Type: "A", Price: 320.0, Reserved: false},
			{ID: "A7", Type: "A", Price: 320.0, Reserved: false},
			{ID: "B1", Type: "B", Price: 280.0, Reserved: false},
			{ID: "B2", Type: "B", Price: 280.0, Reserved: false},
			{ID: "B3", Type: "B", Price: 280.0, Reserved: false},
			{ID: "B4", Type: "B", Price: 280.0, Reserved: false},
			{ID: "B5", Type: "B", Price: 280.0, Reserved: false},
			{ID: "B6", Type: "B", Price: 280.0, Reserved: false},
			{ID: "B7", Type: "B", Price: 280.0, Reserved: false},
			{ID: "B8", Type: "B", Price: 280.0, Reserved: false},
			{ID: "C1", Type: "C", Price: 240.0, Reserved: false},
			{ID: "C2", Type: "C", Price: 240.0, Reserved: false},
			{ID: "C3", Type: "C", Price: 240.0, Reserved: false},
			{ID: "C4", Type: "C", Price: 240.0, Reserved: false},
			{ID: "C5", Type: "C", Price: 240.0, Reserved: false},
			{ID: "C6", Type: "C", Price: 240.0, Reserved: false},
			{ID: "C7", Type: "C", Price: 240.0, Reserved: false},
			{ID: "C8", Type: "C", Price: 240.0, Reserved: false},
			{ID: "C9", Type: "C", Price: 240.0, Reserved: false},
		},
	}

	Taxes = map[string]float32{
		"Service Tax":        14,
		"Swachh Bharat Cess": 0.5,
		"Krishi Kalyan Cess": 0.5,
	}

	TotalSales = []Sales{}
)
