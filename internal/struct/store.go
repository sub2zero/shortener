package store

type UrlShort struct {
	Url      string
	ShortUrl string
}

type UrlStore struct {
	url []UrlShort
}

func (m *UrlStore) Create() {
	m.url = []UrlShort{}
}

func (m *UrlStore) Add(mUrl string, mShortUrl string) {
	m.url = append(m.url, UrlShort{
		Url:      mUrl,
		ShortUrl: mShortUrl,
	})
}
func (m *UrlStore) AddUrl(mUrl string, mShortUrl string) error {

	m.Add(mUrl, mShortUrl)
	// log.Printf("saved gauge:/%v/%v\n", mUrl, *mShortUrl)
	// m.url[mUrl] = *mShortUrl
	return nil
}

// func (m *UrlStore) Update(mUrl string, mShortUrl string) error {
// 	log.Printf("saved: /%v/%v\n", mUrl, mShortUrl)

// 	m.url[mUrl] = mShortUrl

//		return nil
//	}
func (m *UrlStore) GetAll() string {
	var res string
	for _, v := range m.url {
		res += "url:" + v.Url + " shotUrl" + v.ShortUrl + "\n"
	}
	// for l, m := range m.Counter {
	// 	res += "counter/" + l + "/" + strconv.FormatInt(m, 10) + "\n"
	// }
	return res
}
