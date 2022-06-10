package transformer

type Resource struct {
	Transformer Transformer
}

type Resourcer interface {
	GetTransformer() interface{}
	SetTransformer(transformer Transformer)
}

func (r *Resource) GetTransformer() interface{} {
	return r.Transformer
}

func (r *Resource) SetTransformer(transformer Transformer) {
	r.Transformer = transformer
}

type Item struct {
	Resource
	Data interface{}
}

type Collection struct {
	Resource
	Total        int
	Page         int
	ItemsPerPage int
	Data         []interface{}
}

type Error struct {
	Resource
	Data interface{}
}

func (i *Item) GetData() interface{} {
	return i.Data
}

func (i *Item) SetData(data interface{}) {
	i.Data = data
}

func (c *Collection) GetData() []interface{} {
	return c.Data
}

func (c *Collection) SetData(data []interface{}) {
	c.Data = data
}

func (c *Collection) GetTotal() int {
	return c.Total
}

func (c *Collection) SetTotal(total int) {
	c.Total = total
}

func (c *Collection) GetPage() int {
	return c.Page
}

func (c *Collection) SetPage(page int) {
	c.Page = page
}

func (c *Collection) GetItemsPerPage() int {
	return c.ItemsPerPage
}

func (c *Collection) SetItemsPerPage(itemsPerPage int) {
	c.ItemsPerPage = itemsPerPage
}

func (e *Error) GetData() interface{} {
	return e.Data
}

func (e *Error) SetData(data interface{}) {
	e.Data = data
}

func NewItem(data interface{}, transformer Transformer) Item {
	i := Item{}
	i.SetData(data)
	i.SetTransformer(transformer)
	return i
}

func NewCollection(data []interface{}, transformer Transformer) Collection {
	c := Collection{}
	c.SetData(data)
	c.SetTransformer(transformer)
	return c
}

func NewError(data interface{}) Error {
	e := Error{}
	e.SetData(data)
	e.SetTransformer(nil)
	return e
}
