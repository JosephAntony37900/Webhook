package value_objects

type PullRequestEvent struct {
	Action      string      `json:"action"`
	PullRequest PullRequest `json:"pull_request"`
	Repository  Repository  `json:"repository"`
}

type Repository struct {
	ID       int    `json:"id"`
	FullName string `json:"full_name"`
}
type PullRequest struct {
	ID      int    `json:"id"`
	Title   string `json:"title"`
	User    User   `json:"user"`
	Head    Branch `json:"head"`
	Base    Branch `json:"base"`
	URL     string `json:"url"`
	HTMLURL string `json:"html_url"`
}

type User struct {
	Login string `json:"login"`
}

type Branch struct {
	Ref string `json:"ref"`
	SHA string `json:"sha"`
}
