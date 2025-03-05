# YouTube Transcript Summarizer (YTS) CLI

A powerful command-line tool that leverages AI to turn YouTube video transcripts into concise, well-structured summaries. Perfect for researchers, content creators, and anyone who wants to quickly understand video content without watching the full video.

## ✨ Features

- 🎯 Extract transcripts from any YouTube video with available captions
- 🤖 Generate AI-powered summaries using your choice of local or cloud LLMs
- 🔍 Ask specific questions about video content with the query mode
- 📝 Multiple summary formats (concise or detailed analysis)
- 🌍 Support for videos with auto-generated captions
- ⚡ Real-time streaming output as summaries are generated
- 💾 Save summaries and transcripts to files
- ⚙️ Extensive configuration options
- 🔒 Secure API key management

## 🚀 Quick Start

1. Install YTS:

   ```bash
   # Linux/macOS
   curl -sSL https://raw.githubusercontent.com/conormkelly/yts-cli/main/install.sh | bash
   ```

   ```powershell
   # Windows
   irm https://raw.githubusercontent.com/conormkelly/yts-cli/main/install.ps1 | iex
   ```

2. Choose and set up a provider:

   ```bash
   # For local providers (free, runs on your machine):
   yts config set provider lmstudio  # or ollama
   ```

   ```bash
   # For cloud providers (requires API key):
   yts config set provider claude
   yts apikey set claude your-api-key
   ```

3. Generate your first summary:

   ```bash
   yts https://www.youtube.com/watch?v=dQw4w9WgXcQ
   ```

## 🔧 Dependencies

Choose at least one of these LLM providers:

### Local Providers

These run models on your own machine:

#### LM Studio

- Download from [lmstudio.ai](https://lmstudio.ai/)
- Compatible with most GGUF models
- Free to use, runs locally

#### Ollama

- Download from [ollama.ai](https://ollama.ai/)
- Works with llama2, codellama, mistral, etc.
- Free to use, runs locally

### Cloud Providers

These require API keys:

#### Claude (Anthropic)

- Get API key from [anthropic.com/api](https://www.anthropic.com/api)

#### OpenAI

- Get API key from [platform.openai.com](https://platform.openai.com/api-keys)

## 📦 Installation

### Option 1: Prebuilt Binaries (Recommended)

#### Unix-like Systems (Linux/macOS)

```bash
curl -sSL https://raw.githubusercontent.com/conormkelly/yts-cli/main/install.sh | bash
```

#### Windows

```powershell
irm https://raw.githubusercontent.com/conormkelly/yts-cli/main/install.ps1 | iex
```

### Option 2: Build from Source

Prerequisites:

- Go >=1.23

```bash
# Clone the repository
git clone https://github.com/conormkelly/yts-cli
cd yts-cli

# Build the binary
make build

# Optional: Install globally
make install
```

## 🛠️ Basic Usage

### Generate Summaries

```bash
# Basic summary (default: short)
yts https://www.youtube.com/watch?v=video_id

# Long summary
yts -l https://www.youtube.com/watch?v=video_id

# Save to file
yts https://www.youtube.com/watch?v=video_id -o summary.txt
```

### Query Video Content

Ask specific questions about a video's content:

```bash
# Ask a question about the video
yts https://www.youtube.com/watch?v=video_id -q "Does this video explain quantum computing?"

# Save the query result to a file
yts https://www.youtube.com/watch?v=video_id -q "What are the main points about climate change?" -o answer.txt
```

Query mode uses the video's title and transcript as context to provide accurate answers based solely on the video content. This is particularly useful for:

- Fact-checking clickbait titles
- Finding specific information in long videos
- Evaluating video content before watching
- Extracting technical details from educational content

#### Example Query Output

```txt
Title: Understanding Climate Feedback Loops

Question: What are the three main positive feedback loops mentioned?

The three main positive feedback loops mentioned in the video are:

1. Ice-Albedo Feedback Loop: As ice melts due to warming, dark ocean water is exposed, which absorbs more heat than reflective ice, leading to further warming and more ice melt.

2. Water Vapor Feedback Loop: As temperatures rise, more water evaporates into the atmosphere. Since water vapor is a greenhouse gas, this increases warming, creating a self-reinforcing cycle.

3. Permafrost Methane Release: As permafrost thaws in Arctic regions, it releases trapped methane, which is a potent greenhouse gas that causes additional warming, leading to more permafrost thaw.

The video emphasizes that these positive feedback loops have the potential to accelerate climate change beyond current predictions if they reach tipping points.
```

### Example Output

#### Concise Summary

```txt
Title: Understanding Quantum Computing Basics

Core Message: Quantum computing harnesses quantum mechanical phenomena to solve 
specific problems exponentially faster than classical computers.

Key Points:
1. Qubits can exist in multiple states simultaneously through superposition
2. Quantum entanglement enables powerful parallel processing capabilities
3. Current practical limitations include decoherence and error correction
4. Most promising applications include cryptography and molecular simulation

Call to Action: Researchers encouraged to explore IBM's quantum computing cloud platform.
```

#### Detailed Analysis

```txt
Title: Understanding Quantum Computing Basics

1. Executive Summary
Comprehensive introduction to quantum computing fundamentals, explaining how quantum 
mechanics enables new computing paradigms. The presentation covers basic principles, 
current challenges, and practical applications.

2. Key Concepts Covered
- Quantum superposition and its role in computation
- Entanglement as a computational resource
- Quantum gates and circuit model
- Error correction challenges and solutions

[... continues with more detailed analysis ...]
```

### Transcript Formatting

Get video transcripts with various formatting options:

```bash
# Display formatted transcript (default)
yts transcript https://www.youtube.com/watch?v=dQw4w9WgXcQ

# Display raw unformatted transcript
yts transcript --raw https://www.youtube.com/watch?v=dQw4w9WgXcQ

# Include timestamps in the output
yts transcript --timestamps https://www.youtube.com/watch?v=dQw4w9WgXcQ

# Combine flags (raw output with timestamps)
yts transcript --raw --timestamps https://www.youtube.com/watch?v=dQw4w9WgXcQ

# Save to file (works with any flag combination)
yts transcript -o transcript.txt https://www.youtube.com/watch?v=dQw4w9WgXcQ
```

Flags:

- `-o, --output`: Save transcript to a file
- `-r, --raw`: Output raw transcript without AI formatting
- `-t, --timestamps`: Include timestamps in the output

#### Transcript Formatting Example

Before (raw transcript):

```txt
hey guys today were gonna talk about quantum computing its pretty cool and it uses
these things called qubits which are different from regular bits
```

After:

```txt
Hey guys! Today we're gonna talk about quantum computing. It's pretty cool and it 
uses these things called qubits, which are different from regular bits.
```

## ⚙️ Configuration

### Provider Selection

```bash
# Set default provider
yts config set provider lmstudio  # or: ollama, claude, openai

# Override for single command
yts -p ollama https://youtube.com/watch?v=video_id
```

### API Key Management

For cloud providers, securely store your API keys:

```bash
# Store API keys in system keyring
yts apikey set claude your-api-key
yts apikey set openai your-api-key

# Remove stored keys
yts apikey delete claude
yts apikey delete openai

# Verify key status
yts config view
```

### Configuration Management

```bash
# View current settings
yts config view

# Edit configuration file
yts config edit

# Set individual values
yts config set providers.claude.temperature 0.7
```

### Valid Configuration Paths

```bash
# Global
provider                           # Active provider selection
version                            # Configuration version

# Query Settings
queries.system_prompt              # Template for answering questions about videos

# LM Studio Settings
providers.lmstudio.base_url       # API endpoint
providers.lmstudio.model          # Model name

# Ollama Settings
providers.ollama.base_url         # API endpoint
providers.ollama.model            # Model name

# Claude Settings
providers.claude.model            # Model name
providers.claude.temperature      # Generation temperature (0.0-1.0)
providers.claude.max_tokens       # Maximum response tokens
providers.claude.timeout_seconds  # API timeout
providers.claude.max_retries      # Retry attempts

# OpenAI Settings
providers.openai.model           # Model name
providers.openai.temperature     # Generation temperature
providers.openai.max_tokens      # Maximum response tokens
providers.openai.timeout_seconds # API timeout
providers.openai.max_retries     # Retry attempts
providers.openai.organization_id # Optional org ID
```

### Configuration File Location

- Linux: `~/.config/yts/config.json`
- macOS: `~/Library/Application Support/yts/config.json`
- Windows: `%AppData%\yts\config.json`

### Environment Variables

Override settings using environment variables:

```bash
# Provider Selection
export YTS_PROVIDER=claude

# LM Studio
export YTS_LMSTUDIO_URL=http://localhost:1234
export YTS_LMSTUDIO_MODEL=llama-2-13b-chat

# Ollama
export YTS_OLLAMA_URL=http://localhost:11434
export YTS_OLLAMA_MODEL=mistral

# Claude
export YTS_CLAUDE_MODEL=claude-3-sonnet-20240229
export YTS_CLAUDE_TEMPERATURE=0.7
export YTS_CLAUDE_MAX_TOKENS=4096
export YTS_CLAUDE_TIMEOUT=120

# OpenAI
export YTS_OPENAI_MODEL=gpt-4
export YTS_OPENAI_TEMPERATURE=0.7
export YTS_OPENAI_MAX_TOKENS=4096
export YTS_OPENAI_TIMEOUT=120
export YTS_OPENAI_ORG_ID=org-...
```

## ❗ Troubleshooting

### Common Issues

1. "No transcript found"
   - Verify the video has captions enabled
   - For non-English videos, auto-translated captions may not be available
   - Some videos have disabled transcripts - try another video

2. "Provider not responding"
   - Local providers (LM Studio/Ollama):
     - Verify the service is running
     - Check the correct port is set
     - Ensure model is properly loaded
   - Cloud providers (Claude/OpenAI):
     - Verify API key is correct
     - Check internet connectivity
     - Confirm API service status

3. "Rate limiting/Quota exceeded"
   - Cloud providers: Check your API quota and limits
   - Consider switching to local providers for high-volume use
   - Implement exponential backoff in scripts

4. Performance Considerations
   - Large videos (>1 hour) may take longer to process
   - Local providers are generally slower but free
   - Cloud providers offer faster processing but incur costs
   - Network speed affects transcript download time

## 🔍 How It Works

1. **Transcript Fetching**
   - Extracts captions directly from YouTube
   - Supports both manual and auto-generated captions
   - Handles multiple languages and formats
   - Processes raw XML into clean text

2. **AI Processing**
   - Sends raw transcript to chosen LLM
   - Uses optimized prompts for different summary types
   - Streams completions for real-time feedback
   - Implements retry logic and error handling

3. **Output Handling**
   - Real-time streaming to terminal
   - Optional file output
   - Proper text formatting and sanitization
   - Error handling and logging

## 🛡️ Technical Details

### Performance

- Transcript fetching: 1-3 seconds
- Summary generation:
  - Local providers: 30-120 seconds
  - Cloud providers: 10-30 seconds

### Limitations

- Maximum video length: None (but longer videos = longer processing)
- Transcript availability depends on YouTube
- Rate limits apply for cloud providers
- Local processing speed depends on hardware

## 📚 Contributing

1. Fork the repository
2. Create your feature branch: `git checkout -b feature/amazing-feature`
3. Commit your changes: `git commit -am 'Add some feature'`
4. Push to the branch: `git push origin feature/amazing-feature`
5. Open a Pull Request

## ⚖️ Legal Notice

This tool accesses publicly available YouTube video transcripts. While I believe this falls under fair use:

- Review YouTube's Terms of Service
- Use responsibly and respect rate limits
- Consider YouTube's official API for commercial applications
- Don't use for mass data collection

## 📄 License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## 🙏 Acknowledgments

- Transcript handling inspired by [youtube-transcript-api](https://github.com/jdepoix/youtube-transcript-api)
- Built with [Cobra](https://github.com/spf13/cobra) CLI framework
- Configuration managed with [Viper](https://github.com/spf13/viper)
- Secure key storage by [go-keyring](https://github.com/zalando/go-keyring)
