package database

import (
	"io"
	"net/http"
)

type Scramble struct {
	Instruction string
	ShortName   string
	ShowName    string
}

func Tnoodle(scramble string) (s Scramble) {
	switch scramble {
	case "2":
	case "222":
		s = Scramble{"2", "222", "2阶"}
	case "333":
	case "3":
		s = Scramble{"3", "333", "3阶"}
	case "444":
	case "4":
		s = Scramble{"4", "444", "4阶"}
	case "555":
	case "5":
		s = Scramble{"5", "555", "5阶"}
	case "666":
	case "6":
		s = Scramble{"6", "666", "6阶"}
	case "777":
	case "7":
		s = Scramble{"7", "777", "7阶"}
	case "py":
	case "pyram":
		s = Scramble{"py", "pyram", "pyram"}
	case "skewb":
	case "sk":
		s = Scramble{"sk", "skewb", "skewb"}
	case "sq1":
	case "sq":
		s = Scramble{"sq", "sq1", "sq1"}
	case "clock":
	case "cl":
		s = Scramble{"cl", "clock", "clock"}
	case "minx":
	case "mx":
		s = Scramble{"mx", "minx", "minx"}
	case "333fm":
	case "fm":
		s = Scramble{"fm", "333fm", "333fm"}
	default:
		s = Scramble{"instruction", "shortName", "showName"}
	}
	return
}

func GetScramble(s string) string {
	resp, err := http.Get("http://localhost:2014/scramble/.txt?=" + s)
	if err != nil {
		return "获取失败"
	}
	v, _ := io.ReadAll(resp.Body)
	return string(v)
}
