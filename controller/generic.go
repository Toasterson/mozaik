package controller


type ListController struct {
	Controller
	ItemList interface{}
	ItemName string
	//Storer 		StoreInterface
}

func (this *ListController) Get() error {
	this.Data[this.ItemName] = this.ItemList
	return this.Render()
}