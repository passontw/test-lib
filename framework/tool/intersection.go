package tool

/**
 * PickFromObjectList
 * 从对象列表中挑选符合条件的列表
 *
 * @param srcList  - 原始列表
 * @param strings []string - 比较列表
 * @param keyFunc  - 返回获取比较元素
 * @return []T  - 返回值选取列表
 */

func PickFromObjectList[T any](srcList []T, strings []string, keyFunc func(obj T) string) []T {
	stringMap := make(map[string]bool)
	for _, str := range strings {
		stringMap[str] = true
	}

	var result []T
	for _, obj := range srcList {
		if stringMap[keyFunc(obj)] {
			result = append(result, obj)
		}
	}

	return result
}
