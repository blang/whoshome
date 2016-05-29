package whoshome

import (
	"bufio"
	"io"
	"os"
	"strings"
)

type ARPProvider struct {
	arpFile  string
	mac2name map[string]string
}

func NewARPProvider(arpFile string, mac2name map[string]string) *ARPProvider {
	return &ARPProvider{
		arpFile:  arpFile,
		mac2name: mac2name,
	}
}

func (p *ARPProvider) Present() ([]string, error) {
	f, err := os.Open(p.arpFile)
	if err != nil {
		return nil, err
	}
	var ls []string
	br := bufio.NewReader(f)
	// skip first line
	br.ReadString('\n')
	for {
		line, err := br.ReadString('\n')
		if err != nil {
			if err != io.EOF {
				return nil, err
			}
			break
		}
		parts := strings.Fields(line)
		if len(parts) != 6 {
			continue
		}
		if parts[2] != "0x2" {
			continue
		}
		// Only selected devices
		if name, ok := p.mac2name[parts[3]]; ok {
			ls = append(ls, name)
		}
	}
	return ls, nil
}
