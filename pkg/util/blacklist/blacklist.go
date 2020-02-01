package blacklist

import (
	"io/ioutil"

	yaml "gopkg.in/yaml.v2"
)

type Blacklist struct {
	Enabled         bool             `yaml:"enabled"`
	List            []*BlacklistItem `yaml:"blacklist"`
	listBySourceUrl map[string]*BlacklistItem
	listByAliasUrl  map[string]*BlacklistItem
}

type BlacklistItem struct {
	Source string
	Alias  string
}

func NewBlacklist(path string) (*Blacklist, error) {
	bl := &Blacklist{
		List:            []*BlacklistItem{},
		listBySourceUrl: make(map[string]*BlacklistItem),
		listByAliasUrl:  make(map[string]*BlacklistItem),
	}
	if len(path) > 0 {
		err := bl.LoadFromFile(path)
		if err != nil {
			return nil, err
		}
		bl.FillLists()
	}
	return bl, nil
}

func (b *Blacklist) FillLists() {
	for _, item := range b.List {
		b.listBySourceUrl[item.Source] = item
		b.listByAliasUrl[item.Alias] = item
	}
}

func (b *Blacklist) LoadFromFile(path string) error {
	configBytes, err := ioutil.ReadFile(path)
	if err != nil {
		return err
	}
	var bl Blacklist
	if err := yaml.Unmarshal(configBytes, &bl); err != nil {
		return err
	}
	b.List = bl.List
	b.Enabled = bl.Enabled
	return nil
}

func (b *Blacklist) IsAlias(url string) bool {
	if len(b.List) == 0 {
		return false
	}
	_, ok := b.listByAliasUrl[url]
	return ok
}

func (b *Blacklist) IsBlocked(url string) bool {
	if len(b.List) == 0 {
		return false
	}
	_, ok := b.listBySourceUrl[url]
	return ok
}

func (b *Blacklist) GetAlias(url string) string {
	if len(b.List) == 0 {
		return url
	}
	item, ok := b.listBySourceUrl[url]
	if ok {
		return item.Alias
	}
	return url
}

func (b *Blacklist) GetSource(url string) string {
	if len(b.List) == 0 {
		return url
	}
	item, ok := b.listByAliasUrl[url]
	if ok {
		return item.Source
	}
	return url
}
