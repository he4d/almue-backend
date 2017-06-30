package model

//Floor represents the database object of a floor
type Floor struct {
	Base
	Description *string `json:"description"`
}

//DeepCopy creates a deep copy of a Lighting
func (f *Floor) DeepCopy() *Floor {
	if f == nil {
		return nil
	}
	descr := *f.Description
	copy := &Floor{
		Base:        f.Base,
		Description: &descr,
	}
	return copy
}
