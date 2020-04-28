package metadata

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
)

func TestMetadata(t *testing.T) {
	m := http.NewServeMux()
	s := httptest.NewServer(m)
	defer s.Close()
	client := NewClient()
	client.SetBaseURL(s.URL)

	m.HandleFunc("/v1.json", func(writer http.ResponseWriter, request *http.Request) {
		response := `{
    "bgp": {
        "ipv4": {
            "my-address": "192.0.2.2",
            "my-asn": "64496",
            "peer-address": "169.254.169.254",
            "peer-asn": "64515"
        },
        "ipv6": {
            "my-address": "2001:0db8:7400:7ccf:5428:d5ff:fe28:1910",
            "my-asn": "64496",
            "peer-address": "2001:19f0:ffff::1",
            "peer-asn": "64515"
        }
    },
    "hostname": "vultr-guest",
    "instanceid": "a747bfz6385e",
    "public-keys": "ssh-rsaroot@example2",
    "region": {
        "regioncode": "EWR"
    }
}`
		fmt.Fprint(writer, response)
	})

	res, err := client.Metadata()
	if err != nil {
		t.Errorf("Metadata returned error : %v", err)
	}

	expected := &MetaData{
		Hostname:   "vultr-guest",
		InstanceID: "a747bfz6385e",
		PublicKeys: "ssh-rsaroot@example2",
		Region: struct {
			RegionCode string `json:"regioncode,omitempty"`
		}{"EWR"},
		BGP: struct {
			IPv4 struct {
				MyAddress   string `json:"my-address,omitempty"`
				MyASN       string `json:"my-asn,omitempty"`
				PeerAddress string `json:"peer-address,omitempty"`
				PeerASN     string `json:"peer-asn,omitempty"`
			} `json:"ipv4,omitempty"`
			IPv6 struct {
				MyAddress   string `json:"my-address,omitempty"`
				MyASN       string `json:"my-asn,omitempty"`
				PeerAddress string `json:"peer-address,omitempty"`
				PeerASN     string `json:"peer-asn,omitempty"`
			} `json:"ipv6,omitempty"`
		}{
			IPv4: struct {
				MyAddress   string `json:"my-address,omitempty"`
				MyASN       string `json:"my-asn,omitempty"`
				PeerAddress string `json:"peer-address,omitempty"`
				PeerASN     string `json:"peer-asn,omitempty"`
			}{"192.0.2.2", "64496", "169.254.169.254", "64515"},
			IPv6: struct {
				MyAddress   string `json:"my-address,omitempty"`
				MyASN       string `json:"my-asn,omitempty"`
				PeerAddress string `json:"peer-address,omitempty"`
				PeerASN     string `json:"peer-asn,omitempty"`
			}{"2001:0db8:7400:7ccf:5428:d5ff:fe28:1910", "64496", "2001:19f0:ffff::1", "64515"},
		},
	}

	if !reflect.DeepEqual(expected, res) {
		t.Errorf("Expected - %v : actual %v", expected, res)
	}

}

func TestRegionCodeToID(t *testing.T) {
	regionID := RegionCodeToID("EWR")

	if regionID != "1" {
		t.Error("regionID does not match up")
	}
}
