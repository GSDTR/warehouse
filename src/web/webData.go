package Web

var warehouseData []byte

func getJsonData() []byte{
	return warehouseData
}

func SetJsonData(data []byte) {
	warehouseData = data
}
