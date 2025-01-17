package main

import (
	"context"
	"cryoutilities/internal"
	"errors"
	"github.com/cristalhq/acmd"
	"log"
	"os"
	"strconv"
	"strings"
)

func main() {
	// Delete old log file
	_ = os.Remove(internal.LogFilePath)
	// Create a log file
	logFile, err := os.OpenFile(internal.LogFilePath, os.O_RDWR|os.O_CREATE, 0666)
	if err != nil {
		log.Panic(err)
	}
	defer func(logFile *os.File) {
		err := logFile.Close()
		if err != nil {

		}
	}(logFile)
	log.SetOutput(logFile)

	// Create loggers
	internal.CryoUtils.InfoLog = log.New(logFile, "INFO\t", log.Ldate|log.Ltime)
	internal.CryoUtils.ErrorLog = log.New(logFile, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)

	// Print the current version as a test
	internal.CryoUtils.InfoLog.Println("Current Version:", internal.CurrentVersionNumber)

	// Provide a command structure for parsing
	cmds := []acmd.Command{
		{
			Name:        "gui",
			Description: "Run in GUI mode",
			ExecFunc: func(ctx context.Context, args []string) error {
				internal.InitUI()
				return nil
			},
		},
		{
			Name:        "swap",
			Description: "Change swap file size in increments of 1GB.",
			ExecFunc: func(ctx context.Context, args []string) error {
				internal.CryoUtils.InfoLog.Println("Starting swap file resize...")
				size, err := strconv.Atoi(args[0])
				if err != nil {
					return err
				}
				err = internal.ChangeSwapSizeCLI(size, false)
				if err != nil {
					return err
				}
				internal.CryoUtils.InfoLog.Println("Success!")
				return nil
			},
		},
		{
			Name:        "swappiness",
			Description: "Change swappiness to the specified value 0-200.",
			ExecFunc: func(ctx context.Context, args []string) error {
				internal.CryoUtils.InfoLog.Println("Starting swappiness change...")
				swappiness := args[0]
				swappinessInt, err := strconv.Atoi(swappiness)
				if err != nil || swappinessInt < 0 || swappinessInt > 200 {
					return errors.New("invalid swappiness value")
				}
				err = internal.ChangeSwappiness(swappiness)
				if err != nil {
					return err
				}
				internal.CryoUtils.InfoLog.Println("Success!")
				return nil
			},
		},
		{
			Name:        "hugepages",
			Description: "Enable or disable hugepages. Accepts 'true', 'false', 'enable' or 'disable'.\n\tRecommended: Enabled",
			ExecFunc: func(ctx context.Context, args []string) error {
				arg := strings.ToLower(args[0])
				if arg == "true" || arg == "enable" {
					internal.CryoUtils.InfoLog.Println("Enabling HugePages...")
					err := internal.SetHugePages()
					if err != nil {
						return err
					}
				} else if arg == "false" || arg == "disable" {
					internal.CryoUtils.InfoLog.Println("Disabling HugePages...")
					err := internal.RevertHugePages()
					if err != nil {
						return err
					}
				} else {
					return errors.New("invalid argument provided")
				}
				return nil
			},
		},
		{
			Name:        "compaction_proactiveness",
			Description: "Set or revert compaction proactiveness. Accepts 'recommended' or 'stock'.",
			ExecFunc: func(ctx context.Context, args []string) error {
				arg := strings.ToLower(args[0])
				if arg == "recommended" {
					internal.CryoUtils.InfoLog.Println("Setting Compaction Proactiveness...")
					err := internal.SetCompactionProactiveness()
					if err != nil {
						return err
					}
				} else if arg == "stock" {
					internal.CryoUtils.InfoLog.Println("Reverting Compaction Proactiveness...")
					err := internal.RevertCompactionProactiveness()
					if err != nil {
						return err
					}
				} else {
					return errors.New("invalid argument provided")
				}
				return nil
			},
		},
		{
			Name:        "defrag",
			Description: "Enable or disable hugepage defrag. Accepts 'true', 'false', 'enable' or 'disable'.\n\tRecommended: Disabled",
			ExecFunc: func(ctx context.Context, args []string) error {
				arg := strings.ToLower(args[0])
				if arg == "true" || arg == "enable" {
					internal.CryoUtils.InfoLog.Println("Enabling HugePAge Defrag...")
					err := internal.RevertDefrag()
					if err != nil {
						return err
					}
				} else if arg == "false" || arg == "disable" {
					internal.CryoUtils.InfoLog.Println("Revert Compaction Proactiveness...")
					err := internal.SetDefrag()
					if err != nil {
						return err
					}
				} else {
					return errors.New("invalid argument provided")
				}
				return nil
			},
		},
		{
			Name:        "page_lock_unfairness",
			Description: "Set or revert page lock unfairness. Accepts 'recommended' or 'stock'.",
			ExecFunc: func(ctx context.Context, args []string) error {
				arg := strings.ToLower(args[0])
				if arg == "recommended" {
					internal.CryoUtils.InfoLog.Println("Setting Page Lock Unfairness...")
					err := internal.SetPageLockUnfairness()
					if err != nil {
						return err
					}
				} else if arg == "stock" {
					internal.CryoUtils.InfoLog.Println("Reverting Page Lock Unfairness...")
					err := internal.RevertPageLockUnfairness()
					if err != nil {
						return err
					}
				} else {
					return errors.New("invalid argument provided")
				}
				return nil
			},
		},
		{
			Name:        "shmem",
			Description: "Enable or disable shared memory. Accepts 'true', 'false', 'enable' or 'disable'.\n\tRecommended: Enabled",
			ExecFunc: func(ctx context.Context, args []string) error {
				arg := strings.ToLower(args[0])
				if arg == "true" || arg == "enable" {
					internal.CryoUtils.InfoLog.Println("Setting Shared Memory...")
					err := internal.SetShMem()
					if err != nil {
						return err
					}
				} else if arg == "false" || arg == "disable" {
					internal.CryoUtils.InfoLog.Println("Reverting Shared Memory...")
					err := internal.RevertShMem()
					if err != nil {
						return err
					}
				} else {
					return errors.New("invalid argument provided")
				}
				return nil
			},
		},
		{
			Name:        "recommended",
			Description: "Set all values to Cryo's recommendations.",
			ExecFunc: func(ctx context.Context, args []string) error {
				err := internal.UseRecommendedSettings()
				if err != nil {
					return err
				}
				return nil
			},
		},
		{
			Name:        "stock",
			Description: "Set all values to Valve defaults.",
			ExecFunc: func(ctx context.Context, args []string) error {
				err := internal.UseStockSettings()
				if err != nil {
					return err
				}
				return nil
			},
		},
	}

	// If no args are passed, assume "gui"
	if len(os.Args) <= 1 {
		os.Args = []string{"", "gui"}
	}

	// Basic program metadata
	r := acmd.RunnerOf(cmds, acmd.Config{
		AppName:         "cryoutilities",
		AppDescription:  "CryoByte33's Steam Deck utility script.",
		PostDescription: "NOTE: You NEED to run this with sudo if not using GUI mode.",
		Version:         internal.CurrentVersionNumber,
	})

	// Run the command parser
	if err := r.Run(); err != nil {
		internal.CryoUtils.ErrorLog.Println(err)
		os.Exit(1)
	}
}
