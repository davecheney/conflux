package recon

import (
	"fmt"
	"github.com/stvp/go-toml-config"
	"net"
	"strings"
)

type Settings struct {
	Version                     string
	LogName                     string
	HttpPort                    int
	ReconPort                   int
	Partners                    []string
	Filters                     []string
	ThreshMult                  int
	BitQuantum                  int
	MBar                        int
	SplitThreshold              int
	JoinThreshold               int
	NumSamples                  int
	GossipIntervalSecs          int
	MaxOutstandingReconRequests int
}

func NewSettings() *Settings {
	s := &Settings{
		Version:                     "experimental",
		LogName:                     "conflux.recon",
		HttpPort:                    11371,
		ReconPort:                   11370,
		ThreshMult:                  DefaultThreshMult,
		BitQuantum:                  DefaultBitQuantum,
		MBar:                        DefaultMBar,
		GossipIntervalSecs:          60,
		MaxOutstandingReconRequests: 100}
	s.updateDerived()
	return s
}

func (s *Settings) Config() map[string]string {
	m := make(map[string]string)
	m["version"] = s.Version
	m["http port"] = fmt.Sprintf("%d", s.HttpPort)
	m["bitquantum"] = fmt.Sprintf("%d", s.BitQuantum)
	m["mbar"] = fmt.Sprintf("%d", s.MBar)
	m["filters"] = strings.Join(s.Filters, ",")
	return m
}

func (s *Settings) updateDerived() {
	s.SplitThreshold = s.ThreshMult * s.MBar
	s.JoinThreshold = s.SplitThreshold / 2
	s.NumSamples = s.MBar + 1
}

func LoadSettings(path string) *Settings {
	s := NewSettings()
	c := config.NewConfigSet("conflux.recon", config.ExitOnError)
	version := c.String("Version", s.Version)
	logName := c.String("LogName", s.LogName)
	httpPort := c.Int("HttpPort", s.HttpPort)
	reconPort := c.Int("ReconPort", s.ReconPort)
	threshMult := c.Int("ThreshMult", s.ThreshMult)
	bitQuantum := c.Int("BitQuantum", s.BitQuantum)
	mBar := c.Int("MBar", s.MBar)
	gossipIntervalSecs := c.Int("GossipIntervalSecs", s.GossipIntervalSecs)
	maxOutstandingReconRequests := c.Int("MaxOutstandingReconRequests", s.MaxOutstandingReconRequests)
	partners := c.String("Partners", "")
	filters := c.String("Filters", "")
	c.Parse(path)
	s.Version = *version
	s.LogName = *logName
	s.HttpPort = *httpPort
	s.ReconPort = *reconPort
	s.ThreshMult = *threshMult
	s.BitQuantum = *bitQuantum
	s.MBar = *mBar
	s.GossipIntervalSecs = *gossipIntervalSecs
	s.MaxOutstandingReconRequests = *maxOutstandingReconRequests
	for _, partner := range strings.Split(*partners, ",") {
		s.Partners = append(s.Partners, partner)
	}
	for _, filter := range strings.Split(*filters, ",") {
		s.Filters = append(s.Filters, filter)
	}
	s.updateDerived()
	return s
}

func (s *Settings) PartnerAddrs() (addrs []net.Addr, err error) {
	for _, partner := range s.Partners {
		addr, err := net.ResolveTCPAddr("tcp", partner)
		if err != nil {
			return nil, err
		}
		addrs = append(addrs, addr)
	}
	return
}