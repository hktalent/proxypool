package getter

import (
	"fmt"
	"io/ioutil"
	"net/http"

	//"fmt"
	clog "unknwon.dev/clog/v2"

	"strings"

	"github.com/hktalent/proxypool/pkg/models"
)

// IP89 get ip from www.89ip.cn
func IP89() (result []*models.IP) {
	clog.Info("89IP] start test")
	// https://www.89ip.cn/tqdl.html?api=1&num=9999&port=&address=日本&isp=
	for _, x := range []string{"日本", "%E7%BE%8E%E5%9B%BD", "新加坡", "加拿大"} {
		pollURL := fmt.Sprintf("http://www.89ip.cn/tqdl.html?api=1&num=9999&port=&address=%s&isp=", x)

		resp, err := http.Get(pollURL)
		if err != nil {
			clog.Warn(err.Error())
			return
		}

		if resp.StatusCode != 200 {
			clog.Warn(err.Error())
			return
		}
		defer resp.Body.Close()
		body, _ := ioutil.ReadAll(resp.Body)
		bodyIPs := string(body)
		ips := ExprIP.FindAllString(bodyIPs, 100)

		for index := 0; index < len(ips); index++ {
			ip := models.NewIP()
			ip.Data = strings.TrimSpace(ips[index])
			ip.Type1 = "http"
			ip.Source = "89ip"
			clog.Info("[89IP] ip = %s, type = %s", ip.Data, ip.Type1)
			result = append(result, ip)
		}

		clog.Info("89IP done.")
	}
	return
}
