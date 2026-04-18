package main

import (
	"encoding/json"
	"fmt"

	"github.com/spf13/cobra"
)

func writeJSON(cmd *cobra.Command, v any) error {
	data, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		return fmt.Errorf("marshal JSON: %w", err)
	}
	if _, werr := fmt.Fprintln(cmd.OutOrStdout(), string(data)); werr != nil {
		return fmt.Errorf("write JSON: %w", werr)
	}
	return nil
}
