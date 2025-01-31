package main

import (
	"log"
	"runtime"
	"sync"
	"time"

	"github.com/hktalent/proxypool/api"
	"github.com/hktalent/proxypool/getter"
	"github.com/hktalent/proxypool/pkg/initial"
	"github.com/hktalent/proxypool/pkg/models"
	"github.com/hktalent/proxypool/pkg/storage"
)

func main() {

	//init the database
	initial.GlobalInit()

	runtime.GOMAXPROCS(runtime.NumCPU())
	ipChan := make(chan *models.IP, 2000)

	// Start HTTP
	go func() {
		api.Run()
	}()

	// Check the IPs in DB
	go func() {
		storage.CheckProxyDB()
	}()

	// Check the IPs in channel
	for i := 0; i < 50; i++ {
		go func() {
			for {
				storage.CheckProxy(<-ipChan)
			}
		}()
	}

	// Start getters to scraper IP and put it in channel
	for {
		n := models.CountIPs()
		log.Printf("Chan: %v, IP: %v\n", len(ipChan), n)
		if len(ipChan) < 10000 {
			go run(ipChan)
		}
		time.Sleep(10 * time.Minute)
	}
}

func run(ipChan chan<- *models.IP) {
	var wg sync.WaitGroup
	funs := []func() []*models.IP{
		getter.LoadStatic,
		getter.FQDL,  //新代理
		getter.PZZQZ, //新代理
		getter.Fatezero,
		getter.Freeproxylists,
		getter.FreeProxyList,
		getter.Geonode,
		//getter.Data5u,
		//getter.Feiyi,
		//getter.IP66, //need to remove it
		getter.IP3306,
		getter.KDL,
		//getter.GBJ,	//因为网站限制，无法正常下载数据
		//getter.Xici,
		//getter.XDL,
		//getter.IP181,  // 已经无法使用
		//getter.YDL,	//失效的采集脚本，用作系统容错实验
		// getter.PLP, //need to remove it
		getter.PLPSSL,
		getter.IP89,
	}
	for _, f := range funs {
		wg.Add(1)
		go func(f func() []*models.IP) {
			defer func() {
				if err := recover(); nil != err {
					log.Println(err)
				}
			}()
			defer wg.Done()
			temp := f()
			//log.Println("[run] get into loop")
			for _, v := range temp {
				//log.Println("[run] len of ipChan %v",v)
				ipChan <- v
			}
		}(f)
	}
	wg.Wait()
	log.Println("All getters finished.")
}
