# YouTube Transcript Summarizer (YTS) CLI

A command-line tool that fetches YouTube video transcripts and generates concise, well-structured summaries using AI. Perfect for quickly understanding video content without watching the full video.

## Features

- 🎥 Fetch transcripts from any YouTube video with available captions
- 🤖 Generate AI-powered summaries using local LLM models
- 📝 Multiple summary types (short, medium, long)
- 💾 Save summaries to an output file
- 🌐 Support for videos with auto-generated captions
- ⚡ Streaming output for real-time summary generation
- 📄 Output formatted transcript / save to file

## Dependencies

- Python 3.8+
- LM Studio installation
- Any model on LM studio e.g. `llama-3.2-3b-instruct`

## Installation

If you have Python 3.8+, LM studio and a model installed you should be good to go!

You can download the prebuilt binaries.

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

The first time you run YTS, it will automatically set up a Python virtual environment and install required dependencies.

## Usage

### Basic Usage

Summarize a YouTube video:

```bash
yts https://www.youtube.com/watch?v=dQw4w9WgXcQ
```

### Summary Types

Generate different summary lengths (default is medium):

```bash
# Short summary
yts https://www.youtube.com/watch?v=dQw4w9WgXcQ -s

# Long summary
yts https://www.youtube.com/watch?v=dQw4w9WgXcQ -l
```

### Save to File

Save the summary to a Markdown file:

```bash
yts https://www.youtube.com/watch?v=dQw4w9WgXcQ --output summary.txt
```

### Transcript Formatting

To get a formatted version of the raw transcript without summarization:

```bash
# Display formatted transcript
yts transcript https://www.youtube.com/watch?v=dQw4w9WgXcQ

# Save formatted transcript to file
yts transcript https://www.youtube.com/watch?v=dQw4w9WgXcQ --output my-transcript.txt
```

### Configuration

YTS uses a configuration file located in a platform-specific directory.

You can customize:

- LLM base URL / port
- Model selection
- Output format
- Default summary type (short, medium [default], long)
- The short, medium, and long summary system prompts
- Transcript system prompt

You can view the current config file path and values with `yts config`.

### Environment Variables

You can override configuration settings using environment variables:

- `YTS_LLM_URL`: Override the LLM API endpoint
- `YTS_MODEL`: Override the model selection

## How It Works

1. **Transcript Fetching**: YTS uses the `youtube-transcript-api` Python package to fetch video transcripts. It automatically handles auto-generated captions.

2. **AI Processing**: The transcript is processed using a local LLM (default: llama-3.2-3b-instruct) to generate a coherent summary. The summary is streamed in real-time to provide immediate feedback.

3. **Output Generation**: The summary is formatted in Markdown and displayed in the terminal, optionally being saved to a file.

## Development

### Project Structure

```txt
.
├── cmd/                  # Command implementations
├── internal/
│   ├── config/           # Configuration management
│   ├── llm/              # LLM client and processing
│   └── transcript/       # Transcript fetching
├── main.go               # Entry point
└── Makefile              # Build and development commands
```

### Building

- Build for current platform: `make build`
- Create release binaries: `make release`
- Install globally: `make install`

## Contributing

1. Fork the repository
2. Create your feature branch: `git checkout -b feature/amazing-feature`
3. Open a pull request

## License

This project is licensed under the MIT License - see the LICENSE file for details.

## Acknowledgments

- Uses [youtube-transcript-api](https://github.com/jdepoix/youtube-transcript-api) for transcript fetching
- Built with [Cobra](https://github.com/spf13/cobra) CLI framework
- Configuration managed with [Viper](https://github.com/spf13/viper)
