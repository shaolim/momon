package receipt

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github/shaolim/momon/internal/receipt/model"
	"os"
	"path/filepath"
	"strings"

	"github.com/openai/openai-go/v3"
)

type Receipt struct {
	client *openai.Client
	prompt string
}

func New(client *openai.Client) *Receipt {
	prompt := `You are a receipt information extraction assistant. Your task is to analyze the uploaded image and extract structured receipt data.

VALIDATION RULES:
1. Verify the image is a valid receipt (must contain: merchant name, date, items with prices, and total)
2. If the image is NOT a receipt (e.g., random photo, document, etc.), return an error response
3. If the receipt is too blurry or text is unreadable, return an error response

OUTPUT FORMAT - VALID RECEIPT:
Return ONLY valid JSON (no comments, no additional text):
{
    "shop": "Name of the merchant or store",
    "transactionDate": "YYYY-MM-DD HH:MM format (use 24-hour time)",
    "items": [
        {
            "name": "Item name or description",
            "quantity": 1,
            "price": 1000,
            "tax": 0,
            "totalPrice": 1000
        }
    ],
    "tax": 0,
    "total": 1000,
    "isValid": true
}

FIELD DESCRIPTIONS:
- shop: Merchant/store name as shown on receipt
- transactionDate: Date and time in YYYY-MM-DD HH:MM format (if time not visible, use 00:00)
- items: Array of all purchased items
  - name: Product name/description
  - quantity: Number of units purchased (default: 1 if not specified)
  - price: Unit price per item (not total)
  - tax: Tax amount for this specific item (0 if not itemized)
  - totalPrice: Calculated as (quantity × price) + tax
- tax: Total tax amount for entire receipt (sum all item taxes, or use receipt total tax)
- total: Final total amount paid (must match receipt total)
- isValid: Must be true for valid receipts

OUTPUT FORMAT - INVALID RECEIPT:
Return ONLY valid JSON:
{
    "isValid": false,
    "message": "Descriptive error message explaining why the receipt is invalid"
}

CALCULATION REQUIREMENTS:
- Verify that sum of all item totalPrices matches the receipt subtotal
- Verify that subtotal + tax = total on receipt
- If calculations don't match receipt within 1% tolerance, still extract visible data but note any discrepancy
- All monetary values should be in the smallest currency unit (e.g., cents, not dollars)

CRITICAL REQUIREMENTS:
- Return ONLY the JSON object, no additional text or explanation
- Ensure all JSON is properly formatted and valid
- Use null for missing optional fields, not empty strings
- Always return either the valid receipt format OR the error format, never both`

	return &Receipt{
		client: client,
		prompt: prompt,
	}
}

func getMimeType(path string) string {
	ext := strings.ToLower(filepath.Ext(path))
	switch ext {
	case ".jpg", ".jpeg":
		return "image/jpeg"
	case ".png":
		return "image/png"
	case ".gif":
		return "image/gif"
	case ".webp":
		return "image/webp"
	default:
		return "image/jpeg" // default fallback
	}
}

// cleanJSONResponse removes markdown code blocks and other formatting from the API response
func cleanJSONResponse(content string) string {
	// Trim whitespace
	content = strings.TrimSpace(content)

	// Remove markdown code blocks (```json ... ``` or ``` ... ```)
	if strings.HasPrefix(content, "```") {
		// Find the first newline after ```
		lines := strings.Split(content, "\n")
		if len(lines) > 0 {
			// Remove first line (```json or ```)
			lines = lines[1:]
		}
		// Remove last line if it's ```
		if len(lines) > 0 && strings.TrimSpace(lines[len(lines)-1]) == "```" {
			lines = lines[:len(lines)-1]
		}
		content = strings.Join(lines, "\n")
	}

	// Trim again after removing code blocks
	content = strings.TrimSpace(content)

	return content
}

func (r *Receipt) ReadReceipt(ctx context.Context, path string) (*model.Receipt, error) {
	imgBytes, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read file: %w", err)
	}
	mimeType := getMimeType(path)
	dataURL := fmt.Sprintf("data:%s;base64,%s", mimeType, base64.StdEncoding.EncodeToString(imgBytes))

	messages := []openai.ChatCompletionMessageParamUnion{
		{
			OfUser: &openai.ChatCompletionUserMessageParam{
				Content: openai.ChatCompletionUserMessageParamContentUnion{
					OfArrayOfContentParts: []openai.ChatCompletionContentPartUnionParam{
						{
							OfText: &openai.ChatCompletionContentPartTextParam{
								Text: r.prompt,
							},
						},
						{
							OfImageURL: &openai.ChatCompletionContentPartImageParam{
								ImageURL: openai.ChatCompletionContentPartImageImageURLParam{
									URL: dataURL,
								},
							},
						},
					},
				},
			},
		},
	}

	req := openai.ChatCompletionNewParams{
		Messages: messages,
		Model:    openai.ChatModelGPT4o,
	}

	resp, err := r.client.Chat.Completions.New(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("failed to get response from OpenAI: %w", err)
	}
	if len(resp.Choices) == 0 {
		return nil, fmt.Errorf("no choices returned from OpenAI")
	}
	content := resp.Choices[0].Message.Content

	// Clean the response to remove markdown code blocks
	cleanedContent := cleanJSONResponse(content)

	var result model.Receipt
	err = json.Unmarshal([]byte(cleanedContent), &result)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal receipt: %w, content: %s", err, content)
	}

	return &result, nil
}
