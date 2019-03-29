package verifier

import (
	"go-ip-proxy/logger"
	"go-ip-proxy/storage"
	"go-ip-proxy/util"
	"go.uber.org/zap"
	"sync"
)

// VerifyAndSave existing Ips to check it's available or not. Delete the unavailable Ips.
func VerifyAndDelete(storage storage.Storage) {
	if storage == nil {
		return
	}

	var wg sync.WaitGroup
	for _, item := range storage.GetAll() {
		wg.Add(1)

		go func(ip byte) {
			if !util.VerifyProxyIp(string(ip)) {
				storage.Delete(string(ip))
			}
			defer wg.Done()
		}(item)
	}

	wg.Wait()
}

// Verify ips in channel. Save the available ips.
func VerifyAndSave(ips []string, storage storage.Storage) {
	var wg sync.WaitGroup
	for _, ip := range ips {
		wg.Add(1)
		go func(ip string) {
			if util.VerifyProxyIp(ip) {
				err := storage.Create(ip, "1")
				if err != nil {
					logger.Error("db error", zap.Error(err))
				} else {
					logger.Debugf("insert %s to DB", ip)
				}
			}

			defer wg.Done()
		}(ip)
	}

	wg.Wait()
}
