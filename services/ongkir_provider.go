package services

type OngkirRequest struct {
	OriginLat    float64
	OriginLng    float64
	OriginPostal string
	DestLat      float64
	DestLng      float64
	DestPostal   string
	Items        []OngkirItem
}

type OngkirItem struct {
	Name   string
	Weight int
	Length int
	Width  int
	Height int
	Value  int
}

type OngkirResult struct {
	EkspedisiKode string
	NamaEkspedisi string
	Layanan       []OngkirLayanan
}

type OngkirLayanan struct {
	KodeLayanan string
	NamaLayanan string
	Deskripsi   string
	Harga       int
	EstimasiMin int
	EstimasiMax int
	Durasi      string
}

type OngkirProvider interface {
	CekOngkir(req OngkirRequest) ([]OngkirResult, error)
}
