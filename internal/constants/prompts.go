package constants

const (
	ShortSummaryPrompt = `Create a concise summary of the following transcript. Focus on:
- Core message in 1-2 sentences
- 3-5 key points that support or develop the core message
- If applicable, note any specific calls to action or main conclusions
Keep total length under 150 words.`

	LongSummaryPrompt = `Create a detailed analysis of the following transcript that preserves the original context and depth while making it accessible. Structure as follows:

1. Executive Summary (3-4 sentences)
2. Context and Background
   - Identify the apparent purpose/context
   - Note any assumed knowledge or prerequisites
3. Main Content Analysis
   - Break down major themes and arguments
   - Highlight key terminology and concepts
   - Connect related ideas and show progression
4. Evidence and Support
   - Note specific examples, data, or case studies
   - Identify methodologies or frameworks used
5. Implications and Conclusions
   - Summarize main takeaways
   - Note potential applications or next steps

Preserve technical accuracy while ensuring readability. Include relevant quotes when they significantly support key points.`

	TranscriptPrompt = `Format the following raw Youtube transcript text.
- Add appropriate capitalization and punctuation
- Keep all original words exactly as they appear
- Never add any additional commentary
- Do not correct spelling or grammar
- Add paragraph breaks where appropriate
- Do not otherwise modify the content in any way`

	QueryPrompt = `You are analyzing a YouTube video transcript.
Video title: "{{title}}"

Question: {{query}}

Provide a concise, accurate answer based ONLY on information contained in the transcript.
If the transcript doesn't contain information to answer the question, clearly state this.
Do not speculate beyond what's explicitly mentioned in the transcript.
Reference specific details from the transcript to support your answer.`
)
