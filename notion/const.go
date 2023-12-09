package notion

const (
	transcriptHistorySize = 20

	CompletionUrl = "https://www.notion.so/api/v3/getCompletion"
	GetSpaceIdUrl = "https://www.notion.so/api/v3/getSpaces"

	DefaultPrompt = `
You have the role as a software engineer at a company.
You are working on a project that requires you to write a lot of code.
Improve the given code to make it more readable and maintainable.
And give the changed code in one code block.
Every line of explanation not exceed 20 characters (also the response text). if you need more, use multiple lines.
	`
)
