package priority

type Priority uint8

const (
	Low Priority = iota + 1
	Normal
	Medium
	High
)

func New(value uint8) Priority {
	switch value {
	case 1:
		return Low
	case 2:
		return Normal
	case 3:
		return Medium
	case 4:
		return High
	default:
		return Normal
	}
}

func GetIndex(value string) Priority {
	switch value {
	case "L":
		return Low
	case "":
		return Normal
	case "M":
		return Medium
	case "H":
		return High
	default:
		return Normal
	}
}

func (p Priority) String() string {
	switch p {
	case Low:
		return "L"
	case Normal:
		return ""
	case Medium:
		return "M"
	case High:
		return "H"
	default:
		return ""
	}
}
