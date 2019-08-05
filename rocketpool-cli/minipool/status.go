package minipool

import (
    "fmt"

    "gopkg.in/urfave/cli.v1"

    "github.com/rocket-pool/smartnode/shared/services"
    "github.com/rocket-pool/smartnode/shared/services/rocketpool/minipool"
    "github.com/rocket-pool/smartnode/shared/services/rocketpool/node"
    "github.com/rocket-pool/smartnode/shared/utils/eth"
)


// Get the node's minipool statuses
func getMinipoolStatus(c *cli.Context) error {

    // Initialise services
    p, err := services.NewProvider(c, services.ProviderOpts{
        AM: true,
        CM: true,
        LoadContracts: []string{"rocketPoolToken", "utilAddressSetStorage"},
        LoadAbis: []string{"rocketMinipool"},
        WaitClientSync: true,
        WaitRocketStorage: true,
    })
    if err != nil {
        return err
    }

    // Get minipool addresses
    nodeAccount, _ := p.AM.GetNodeAccount()
    minipoolAddresses, err := node.GetMinipoolAddresses(nodeAccount.Address, p.CM)
    if err != nil {
        return err
    }
    minipoolCount := len(minipoolAddresses)

    // Get minipool details
    detailsChannels := make([]chan *minipool.Details, minipoolCount)
    errorChannel := make(chan error)
    for mi := 0; mi < minipoolCount; mi++ {
        detailsChannels[mi] = make(chan *minipool.Details)
        go (func(mi int) {
            if details, err := minipool.GetDetails(p.CM, minipoolAddresses[mi]); err != nil {
                errorChannel <- err
            } else {
                detailsChannels[mi] <- details
            }
        })(mi)
    }

    // Receive minipool details
    minipoolDetails := make([]*minipool.Details, minipoolCount)
    for mi := 0; mi < minipoolCount; mi++ {
        select {
            case details := <-detailsChannels[mi]:
                minipoolDetails[mi] = details
            case err := <-errorChannel:
                return err
        }
    }

    // Log status & return
    fmt.Fprintln(p.Output, fmt.Sprintf("Node has %d minipools", minipoolCount))
    for _, details := range minipoolDetails {
        fmt.Fprintln(p.Output, "--------")
        fmt.Fprintln(p.Output, "Address:               ", details.Address.Hex())
        fmt.Fprintln(p.Output, "Status:                ", details.StatusType)
        fmt.Fprintln(p.Output, "Status Updated Time:   ", details.StatusTime.Format("2006-01-02, 15:04 -0700 MST"))
        fmt.Fprintln(p.Output, "Staking Duration:      ", details.StakingDurationId)
        fmt.Fprintln(p.Output, "Node ETH Deposited:    ", fmt.Sprintf("%.2f", eth.WeiToEth(details.NodeEtherBalanceWei)))
        fmt.Fprintln(p.Output, "Node RPL Deposited:    ", fmt.Sprintf("%.2f", eth.WeiToEth(details.NodeRplBalanceWei)))
        fmt.Fprintln(p.Output, "Deposit Count:         ", details.DepositCount.String())
        fmt.Fprintln(p.Output, "User Deposit Capacity: ", fmt.Sprintf("%.2f", eth.WeiToEth(details.UserDepositCapacityWei)))
        fmt.Fprintln(p.Output, "User Deposit Total:    ", fmt.Sprintf("%.2f", eth.WeiToEth(details.UserDepositTotalWei)))
    }
    return nil

}

