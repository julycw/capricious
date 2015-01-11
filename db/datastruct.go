package db

import ()

type DataStruct map[string]interface{}

func NewDataStruct(data map[string]interface{}) DataStruct {
	var size int = 0
	if data != nil {
		size = len(data)
	}
	d := make(DataStruct, size)

	if data != nil {
		for key, value := range data {
			d[key] = value
		}
	}

	return d
}

func clearData(data *DataStruct) {
	delete(*data, "_uuid")
	delete(*data, "_update_at")
	delete(*data, "_insert_at")
}
