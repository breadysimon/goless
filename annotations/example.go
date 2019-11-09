package annotations

// @Api(router="/aaa/aa")
type Aaa struct {
	x string `gorm:"size=10"`
	Y int    `json:"y"`
}
