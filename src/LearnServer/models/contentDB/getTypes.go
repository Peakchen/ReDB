package contentDB

// GetTypeInfo 获取某个题目的题型信息, data: 接收信息的指针
func GetTypeInfo(problemID string, subIdx int, data interface{}) error {
	db := GetDB()
	err := db.QueryRowx(`SELECT distinct p.typeName as typename, t.category, t.priority, t.chapNum as typeChapter, t.sectNum as typeSection
						 FROM probtypes as p, typenames as t
						 WHERE p.typeName = t.name and p.problemID = ? and p.subIdx = ?;`, problemID, subIdx).StructScan(data)
	return err
}
