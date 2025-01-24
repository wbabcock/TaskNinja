package main

func getPriorityIndex(value string) int {
	p := 2
	switch value {
	case "L":
		p = 1
	case "":
		p = 2
	case "M":
		p = 3
	case "H":
		p = 4
	}
	return p
}

func getPriorityValue(value uint64) string {
	p := ""
	switch value {
	case 1:
		p = "L"
	case 2:
		p = ""
	case 3:
		p = "M"
	case 4:
		p = "H"
	}
	return p
}
