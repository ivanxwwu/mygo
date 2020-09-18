package data

var datas []string
var datas2 []string
var Len1 int64
var Len2 int64

func Add(str string) string {
	data := []byte(str)
	sData := string(data)
	Len1 += int64(len(str))
	datas = append(datas, sData)

	return sData
}

func Add2(str string) string {
	return Add3(str)
}

func Add3(str string) string {
	data := []byte(str)
	sData := string(data)
	Len2 += int64(len(str))
	datas2 = append(datas2, sData)

	return sData
}

