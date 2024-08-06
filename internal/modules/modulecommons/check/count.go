package check

var count int

func AddCheckCount() {
	count++
}

func GetCheckCount() int {
	return count
}
