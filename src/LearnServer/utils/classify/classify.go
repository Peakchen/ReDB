package classify

// List 等待归类的元素构成的列表
type List interface {
	// InSameClass 给定两个元素，返回两个元素是否在同一类
	InSameClass(int, int) bool
	// Length 得到List长度
	Length() int
}

// Classify 进行归类，返回的列表中，每个元素是归为一类的元素下标构成的列表
func Classify(toClassify List) [][]int {
	resultIndex := [][]int{}
	length := toClassify.Length()
	for i := 0; i < length; i++ {
		foundClass := false
		for ri, class := range resultIndex {
			if toClassify.InSameClass(i, class[0]) {
				resultIndex[ri] = append(resultIndex[ri], i)
				foundClass = true
				break
			}
		}
		if !foundClass {
			resultIndex = append(resultIndex, []int{i})
		}
	}

	return resultIndex
}
