// Copyright (c) 2020, 2021, Oracle and/or its affiliates.
// Licensed under the Universal Permissive License v 1.0 as shown at https://oss.oracle.com/licenses/upl.

package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
)

func init() {
	clustersCmd.AddCommand(clusterListCmd)
}

var clusterListCmd = &cobra.Command{
	Use: "list",
	Short: "List the clusters",
	Long: "List the clusters",
	RunE: func(cmd *cobra.Command, args []string) error {
		if err := listClusters(args); err != nil {
			return err
		}
		return nil
	},
}

func listClusters(args []string) error {
	fmt.Println("list clusters...")
	return nil
}