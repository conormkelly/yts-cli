package transcript

import (
	"bytes"
	"encoding/json"
	"encoding/xml"
	"errors"
	"fmt"
	"html"
	"io"
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"time"
)

// Custom error types
type ErrTranscriptsDisabled struct{ VideoID string }
type ErrNoTranscriptFound struct{ VideoID string }

func (e ErrTranscriptsDisabled) Error() string {
	return fmt.Sprintf("transcripts are disabled for video: %s", e.VideoID)
}

func (e ErrNoTranscriptFound) Error() string {
	return fmt.Sprintf("no transcript found for video: %s", e.VideoID)
}

// InnerTubeContext represents the YouTube InnerTube API context
type InnerTubeContext struct {
	Client struct {
		ClientName    string `json:"clientName"`
		ClientVersion string `json:"clientVersion"`
	} `json:"client"`
}

// InnerTubeRequest represents the request to YouTube's InnerTube API
type InnerTubeRequest struct {
	Context InnerTubeContext `json:"context"`
	VideoID string           `json:"videoId"`
}

// InnerTubeResponse represents the response from YouTube's InnerTube API
type InnerTubeResponse struct {
	Captions struct {
		PlayerCaptionsTracklistRenderer struct {
			CaptionTracks []struct {
				BaseURL      string `json:"baseUrl"`
				Name         struct {
					Runs []struct {
						Text string `json:"text"`
					} `json:"runs"`
				} `json:"name"`
				LanguageCode string `json:"languageCode"`
				Kind         string `json:"kind,omitempty"`
			} `json:"captionTracks"`
		} `json:"playerCaptionsTracklistRenderer"`
	} `json:"captions"`
	PlayabilityStatus struct {
		Status string `json:"status"`
		Reason string `json:"reason,omitempty"`
	} `json:"playabilityStatus"`
}

// CaptionsData represents the YouTube captions JSON structure
type CaptionsData struct {
	CaptionTracks []struct {
		BaseURL string `json:"baseUrl"`
		Name    struct {
			SimpleText string `json:"simpleText"`
		} `json:"name"`
		LanguageCode string `json:"languageCode"`
		Kind         string `json:"kind"`
	} `json:"captionTracks"`
}

// TranscriptResponse represents a single caption entry
type TranscriptResponse struct {
	Text     string  `json:"text"`
	Start    float64 `json:"start"`
	Duration float64 `json:"duration"`
}

// TranscriptFetcher handles fetching transcripts from YouTube
type TranscriptFetcher struct {
	httpClient *http.Client
}

func NewTranscriptFetcher() *TranscriptFetcher {
	return &TranscriptFetcher{
		httpClient: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}

func extractVideoID(url string) (string, error) {
	patterns := []string{
		`(?:v=|\/)([0-9A-Za-z_-]{11}).*`,
		`(?:youtu\.be\/)([0-9A-Za-z_-]{11})`,
	}

	for _, pattern := range patterns {
		re := regexp.MustCompile(pattern)
		if matches := re.FindStringSubmatch(url); len(matches) > 1 {
			return matches[1], nil
		}
	}

	return "", fmt.Errorf("could not extract video ID from URL: %s", url)
}

func (f *TranscriptFetcher) Fetch(videoURL string) (string, []TranscriptResponse, error) {
	// 1. Extract video ID
	videoID, err := extractVideoID(videoURL)
	if err != nil {
		return "", nil, fmt.Errorf("invalid video ID: %w", err)
	}

	// 2. Fetch video page to get title and API key
	watchURL := fmt.Sprintf("https://www.youtube.com/watch?v=%s", videoID)
	req, err := http.NewRequest("GET", watchURL, nil)
	if err != nil {
		return "", nil, fmt.Errorf("failed to create request: %w", err)
	}

	// Add headers to mimic browser request
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/91.0.4472.124 Safari/537.36")
	req.Header.Set("Accept-Language", "en-US,en;q=0.9")

	resp, err := f.httpClient.Do(req)
	if err != nil {
		return "", nil, fmt.Errorf("failed to fetch video page: %w", err)
	}
	defer resp.Body.Close()

	// 3. Read and parse HTML
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", nil, fmt.Errorf("failed to read response body: %w", err)
	}

	htmlBody := string(body)

	// Extract title
	title, err := extractTitle(htmlBody)
	if err != nil {
		return "", nil, fmt.Errorf("failed to extract video title: %w", err)
	}

	// Extract InnerTube API key
	apiKey, err := extractInnerTubeAPIKey(htmlBody, videoID)
	if err != nil {
		return "", nil, fmt.Errorf("failed to extract API key: %w", err)
	}

	// 4. Fetch captions using InnerTube API
	innerTubeResp, err := f.fetchInnerTubeData(videoID, apiKey)
	if err != nil {
		return "", nil, fmt.Errorf("failed to fetch InnerTube data: %w", err)
	}

	// Check playability status
	if innerTubeResp.PlayabilityStatus.Status != "OK" {
		return "", nil, fmt.Errorf("video is not playable: %s - %s", innerTubeResp.PlayabilityStatus.Status, innerTubeResp.PlayabilityStatus.Reason)
	}

	// Extract caption tracks
	captionTracks := innerTubeResp.Captions.PlayerCaptionsTracklistRenderer.CaptionTracks
	if len(captionTracks) == 0 {
		return "", nil, &ErrNoTranscriptFound{VideoID: videoID}
	}

	// 5. Fetch and parse transcript (using first available track)
	// Remove &fmt=srv3 from the URL as the Python library does
	baseURL := strings.Replace(captionTracks[0].BaseURL, "&fmt=srv3", "", 1)
	transcript, err := f.fetchTranscriptFromURL(baseURL)
	if err != nil {
		return "", nil, fmt.Errorf("failed to fetch transcript: %w", err)
	}

	return title, transcript, nil
}

func extractTitle(htmlText string) (string, error) {
	parts := strings.Split(htmlText, "<title>")
	if len(parts) < 2 {
		return "", errors.New("could not find <title> tag in html")
	}
	// Find the end of the title
	rawTitle := strings.Split(parts[1], " - YouTube</title>")[0]
	return html.UnescapeString(rawTitle), nil
}

func extractInnerTubeAPIKey(html string, videoID string) (string, error) {
	// Pattern to extract the InnerTube API key from the HTML
	pattern := `"INNERTUBE_API_KEY":\s*"([a-zA-Z0-9_-]+)"`
	re := regexp.MustCompile(pattern)
	matches := re.FindStringSubmatch(html)

	if len(matches) > 1 {
		return matches[1], nil
	}

	// Check if there's a reCAPTCHA (IP blocked)
	if strings.Contains(html, `class="g-recaptcha"`) {
		return "", fmt.Errorf("IP blocked by YouTube (reCAPTCHA detected)")
	}

	return "", fmt.Errorf("could not extract InnerTube API key for video: %s", videoID)
}

func (f *TranscriptFetcher) fetchInnerTubeData(videoID string, apiKey string) (*InnerTubeResponse, error) {
	// Create InnerTube API request
	innerTubeReq := InnerTubeRequest{
		Context: InnerTubeContext{
			Client: struct {
				ClientName    string `json:"clientName"`
				ClientVersion string `json:"clientVersion"`
			}{
				ClientName:    "ANDROID",
				ClientVersion: "20.10.38",
			},
		},
		VideoID: videoID,
	}

	jsonData, err := json.Marshal(innerTubeReq)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal InnerTube request: %w", err)
	}

	// Make POST request to InnerTube API
	url := fmt.Sprintf("https://www.youtube.com/youtubei/v1/player?key=%s", apiKey)
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, fmt.Errorf("failed to create InnerTube request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept-Language", "en-US")

	resp, err := f.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch InnerTube data: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("InnerTube API returned status: %d", resp.StatusCode)
	}

	// Parse response
	var innerTubeResp InnerTubeResponse
	if err := json.NewDecoder(resp.Body).Decode(&innerTubeResp); err != nil {
		return nil, fmt.Errorf("failed to parse InnerTube response: %w", err)
	}

	return &innerTubeResp, nil
}

func (f *TranscriptFetcher) fetchTranscriptFromURL(url string) ([]TranscriptResponse, error) {
	resp, err := f.httpClient.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	// Parse XML response
	decoder := xml.NewDecoder(resp.Body)
	var transcript []TranscriptResponse

	for {
		token, err := decoder.Token()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, err
		}

		switch se := token.(type) {
		case xml.StartElement:
			if se.Name.Local == "text" {
				var entry TranscriptResponse
				// Get start and duration attributes
				for _, attr := range se.Attr {
					switch attr.Name.Local {
					case "start":
						entry.Start, _ = strconv.ParseFloat(attr.Value, 64)
					case "dur":
						entry.Duration, _ = strconv.ParseFloat(attr.Value, 64)
					}
				}
				// Get text content
				if textToken, err := decoder.Token(); err == nil {
					if text, ok := textToken.(xml.CharData); ok {
						entry.Text = html.UnescapeString(string(text))
					}
				}
				transcript = append(transcript, entry)
			}
		}
	}

	return transcript, nil
}
