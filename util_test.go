package logger

type Object struct {
	FirstName   string `json:"firstName"`
	LastName    string `json:"lastName"`
	PIN         string `json:"pin"         mask:"pin"`
	FullName    string `json:"fullName"    mask:"name"`
	PhoneNumber string `json:"phoneNumber" mask:"phone"`
	Address     string `json:"address"     mask:"any"`
}

type Main struct {
	Object
	MapString          map[string]string  `json:"mapString"`
	MapInteger         map[int]int        `json:"mapInt"`
	MapObject          map[string]Object  `json:"mapObject"`
	MapObjectPointer   map[string]*Object `json:"mapObjectPointer"`
	SliceString        []string           `json:"sliceString"`
	SliceInteger       []int              `json:"sliceInteger"`
	SliceObject        []Object           `json:"sliceObject"`
	SliceObjectPointer []*Object          `json:"sliceObjectPointer"`
	Integer            int                `json:"integer"`
	Float              float64            `json:"float64"`
	Boolean            bool               `json:"boolean"`
}
