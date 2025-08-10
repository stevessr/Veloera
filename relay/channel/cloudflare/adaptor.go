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

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"
	"veloera/dto"
	"veloera/relay/channel"
	relaycommon "veloera/relay/common"
	"veloera/relay/constant"

	"github.com/gin-gonic/gin"
)

type Adaptor struct {
}

func (a *Adaptor) ConvertClaudeRequest(*gin.Context, *relaycommon.RelayInfo, *dto.ClaudeRequest) (any, error) {
	//TODO implement me
	panic("implement me")
	return nil, nil
}

func (a *Adaptor) Init(info *relaycommon.RelayInfo) {
}

func (a *Adaptor) GetRequestURL(info *relaycommon.RelayInfo) (string, error) {
	// Check if this is an AI Gateway URL vs Workers AI
	if strings.HasPrefix(info.BaseUrl, "https://gateway.ai.cloudflare.com") {
		// AI Gateway format: https://gateway.ai.cloudflare.com/v1/{account_id}/{gateway_id}
		// The gateway_id should be stored in the Other field
		var gatewayId string
		if info.Other != nil {
			if gid, exists := info.Other["gateway_id"]; exists {
				if gatewayIdStr, ok := gid.(string); ok {
					gatewayId = gatewayIdStr
				}
			}
		}
		if gatewayId == "" {
			return "", errors.New("gateway_id is required for Cloudflare AI Gateway")
		}
		
		// Determine the provider endpoint based on the model name
		_, endpoint := a.getProviderEndpoint(info.UpstreamModelName, info.RelayMode)
		return fmt.Sprintf("%s/v1/%s/%s%s", info.BaseUrl, info.ApiVersion, gatewayId, endpoint), nil
	}
	
	// Workers AI format (legacy)
	switch info.RelayMode {
	case constant.RelayModeChatCompletions:
		return fmt.Sprintf("%s/client/v4/accounts/%s/ai/v1/chat/completions", info.BaseUrl, info.ApiVersion), nil
	case constant.RelayModeEmbeddings:
		return fmt.Sprintf("%s/client/v4/accounts/%s/ai/v1/embeddings", info.BaseUrl, info.ApiVersion), nil
	default:
		return fmt.Sprintf("%s/client/v4/accounts/%s/ai/run/%s", info.BaseUrl, info.ApiVersion, info.UpstreamModelName), nil
	}
}

func (a *Adaptor) SetupRequestHeader(c *gin.Context, req *http.Header, info *relaycommon.RelayInfo) error {
	channel.SetupApiRequestHeader(info, c, req)
	
	// For AI Gateway, we use the API key as the CF-Authorization header
	// For Workers AI direct, we use Bearer token
	if strings.HasPrefix(info.BaseUrl, "https://gateway.ai.cloudflare.com") {
		req.Set("CF-Authorization", fmt.Sprintf("Bearer %s", info.ApiKey))
	} else {
		req.Set("Authorization", fmt.Sprintf("Bearer %s", info.ApiKey))
	}
	
	return nil
}

func (a *Adaptor) ConvertOpenAIRequest(c *gin.Context, info *relaycommon.RelayInfo, request *dto.GeneralOpenAIRequest) (any, error) {
	if request == nil {
		return nil, errors.New("request is nil")
	}
	
	// Check if this is AI Gateway vs Workers AI
	if strings.HasPrefix(info.BaseUrl, "https://gateway.ai.cloudflare.com") {
		// For AI Gateway, convert the model name to remove provider prefix
		provider, _ := a.getProviderEndpoint(info.UpstreamModelName, info.RelayMode)
		switch provider {
		case "openai":
			// Remove openai/ prefix for the actual request
			if strings.HasPrefix(info.UpstreamModelName, "openai/") {
				request.Model = strings.TrimPrefix(info.UpstreamModelName, "openai/")
			}
			return request, nil
		case "anthropic":
			// For Anthropic through gateway, we need to convert OpenAI format to Claude format
			return a.convertToClaudeRequest(*request), nil
		case "workers-ai":
			// For Workers AI through gateway, use the original model name format
			if strings.HasPrefix(info.UpstreamModelName, "workers-ai/") {
				modelName := strings.TrimPrefix(info.UpstreamModelName, "workers-ai/")
				// Map friendly names back to full Workers AI model names
				request.Model = a.mapToWorkersAIModel(modelName)
			}
			return request, nil
		default:
			return request, nil
		}
	}
	
	// Workers AI direct (legacy)
	switch info.RelayMode {
	case constant.RelayModeCompletions:
		return convertCf2CompletionsRequest(*request), nil
	default:
		return request, nil
	}
}

func (a *Adaptor) ConvertOpenAIResponsesRequest(c *gin.Context, info *relaycommon.RelayInfo, request dto.OpenAIResponsesRequest) (any, error) {
	// TODO implement me
	return nil, errors.New("not implemented")
}

func (a *Adaptor) DoRequest(c *gin.Context, info *relaycommon.RelayInfo, requestBody io.Reader) (any, error) {
	return channel.DoApiRequest(a, c, info, requestBody)
}

func (a *Adaptor) ConvertRerankRequest(c *gin.Context, relayMode int, request dto.RerankRequest) (any, error) {
	return request, nil
}

func (a *Adaptor) ConvertEmbeddingRequest(c *gin.Context, info *relaycommon.RelayInfo, request dto.EmbeddingRequest) (any, error) {
	return request, nil
}

func (a *Adaptor) ConvertAudioRequest(c *gin.Context, info *relaycommon.RelayInfo, request dto.AudioRequest) (io.Reader, error) {
	// 添加文件字段
	file, _, err := c.Request.FormFile("file")
	if err != nil {
		return nil, errors.New("file is required")
	}
	defer file.Close()
	// 打开临时文件用于保存上传的文件内容
	requestBody := &bytes.Buffer{}

	// 将上传的文件内容复制到临时文件
	if _, err := io.Copy(requestBody, file); err != nil {
		return nil, err
	}
	return requestBody, nil
}

func (a *Adaptor) ConvertImageRequest(c *gin.Context, info *relaycommon.RelayInfo, request dto.ImageRequest) (any, error) {
	//TODO implement me
	return nil, errors.New("not implemented")
}

func (a *Adaptor) DoResponse(c *gin.Context, resp *http.Response, info *relaycommon.RelayInfo) (usage any, err *dto.OpenAIErrorWithStatusCode) {
	switch info.RelayMode {
	case constant.RelayModeEmbeddings:
		fallthrough
	case constant.RelayModeChatCompletions:
		if info.IsStream {
			err, usage = cfStreamHandler(c, resp, info)
		} else {
			err, usage = cfHandler(c, resp, info)
		}
	case constant.RelayModeAudioTranslation:
		fallthrough
	case constant.RelayModeAudioTranscription:
		err, usage = cfSTTHandler(c, resp, info)
	}
	return
}

func (a *Adaptor) GetModelList() []string {
	return ModelList
}

func (a *Adaptor) GetChannelName() string {
	return ChannelName
}

// getProviderEndpoint determines the provider and endpoint based on the model name
func (a *Adaptor) getProviderEndpoint(modelName string, relayMode int) (string, string) {
	// Extract provider prefix from model name
	parts := strings.Split(modelName, "/")
	if len(parts) < 2 {
		// Default to OpenAI if no provider specified
		return "openai", a.getOpenAIEndpoint(relayMode)
	}
	
	provider := parts[0]
	switch provider {
	case "openai":
		return "openai", a.getOpenAIEndpoint(relayMode)
	case "anthropic":
		return "anthropic", a.getAnthropicEndpoint(relayMode)
	case "workers-ai", "cf":
		return "workers-ai", a.getWorkersAIEndpoint(relayMode)
	default:
		// Default to OpenAI format
		return "openai", a.getOpenAIEndpoint(relayMode)
	}
}

func (a *Adaptor) getOpenAIEndpoint(relayMode int) string {
	switch relayMode {
	case constant.RelayModeChatCompletions:
		return "/v1/chat/completions"
	case constant.RelayModeEmbeddings:
		return "/v1/embeddings" 
	case constant.RelayModeImagesGenerations:
		return "/v1/images/generations"
	case constant.RelayModeAudioTranscription:
		return "/v1/audio/transcriptions"
	case constant.RelayModeAudioTranslation:
		return "/v1/audio/translations"
	default:
		return "/v1/chat/completions"
	}
}

func (a *Adaptor) getAnthropicEndpoint(relayMode int) string {
	switch relayMode {
	case constant.RelayModeChatCompletions:
		return "/v1/messages"
	default:
		return "/v1/messages"
	}
}

func (a *Adaptor) getWorkersAIEndpoint(relayMode int) string {
	switch relayMode {
	case constant.RelayModeChatCompletions:
		return "/v1/chat/completions"
	case constant.RelayModeEmbeddings:
		return "/v1/embeddings"
	default:
		return "/v1/chat/completions"
	}
}

// mapToWorkersAIModel maps friendly model names to full Workers AI model identifiers
func (a *Adaptor) mapToWorkersAIModel(modelName string) string {
	modelMap := map[string]string{
		"llama-3.1-8b-instruct": "@cf/meta/llama-3.1-8b-instruct",
		"llama-2-7b-chat":       "@cf/meta/llama-2-7b-chat-fp16", 
		"mistral-7b-instruct":   "@cf/mistral/mistral-7b-instruct-v0.1",
	}
	
	if fullName, exists := modelMap[modelName]; exists {
		return fullName
	}
	return modelName // Return as-is if no mapping found
}

// convertToClaudeRequest converts OpenAI format to Claude Messages format
func (a *Adaptor) convertToClaudeRequest(request dto.GeneralOpenAIRequest) map[string]interface{} {
	claudeRequest := map[string]interface{}{
		"model":      request.Model,
		"max_tokens": request.GetMaxTokens(),
	}
	
	if request.Temperature != nil {
		claudeRequest["temperature"] = *request.Temperature
	}
	if request.TopP != 0 {
		claudeRequest["top_p"] = request.TopP
	}
	if request.Stream {
		claudeRequest["stream"] = true
	}
	
	// Convert messages
	var messages []map[string]interface{}
	var systemMessage string
	
	for _, msg := range request.Messages {
		if msg.Role == "system" {
			systemMessage = msg.StringContent()
		} else {
			claudeMsg := map[string]interface{}{
				"role":    msg.Role,
				"content": msg.StringContent(),
			}
			messages = append(messages, claudeMsg)
		}
	}
	
	if systemMessage != "" {
		claudeRequest["system"] = systemMessage
	}
	claudeRequest["messages"] = messages
	
	return claudeRequest
}
