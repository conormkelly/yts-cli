# YouTube Transcript Summarizer (YTS) CLI

A command-line tool that fetches YouTube video transcripts and generates concise, well-structured summaries using AI. Perfect for quickly understanding video content without watching the full video.

## Features

- 🎥 Fetch transcripts from any YouTube video with available captions
- 🤖 Generate AI-powered summaries using local LLM models
- 🔄 Support for multiple LLM providers (LM Studio, Ollama)
- 📝 Multiple summary types (short, long)
- 💾 Save summaries to an output file
- 🌐 Support for videos with auto-generated captions
- ⚡ Streaming output for real-time summary generation
- 📄 Output formatted transcripts

## Dependencies

You'll need one of the following LLM providers:

### LM Studio

- [LM Studio](https://lmstudio.ai/) installation
- Any compatible model

### Ollama

- [Ollama](https://ollama.ai/) installation
- Any compatible model (e.g., llama2, codellama, mistral)

## Installation

### Prebuilt binaries

You can download the latest release from the [releases page](https://github.com/conormkelly/yts-cli/releases/latest).

### From Source

#### Prerequisites

- Go >=1.23

1. Clone the repository:

   ```bash
   git clone https://github.com/conormkelly/yts-cli
   cd yts-cli
   ```

2. Build the binary:

   ```bash
   make build
   ```

3. Install globally (optional):

   ```bash
   make install
   ```

## Usage

### Basic Usage

Summarize a YouTube video:

```bash
yts https://www.youtube.com/watch?v=dQw4w9WgXcQ
```

### LLM Provider Selection

YTS supports multiple LLM providers. Configure your preferred provider:

```bash
# Use LM Studio (default)
yts https://www.youtube.com/watch?v=dQw4w9WgXcQ

# Use Ollama
yts --provider ollama https://www.youtube.com/watch?v=dQw4w9WgXcQ

# Or ...
yts -p ollama https://www.youtube.com/watch?v=dQw4w9WgXcQ
```

### Summary Types

Generate different summary lengths:

```bash
# Short summary (default)
yts https://www.youtube.com/watch?v=dQw4w9WgXcQ

# Long summary
yts https://www.youtube.com/watch?v=dQw4w9WgXcQ -l
```

### Save to File

Save the summary to a file:

```bash
yts https://www.youtube.com/watch?v=dQw4w9WgXcQ --output summary.txt
```

### Transcript Formatting

To get a formatted version of the raw transcript without summarization:

```bash
# Display formatted transcript
yts transcript https://www.youtube.com/watch?v=dQw4w9WgXcQ

# Save formatted transcript to file
yts transcript https://www.youtube.com/watch?v=dQw4w9WgXcQ -o my-transcript.txt
```

### Configuration

YTS uses a configuration file located in a platform-specific directory.

You can customize:

- LLM provider selection (lmstudio, ollama)
- Provider-specific settings:
  - LM Studio:
    - Base URL
    - Model selection
  - Ollama:
    - Base URL
    - Model selection
- Output format
- Summary system prompts
- Transcript system prompt

View current config:

```bash
yts config
```

### Environment Variables

Override configuration settings using environment variables:

- Provider Selection:
  - `YTS_PROVIDER`: Select LLM provider ("lmstudio" or "ollama")

- LM Studio Settings:
  - `YTS_LMSTUDIO_URL`: Override the LM Studio API endpoint
  - `YTS_LMSTUDIO_MODEL`: Override the model selection

- Ollama Settings:
  - `YTS_OLLAMA_URL`: Override the Ollama API endpoint
  - `YTS_OLLAMA_MODEL`: Override the model selection

## How It Works

1. **Transcript Fetching**: YTS fetches transcripts directly from YouTube using a pure Go implementation. It handles both manual and auto-generated captions, supporting multiple languages and formats. The fetcher:
   - Extracts caption metadata from the video page
   - Downloads the raw transcript XML
   - Processes and formats the captions into clean text

2. **AI Processing**: The transcript is processed using your chosen LLM provider to generate a coherent summary. The processing pipeline:
   - Sends a system prompt to the LLM based on the selected summary type (short/medium/long)
   - Streams completions for responsive feedback

3. **Output Generation**: The summary is displayed in the terminal, with:
   - Proper text formatting and sanitization
   - Real-time streaming output
   - Optional file saving

## Development

### Project Structure

```txt
.
├── cmd/                   # Command implementations
│   ├── config.go          # Configuration subcommand
│   ├── root.go            # Main command logic
│   ├── transcript.go      # Transcript subcommand
│   └── version.go         # Version subcommand
├── internal/
│   ├── config/            # Configuration management
│   │   └── config.go      # Config types and loading
│   ├── llm/               # LLM provider integration
│   │   ├── lmstudio.go    # LM Studio provider implementation
│   │   ├── ollama.go      # Ollama provider implementation
│   │   └── provider.go    # Provider interface and factory
│   └── transcript/        # Transcript processing
│       └── fetcher.go     # YouTube transcript fetching
├── .github/               # GitHub specific files
│   └── workflows/         # GitHub Actions workflows
├── main.go                # Entry point
├── Makefile               # Build and development commands
└── go.mod                 # Go module definition
```

### Building

- Build for current platform: `make build`
- Create release binaries: `make release`
- Install globally: `make install`

## Contributing

1. Fork the repository
2. Create your feature branch: `git checkout -b feature/amazing-feature`
3. Open a pull request

## Legal Notice

This tool accesses publicly available YouTube video transcripts. While I believe this falls within fair use, users should:

- Review YouTube's Terms of Service
- Use the tool responsibly
- Consider YouTube's official API for commercial applications

## License

This project is licensed under the MIT License - see the LICENSE file for details.

## Acknowledgments

- Transcript fetching approach inspired by [youtube-transcript-api](https://github.com/jdepoix/youtube-transcript-api) (reimplemented in pure Go)
- Built with [Cobra](https://github.com/spf13/cobra) CLI framework
- Configuration managed with [Viper](https://github.com/spf13/viper)
