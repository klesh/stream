package trellogithub

import (
	"github.com/merico-dev/stream/pkg/util/log"
)

// Update remove and set up trello-github-integ workflows.
func Update(options map[string]interface{}) (map[string]interface{}, error) {
	tg, err := NewTrelloGithub(options)
	if err != nil {
		return nil, err
	}

	api := tg.GetApi()
	log.Infof("API is %s.", api.Name)
	ws := defaultWorkflows.GetWorkflowByNameVersionTypeString(api.Name)

	for _, w := range ws {
		err := tg.client.DeleteWorkflow(w, tg.options.Branch)
		if err != nil {
			return nil, err
		}

		if err := tg.renderTemplate(w); err != nil {
			return nil, err
		}
		err = tg.client.AddWorkflow(w, tg.options.Branch)
		if err != nil {
			return nil, err
		}
	}
	log.Success("Adding workflow file succeeded.")
	trelloIds, err := tg.CreateTrelloItems()
	if err != nil {
		return nil, err
	}
	log.Success("Creating trello board succeeded.")
	if err := tg.AddTrelloIdSecret(trelloIds); err != nil {
		return nil, err
	}

	log.Success("Adding secret keys for trello succeeded.")

	return buildState(tg, trelloIds), nil
}
