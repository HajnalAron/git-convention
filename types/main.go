package types

/*
{
   "type": "revert",
   "description": "Reverts a previous commit",
}
*/

type Branch struct {
	Type string `json:"type"`
	Desc string `json:"description"`
}

func (b Branch) Title() string       { return b.Type }
func (b Branch) Description() string { return b.Desc }
func (b Branch) FilterValue() string { return b.Type }

/*
{
"type": "fix",
"description": "A bug fix",
"emoji": "🐛"
}
*/

type Commit struct {
	Type  string `json:"type"`
	Desc  string `json:"description"`
	Emoji string `json:"emoji"`
}

func (c Commit) Title() string       { return c.Type }
func (c Commit) Description() string { return c.Desc }
func (c Commit) FilterValue() string { return c.Type }
