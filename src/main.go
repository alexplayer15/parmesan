package main

import (
    "fmt"
    "os"

    "github.com/spf13/cobra"
    "github.com/alexp/parmesan/src/commands"
)

var rootCmd = &cobra.Command{
    Use:   "parmesan",
    Short: "CLI tool to generate requests based off your OAS",
}

func main() {
    if err := rootCmd.Execute(); err != nil {
        fmt.Println(err)
        os.Exit(1)
    }

    rootCmd.AddCommand(commands.GenerateRequestCmd)  // Use GenerateRequestCmd from the commands package
    if err := rootCmd.Execute(); err != nil {
        fmt.Println(err)
        os.Exit(1)
    }
}
