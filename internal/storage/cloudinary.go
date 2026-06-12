package storage

import (
	"bytes"
	"crypto/sha1"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"sort"
	"strconv"
	"strings"
	"time"
)

// Client uploads task attachments to Cloudinary (free tier compatible).
type Client struct {
	cloudName string
	apiKey    string
	apiSecret string
	folder    string
	http      *http.Client
}

func NewCloudinary(cloudName, apiKey, apiSecret, folder string) *Client {
	if folder == "" {
		folder = "tasktheteam"
	}
	return &Client{
		cloudName: cloudName,
		apiKey:    apiKey,
		apiSecret: apiSecret,
		folder:    folder,
		http:      &http.Client{Timeout: 60 * time.Second},
	}
}

type uploadResponse struct {
	SecureURL    string `json:"secure_url"`
	PublicID     string `json:"public_id"`
	ResourceType string `json:"resource_type"`
	Bytes        int64  `json:"bytes"`
	Error        *struct {
		Message string `json:"message"`
	} `json:"error"`
}

func (c *Client) Upload(file *multipart.FileHeader) (*UploadResult, error) {
	src, err := file.Open()
	if err != nil {
		return nil, err
	}
	defer src.Close()

	timestamp := strconv.FormatInt(time.Now().Unix(), 10)
	signature := c.sign(map[string]string{
		"folder":    c.folder,
		"timestamp": timestamp,
	})

	var body bytes.Buffer
	writer := multipart.NewWriter(&body)

	_ = writer.WriteField("api_key", c.apiKey)
	_ = writer.WriteField("timestamp", timestamp)
	_ = writer.WriteField("folder", c.folder)
	_ = writer.WriteField("signature", signature)

	part, err := writer.CreateFormFile("file", file.Filename)
	if err != nil {
		return nil, err
	}
	if _, err := io.Copy(part, src); err != nil {
		return nil, err
	}
	if err := writer.Close(); err != nil {
		return nil, err
	}

	url := fmt.Sprintf("https://api.cloudinary.com/v1_1/%s/auto/upload", c.cloudName)
	req, err := http.NewRequest(http.MethodPost, url, &body)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", writer.FormDataContentType())

	res, err := c.http.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	raw, _ := io.ReadAll(res.Body)
	var parsed uploadResponse
	if err := json.Unmarshal(raw, &parsed); err != nil {
		return nil, fmt.Errorf("cloudinary: invalid response")
	}
	if parsed.Error != nil {
		return nil, fmt.Errorf("cloudinary: %s", parsed.Error.Message)
	}
	if parsed.SecureURL == "" {
		return nil, fmt.Errorf("cloudinary: upload failed (status %d)", res.StatusCode)
	}

	return &UploadResult{
		URL:          parsed.SecureURL,
		PublicID:     parsed.PublicID,
		ResourceType: parsed.ResourceType,
		Bytes:        parsed.Bytes,
	}, nil
}

func (c *Client) Delete(publicID, resourceType string) error {
	if publicID == "" {
		return nil
	}
	if resourceType == "" {
		resourceType = "image"
	}

	timestamp := strconv.FormatInt(time.Now().Unix(), 10)
	signature := c.sign(map[string]string{
		"public_id": publicID,
		"timestamp": timestamp,
	})

	form := urlValues(map[string]string{
		"api_key":   c.apiKey,
		"public_id": publicID,
		"timestamp": timestamp,
		"signature": signature,
	})

	url := fmt.Sprintf("https://api.cloudinary.com/v1_1/%s/%s/destroy", c.cloudName, resourceType)
	req, err := http.NewRequest(http.MethodPost, url, strings.NewReader(form))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	res, err := c.http.Do(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()
	return nil
}

func (c *Client) sign(params map[string]string) string {
	keys := make([]string, 0, len(params))
	for k := range params {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	var parts []string
	for _, k := range keys {
		parts = append(parts, k+"="+params[k])
	}
	toSign := strings.Join(parts, "&") + c.apiSecret
	sum := sha1.Sum([]byte(toSign))
	return hex.EncodeToString(sum[:])
}

func urlValues(values map[string]string) string {
	keys := make([]string, 0, len(values))
	for k := range values {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	var parts []string
	for _, k := range keys {
		parts = append(parts, k+"="+values[k])
	}
	return strings.Join(parts, "&")
}
