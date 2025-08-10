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
package cloudflare_gateway

var ModelList = []string{
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
	
	// Workers AI models through gateway
	"workers-ai/llama-3.1-8b-instruct",
	"workers-ai/llama-2-7b-chat", 
	"workers-ai/mistral-7b-instruct",
}

var ChannelName = "cloudflare_gateway"