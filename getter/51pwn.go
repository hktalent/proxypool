package getter

import (
	"bufio"
	"fmt"
	"github.com/hktalent/htmlquery"
	"github.com/hktalent/proxypool/pkg/models"
	"github.com/tidwall/gjson"
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

const (
	MacLineSize = 10 * 1024 * 1024 // 10M
)

var ExprIP = regexp.MustCompile(`((25[0-5]|2[0-4]\d|((1\d{2})|([1-9]?\d)))\.){3}(25[0-5]|2[0-4]\d|((1\d{2})|([1-9]?\d)))\:([0-9]+)`)

func LoadStatic() (result []*models.IP) {
	if szPwd, err := os.Getwd(); nil == err {
		szPath := "/conf/"
		for _, x := range []string{"ALIILAPRO_Proxy", "openproxylist", "proxylist-update-every-minute"} {
			szCur := szPwd + szPath + x + "/"
			var Visit = func(path string, d fs.DirEntry, err error) error {
				if strings.HasSuffix(path, ".txt") {

				}
				return nil
			}

			if err := filepath.WalkDir(szCur, Visit); nil != err {
				fmt.Printf("filepath.WalkDir() returned %v\n", err)
			}
		}
	}
	return
}

func FreeProxyList() (result []*models.IP) {
	pollURL := "https://free-proxy-list.net"
	if doc, err := htmlquery.LoadURL(pollURL); nil != err {
		log.Println("free-proxy-list", err)
	} else {
		// ips := ExprIP.FindAllString(bodyIPs, 100)
		trNode := htmlquery.Find(doc, "#list > div > div.table-responsive > div > table > tr")
		for i := 0; i < len(trNode); i++ {
			tdNode := htmlquery.Find(trNode[i], "//td")
			ip := extractTextFromDivNode(tdNode[0])
			if 0 == len(ExprIP.FindAllString(ip, -1)) {
				continue
			}
			port := extractTextFromDivNode(tdNode[1])
			Type := "http"
			speed := htmlquery.InnerText(tdNode[7])
			s1 := " secs ago"
			if len(speed) > len(s1) {
				speed = speed[0 : len(speed)-len(s1)]
			}
			if "yes" == htmlquery.InnerText(tdNode[6]) {
				Type = "https"
			}

			IP := models.NewIP()
			IP.Data = ip + ":" + port
			IP.Type1 = Type
			IP.Source = "FreeProxyList"
			IP.Speed = extractSpeed(speed)
			result = append(result, IP)
		}
	}
	return
}

func Freeproxylists() (result []*models.IP) {
	pollURL := "https://www.freeproxylists.net/?c=&pt=&pr=&a%5B%5D=0&a%5B%5D=1&a%5B%5D=2&u=0"
	if doc, err := htmlquery.LoadURL(pollURL); nil != err {
		log.Println("Freeproxylists", err)
	} else {
		// ips := ExprIP.FindAllString(bodyIPs, 100)
		trNode := htmlquery.Find(doc, "body > div:nth-child(3) > div:nth-child(2) > table > tr")
		for i := 0; i < len(trNode); i++ {
			tdNode := htmlquery.Find(trNode[i], "//td")
			ip := extractTextFromDivNode(tdNode[0])
			if 0 == len(ExprIP.FindAllString(ip, -1)) {
				continue
			}
			port := extractTextFromDivNode(tdNode[1])
			Type := "http"
			speed := htmlquery.InnerText(tdNode[8])

			IP := models.NewIP()
			IP.Data = ip + ":" + port
			IP.Type1 = Type
			IP.Source = "freeproxylists"
			IP.Speed = extractSpeed(speed)
			result = append(result, IP)
		}
	}
	return
}

func Geonode() (result []*models.IP) {
	for i := 1; i < 500; i++ {
		pollURL := fmt.Sprintf("https://proxylist.geonode.com/api/proxy-list?limit=500&page=%d&sort_by=lastChecked&sort_type=desc&protocols=socks5", i)
		if doc, err := htmlquery.LoadURL(pollURL); nil != err {
			log.Println("geonode", err)
		} else {
			oJ := gjson.Parse(doc.Data)
			oD1 := oJ.Get("data")
			oA1 := oD1.Array()
			if 0 == len(oA1) {
				break
			}
			for _, x := range oA1 {
				ip := x.Get("ip")
				port := x.Get("port")
				IP := models.NewIP()
				IP.Data = fmt.Sprintf("%v:%v", ip, port)
				IP.Type1 = fmt.Sprintf("%v", x.Get("type"))
				IP.Source = "geonode"
				IP.Speed = extractSpeed(fmt.Sprintf("%v", x.Get("responseTime")))
				result = append(result, IP)
			}
		}
	}
	return
}

func Fatezero() (result []*models.IP) {
	pollURL := "http://proxylist.fatezero.org/proxy.list"
	if doc, err := htmlquery.LoadURL(pollURL); nil != err {
		log.Println("fatezero", err)
	} else {
		scanner := bufio.NewScanner(strings.NewReader(doc.Data))
		scanner.Buffer(make([]byte, MacLineSize), MacLineSize)
		for scanner.Scan() {
			if value := strings.TrimSpace(scanner.Text()); 0 < len(value) {
				oJ := gjson.Parse(value)
				ip := oJ.Get("host")
				port := oJ.Get("port")
				IP := models.NewIP()
				IP.Data = fmt.Sprintf("%v:%v", ip, port)
				IP.Type1 = fmt.Sprintf("%v", oJ.Get("type"))
				IP.Source = "fatezero"
				IP.Speed = 100
				result = append(result, IP)
			}
		}
	}
	return
}
