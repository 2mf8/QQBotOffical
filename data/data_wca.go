package database

type Info struct {
	Msg                  string `json:"msg"`
	PageLast             bool   `json:"pageLast"`
	PageEmpty            bool   `json:"pageEmpty"`
	Data                 []PeopleInfo
	TotalPages           int  `json:"totalPages"`
	PageFirst            bool `json:"pageFirst"`
	PageSize             int  `json:"pageSize"`
	PageNumberOfElements int  `json:"pageNumberOfElements"`
	PageNum              int  `json:"pageNum"`
	Retcode              int  `json:"retcode"`
	TotalElements        int  `json:"totalElements"`
}

type PeopleInfo struct {
	CountryId string `json:"countryId"`
	Gender    string `json:"gender"`
	Id        string `json:"id"`
	Name      string `json:"name"`
	SubId     int    `json:"subId"`
}

type RankInfo struct {
	Msg                  string `json:"msg"`
	PageLast             bool   `json:"pageLast"`
	PageEmpty            bool   `json:"pageEmpty"`
	Data                 []PeopleRankInfo
	TotalPages           int  `json:"totalPages"`
	PageFirst            bool `json:"pageFirst"`
	PageSize             int  `json:"pageSize"`
	PageNumberOfElements int  `json:"pageNumberOfElements"`
	PageNum              int  `json:"pageNum"`
	Retcode              int  `json:"retcode"`
	TotalElements        int  `json:"totalElements"`
}

type PeopleRankInfo struct {
	Best          int    `json:"best"`
	ContinentRank int    `json:"continentRank"`
	CountryRank   string `json:"countryRank"`
	EventId       string `json:"eventId"`
	PersonId      string `json:"personId"`
	WorldRank     int    `json:"worldRank"`
}

/*func convertTime(result int, eventId string) string {
	if result == -1 {
		return "DNF"
	} else if result == -2 {
		return "DNS"
	} else if result == 0 {
		return ""
	} else if eventId == "333fm" {
		if result > 1000 {
			return fmt.Sprintf("%.2f", float64(result)/float64(100))
		} else {
			return strconv.Itoa(result)
		}
	} else if eventId == "333mbf" {
		mbfDifference := 99 - result/10000000
		mbfMissed := result % 100
		mbfSolved := mbfDifference + mbfMissed
		mbfAttempted := mbfSolved + mbfMissed
		mbfTime := result / 100
		mbfSec := mbfTime % 10000
		mbfMin := mbfSec / 60
		mbfSec = mbfSec % 60
		return fmt.Sprintf("%d/%d %d:%02d", mbfSolved, mbfAttempted, mbfMin, mbfSec)
	} else {
		sec := result / 100
		msec := result % 100
		if sec > 59 {
			min := sec / 60
			sec := sec % 60
			return fmt.Sprintf("%d:%02d.%02d", min, sec, msec)
		}
		return fmt.Sprintf("%d.%02d", sec, msec)
	}
}

func wcaPersonHandler(s string) string {
	singleUrl := "http://www.2mf8.cn:8083/wcaSingle/findBestResultsByPersonId?personId=" + url.QueryEscape(s)
	averageUrl := "http://www.2mf8.cn:8083/wcaAverage/findBestResultsByPersonId?personId=" + url.QueryEscape(s)
	resp1, _ := http.Get(singleUrl)
	body1, _ := io.ReadAll(resp1.Body)
	resp1.Body.Close()
	resp2, _ := http.Get(averageUrl)
	body2, _ := io.ReadAll(resp2.Body)
	resp2.Body.Close()
	sr := RankInfo{}
	ar := RankInfo{}
	result := ""
	json.Unmarshal([]byte(body1), &sr)
	json.Unmarshal([]byte(body2), &ar)
	for _, i := range sr.Data {
		s_i := i.EventId + " " + convertTime(i.Best, i.EventId)
		result += "\n" + s_i
		for _, j := range ar.Data {
			if i.EventId == j.EventId {
				s_j := convertTime(j.Best, j.EventId)
				result += " | " + s_j
			}
		}
	}
	return result
}*/

/*func main() {
	url := "http://www.2mf8.cn:8083/wcaPerson/searchPeople?q=" + url.QueryEscape("2017WANY29")
	resp, _ := http.Get(url)
	s := Info{}
	gen := ""
	body, _ := ioutil.ReadAll(resp.Body)
	resp.Body.Close()
	json.Unmarshal([]byte(body), &s)
	//fmt.Println(fmt.Sprintf("%+v", s))
	//fmt.Println(len(s.Data))
	if s.TotalElements == 1 {
		//fmt.Println(s.Data[0].Id)
		if s.Data[0].Gender == "m" {
			gen = "Male"
		} else {
			gen = "Female"
		}
		s_r := s.Data[0].Name + "\n" + s.Data[0].Id + "," + s.Data[0].CountryId + "," + gen + wcaPersonHandler(s.Data[0].Id)
		//fmt.Println(s_r)
	}
}*/
