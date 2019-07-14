package machine

type Item struct {
	Name  string
	Price int
}

type Inventory struct {
	Item
	Stock int
}
