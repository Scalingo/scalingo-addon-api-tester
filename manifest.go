package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/url"
	"os"
	"reflect"
	"strings"
	"sync"
)

type Plan struct {
	Name        string            `json:"name"`
	DisplayName string            `json:"display_name"`
	Price       float64           `json:"price"`
	Description string            `json:"description"`
	Options     map[string]string `json:"options"`
}

type Manifest struct {
	Name             string   `json:"name"`
	Username         string   `json:"username"`
	Password         string   `json:"password"`
	SSOSalt          string   `json:"sso_salt"`
	LogoURL          string   `json:"logo_url"`
	ShortDescription string   `json:"short_description"`
	Description      string   `json:"description"`
	ConfigVars       []string `json:"config_vars"`
	LogDrain         bool     `json:"log_drain"`
	Production       struct {
		BaseURL string `json:"base_url"`
		SSOURL  string `json:"sso_url"`
	} `json:"production"`
	Test struct {
		BaseURL string `json:"base_url"`
		SSOURL  string `json:"sso_url"`
	} `json:"test"`
	Plans []Plan `json:"plans"`
}

func readManifest(path string) (*Manifest, error) {
	fd, err := os.OpenFile(path, os.O_RDONLY, 0644)
	if err != nil {
		return nil, err
	}
	defer fd.Close()

	manifest = &Manifest{}
	err = json.NewDecoder(fd).Decode(&manifest)
	if err != nil {
		return nil, fmt.Errorf("fail to parse JSON: %v", err)
	}

	return manifest, nil
}

func (m *Manifest) PlanExist(p string) bool {
	for _, plan := range m.Plans {
		if plan.Name == p {
			return true
		}
	}
	return false
}

func (m *Manifest) PlanOptions(p string) map[string]string {
	for _, plan := range m.Plans {
		if plan.Name == p {
			return plan.Options
		}
	}
	return map[string]string{}
}

func (m *Manifest) CheckAddonConfig(config map[string]string) error {
	var i int
	for name, _ := range config {
		i = 0
		for _, v := range m.ConfigVars {
			if name == v {
				break
			}
			i++
		}
		if i == len(m.ConfigVars) {
			return fmt.Errorf("%v is not part of the config_vars array of the manifest", name)
		}
	}
	return nil
}

func (m *Manifest) Check() error {
	var errorMessages []string
	errChan := make(chan error)
	wg := &sync.WaitGroup{}
	wg.Add(12)

	go checkURL(wg, errChan, m.Production.BaseURL, "production.base_url")
	go checkURL(wg, errChan, m.Production.SSOURL, "production.sso_url")
	go checkURL(wg, errChan, m.Test.BaseURL, "test.base_url")
	go checkURL(wg, errChan, m.Test.SSOURL, "test.sso_url")
	go checkURL(wg, errChan, m.LogoURL, "logo_url")
	go checkString(wg, errChan, m.Username, "username")
	go checkString(wg, errChan, m.Password, "password")
	go checkString(wg, errChan, m.SSOSalt, "sso_salt")
	go checkString(wg, errChan, m.ShortDescription, "short_description")
	go checkString(wg, errChan, m.Description, "description")
	go checkArray(wg, errChan, m.Plans, "plans", "plan")
	go checkArray(wg, errChan, m.ConfigVars, "config_vars", "environment variable")

	go func() {
		wg.Wait()
		close(errChan)
	}()

	for err := range errChan {
		errorMessages = append(errorMessages, err.Error())
	}

	if len(errorMessages) > 0 {
		return errors.New("  → " + strings.Join(errorMessages, "\n  → "))
	} else {
		return nil
	}
}

func checkURL(wg *sync.WaitGroup, c chan error, raw string, name string) {
	defer wg.Done()
	_, err := url.Parse(raw)
	if err != nil {
		c <- fmt.Errorf("%v is not a valid URL: %v", name, err)
	}
}

func checkString(wg *sync.WaitGroup, c chan error, s string, name string) {
	defer wg.Done()
	if s == "" {
		c <- fmt.Errorf("%v can't be blank", name)
	}
}

func checkArray(wg *sync.WaitGroup, c chan error, array interface{}, name string, elemName string) {
	defer wg.Done()
	if array == nil {
		c <- fmt.Errorf("%v should be defined", name)
	} else if reflect.ValueOf(array).Len() == 0 {
		c <- fmt.Errorf("%v should have at least one element %v", name, elemName)
	}
}
