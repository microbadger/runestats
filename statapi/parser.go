package statapi

import (
	"html/template"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

const (
	_oldSchoolURL    = "http://services.runescape.com/m=hiscore_oldschool/hiscorepersonal.ws?user1="
	_oldSchoolAPIURL = "http://services.runescape.com/m=hiscore_oldschool/index_lite.ws?player="
)

func OldSchoolAPIHandler(p string) []string {
	res, err := http.Get(_oldSchoolAPIURL + p)
	if err != nil {
		return nil
	}
	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil
	}
	var stats []string
	for _, row := range strings.Split(string(body), "\n") {
		stat := newStatFromAPI(row)
		if stat != "" {
			stats = append(stats, stat)
		}
	}
	return stats
}

func newStatFromAPI(row string) string {
	if row == "" {
		return ""
	}
	stats := strings.Split(row, ",")
	if len(stats) > 0 {
		return stats[1]
	}
	return ""
}

func OldSchoolHandler(p string) []Stat {
	if p == "" {
		return nil
	}
	doc, err := goquery.NewDocument(_oldSchoolURL + p[1:])
	if err != nil {
		return nil
	}
	var rows []string
	doc.Find("td").Each(func(i int, s *goquery.Selection) {
		if i > 12 && i < 131 {
			res, _ := s.Html()
			rows = append(rows, res)
		}
	})
	if rows == nil {
		return nil
	}
	var stats []Stat
	stats = append(stats, newStat(rows[0], "", rows[2]))
	for i := 5; i < 118; i = i + 5 {
		stats = append(stats, newStat(rows[i], rows[i-1], rows[i+2]))
	}

	return stats
}

func newStat(t string, p string, v string) Stat {
	// TODO: Once painter is functional, change new stats to only
	//	     keep the type of stat and the value, since that's all i need
	s := Stat{
		Type:    template.HTML(strings.Replace(t, "\n", "", -1)),
		Picture: template.HTML(p),
		Value:   v}
	s.Type = template.HTML(strings.Replace(string(s.Type), "overall.ws",
		"http://services.runescape.com/m=hiscore_oldschool/overall.ws", -1))
	s.Picture = template.HTML(strings.Replace(string(s.Picture),
		"http://www.runescape.com/img/rsp777/hiscores",
		"templates/static/images/", -1))
	return s
}
