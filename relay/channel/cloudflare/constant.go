// Copyright (c) 2025 Tethys Plex
//
// This file is part of Veloera.
//
// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU General Public License for more details.
//
// You should have received a copy of the GNU General Public License
// along with this program. If not, see <https://www.gnu.org/licenses/>.
package cloudflare

var ModelList = []string{
	// Cloudflare Workers AI models (legacy format)
	"@cf/meta/llama-3.1-8b-instruct",
	"@cf/meta/llama-2-7b-chat-fp16",
	"@cf/meta/llama-2-7b-chat-int8",
	"@cf/mistral/mistral-7b-instruct-v0.1",
	"@hf/thebloke/deepseek-coder-6.7b-base-awq",
	"@hf/thebloke/deepseek-coder-6.7b-instruct-awq",
	"@cf/deepseek-ai/deepseek-math-7b-base",
	"@cf/deepseek-ai/deepseek-math-7b-instruct",
	"@cf/thebloke/discolm-german-7b-v1-awq",
	"@cf/tiiuae/falcon-7b-instruct",
	"@cf/google/gemma-2b-it-lora",
	"@hf/google/gemma-7b-it",
	"@cf/google/gemma-7b-it-lora",
	"@hf/nousresearch/hermes-2-pro-mistral-7b",
	"@hf/thebloke/llama-2-13b-chat-awq",
	"@cf/meta-llama/llama-2-7b-chat-hf-lora",
	"@cf/meta/llama-3-8b-instruct",
	"@hf/thebloke/llamaguard-7b-awq",
	"@hf/thebloke/mistral-7b-instruct-v0.1-awq",
	"@hf/mistralai/mistral-7b-instruct-v0.2",
	"@cf/mistral/mistral-7b-instruct-v0.2-lora",
	"@hf/thebloke/neural-chat-7b-v3-1-awq",
	"@cf/openchat/openchat-3.5-0106",
	"@hf/thebloke/openhermes-2.5-mistral-7b-awq",
	"@cf/microsoft/phi-2",
	"@cf/qwen/qwen1.5-0.5b-chat",
	"@cf/qwen/qwen1.5-1.8b-chat",
	"@cf/qwen/qwen1.5-14b-chat-awq",
	"@cf/qwen/qwen1.5-7b-chat-awq",
	"@cf/defog/sqlcoder-7b-2",
	"@hf/nexusflow/starling-lm-7b-beta",
	"@cf/tinyllama/tinyllama-1.1b-chat-v1.0",
	"@hf/thebloke/zephyr-7b-beta-awq",
	
	// OpenAI models through AI Gateway
	"openai/gpt-4o",
	"openai/gpt-4o-mini",
	"openai/gpt-4",
	"openai/gpt-4-turbo",
	"openai/gpt-3.5-turbo",
	"openai/text-embedding-3-small",
	"openai/text-embedding-3-large",
	"openai/text-embedding-ada-002",
	"openai/whisper-1",
	"openai/dall-e-3",
	"openai/dall-e-2",
	
	// Anthropic models through AI Gateway  
	"anthropic/claude-3-5-sonnet-20241022",
	"anthropic/claude-3-5-haiku-20241022",
	"anthropic/claude-3-opus-20240229",
	"anthropic/claude-3-sonnet-20240229",
	"anthropic/claude-3-haiku-20240307",
	
	// Workers AI models with new prefix for gateway routing
	"workers-ai/llama-3.1-8b-instruct",
	"workers-ai/llama-2-7b-chat", 
	"workers-ai/mistral-7b-instruct",
}

var ChannelName = "cloudflare"
