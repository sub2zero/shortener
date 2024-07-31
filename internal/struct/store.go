package store

type UrlStore struct {
	Gauge   map[string]float64
	Counter map[string]int64
}

func (m *UrlStore) Create() {
	m.Gauge = map[string]float64{}
	m.Counter = map[string]int64{}
}

func (m *UrlStore) GetAll() string {
	var res string = "Hellow all"
	// for k, v := range m.Gauge {
	// 	res += "gauge/" + k + "/" + fmt.Sprintf("%v", v) + "\n"
	// }
	// for l, m := range m.Counter {
	// 	res += "counter/" + l + "/" + strconv.FormatInt(m, 10) + "\n"
	// }
	return res
}
