package service

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
)

func (s *Service) loginQBit() (err error) {
	client := &http.Client{}
	loginData := []byte(fmt.Sprintf("username=%s&password=%s", s.settings.QBitUser, s.settings.QBitPass))

	req, err := http.NewRequest("POST", s.settings.QBitURL+"/api/v2/auth/login", bytes.NewBuffer(loginData))
	if err != nil {
		return ErrQBitLogin
	}

	req.Header.Set("Referer", s.settings.QBitURL)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp, err := client.Do(req)
	if err != nil {
		s.logger.Info("login " + err.Error())
		return ErrQBitLogin
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		s.logger.Info("http err " + resp.Status)
		return ErrQBitLogin
	}

	s.qBitClient = client
	s.qBitCookie = resp.Header.Get("Set-Cookie")
	return nil
}

func (s *Service) updateQBitPort(port uint16) (err error) {
	if s.settings.QBitURL == "" {
		return nil
	}

	if s.qBitClient == nil || s.qBitCookie == "" {
		err = s.loginQBit()
		if err != nil {
			return err
		}
	}

	preferences := map[string]interface{}{
		"listen_port": port,
		"upnp":        false,
		"random_port": false,
	}
	jsonData, err := json.Marshal(preferences)
	if err != nil {
		return err
	}
	data := []byte("json=" + url.QueryEscape(string(jsonData)))

	req, err := http.NewRequest("POST", s.settings.QBitURL+"/api/v2/app/setPreferences", bytes.NewBuffer(data))
	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Content-Length", fmt.Sprintf("%d", len(data)))
	req.Header.Set("Cookie", s.qBitCookie)
	// req.PostForm = form

	resp, err := s.qBitClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to change listen port")
	}

	return nil
}
