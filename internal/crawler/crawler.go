package crawler

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/barasher/dep-carto/internal/model"
	"github.com/barasher/dep-carto/internal/parser"
	"github.com/rs/zerolog/log"
)

const (
	xmlExtension  = ".xml"
	jsonExtension = ".json"
	urlSuffix     = "server"
	slash         = "/"
)

type Crawler struct {
	url      string
	hostname string
	ips      []string
	re       parser.RefExtractor
	date     time.Time
}

func NewCrawler(u string, re parser.RefExtractor) (Crawler, error) {
	c := Crawler{
		re:   re,
		date: time.Now(),
	}

	if strings.HasSuffix(u, slash) {
		c.url = fmt.Sprintf("%v%v", u, urlSuffix)
	} else {
		c.url = fmt.Sprintf("%v%v%v", u, slash, urlSuffix)
	}

	h, err := hostname()
	if err != nil {
		return c, err
	}
	c.hostname = h

	i, err := ips()
	if err != nil {
		return c, err
	}
	c.ips = i

	return c, nil
}

func hostname() (string, error) {
	h, err := os.Hostname()
	if err != nil {
		return h, fmt.Errorf("error while getting hostname: %w", err)
	}
	return h, nil
}

func ips() ([]string, error) {
	var ips []string
	interfaces, err := net.Interfaces()
	if err != nil {
		return ips, fmt.Errorf("error while getting interfaces: %w", err)
	}
	for _, i := range interfaces {
		addresses, err := i.Addrs()
		if err != nil {
			return ips, fmt.Errorf("error while getting adresses for interface %v: %w", i.Name, err)
		}
		for _, address := range addresses {
			var ip net.IP
			switch v := address.(type) {
			case *net.IPNet:
				ip = v.IP
			case *net.IPAddr:
				ip = v.IP
			}
			if !ip.IsLoopback() {
				ips = append(ips, ip.String())
			}
		}
	}
	return ips, nil
}

func (c Crawler) Crawl(inputs []string) error {
	var deps []string
	for _, curInput := range inputs {
		err := filepath.Walk(curInput, func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}
			if !info.IsDir() {
				d, err := c.parseFile(path)
				if err != nil {
					return fmt.Errorf("error while parsing %v: %w", path, err)
				}
				if len(d) > 0 {
					deps = append(deps, d...)
				}
			}
			return nil
		})
		if err != nil {
			return fmt.Errorf("error while browsing %v: %w", curInput, err)
		}
	}

	if len(deps) > 0 {
		var dDeps []string
		m := make(map[string]bool)
		for _, curD := range deps {
			if _, found := m[curD]; !found {
				m[curD] = true
				dDeps = append(dDeps, curD)
			}
		}
		s := c.buildServer(dDeps)
		if err := c.pushServer(s); err != nil {
			return fmt.Errorf("error while pushing server: %w", err)
		}
	}

	return nil
}

func (c Crawler) parseFile(f string) ([]string, error) {
	var p parser.Parser
	switch ext := filepath.Ext(f); ext {
	case xmlExtension:
		p = parser.NewXmlParser(&c.re)
	case jsonExtension:
		p = parser.NewJsonParser(&c.re)
	default:
		log.Info().Str("file", f).Msgf("Unsupported extension (%v)", ext)
		return nil, nil
	}

	r, err := os.Open(f)
	if err != nil {
		return nil, fmt.Errorf("error while opening %v: %w", f, err)
	}
	defer r.Close()
	return p.Parse(r)
}

func (c Crawler) buildServer(deps []string) model.Server {
	d := make([]model.Dependency, len(deps))
	for i, cur := range deps {
		d[i] = model.Dependency{Resource: cur}
	}
	return model.Server{
		Hostname:     c.hostname,
		IPs:          c.ips,
		Dependencies: d,
		LastUpdate:   c.date,
	}
}

func (c Crawler) pushServer(s model.Server) error {
	buf := bytes.Buffer{}
	if err := json.NewEncoder(&buf).Encode(s); err != nil {
		return fmt.Errorf("error while marshalling server: %w", err)
	}
	resp, err := http.Post(c.url, "application/json", &buf)
	if err != nil {
		return fmt.Errorf("error while posting server: %w", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusNoContent {
		return fmt.Errorf("unexpected http status (%v)", resp.StatusCode)
	}
	return nil
}
