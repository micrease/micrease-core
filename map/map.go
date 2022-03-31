package _map

type Map struct {
	data map[string]interface{}
}

func NewMap() *Map {
	m := new(Map)
	m.data = make(map[string]interface{})
	return m
}

func (m *Map) NewMap(data map[string]interface{}) *Map {
	m.data = data
	return m
}

func (c *Map) Set(key string, value interface{}) {
	c.data[key] = value
}

func (c *Map) Get(key string) interface{} {
	if !c.Exist(key) {
		return nil
	}
	return c.data[key]
}

func (c *Map) Exist(key string) bool {
	if _, ok := c.data[key]; ok {
		return true
	}
	return false
}

func (c *Map) GetString(key string) string {
	if c.Exist(key) {
		return c.Get(key).(string)
	}
	return ""
}
