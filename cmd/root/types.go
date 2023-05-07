package root

import "github.com/bilalcaliskan/s3-manager/internal/prompt"

var (
	selectRunner    prompt.SelectRunner = prompt.GetSelectRunner("Select operation", []string{"search", "clean"})
	accessKeyRunner prompt.PromptRunner = prompt.GetPromptRunner("Provide AWS Access Key", nil)
	secretKeyRunner prompt.PromptRunner = prompt.GetPromptRunner("Provide AWS Secret Key", nil)
	regionRunner    prompt.PromptRunner = prompt.GetPromptRunner("Provide AWS Region", nil)
	bucketRunner    prompt.PromptRunner = prompt.GetPromptRunner("Provide AWS Bucket Name", nil)
)
