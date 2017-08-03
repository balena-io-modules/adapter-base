package wifi

import (
	"context"
	"fmt"
	"io/ioutil"
	"net/http"
	"os/exec"

	log "github.com/Sirupsen/logrus"
	"github.com/lair-framework/go-nmap"
	"github.com/parnurzeal/gorequest"
)

// TODO: these should be configurable
const (
	GATEWAY_IP = "192.168.1.112"
	SCAN_RANGE = "192.168.1.*"
)

type Host struct {
	Name string
	Ip   string
	Mac  string
}

func PostForm(url, filePath string) error {
	req := gorequest.New()
	req.Post(url)
	req.Type("multipart")
	req.SendFile(filePath, "firmware.bin", "image")

	log.WithFields(log.Fields{
		"URL":    req.Url,
		"Method": req.Method,
	}).Info("Posting form")

	resp, _, errs := req.End()
	return handleResp(resp, errs, http.StatusOK)
}

func Scan(ctx context.Context) ([]Host, error) {
	cmd := exec.CommandContext(ctx, "bash", "-c", fmt.Sprintf("nmap -sP %s -oX /tmp/scan.txt", SCAN_RANGE))

	if err := cmd.Run(); err != nil {
		return nil, err
	}

	file, err := ioutil.ReadFile("/tmp/scan.txt")
	if err != nil {
		return nil, err
	}

	nmap, err := nmap.Parse(file)
	if err != nil {
		return nil, err
	}

	var hosts []Host
	for _, host := range nmap.Hosts {
		h := Host{}

		for _, address := range host.Addresses {
			if address.AddrType == "mac" {
				h.Mac = address.Addr
			} else {
				h.Ip = address.Addr
			}
		}

		// Ignore the gateway device
		if h.Ip == GATEWAY_IP {
			continue
		}

		url := "http://" + h.Ip + "/id"
		resp, body, errs := gorequest.New().Get(url).End()
		if err := handleResp(resp, errs, 200); err != nil {
			log.WithFields(log.Fields{
				"Error": err,
				"URL":   url,
				"IP":    h.Ip,
				"MAC":   h.Mac,
			}).Warn("Unable to get device ID")
			continue
		}
		h.Name = body

		hosts = append(hosts, h)
	}

	return hosts, nil
}

func handleResp(resp gorequest.Response, errs []error, statusCode int) error {
	if errs != nil {
		return errs[0]
	}

	if resp.StatusCode != statusCode {
		return fmt.Errorf("Invalid response received: %s", resp.Status)
	}

	log.WithFields(log.Fields{
		"Response": resp.Status,
	}).Debug("Valid response received")

	return nil
}
