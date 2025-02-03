package transcript

import (
	"encoding/json"
	"encoding/xml"
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

func (f *TranscriptFetcher) Fetch(videoURL string) ([]TranscriptResponse, error) {
	// 1. Extract video ID
	videoID, err := extractVideoID(videoURL)
	if err != nil {
		return nil, fmt.Errorf("invalid video ID: %w", err)
	}

	// 2. Fetch video page
	watchURL := fmt.Sprintf("https://www.youtube.com/watch?v=%s", videoID)
	req, err := http.NewRequest("GET", watchURL, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	// Add headers to mimic browser request
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/91.0.4472.124 Safari/537.36")
	req.Header.Set("Accept-Language", "en-US,en;q=0.9")

	resp, err := f.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch video page: %w", err)
	}
	defer resp.Body.Close()

	// 3. Extract captions JSON
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	captionsData, err := extractCaptionsJSON(string(body), videoID)
	if err != nil {
		return nil, fmt.Errorf("failed to extract captions: %w", err)
	}

	// Debug: Print available captions
	// fmt.Println("Available caption tracks:")
	// for _, track := range captionsData.CaptionTracks {
	// 	fmt.Printf("- Language: %s (%s)\n", track.Name.SimpleText, track.LanguageCode)
	// }

	if len(captionsData.CaptionTracks) == 0 {
		return nil, &ErrNoTranscriptFound{VideoID: videoID}
	}

	// 4. Fetch and parse transcript (using first available track)
	transcript, err := f.fetchTranscriptFromURL(captionsData.CaptionTracks[0].BaseURL)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch transcript: %w", err)
	}

	return transcript, nil
}

func extractCaptionsJSON(html string, videoID string) (*CaptionsData, error) {
	parts := strings.Split(html, `"captions":`)
	if len(parts) < 2 {
		return nil, &ErrTranscriptsDisabled{VideoID: videoID}
	}

	// Find the end of the captions JSON object
	jsonPart := strings.Split(parts[1], `,"videoDetails"`)[0]
	jsonPart = strings.TrimSpace(jsonPart)

	var captionsData struct {
		PlayerCaptionsTracklistRenderer CaptionsData `json:"playerCaptionsTracklistRenderer"`
	}

	if err := json.Unmarshal([]byte(jsonPart), &captionsData); err != nil {
		return nil, fmt.Errorf("failed to parse captions JSON: %w", err)
	}

	return &captionsData.PlayerCaptionsTracklistRenderer, nil
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
