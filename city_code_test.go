package common

import (
	"log"

	"testing"
)

// just used to generate city-code map
func _TestGenCityCode(t *testing.T) {
	NAME2CODE := map[string]string{}
	for k, v := range CODE2NAME {
		NAME2CODE[v] = k
	}
	log.Printf("NAME2CODE:%#v", NAME2CODE)

	MDCG_CITY_CODE := map[string]string{}
	for k, v := range MDCG_CITY {
		MDCG_CITY_CODE[v] = k
	}
	log.Printf("MDCG_CITY_CODE:%#v", MDCG_CITY_CODE)

	OTHER_CITY_CODE := map[string]string{}
	for k, v := range OTHER_CITY {
		OTHER_CITY_CODE[v] = k
	}
	log.Printf("OTHER_CITY_CODE:%#v", OTHER_CITY_CODE)
}

func TestGetCodeFromCity(t *testing.T) {
	city := "白沙黎族自治县"
	code := GetCodeFromCity(city)
	if code != "469025" {
		t.Errorf("TestGetCodeFromCity err should be 469025")
	}
	city = "白沙黎族自治县22"
	code = GetCodeFromCity(city)
	if code != "" {
		t.Logf("TestGetCodeFromCity err should be empty", city, code)
	}
}
