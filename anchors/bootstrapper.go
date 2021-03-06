package anchors

import (
	"errors"

	"github.com/centrifuge/go-centrifuge/bootstrap"
	"github.com/centrifuge/go-centrifuge/config"
	"github.com/centrifuge/go-centrifuge/ethereum"
	"github.com/centrifuge/go-centrifuge/queue"
)

type Bootstrapper struct {
}

// Bootstrap initializes the AnchorRepositoryContract as well as the AnchoringConfirmationTask that depends on it.
// the AnchoringConfirmationTask is added to be registered on the Queue at queue.Bootstrapper
func (*Bootstrapper) Bootstrap(context map[string]interface{}) error {
	if _, ok := context[bootstrap.BootstrappedConfig]; !ok {
		return errors.New("config hasn't been initialized")
	}
	if _, ok := context[bootstrap.BootstrappedEthereumClient]; !ok {
		return errors.New("ethereum client hasn't been initialized")
	}

	client := ethereum.GetConnection()
	repositoryContract, err := NewEthereumAnchorRepositoryContract(config.Config.GetContractAddress("anchorRepository"), client.GetClient())
	if err != nil {
		return err
	}

	anchorRepo := NewEthereumAnchorRepository(config.Config, repositoryContract)
	setAnchorRepository(anchorRepo)
	if err != nil {
		return err
	}
	return queue.InstallQueuedTask(context, NewAnchoringConfirmationTask(&repositoryContract.EthereumAnchorRepositoryContractFilterer, ethereum.DefaultWaitForTransactionMiningContext))
}
