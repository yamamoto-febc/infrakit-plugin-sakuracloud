package sacloud

import (
	"fmt"
	"reflect"
)

type VPCRouterSetting struct {
	Interfaces         []*VPCRouterInterface        `json:",omitempty"`
	StaticNAT          *VPCRouterStaticNAT          `json:",omitempty"`
	PortForwarding     *VPCRouterPortForwarding     `json:",omitempty"`
	Firewall           *VPCRouterFirewall           `json:",omitempty"`
	DHCPServer         *VPCRouterDHCPServer         `json:",omitempty"`
	DHCPStaticMapping  *VPCRouterDHCPStaticMapping  `json:",omitempty"`
	L2TPIPsecServer    *VPCRouterL2TPIPsecServer    `json:",omitempty"`
	PPTPServer         *VPCRouterPPTPServer         `json:",omitempty"`
	RemoteAccessUsers  *VPCRouterRemoteAccessUsers  `json:",omitempty"`
	SiteToSiteIPsecVPN *VPCRouterSiteToSiteIPsecVPN `json:",omitempty"`
	StaticRoutes       *VPCRouterStaticRoutes       `json:",omitempty"`
	VRID               *int                         `json:",omitempty"`
}

type VPCRouterInterface struct {
	IPAddress        []string `json:",omitempty"`
	NetworkMaskLen   int      `json:",omitempty"`
	VirtualIPAddress string   `json:",omitempty"`
	IPAliases        []string `json:",omitempty"`
}

func (s *VPCRouterSetting) AddInterface(vip string, ipaddress []string, maskLen int) {
	if s.Interfaces == nil {
		s.Interfaces = []*VPCRouterInterface{nil}
	}
	s.Interfaces = append(s.Interfaces, &VPCRouterInterface{
		VirtualIPAddress: vip,
		IPAddress:        ipaddress,
		NetworkMaskLen:   maskLen,
	})
}

type VPCRouterStaticNAT struct {
	Config  []*VPCRouterStaticNATConfig `json:",omitempty"`
	Enabled string                      `json:",omitempty"`
}
type VPCRouterStaticNATConfig struct {
	GlobalAddress  string `json:",omitempty"`
	PrivateAddress string `json:",omitempty"`
}

func (s *VPCRouterSetting) AddStaticNAT(globalAddress string, privateAddress string) {
	if s.StaticNAT == nil {
		s.StaticNAT = &VPCRouterStaticNAT{
			Enabled: "True",
		}
	}

	if s.StaticNAT.Config == nil {
		s.StaticNAT.Config = []*VPCRouterStaticNATConfig{}
	}

	s.StaticNAT.Config = append(s.StaticNAT.Config, &VPCRouterStaticNATConfig{
		GlobalAddress:  globalAddress,
		PrivateAddress: privateAddress,
	})
}

func (s *VPCRouterSetting) RemoveStaticNAT(globalAddress string, privateAddress string) {
	if s.StaticNAT == nil {
		return
	}

	if s.StaticNAT.Config == nil {
		s.StaticNAT.Enabled = "False"
		return
	}

	dest := []*VPCRouterStaticNATConfig{}
	for _, c := range s.StaticNAT.Config {
		if c.GlobalAddress != globalAddress || c.PrivateAddress != privateAddress {
			dest = append(dest, c)
		}
	}
	s.StaticNAT.Config = dest
	if len(s.StaticNAT.Config) == 0 {
		s.StaticNAT.Enabled = "False"
		s.StaticNAT.Config = nil
		return
	}
	s.StaticNAT.Enabled = "True"
}

func (s *VPCRouterSetting) FindStaticNAT(globalAddress string, privateAddress string) *VPCRouterStaticNATConfig {
	for _, c := range s.StaticNAT.Config {
		if c.GlobalAddress == globalAddress && c.PrivateAddress == privateAddress {
			return c
		}
	}
	return nil
}

type VPCRouterPortForwarding struct {
	Config  []*VPCRouterPortForwardingConfig `json:",omitempty"`
	Enabled string                           `json:",omitempty"`
}
type VPCRouterPortForwardingConfig struct {
	Protocol       string `json:",omitempty"` // tcp/udp only
	GlobalPort     string `json:",omitempty"`
	PrivateAddress string `json:",omitempty"`
	PrivatePort    string `json:",omitempty"`
}

func (s *VPCRouterSetting) AddPortForwarding(protocol string, globalPort string, privateAddress string, privatePort string) {
	if s.PortForwarding == nil {
		s.PortForwarding = &VPCRouterPortForwarding{
			Enabled: "True",
		}
	}

	if s.PortForwarding.Config == nil {
		s.PortForwarding.Config = []*VPCRouterPortForwardingConfig{}
	}

	s.PortForwarding.Config = append(s.PortForwarding.Config, &VPCRouterPortForwardingConfig{
		Protocol:       protocol,
		GlobalPort:     globalPort,
		PrivateAddress: privateAddress,
		PrivatePort:    privatePort,
	})
}

func (s *VPCRouterSetting) RemovePortForwarding(protocol string, globalPort string, privateAddress string, privatePort string) {
	if s.PortForwarding == nil {
		return
	}

	if s.PortForwarding.Config == nil {
		s.PortForwarding.Enabled = "False"
		return
	}

	dest := []*VPCRouterPortForwardingConfig{}
	for _, c := range s.PortForwarding.Config {
		if c.Protocol != protocol || c.GlobalPort != globalPort ||
			c.PrivateAddress != privateAddress || c.PrivatePort != privatePort {
			dest = append(dest, c)
		}
	}
	s.PortForwarding.Config = dest
	if len(s.PortForwarding.Config) == 0 {
		s.PortForwarding.Enabled = "False"
		s.PortForwarding.Config = nil
		return
	}
	s.PortForwarding.Enabled = "True"
}
func (s *VPCRouterSetting) FindPortForwarding(protocol string, globalPort string, privateAddress string, privatePort string) *VPCRouterPortForwardingConfig {
	for _, c := range s.PortForwarding.Config {
		if c.Protocol == protocol && c.GlobalPort == globalPort &&
			c.PrivateAddress == privateAddress && c.PrivatePort == privatePort {
			return c
		}
	}
	return nil
}

type VPCRouterFirewall struct {
	Config  []*VPCRouterFirewallSetting `json:",omitempty"`
	Enabled string                      `json:",omitempty"`
}
type VPCRouterFirewallSetting struct {
	Receive []*VPCRouterFirewallRule `json:",omitempty"`
	Send    []*VPCRouterFirewallRule `json:",omitempty"`
}
type VPCRouterFirewallRule struct {
	Action             string `json:",omitempty"`
	Protocol           string `json:",omitempty"`
	SourceNetwork      string `json:",omitempty"`
	SourcePort         string `json:",omitempty"`
	DestinationNetwork string `json:",omitempty"`
	DestinationPort    string `json:",omitempty"`
}

func (s *VPCRouterSetting) addFirewallRule(direction string, rule *VPCRouterFirewallRule) {
	if s.Firewall == nil {
		s.Firewall = &VPCRouterFirewall{
			Enabled: "True",
		}
	}
	if s.Firewall.Config == nil || len(s.Firewall.Config) == 0 {
		s.Firewall.Config = []*VPCRouterFirewallSetting{
			{
				Receive: []*VPCRouterFirewallRule{},
				Send:    []*VPCRouterFirewallRule{},
			},
		}
	}

	switch direction {
	case "send":
		s.Firewall.Config[0].Send = append(s.Firewall.Config[0].Send, rule)
	case "receive":
		s.Firewall.Config[0].Receive = append(s.Firewall.Config[0].Receive, rule)
	}
}

func (s *VPCRouterSetting) removeFirewallRule(direction string, rule *VPCRouterFirewallRule) {

	if s.Firewall == nil {
		return
	}

	if s.Firewall.Config == nil {
		s.Firewall.Enabled = "False"
		return
	}

	switch direction {
	case "send":
		dest := []*VPCRouterFirewallRule{}
		for _, c := range s.Firewall.Config[0].Send {
			if c.Action != rule.Action || c.Protocol != rule.Protocol ||
				c.SourceNetwork != rule.SourceNetwork || c.SourcePort != rule.SourcePort ||
				c.DestinationNetwork != rule.DestinationNetwork || c.DestinationPort != rule.DestinationPort {
				dest = append(dest, c)
			}
		}
		s.Firewall.Config[0].Send = dest
	case "receive":
		dest := []*VPCRouterFirewallRule{}
		for _, c := range s.Firewall.Config[0].Receive {
			if c.Action != rule.Action || c.Protocol != rule.Protocol ||
				c.SourceNetwork != rule.SourceNetwork || c.SourcePort != rule.SourcePort ||
				c.DestinationNetwork != rule.DestinationNetwork || c.DestinationPort != rule.DestinationPort {
				dest = append(dest, c)
			}
		}
		s.Firewall.Config[0].Receive = dest
	}

	if len(s.Firewall.Config) == 0 {
		s.Firewall.Enabled = "False"
		s.Firewall.Config = nil
		return
	}

	if len(s.Firewall.Config[0].Send) == 0 && len(s.Firewall.Config[0].Send) == 0 {
		s.Firewall.Enabled = "False"
		s.Firewall.Config = nil
		return
	}

	s.PortForwarding.Enabled = "True"

}

func (s *VPCRouterSetting) findFirewallRule(direction string, rule *VPCRouterFirewallRule) *VPCRouterFirewallRule {
	switch direction {
	case "send":
		for _, c := range s.Firewall.Config[0].Send {
			if c.Action == rule.Action && c.Protocol == rule.Protocol &&
				c.SourceNetwork == rule.SourceNetwork && c.SourcePort == rule.SourcePort &&
				c.DestinationNetwork == rule.DestinationNetwork && c.DestinationPort == rule.DestinationPort {
				return c
			}
		}
	case "receive":
		for _, c := range s.Firewall.Config[0].Receive {
			if c.Action == rule.Action && c.Protocol == rule.Protocol &&
				c.SourceNetwork == rule.SourceNetwork && c.SourcePort == rule.SourcePort &&
				c.DestinationNetwork == rule.DestinationNetwork && c.DestinationPort == rule.DestinationPort {
				return c
			}
		}
	}

	return nil

}

func (s *VPCRouterSetting) AddFirewallRuleSend(isAllow bool, protocol string, sourceNetwork string, sourcePort string, destNetwork string, destPort string) {
	action := "deny"
	if isAllow {
		action = "allow"
	}
	rule := &VPCRouterFirewallRule{
		Action:             action,
		Protocol:           protocol,
		SourceNetwork:      sourceNetwork,
		SourcePort:         sourcePort,
		DestinationNetwork: destNetwork,
		DestinationPort:    destPort,
	}

	s.addFirewallRule("send", rule)
}

func (s *VPCRouterSetting) RemoveFirewallRuleSend(isAllow bool, protocol string, sourceNetwork string, sourcePort string, destNetwork string, destPort string) {
	action := "deny"
	if isAllow {
		action = "allow"
	}
	rule := &VPCRouterFirewallRule{
		Action:             action,
		Protocol:           protocol,
		SourceNetwork:      sourceNetwork,
		SourcePort:         sourcePort,
		DestinationNetwork: destNetwork,
		DestinationPort:    destPort,
	}

	s.removeFirewallRule("send", rule)
}

func (s *VPCRouterSetting) FindFirewallRuleSend(isAllow bool, protocol string, sourceNetwork string, sourcePort string, destNetwork string, destPort string) *VPCRouterFirewallRule {
	action := "deny"
	if isAllow {
		action = "allow"
	}
	rule := &VPCRouterFirewallRule{
		Action:             action,
		Protocol:           protocol,
		SourceNetwork:      sourceNetwork,
		SourcePort:         sourcePort,
		DestinationNetwork: destNetwork,
		DestinationPort:    destPort,
	}

	return s.findFirewallRule("send", rule)
}

func (s *VPCRouterSetting) AddFirewallRuleReceive(isAllow bool, protocol string, sourceNetwork string, sourcePort string, destNetwork string, destPort string) {
	action := "deny"
	if isAllow {
		action = "allow"
	}
	rule := &VPCRouterFirewallRule{
		Action:             action,
		Protocol:           protocol,
		SourceNetwork:      sourceNetwork,
		SourcePort:         sourcePort,
		DestinationNetwork: destNetwork,
		DestinationPort:    destPort,
	}

	s.addFirewallRule("receive", rule)
}

func (s *VPCRouterSetting) RemoveFirewallRuleReceive(isAllow bool, protocol string, sourceNetwork string, sourcePort string, destNetwork string, destPort string) {
	action := "deny"
	if isAllow {
		action = "allow"
	}
	rule := &VPCRouterFirewallRule{
		Action:             action,
		Protocol:           protocol,
		SourceNetwork:      sourceNetwork,
		SourcePort:         sourcePort,
		DestinationNetwork: destNetwork,
		DestinationPort:    destPort,
	}

	s.removeFirewallRule("receive", rule)
}

func (s *VPCRouterSetting) FindFirewallRuleReceive(isAllow bool, protocol string, sourceNetwork string, sourcePort string, destNetwork string, destPort string) *VPCRouterFirewallRule {
	action := "deny"
	if isAllow {
		action = "allow"
	}
	rule := &VPCRouterFirewallRule{
		Action:             action,
		Protocol:           protocol,
		SourceNetwork:      sourceNetwork,
		SourcePort:         sourcePort,
		DestinationNetwork: destNetwork,
		DestinationPort:    destPort,
	}

	return s.findFirewallRule("receive", rule)
}

type VPCRouterDHCPServer struct {
	Config  []*VPCRouterDHCPServerConfig `json:",omitempty"`
	Enabled string                       `json:",omitempty"`
}
type VPCRouterDHCPServerConfig struct {
	Interface  string `json:",omitempty"`
	RangeStart string `json:",omitempty"`
	RangeStop  string `json:",omitempty"`
}

func (s *VPCRouterSetting) AddDHCPServer(nicIndex int, rangeStart string, rangeStop string) {
	if s.DHCPServer == nil {
		s.DHCPServer = &VPCRouterDHCPServer{
			Enabled: "True",
		}
	}
	if s.DHCPServer.Config == nil {
		s.DHCPServer.Config = []*VPCRouterDHCPServerConfig{}
	}

	nic := fmt.Sprintf("eth%d", nicIndex)
	s.DHCPServer.Config = append(s.DHCPServer.Config, &VPCRouterDHCPServerConfig{
		Interface:  nic,
		RangeStart: rangeStart,
		RangeStop:  rangeStop,
	})

}

func (s *VPCRouterSetting) RemoveDHCPServer(nicIndex int, rangeStart string, rangeStop string) {
	if s.DHCPServer == nil {
		return
	}

	if s.DHCPServer.Config == nil {
		s.DHCPServer.Enabled = "False"
		return
	}

	dest := []*VPCRouterDHCPServerConfig{}
	for _, c := range s.DHCPServer.Config {
		if c.Interface != fmt.Sprintf("eth%d", nicIndex) || c.RangeStart != rangeStart || c.RangeStop != rangeStop {
			dest = append(dest, c)
		}
	}
	s.DHCPServer.Config = dest

	if len(s.DHCPServer.Config) == 0 {
		s.DHCPServer.Enabled = "False"
		s.DHCPServer.Config = nil
		return
	}
	s.DHCPServer.Enabled = "True"

}

func (s *VPCRouterSetting) FindDHCPServer(nicIndex int, rangeStart string, rangeStop string) *VPCRouterDHCPServerConfig {
	for _, c := range s.DHCPServer.Config {
		if c.Interface == fmt.Sprintf("eth%d", nicIndex) && c.RangeStart == rangeStart && c.RangeStop == rangeStop {
			return c
		}
	}
	return nil
}

type VPCRouterDHCPStaticMapping struct {
	Config  []*VPCRouterDHCPStaticMappingConfig `json:",omitempty"`
	Enabled string                              `json:",omitempty"`
}
type VPCRouterDHCPStaticMappingConfig struct {
	IPAddress  string `json:",omitempty"`
	MACAddress string `json:",omitempty"`
}

func (s *VPCRouterSetting) AddDHCPStaticMapping(ipAddress string, macAddress string) {
	if s.DHCPStaticMapping == nil {
		s.DHCPStaticMapping = &VPCRouterDHCPStaticMapping{
			Enabled: "True",
		}
	}
	if s.DHCPStaticMapping.Config == nil {
		s.DHCPStaticMapping.Config = []*VPCRouterDHCPStaticMappingConfig{}
	}

	s.DHCPStaticMapping.Config = append(s.DHCPStaticMapping.Config, &VPCRouterDHCPStaticMappingConfig{
		IPAddress:  ipAddress,
		MACAddress: macAddress,
	})
}

func (s *VPCRouterSetting) RemoveDHCPStaticMapping(ipAddress string, macAddress string) {
	if s.DHCPStaticMapping == nil {
		return
	}

	if s.DHCPStaticMapping.Config == nil {
		s.DHCPStaticMapping.Enabled = "False"
		return
	}

	dest := []*VPCRouterDHCPStaticMappingConfig{}
	for _, c := range s.DHCPStaticMapping.Config {
		if c.IPAddress != ipAddress || c.MACAddress != macAddress {
			dest = append(dest, c)
		}
	}
	s.DHCPStaticMapping.Config = dest

	if len(s.DHCPStaticMapping.Config) == 0 {
		s.DHCPStaticMapping.Enabled = "False"
		s.DHCPStaticMapping.Config = nil
		return
	}
	s.DHCPStaticMapping.Enabled = "True"

}

func (s *VPCRouterSetting) FindDHCPStaticMapping(ipAddress string, macAddress string) *VPCRouterDHCPStaticMappingConfig {
	for _, c := range s.DHCPStaticMapping.Config {
		if c.IPAddress == ipAddress && c.MACAddress == macAddress {
			return c
		}
	}
	return nil
}

type VPCRouterL2TPIPsecServer struct {
	Config  *VPCRouterL2TPIPsecServerConfig `json:",omitempty"`
	Enabled string                          `json:",omitempty"`
}

type VPCRouterL2TPIPsecServerConfig struct {
	PreSharedSecret string `json:",omitempty"`
	RangeStart      string `json:",omitempty"`
	RangeStop       string `json:",omitempty"`
}

func (s *VPCRouterSetting) EnableL2TPIPsecServer(preSharedSecret string, rangeStart string, rangeStop string) {
	if s.L2TPIPsecServer == nil {
		s.L2TPIPsecServer = &VPCRouterL2TPIPsecServer{
			Enabled: "True",
		}
	}
	s.L2TPIPsecServer.Config = &VPCRouterL2TPIPsecServerConfig{
		PreSharedSecret: preSharedSecret,
		RangeStart:      rangeStart,
		RangeStop:       rangeStop,
	}
}

func (s *VPCRouterSetting) DisableL2TPIPsecServer() {
	if s.L2TPIPsecServer == nil {
		s.L2TPIPsecServer = &VPCRouterL2TPIPsecServer{}
	}
	s.L2TPIPsecServer.Enabled = "False"
	s.L2TPIPsecServer.Config = nil
}

type VPCRouterPPTPServer struct {
	Config  *VPCRouterPPTPServerConfig `json:",omitempty"`
	Enabled string                     `json:",omitempty"`
}
type VPCRouterPPTPServerConfig struct {
	RangeStart string `json:",omitempty"`
	RangeStop  string `json:",omitempty"`
}

func (s *VPCRouterSetting) EnablePPTPServer(rangeStart string, rangeStop string) {
	if s.PPTPServer == nil {
		s.PPTPServer = &VPCRouterPPTPServer{
			Enabled: "True",
		}
	}
	s.PPTPServer.Config = &VPCRouterPPTPServerConfig{
		RangeStart: rangeStart,
		RangeStop:  rangeStop,
	}
}

func (s *VPCRouterSetting) DisablePPTPServer() {
	if s.PPTPServer == nil {
		s.PPTPServer = &VPCRouterPPTPServer{}
	}
	s.PPTPServer.Enabled = "False"
	s.PPTPServer.Config = nil
}

type VPCRouterRemoteAccessUsers struct {
	Config  []*VPCRouterRemoteAccessUsersConfig `json:",omitempty"`
	Enabled string                              `json:",omitempty"`
}
type VPCRouterRemoteAccessUsersConfig struct {
	UserName string `json:",omitempty"`
	Password string `json:",omitempty"`
}

func (s *VPCRouterSetting) AddRemoteAccessUser(userName string, password string) {
	if s.RemoteAccessUsers == nil {
		s.RemoteAccessUsers = &VPCRouterRemoteAccessUsers{
			Enabled: "True",
		}
	}
	if s.RemoteAccessUsers.Config == nil {
		s.RemoteAccessUsers.Config = []*VPCRouterRemoteAccessUsersConfig{}
	}
	s.RemoteAccessUsers.Config = append(s.RemoteAccessUsers.Config, &VPCRouterRemoteAccessUsersConfig{
		UserName: userName,
		Password: password,
	})
}

func (s *VPCRouterSetting) RemoveRemoteAccessUser(userName string, password string) {
	if s.RemoteAccessUsers == nil {
		return
	}

	if s.RemoteAccessUsers.Config == nil {
		s.RemoteAccessUsers.Enabled = "False"
		return
	}

	dest := []*VPCRouterRemoteAccessUsersConfig{}
	for _, c := range s.RemoteAccessUsers.Config {
		if c.UserName != userName || c.Password != password {
			dest = append(dest, c)
		}
	}
	s.RemoteAccessUsers.Config = dest

	if len(s.RemoteAccessUsers.Config) == 0 {
		s.RemoteAccessUsers.Enabled = "False"
		s.RemoteAccessUsers.Config = nil
		return
	}
	s.RemoteAccessUsers.Enabled = "True"
}

func (s *VPCRouterSetting) FindRemoteAccessUser(userName string, password string) *VPCRouterRemoteAccessUsersConfig {
	for _, c := range s.RemoteAccessUsers.Config {
		if c.UserName == userName && c.Password == password {
			return c
		}
	}
	return nil
}

type VPCRouterSiteToSiteIPsecVPN struct {
	Config  []*VPCRouterSiteToSiteIPsecVPNConfig `json:",omitempty"`
	Enabled string                               `json:",omitempty"`
}

type VPCRouterSiteToSiteIPsecVPNConfig struct {
	LocalPrefix     []string `json:",omitempty"`
	Peer            string   `json:",omitempty"`
	PreSharedSecret string   `json:",omitempty"`
	RemoteID        string   `json:",omitempty"`
	Routes          []string `json:",omitempty"`
}

func (s *VPCRouterSetting) AddSiteToSiteIPsecVPN(localPrefix []string, peer string, preSharedSecret string, remoteID string, routes []string) {
	if s.SiteToSiteIPsecVPN == nil {
		s.SiteToSiteIPsecVPN = &VPCRouterSiteToSiteIPsecVPN{
			Enabled: "True",
		}
	}
	if s.SiteToSiteIPsecVPN.Config == nil {
		s.SiteToSiteIPsecVPN.Config = []*VPCRouterSiteToSiteIPsecVPNConfig{}
	}

	s.SiteToSiteIPsecVPN.Config = append(s.SiteToSiteIPsecVPN.Config, &VPCRouterSiteToSiteIPsecVPNConfig{
		LocalPrefix:     localPrefix,
		Peer:            peer,
		PreSharedSecret: preSharedSecret,
		RemoteID:        remoteID,
		Routes:          routes,
	})
}

func (s *VPCRouterSetting) RemoveSiteToSiteIPsecVPN(localPrefix []string, peer string, preSharedSecret string, remoteID string, routes []string) {
	config := &VPCRouterSiteToSiteIPsecVPNConfig{
		LocalPrefix:     localPrefix,
		Peer:            peer,
		PreSharedSecret: preSharedSecret,
		RemoteID:        remoteID,
		Routes:          routes,
	}

	if s.SiteToSiteIPsecVPN == nil {
		return
	}

	if s.SiteToSiteIPsecVPN.Config == nil {
		s.SiteToSiteIPsecVPN.Enabled = "False"
		return
	}

	dest := []*VPCRouterSiteToSiteIPsecVPNConfig{}
	for _, c := range s.SiteToSiteIPsecVPN.Config {
		if !s.isSameSiteToSiteIPsecVPNConfig(c, config) {
			dest = append(dest, c)
		}
	}
	s.SiteToSiteIPsecVPN.Config = dest

	if len(s.SiteToSiteIPsecVPN.Config) == 0 {
		s.SiteToSiteIPsecVPN.Enabled = "False"
		s.SiteToSiteIPsecVPN.Config = nil
		return
	}
	s.SiteToSiteIPsecVPN.Enabled = "True"
}

func (s *VPCRouterSetting) FindSiteToSiteIPsecVPN(localPrefix []string, peer string, preSharedSecret string, remoteID string, routes []string) *VPCRouterSiteToSiteIPsecVPNConfig {
	config := &VPCRouterSiteToSiteIPsecVPNConfig{
		LocalPrefix:     localPrefix,
		Peer:            peer,
		PreSharedSecret: preSharedSecret,
		RemoteID:        remoteID,
		Routes:          routes,
	}

	for _, c := range s.SiteToSiteIPsecVPN.Config {
		if s.isSameSiteToSiteIPsecVPNConfig(c, config) {
			return c
		}
	}
	return nil
}

func (s *VPCRouterSetting) isSameSiteToSiteIPsecVPNConfig(c1 *VPCRouterSiteToSiteIPsecVPNConfig, c2 *VPCRouterSiteToSiteIPsecVPNConfig) bool {
	return reflect.DeepEqual(c1.LocalPrefix, c2.LocalPrefix) &&
		c1.Peer == c2.Peer &&
		c1.PreSharedSecret == c2.PreSharedSecret &&
		c1.RemoteID == c2.RemoteID &&
		reflect.DeepEqual(c1.Routes, c2.Routes)
}

type VPCRouterStaticRoutes struct {
	Config  []*VPCRouterStaticRoutesConfig `json:",omitempty"`
	Enabled string                         `json:",omitempty"`
}
type VPCRouterStaticRoutesConfig struct {
	Prefix  string `json:",omitempty"`
	NextHop string `json:",omitempty"`
}

func (s *VPCRouterSetting) AddStaticRoute(prefix string, nextHop string) {
	if s.StaticRoutes == nil {
		s.StaticRoutes = &VPCRouterStaticRoutes{
			Enabled: "True",
		}
	}
	if s.StaticRoutes.Config == nil {
		s.StaticRoutes.Config = []*VPCRouterStaticRoutesConfig{}
	}
	s.StaticRoutes.Config = append(s.StaticRoutes.Config, &VPCRouterStaticRoutesConfig{
		Prefix:  prefix,
		NextHop: nextHop,
	})
}

func (s *VPCRouterSetting) RemoveStaticRoute(prefix string, nextHop string) {
	if s.StaticRoutes == nil {
		return
	}

	if s.StaticRoutes.Config == nil {
		s.StaticRoutes.Enabled = "False"
		return
	}

	dest := []*VPCRouterStaticRoutesConfig{}
	for _, c := range s.StaticRoutes.Config {
		if c.Prefix != prefix || c.NextHop != nextHop {
			dest = append(dest, c)
		}
	}
	s.StaticRoutes.Config = dest

	if len(s.StaticRoutes.Config) == 0 {
		s.StaticRoutes.Enabled = "False"
		s.StaticRoutes.Config = nil
		return
	}
	s.StaticRoutes.Enabled = "True"
}

func (s *VPCRouterSetting) FindStaticRoute(prefix string, nextHop string) *VPCRouterStaticRoutesConfig {
	for _, c := range s.StaticRoutes.Config {
		if c.Prefix == prefix && c.NextHop == nextHop {
			return c
		}
	}
	return nil
}
