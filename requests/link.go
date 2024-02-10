package requests

type CreateLink struct {
	Title                  string  `schema:"title,required"`
	Description            string  `schema:"description"`
	Latitude               float64 `schema:"latitude"`
	Longitude              float64 `schema:"longitude"`
	HasLocationRestriction bool    `schema:"location_restriction"`
	HasSignature           bool    `schema:"has_signature"`
}
