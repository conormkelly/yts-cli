package transcript

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
)

type Fetcher struct {
	pythonScript string
	venvPath     string
}

type Response struct {
	Transcript string `json:"transcript,omitempty"`
	Error      string `json:"error,omitempty"`
}

// ensureVirtualEnv creates and sets up a virtual environment with required dependencies
func ensureVirtualEnv() (string, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("failed to get home directory: %w", err)
	}

	// Create .yts directory in user's home if it doesn't exist
	ytsDir := filepath.Join(homeDir, ".yts")
	if err := os.MkdirAll(ytsDir, 0755); err != nil {
		return "", fmt.Errorf("failed to create .yts directory: %w", err)
	}

	venvPath := filepath.Join(ytsDir, "venv")
	pythonPath := filepath.Join(venvPath, "bin", "python3")

	// Check if venv already exists and has the package
	if _, err := os.Stat(pythonPath); err == nil {
		// Try importing the package
		cmd := exec.Command(pythonPath, "-c", "import youtube_transcript_api")
		if cmd.Run() == nil {
			return venvPath, nil // Venv exists and package is installed
		}
	}

	fmt.Println("Setting up Python virtual environment...")

	// Create virtual environment
	cmd := exec.Command("python3", "-m", "venv", venvPath)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("failed to create virtual environment: %w", err)
	}

	// Install package in virtual environment
	fmt.Println("Installing youtube_transcript_api...")
	pipPath := filepath.Join(venvPath, "bin", "pip3")
	cmd = exec.Command(pipPath, "install", "youtube-transcript-api")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("failed to install dependencies: %w", err)
	}

	return venvPath, nil
}

// NewFetcher creates a new transcript fetcher and sets up the Python environment
func NewFetcher() (*Fetcher, error) {
	venvPath, err := ensureVirtualEnv()
	if err != nil {
		return nil, err
	}

	script := `
import sys
import json
import re
from youtube_transcript_api import YouTubeTranscriptApi, NoTranscriptFound, TranscriptsDisabled

def extract_video_id(url):
    patterns = [
        r'(?:v=|\/)([0-9A-Za-z_-]{11}).*',
        r'(?:youtu\.be\/)([0-9A-Za-z_-]{11})',
    ]
    for pattern in patterns:
        match = re.search(pattern, url)
        if match:
            return match.group(1)
    return None

def main():
    try:
        video_id = extract_video_id(sys.argv[1])
        if not video_id:
            print(json.dumps({"error": "Invalid YouTube URL"}))
            sys.exit(1)

        transcript = YouTubeTranscriptApi.get_transcript(video_id)
        text = "\n".join(entry["text"] for entry in transcript)
        print(json.dumps({"transcript": text}))

    except NoTranscriptFound:
        print(json.dumps({"error": "No transcript available for this video"}))
        sys.exit(1)
    except TranscriptsDisabled:
        print(json.dumps({"error": "Transcripts are disabled for this video"}))
        sys.exit(1)
    except Exception as e:
        print(json.dumps({"error": str(e)}))
        sys.exit(1)

if __name__ == "__main__":
    main()
`
	tmpDir := os.TempDir()
	scriptPath := filepath.Join(tmpDir, "yt_transcript_fetcher.py")
	if err := os.WriteFile(scriptPath, []byte(script), 0644); err != nil {
		return nil, fmt.Errorf("failed to create script file: %w", err)
	}

	return &Fetcher{
		pythonScript: scriptPath,
		venvPath:     venvPath,
	}, nil
}

// Fetch retrieves the transcript for a given YouTube video URL
func (f *Fetcher) Fetch(videoURL string) (string, error) {
	pythonPath := filepath.Join(f.venvPath, "bin", "python3")
	cmd := exec.Command(pythonPath, f.pythonScript, videoURL)
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	if err := cmd.Run(); err != nil {
		// Check if we have a structured error response
		var response Response
		if err := json.Unmarshal(stdout.Bytes(), &response); err == nil && response.Error != "" {
			return "", fmt.Errorf("transcript error: %s", response.Error)
		}
		// Fall back to command error
		return "", fmt.Errorf("failed to fetch transcript: %w\nStderr: %s", err, stderr.String())
	}

	var response Response
	if err := json.Unmarshal(stdout.Bytes(), &response); err != nil {
		return "", fmt.Errorf("failed to parse transcript response: %w", err)
	}

	if response.Error != "" {
		return "", fmt.Errorf("transcript error: %s", response.Error)
	}

	return response.Transcript, nil
}

// Cleanup removes the temporary Python script
func (f *Fetcher) Cleanup() error {
	return os.Remove(f.pythonScript)
}
