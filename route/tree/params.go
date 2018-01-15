package tree

func (ps params) ByName(name string) string {
	for _, p := range ps {
		if p.key == name {
			return p.value
		}
	}

	return ""
}

func (ps params) asMap() (ret map[string]string) {
	ret = make(map[string]string)
	for _, p := range ps {
		ret[p.key] = p.value
	}
	return
}
