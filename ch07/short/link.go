package short

type Link struct {
	Key string
	URL string
}

var Store = make(map[string]string)

func Set(link *Link) {
	Store[link.Key] = link.URL
}

func Get(key string) *Link {

	if value, ok := Store[key]; ok {
		return &Link{
			Key: key,
			URL: value,
		}
	}
	return nil
}
