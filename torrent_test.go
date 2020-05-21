package main

import (
	"github.com/stretchr/testify/assert"
	"strings"
	"testing"
)

//body := strings.NewReader("d8:intervali900e5:peers2:ip14:23.160.192.1514:porti0e2:ip13:87.123.55.2554:porti0e2:ip14:213.174.18.1874:porti0e2:ip13:172.83.40.1084:porti0e2:ip13:185.21.217.614:porti0e2:ip14:173.20.201.1384:porti0e2:ip13:185.148.3.2464:porti0e2:ip12:65.102.186.24:porti0e2:ip15:199.116.115.1444:porti0e2:ip13:91.121.106.864:porti0e2:ip14:62.210.244.1534:porti0e2:ip12:212.47.231.24:porti0e2:ip14:176.186.51.2464:porti0e2:ip12:91.67.221.674:porti0e2:ip12:193.25.6.2064:porti0e2:ip14:109.123.118.454:porti0e2:ip13:90.150.185.814:porti0e2:ip14:188.98.215.1064:porti0e2:ip11:70.95.30.364:porti0e2:ip15:212.195.128.1344:porti0e2:ip13:71.198.208.744:porti0e2:ip12:87.101.95.994:porti0e2:ip12:89.23.196.374:porti0e2:ip14:82.197.213.1324:porti0e2:ip10:83.91.6.234:porti0e2:ip13:75.127.14.1404:porti0e2:ip14:67.253.193.1544:porti0e2:ip13:172.83.40.1084:porti0e2:ip14:71.172.157.1724:porti0e2:ip13:91.239.36.1474:porti0e2:ip14:90.125.159.2474:porti0e2:ip12:69.250.49.854:porti0e2:ip12:107.202.43.24:porti0e2:ip15:185.172.235.1904:porti0e2:ip13:83.233.178.834:porti0e2:ip13:5.103.184.1014:porti0e2:ip14:89.154.182.2344:porti0e2:ip13:24.153.41.1614:porti0e2:ip13:188.254.61.504:porti0e2:ip12:212.8.50.1424:porti0e2:ip15:141.239.111.1854:porti0e2:ip14:95.129.164.1664:porti0e2:ip14:194.187.249.474:porti0e2:ip13:104.232.116.84:porti0e2:ip11:74.101.2.854:porti0e2:ip13:2.236.113.1664:porti0e2:ip14:99.199.119.1364:porti0e2:ip12:155.4.235.604:porti0e2:ip12:115.186.3.794:porti0e2:ip14:185.45.195.1644:porti0eee")

func Test_func_parse_response_exists(t *testing.T) {
	parseResponse(strings.NewReader(""))
}

func Test_given_bee_code_with_interval_then_parse(t *testing.T) {
	reader := strings.NewReader("d8:intervali900e")

	response := parseResponse(reader)

	assert.Equal(t, 900, response.Interval)
}

func Test_given_bee_code_with_peer_then_parse(t *testing.T) {
	reader := strings.NewReader("d8:intervali900e5:peers2:ip14:23.160.192.1514:porti0e")

	response := parseResponse(reader)

	t.Log(response)

	//assert.Equal(t, "", response)
}
