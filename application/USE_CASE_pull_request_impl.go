package application

import (
	"encoding/json"
	domain "github_wb/domain/value_objects"
	"log"
)

func HandlePullRequestEvent(payload []byte) int {
	var event domain.PullRequestEventPayload

	if err := json.Unmarshal(payload, &event); err != nil {
		return 500
	}

	if event.Action == "closed" {
		baseBranch := event.PullRequest.Base.Ref
		featureBranch := event.PullRequest.Head.Ref
		user := event.PullRequest.User.Login
		prID := event.PullRequest.ID

		log.Printf("Pull Request cerrado:\nID: %d\nBase: %s\nHead: %s\nUsuario: %s", prID, baseBranch, featureBranch, user)
	} else {
		log.Printf("Acción no soportada para Pull Request: %s", event.Action)
	}

	return 200
}
// pureba pullrequ