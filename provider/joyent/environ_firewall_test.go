// Copyright 2014 Canonical Ltd.
// Licensed under the AGPLv3, see LICENCE file for details.
package joyent_test

import (
	"github.com/joyent/gosdc/cloudapi"
	jc "github.com/juju/testing/checkers"
	gc "gopkg.in/check.v1"

	"github.com/juju/juju/core/network"
	"github.com/juju/juju/core/network/firewall"
	"github.com/juju/juju/provider/joyent"
)

type FirewallSuite struct{}

var _ = gc.Suite(&FirewallSuite{})

func (s *FirewallSuite) TestGetIngressRules(c *gc.C) {
	testCases := []struct {
		about    string
		envName  string
		rules    []cloudapi.FirewallRule
		expected firewall.IngressRules
	}{
		{
			"single port model rule",
			"switch",
			[]cloudapi.FirewallRule{{
				"",
				true,
				"FROM tag switch TO tag juju ALLOW tcp PORT 80",
			}},
			firewall.IngressRules{firewall.NewIngressRule(network.MustParsePortRange("80/tcp"), "0.0.0.0/0")},
		},
		{
			"port range model rule",
			"switch",
			[]cloudapi.FirewallRule{{
				"",
				true,
				"FROM tag switch TO tag juju ALLOW tcp (PORT 80 AND PORT 81 AND PORT 82 AND PORT 83)",
			}},
			firewall.IngressRules{firewall.NewIngressRule(network.MustParsePortRange("80-83/tcp"), "0.0.0.0/0")},
		},
	}
	for i, t := range testCases {
		c.Logf("test %d: %s", i, t.about)
		ports, err := joyent.GetPorts(t.envName, t.rules)
		c.Assert(err, jc.ErrorIsNil)
		c.Assert(ports, jc.DeepEquals, t.expected)
	}

}

func (s *FirewallSuite) TestRuleCreation(c *gc.C) {
	testCases := []struct {
		about    string
		rules    firewall.IngressRule
		expected string
	}{{
		"single port firewall rule",
		firewall.NewIngressRule(network.MustParsePortRange("80/tcp")),
		"FROM tag switch TO tag juju ALLOW tcp PORT 80",
	}, {
		"multiple port firewall rule",
		firewall.NewIngressRule(network.MustParsePortRange("80-81/tcp")),
		"FROM tag switch TO tag juju ALLOW tcp ( PORT 80 AND PORT 81 )",
	}}

	for i, t := range testCases {
		c.Logf("test case %d: %s", i, t.about)
		rule := joyent.CreateFirewallRuleAll("switch", t.rules)
		c.Check(rule, gc.Equals, t.expected)
	}
}
