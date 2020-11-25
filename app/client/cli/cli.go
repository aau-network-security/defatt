package cli

import (
	"context"
	"errors"
	"fmt"
	"time"

	pb "github.com/aau-network-security/defat/app/daemon/proto"
	"github.com/spf13/cobra"
)

// will be splitted into different sub files
var (
	UnableToListScenarios = errors.New("Failed to list scenarios")
)

func (c *Client) ListScenarios() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "scenarios",
		Short:   "List available scenarios",
		Example: "defat scenarios list ",
		Run: func(cmd *cobra.Command, args []string) {
			ctx, cancel := context.WithTimeout(context.Background(), time.Second)
			defer cancel()
			r, err := c.rpcClient.ListScenarios(ctx, &pb.EmptyRequest{})
			if err != nil {
				PrintError(err)
				return
			}
			f := formatter{
				header: []string{"SCENARIO ID", "DIFFICULTY", "DURATION", "NUMBER OF NETWORKS", "STORY"},
				fields: []string{"Id", "Difficulty", "Duration", "Networks", "Story"},
			}

			var elements []formatElement
			for _, e := range r.Scenarios {
				elements = append(elements, e)
			}
			table, err := f.AsTable(elements)
			if err != nil {
				PrintError(UnableToListScenarios)
				return
			}
			fmt.Printf(table)
		},
	}
	return cmd
}

func (c *Client) ListChallengesOnScenario() *cobra.Command {
	var scenarioID uint32
	cmd := &cobra.Command{
		Use:     "challenges",
		Short:   "List challenges in given scenario",
		Example: "defat scenario  1 ",
		Run: func(cmd *cobra.Command, args []string) {
			ctx, cancel := context.WithTimeout(context.Background(), time.Second)
			defer cancel()
			resp, err := c.rpcClient.ListScenChals(ctx, &pb.ListScenarioChallengesReq{ScenarioId: scenarioID})
			if err != nil {
				PrintError(err)
				return
			}
			var header []string
			//var fields []string
			var elements []formatElement
			for _, n := range resp.Chals {
				header = append(header, n.Vlan)
				elements = append(elements, n.Challenges)
				//fields = append(fields, strings.Join(n.Challenges, "|"))
			}
			f := formatter{
				header: header,
				fields: header,
			}

			table, err := f.AsTable(elements)
			if err != nil {
				PrintError(UnableToListScenarios)
				return
			}
			fmt.Printf(table)
		},
	}
	return cmd
}
