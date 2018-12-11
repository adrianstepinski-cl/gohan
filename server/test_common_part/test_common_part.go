package test_common_part

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"

	"github.com/cloudwan/gohan/schema"
	. "github.com/onsi/gomega"
)

const (
	BaseURL        = "http://localhost:19090"
	HealthcheckURL = BaseURL + "/gohan/v0.1/healthcheck"
)

const (
	AdminTokenID      = "admin_token"
	MemberTokenID     = "demo_token"
	PowerUserTokenID  = "power_user_token"
	AdminTenantID     = "fc394f2ab2df4114bde39905f800dc57"
	MemberTenantID    = "fc394f2ab2df4114bde39905f800dc57"
	PowerUserTenantID = "acf5662bbff44060b93ac3db3c25a590"
)

func GetNetwork(color string, tenant string) map[string]interface{} {
	return map[string]interface{}{
		"id":                "network" + color,
		"name":              "Network" + color,
		"description":       "The " + color + " Network",
		"tenant_id":         tenant,
		"route_targets":     []string{"1000:10000", "2000:20000"},
		"shared":            false,
		"providor_networks": map[string]interface{}{"segmentation_id": 12, "segmentation_type": "vlan"},
		"config": map[string]interface{}{
			"default_vlan": map[string]interface{}{
				"vlan_id": float64(1),
				"name":    "default_vlan",
			},
			"empty_vlan": map[string]interface{}{},
			"vpn_vlan": map[string]interface{}{
				"name": "vpn_vlan",
			},
		},
	}
}

func GetSubnet(color string, tenant string, parent string) map[string]interface{} {
	return map[string]interface{}{
		"id":          "subnet" + color,
		"name":        "Subnet" + color,
		"description": "The " + color + " Subnet",
		"tenant_id":   tenant,
		"cidr":        "10.0.0.0/24",
		"network_id":  parent,
	}
}

func GetCity(name string) map[string]interface{} {
	return map[string]interface{}{
		"id":   "city" + name,
		"name": "City" + name,
	}
}

func GetSchool(name, cityID string) map[string]interface{} {
	return map[string]interface{}{
		"id":      "school" + name,
		"name":    "School" + name,
		"city_id": cityID,
	}
}

func GetChild(name, schoolID string) map[string]interface{} {
	return map[string]interface{}{
		"id":        name,
		"school_id": schoolID,
	}
}

func GetParent(name, boyID, girlID string) map[string]interface{} {
	return map[string]interface{}{
		"id":      "parent" + name,
		"boy_id":  boyID,
		"girl_id": girlID,
	}
}

func GetNetworkSingularURL(color string) string {
	s, _ := schema.GetManager().Schema("network")
	return BaseURL + s.URL + "/network" + color
}

func GetServerSingularURL(color string) string {
	s, _ := schema.GetManager().Schema("server")
	return BaseURL + s.URL + "/server" + color
}

func GetSubnetSingularURL(color string) string {
	s, _ := schema.GetManager().Schema("subnet")
	return BaseURL + s.URL + "/subnet" + color
}

func GetSubnetFullSingularURL(networkColor, subnetColor string) string {
	return GetSubnetFullPluralURL(networkColor) + "/subnet" + subnetColor
}

func GetSubnetFullPluralURL(networkColor string) string {
	s, _ := schema.GetManager().Schema("network")
	return BaseURL + s.URL + "/network" + networkColor + "/subnets"
}

func TestURL(method, url, token string, postData interface{}, expectedCode int) interface{} {
	data, resp := HttpRequest(method, url, token, postData)
	jsonData, _ := json.MarshalIndent(data, "", "    ")
	ExpectWithOffset(1, resp.StatusCode).To(Equal(expectedCode), string(jsonData))
	return data
}

func TestURLErrorMessage(method, url, token string, postData interface{}, expectedCode int, expectedMessage string) interface{} {
	data, resp := HttpRequest(method, url, token, postData)
	jsonData, _ := json.MarshalIndent(data, "", "    ")
	ExpectWithOffset(1, resp.StatusCode).To(Equal(expectedCode), string(jsonData))
	Expect(data).To(HaveKeyWithValue("error", expectedMessage))
	return data
}

func HttpRequest(method, url, token string, postData interface{}) (interface{}, *http.Response) {
	client := &http.Client{}
	var reader io.Reader
	if postData != nil {
		jsonByte, err := json.Marshal(postData)
		Expect(err).ToNot(HaveOccurred())
		reader = bytes.NewBuffer(jsonByte)
	}
	request, err := http.NewRequest(method, url, reader)
	Expect(err).ToNot(HaveOccurred())
	request.Header.Set("X-Auth-Token", token)
	var data interface{}
	resp, err := client.Do(request)
	Expect(err).ToNot(HaveOccurred())
	defer resp.Body.Close()
	decoder := json.NewDecoder(resp.Body)
	decoder.Decode(&data)
	return data, resp
}
