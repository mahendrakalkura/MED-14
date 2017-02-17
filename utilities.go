package main

import (
	"bytes"
	"database/sql"
	"errors"
	"fmt"
	"github.com/BurntSushi/toml"
	"github.com/getsentry/raven-go"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"gopkg.in/xmlpath.v2"
	"io/ioutil"
	"math/rand"
	"net/http"
	"net/url"
	"strings"
	"time"
)

func get_database(settings *Settings) *sqlx.DB {
	dsn := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		settings.SQLX.Hostname,
		settings.SQLX.Port,
		settings.SQLX.Username,
		settings.SQLX.Password,
		settings.SQLX.Database,
	)
	database := sqlx.MustConnect("postgres", dsn)
	return database
}

func get_http_client(settings *Settings) *http.Client {
	client := &http.Client{}

	timeout := time.Duration(60 * time.Second)
	client.Timeout = timeout

	if len(settings.Proxies.Hostname) > 0 && len(settings.Proxies.Ports) > 0 {
		proxy := get_proxy(settings.Proxies.Hostname, settings.Proxies.Ports)
		proxy_url, err := url.Parse(proxy)
		if err != nil {
			raven.CaptureErrorAndWait(err, nil)
		}
		client.Transport = &http.Transport{Proxy: http.ProxyURL(proxy_url)}
	}

	return client
}

func get_proxy(hostname string, ports []int) string {
	port := get_random_number(ports[0], ports[1]+1)
	return fmt.Sprintf("https://%s:%d", hostname, port)
}

func get_random_number(minimum int, maximum int) int {
	rand.Seed(time.Now().Unix())
	return rand.Intn(maximum-minimum) + minimum
}

func get_results(settings *Settings, address Address) ([]Result, error) {
	var results []Result

	client := get_http_client(settings)

	data := url.Values{}
	data.Add("fin_txt", address.Street)

	request, new_request_err := http.NewRequest(
		"POST",
		"https://www.notariate.zh.ch/deu/ind_sea.php",
		bytes.NewBufferString(data.Encode()),
	)
	if new_request_err != nil {
		raven.CaptureErrorAndWait(new_request_err, nil)
		return results, errors.New("new_request_err")
	}

	request.Header.Add("Content-Type", "application/x-www-form-urlencoded; charset=UTF-8")
	request.Header.Add("DNT", "1")
	request.Header.Add("Host", "www.notariate.zh.ch")
	request.Header.Add("Origin", "https://www.notariate.zh.ch")
	request.Header.Add("Referer", "https://www.notariate.zh.ch/deu/ind_sea.php")
	request.Header.Add("User-Agent", "Go")
	request.Header.Add("X-Requested-With", "XMLHttpRequest")

	response, do_err := client.Do(request)
	if do_err != nil {
		raven.CaptureErrorAndWait(do_err, nil)
		return results, errors.New("do_err")
	}

	defer response.Body.Close()

	body_bytes, read_all_err := ioutil.ReadAll(response.Body)
	if read_all_err != nil {
		raven.CaptureErrorAndWait(read_all_err, nil)
		return results, errors.New("read_all_err")
	}

	body_string := string(body_bytes)

	if strings.Contains(body_string, "Es wurde kein Notariat gefunden.") {
		return results, nil
	}

	reader := strings.NewReader(body_string)
	root, parse_html_err := xmlpath.ParseHTML(reader)
	if parse_html_err != nil {
		raven.CaptureErrorAndWait(parse_html_err, nil)
		return results, errors.New("parse_html_err")
	}

	path := `//p[@class="conbld mt10"]`
	xpath := xmlpath.MustCompile(path)
	iter := xpath.Iter(root)
	for iter.Next() {
		result := Result{}

		result.AddressId = address.Id

		location_path := `./text()`
		location_xpath := xmlpath.MustCompile(location_path)
		location_value, location_ok := location_xpath.String(iter.Node())
		if location_ok {
			result.Location = location_value
			result.Location = strings.Replace(
				result.Location,
				"Mit dem Suchbegriff gefundene Strassen in ",
				"",
				-1,
			)
			result.Location = strings.Replace(result.Location, ":", "", -1)
		}

		office_name_path := `./following-sibling::ul/li/span[@class="conbld"]/a/text()`
		office_name_xpath := xmlpath.MustCompile(office_name_path)
		office_name_value, office_name_ok := office_name_xpath.String(iter.Node())
		if office_name_ok {
			result.OfficeName = office_name_value
		}

		office_url_path := `./following-sibling::ul/li/span[@class="conbld"]/a/@href`
		office_url_xpath := xmlpath.MustCompile(office_url_path)
		office_url_value, office_url_ok := office_url_xpath.String(iter.Node())
		if office_url_ok {
			result.OfficeUrl = office_url_value
			result.OfficeUrl = fmt.Sprintf("https://www.notariate.zh.ch%s", result.OfficeUrl)
		}

		quarter_path := `./following-sibling::ul/li/b/text()`
		quarter_xpath := xmlpath.MustCompile(quarter_path)
		quarter_value, quarter_ok := quarter_xpath.String(iter.Node())
		if quarter_ok {
			result.Quarter = quarter_value
		}

		street_path := `./following-sibling::ul/li`
		street_xpath := xmlpath.MustCompile(street_path)
		street_value, street_ok := street_xpath.String(iter.Node())
		if street_ok {
			result.Street = street_value
			result.Street = strings.TrimSpace(result.Street)
			result.Street = strings.Replace(result.Street, result.Quarter, "", -1)
			streets := strings.Split(result.Street, "\n")
			result.Street = streets[len(streets)-1]
			result.Street = strings.TrimSpace(result.Street)
		}

		results = append(results, result)
	}

	return results, nil
}

func get_settings() *Settings {
	var settings = &Settings{}
	_, err := toml.DecodeFile("settings.toml", settings)
	if err != nil {
		raven.CaptureErrorAndWait(err, nil)
		panic(err)
	}
	return settings
}

func get_text(text sql.NullString) string {
	if text.Valid {
		return text.String
	}
	return ""
}
