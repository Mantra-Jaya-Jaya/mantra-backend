package services

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"math"
	"net/http"
	"os"
	"strconv"
	"strings"
)

type BiteshipAdapter struct {
	apiKey  string
	baseURL string
	client  *http.Client
}

func NewBiteshipAdapter() *BiteshipAdapter {
	apiKey := os.Getenv("BITESHIP_API_KEY")
	baseURL := "https://api.biteship.com"
	if os.Getenv("BITESHIP_MODE") == "production" {
		baseURL = "https://api.biteship.com"
	}

	return &BiteshipAdapter{
		apiKey:  apiKey,
		baseURL: baseURL,
		client:  &http.Client{},
	}
}

type biteshipRateRequest struct {
	OriginPostalCode      string         `json:"origin_postal_code,omitempty"`
	DestinationPostalCode string         `json:"destination_postal_code,omitempty"`
	OriginLatitude        float64        `json:"origin_latitude,omitempty"`
	OriginLongitude       float64        `json:"origin_longitude,omitempty"`
	DestinationLatitude   float64        `json:"destination_latitude,omitempty"`
	DestinationLongitude  float64        `json:"destination_longitude,omitempty"`
	Couriers              string         `json:"couriers"`
	Items                 []biteshipItem `json:"items"`
}

type biteshipItem struct {
	Name     string `json:"name"`
	Weight   int    `json:"weight"`
	Length   int    `json:"length"`
	Width    int    `json:"width"`
	Height   int    `json:"height"`
	Value    int    `json:"value"`
}

type biteshipRateResponse struct {
	Success bool              `json:"success"`
	Object  string            `json:"object"`
	Pricing []biteshipPricing `json:"pricing"`
}

type biteshipPricing struct {
	Company               string `json:"company"`
	CourierName           string `json:"courier_name"`
	CourierCode           string `json:"courier_code"`
	CourierServiceName    string `json:"courier_service_name"`
	CourierServiceCode    string `json:"courier_service_code"`
	Description           string `json:"description"`
	Duration              string `json:"duration"`
	ShipmentDurationRange string `json:"shipment_duration_range"`
	ShipmentDurationUnit  string `json:"shipment_duration_unit"`
	Price                 int    `json:"price"`
}

func (b *BiteshipAdapter) CekOngkir(req OngkirRequest) ([]OngkirResult, error) {
	originPostal := req.OriginPostal
	destPostal := req.DestPostal

	var items []biteshipItem
	for _, it := range req.Items {
		weight := it.Weight
		if weight <= 0 {
			weight = 500
		}
		items = append(items, biteshipItem{
			Name:   it.Name,
			Weight: weight,
			Length: it.Length,
			Width:  it.Width,
			Height: it.Height,
			Value:  it.Value,
		})
	}

	if len(items) == 0 {
		items = []biteshipItem{
			{Name: "Barang", Weight: 1000, Value: 0},
		}
	}

	payload := biteshipRateRequest{
		Couriers: "jne,jnt,sicepat,anteraja,ninja",
		Items:    items,
	}

	originLatStr := os.Getenv("BITESHIP_STORE_COORDINATE_LAT")
	originLngStr := os.Getenv("BITESHIP_STORE_COORDINATE_LONG")
	originLat, _ := strconv.ParseFloat(originLatStr, 64)
	originLng, _ := strconv.ParseFloat(originLngStr, 64)

	if originPostal != "" && destPostal != "" {
		payload.OriginPostalCode = originPostal
		payload.DestinationPostalCode = destPostal
	} else if originLat != 0 || originLng != 0 {
		payload.OriginLatitude = originLat
		payload.OriginLongitude = originLng
		payload.DestinationLatitude = req.DestLat
		payload.DestinationLongitude = req.DestLng
	} else if originPostal != "" {
		payload.OriginPostalCode = originPostal
		payload.DestinationPostalCode = b.getOriginPostalCode("")
	} else {
		payload.OriginPostalCode = b.getOriginPostalCode("")
		payload.DestinationPostalCode = destPostal
	}

	bodyBytes, _ := json.Marshal(payload)

	reqHttp, _ := http.NewRequest("POST", b.baseURL+"/v1/rates/couriers", strings.NewReader(string(bodyBytes)))
	reqHttp.Header.Set("Content-Type", "application/json")
	reqHttp.Header.Set("Authorization", b.apiKey)

	resp, err := b.client.Do(reqHttp)
	if err != nil {
		return nil, fmt.Errorf("biteship request failed: %w", err)
	}
	defer resp.Body.Close()

	respBody, _ := io.ReadAll(resp.Body)

	if resp.StatusCode != http.StatusOK {
		errMsg := string(respBody)

		if strings.Contains(errMsg, "No sufficient balance") || strings.Contains(errMsg, "insufficient") || strings.Contains(errMsg, "top up") {
			return getMockOngkir(req), nil
		}

		return nil, fmt.Errorf("biteship error status %d: %s", resp.StatusCode, errMsg)
	}

	var rateResp biteshipRateResponse
	if err := json.Unmarshal(respBody, &rateResp); err != nil {
		return nil, fmt.Errorf("biteship parse error: %w", err)
	}

	if !rateResp.Success {
		return nil, fmt.Errorf("biteship response not successful")
	}

	if rateResp.Pricing == nil {
		return []OngkirResult{}, nil
	}

	ekspedisiMap := make(map[string]*OngkirResult)
	for _, p := range rateResp.Pricing {
		eks, exists := ekspedisiMap[p.CourierCode]
		if !exists {
			eks = &OngkirResult{
				EkspedisiKode: p.CourierCode,
				NamaEkspedisi: p.CourierName,
			}
			ekspedisiMap[p.CourierCode] = eks
		}

		estimasiMin, estimasiMax := parseDuration(p.ShipmentDurationRange)

		eks.Layanan = append(eks.Layanan, OngkirLayanan{
			KodeLayanan: p.CourierServiceCode,
			NamaLayanan: p.CourierServiceName,
			Deskripsi:   p.Description,
			Harga:       p.Price,
			EstimasiMin: estimasiMin,
			EstimasiMax: estimasiMax,
			Durasi:      p.Duration,
		})
	}

	var results []OngkirResult
	for _, v := range ekspedisiMap {
		results = append(results, *v)
	}

	return results, nil
}

func (b *BiteshipAdapter) getOriginPostalCode(fallback string) string {
	if fallback != "" {
		return fallback
	}
	return os.Getenv("BITESHIP_STORE_POSTAL_CODE")
}

func parseDuration(rangeStr string) (min, max int) {
	if rangeStr == "" {
		return 0, 0
	}
	parts := strings.Split(rangeStr, " - ")
	if len(parts) != 2 {
		return 0, 0
	}
	min, _ = strconv.Atoi(strings.TrimSpace(parts[0]))
	max, _ = strconv.Atoi(strings.TrimSpace(parts[1]))
	return
}

type BiteshipTrackingEvent struct {
	Time        string `json:"time"`
	Description string `json:"description"`
	Status      string `json:"status"`
	City        string `json:"city"`
}

type BiteshipTrackingResponse struct {
	Success bool                    `json:"success"`
	History []BiteshipTrackingEvent `json:"history"`
}

func (b *BiteshipAdapter) TrackShipment(waybill string, courierCode string) ([]BiteshipTrackingEvent, error) {
	reqHttp, _ := http.NewRequest("GET", b.baseURL+fmt.Sprintf("/v1/trackings/%s/couriers/%s", waybill, courierCode), nil)
	reqHttp.Header.Set("Authorization", b.apiKey)

	resp, err := b.client.Do(reqHttp)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		respBody, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("biteship error status %d: %s", resp.StatusCode, string(respBody))
	}

	var result struct {
		Success bool                    `json:"success"`
		History []BiteshipTrackingEvent `json:"history"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}
	return result.History, nil
}

type BiteshipCourierItem struct {
	CourierCode        string `json:"courier_code"`
	CourierName        string `json:"courier_name"`
	CourierServiceName string `json:"courier_service_name"`
	CourierServiceCode string `json:"courier_service_code"`
	Description        string `json:"description"`
}

type BiteshipCouriersResponse struct {
	Success  bool                  `json:"success"`
	Message  string                `json:"message"`
	Couriers []BiteshipCourierItem `json:"couriers"`
}

func (b *BiteshipAdapter) GetCouriers() ([]BiteshipCourierItem, error) {
	reqHttp, _ := http.NewRequest("GET", b.baseURL+"/v1/couriers", nil)
	reqHttp.Header.Set("Authorization", b.apiKey)

	resp, err := b.client.Do(reqHttp)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	respBody, _ := io.ReadAll(resp.Body)
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("biteship error status %d: %s", resp.StatusCode, string(respBody))
	}

	var result BiteshipCouriersResponse
	if err := json.Unmarshal(respBody, &result); err != nil {
		return nil, fmt.Errorf("failed to parse biteship couriers: %w", err)
	}

	if !result.Success {
		return nil, fmt.Errorf("biteship couriers fetch unsuccessful: %s", result.Message)
	}

	return result.Couriers, nil
}

// CreateShipmentRequest is the request payload for creating a Biteship order
type CreateShipmentRequest struct {
	OriginAddress       string  `json:"origin_address"`
	OriginCoordinate    string  `json:"origin_coordinate"`
	DestinationAddress  string  `json:"destination_address"`
	DestinationCoordinate string `json:"destination_coordinate"`
	CourierCode         string  `json:"courier_code"`
	CourierServiceCode  string  `json:"courier_service_code"`
	Items               []CreateShipmentItem `json:"items"`
}

type CreateShipmentItem struct {
	Name          string `json:"name"`
	Weight        int    `json:"weight"`
	Quantity      int    `json:"quantity"`
	Value         int    `json:"value"`
	Length        int    `json:"length,omitempty"`
	Width         int    `json:"width,omitempty"`
	Height        int    `json:"height,omitempty"`
}

type biteshipCreateOrderRequest struct {
	Shipper struct {
		Name    string `json:"name"`
		Phone   string `json:"phone"`
		Address string `json:"address"`
		Coordinate string `json:"coordinate"`
	} `json:"shipper"`
	Destination struct {
		Name         string `json:"name"`
		Phone        string `json:"phone"`
		Address      string `json:"address"`
		Coordinate   string `json:"coordinate"`
	} `json:"destination"`
	Courier struct {
		CourierCode string `json:"courier_code"`
		ServiceCode string `json:"service_code"`
	} `json:"courier"`
	Items []struct {
		Name     string `json:"name"`
		Weight   int    `json:"weight"`
		Quantity int    `json:"quantity"`
		Value    int    `json:"value"`
		Length   int    `json:"length,omitempty"`
		Width    int    `json:"width,omitempty"`
		Height   int    `json:"height,omitempty"`
	} `json:"items"`
}

type biteshipCreateOrderResponse struct {
	Success   bool   `json:"success"`
	Message   string `json:"message"`
	ID        string `json:"id"`
	WaybillID string `json:"waybill_id"`
	Status    string `json:"status"`
	Courier   struct {
		Company string `json:"company"`
		Name    string `json:"name"`
		Phone   string `json:"phone"`
	} `json:"courier"`
}

// CreateShipment creates a shipment order on Biteship and returns waybill ID
func (b *BiteshipAdapter) CreateShipment(req CreateShipmentRequest) (*biteshipCreateOrderResponse, error) {
	orderReq := biteshipCreateOrderRequest{}
	orderReq.Shipper.Name = os.Getenv("BITESHIP_STORE_NAME")
	orderReq.Shipper.Phone = os.Getenv("BITESHIP_STORE_PHONE")
	orderReq.Shipper.Address = req.OriginAddress
	orderReq.Shipper.Coordinate = req.OriginCoordinate

	orderReq.Destination.Name = ""
	orderReq.Destination.Phone = ""
	orderReq.Destination.Address = req.DestinationAddress
	orderReq.Destination.Coordinate = req.DestinationCoordinate

	orderReq.Courier.CourierCode = req.CourierCode
	orderReq.Courier.ServiceCode = req.CourierServiceCode

	for _, it := range req.Items {
		orderReq.Items = append(orderReq.Items, struct {
			Name     string `json:"name"`
			Weight   int    `json:"weight"`
			Quantity int    `json:"quantity"`
			Value    int    `json:"value"`
			Length   int    `json:"length,omitempty"`
			Width    int    `json:"width,omitempty"`
			Height   int    `json:"height,omitempty"`
		}{
			Name:     it.Name,
			Weight:   it.Weight,
			Quantity: it.Quantity,
			Value:    it.Value,
			Length:   it.Length,
			Width:    it.Width,
			Height:   it.Height,
		})
	}

	bodyBytes, _ := json.Marshal(orderReq)
	reqHttp, _ := http.NewRequest("POST", b.baseURL+"/v1/orders", bytes.NewReader(bodyBytes))
	reqHttp.Header.Set("Content-Type", "application/json")
	reqHttp.Header.Set("Authorization", b.apiKey)

	resp, err := b.client.Do(reqHttp)
	if err != nil {
		return nil, fmt.Errorf("biteship create order request failed: %w", err)
	}
	defer resp.Body.Close()

	respBody, _ := io.ReadAll(resp.Body)
	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		return nil, fmt.Errorf("biteship create order error status %d: %s", resp.StatusCode, string(respBody))
	}

	var result biteshipCreateOrderResponse
	if err := json.Unmarshal(respBody, &result); err != nil {
		return nil, fmt.Errorf("biteship create order parse error: %w", err)
	}

	return &result, nil
}

// GetOrderStatus mengambil status order Biteship berdasarkan order ID.
// Dipakai untuk verifikasi / polling status kalau webhook gagal.
func (b *BiteshipAdapter) GetOrderStatus(orderID string) (*biteshipCreateOrderResponse, error) {
	reqHttp, _ := http.NewRequest("GET", b.baseURL+fmt.Sprintf("/v1/orders/%s", orderID), nil)
	reqHttp.Header.Set("Authorization", b.apiKey)

	resp, err := b.client.Do(reqHttp)
	if err != nil {
		return nil, fmt.Errorf("biteship get order failed: %w", err)
	}
	defer resp.Body.Close()

	respBody, _ := io.ReadAll(resp.Body)
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("biteship order error status %d: %s", resp.StatusCode, string(respBody))
	}

	var result biteshipCreateOrderResponse
	if err := json.Unmarshal(respBody, &result); err != nil {
		return nil, fmt.Errorf("biteship order parse error: %w", err)
	}

	return &result, nil
}

// CancelOrder membatalkan order Biteship.
// Dipakai jika shipment perlu dibatalkan sebelum dikirim.
func (b *BiteshipAdapter) CancelOrder(orderID string) error {
	reqHttp, _ := http.NewRequest("POST", b.baseURL+fmt.Sprintf("/v1/orders/%s/cancel", orderID), nil)
	reqHttp.Header.Set("Authorization", b.apiKey)

	resp, err := b.client.Do(reqHttp)
	if err != nil {
		return fmt.Errorf("biteship cancel order failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusAccepted {
		respBody, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("biteship cancel order error status %d: %s", resp.StatusCode, string(respBody))
	}

	return nil
}

func getMockOngkir(req OngkirRequest) []OngkirResult {
	totalWeight := 0
	for _, it := range req.Items {
		w := it.Weight
		if w <= 0 {
			w = 500
		}
		totalWeight += w
	}

	// Estimasi jarak dari toko (Polines) ke tujuan
	jarak := 10
	if req.DestLat != 0 {
		d := Haversine(-7.046389, 110.438333, req.DestLat, req.DestLng)
		jarak = int(math.Round(d))
	}
	if jarak < 1 {
		jarak = 5
	}

	type mockLayanan struct {
		kode string
		nama string
		desc string
		min  int
		max  int
	}

	type mockKurir struct {
		kode    string
		nama    string
		layanan []mockLayanan
	}

	kurirList := []mockKurir{
		{
			kode: "jne", nama: "JNE",
			layanan: []mockLayanan{
				{"jne_reg", "REG", "Reguler", 2, 4},
				{"jne_oke", "OKE", "Ongkos Kirim Ekonomis", 3, 6},
				{"jne_yes", "YES", "Yakin Esok Sampai", 1, 1},
			},
		},
		{
			kode: "jnt", nama: "J&T",
			layanan: []mockLayanan{
				{"jnt_reg", "Reguler", "Reguler", 2, 4},
				{"jnt_ez", "EZ", "Economy", 4, 7},
			},
		},
		{
			kode: "sicepat", nama: "SiCepat",
			layanan: []mockLayanan{
				{"sicepat_best", "BEST", "Besok Sampai Tujuan", 1, 1},
				{"sicepat_reg", "REG", "Reguler", 2, 3},
			},
		},
		{
			kode: "anteraja", nama: "AnterAja",
			layanan: []mockLayanan{
				{"anteraja_reg", "Reguler", "Reguler", 2, 4},
			},
		},
		{
			kode: "ninja", nama: "Ninja Xpress",
			layanan: []mockLayanan{
				{"ninja_reg", "Reguler", "Reguler", 2, 4},
				{"ninja_express", "Express", "Express", 1, 2},
			},
		},
	}

	var results []OngkirResult
	for _, k := range kurirList {
		var layananList []OngkirLayanan
		for _, l := range k.layanan {
			harga := (5000 + totalWeight*50) * jarak / 10
			if harga < 8000 {
				harga = 8000
			}
			layananList = append(layananList, OngkirLayanan{
				KodeLayanan: l.kode,
				NamaLayanan: l.nama,
				Deskripsi:   l.desc,
				Harga:       harga,
				EstimasiMin: l.min,
				EstimasiMax: l.max,
				Durasi:      fmt.Sprintf("%d - %d hari", l.min, l.max),
			})
		}
		results = append(results, OngkirResult{
			EkspedisiKode: k.kode,
			NamaEkspedisi: k.nama,
			Layanan:       layananList,
		})
	}

	return results
}

func Haversine(lat1, lon1, lat2, lon2 float64) float64 {
	const R = 6371
	dLat := (lat2 - lat1) * math.Pi / 180
	dLon := (lon2 - lon1) * math.Pi / 180
	a := math.Sin(dLat/2)*math.Sin(dLat/2) +
		math.Cos(lat1*math.Pi/180)*math.Cos(lat2*math.Pi/180)*
			math.Sin(dLon/2)*math.Sin(dLon/2)
	c := 2 * math.Atan2(math.Sqrt(a), math.Sqrt(1-a))
	return R * c
}
