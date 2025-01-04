package db

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"path"

	"mickamy.com/sampay/config"
)

func Create(ctx context.Context) error {
	db := config.Database()
	common := config.Common()

	// create user
	if common.Env == config.Development {
		cmd := exec.Command("psql", "-U", db.AdminUser, "-h", db.Host, "-d", "postgres", "-f", path.Join(common.PackageRoot, "db", "00_users.sql"))
		cmd.Env = append(os.Environ(), "PGPASSWORD="+db.AdminPassword)

		output, err := cmd.CombinedOutput()
		if err != nil {
			return fmt.Errorf("failed to create user: %w\nOutput: %s", err, string(output))
		}
	}

	// create database
	{
		cmd := exec.Command("psql", "-U", db.AdminUser, "-h", db.Host, "-d", "postgres", "-f", path.Join(common.PackageRoot, "db", "01_database.sql"))
		cmd.Env = append(os.Environ(), "PGPASSWORD="+db.AdminPassword)

		output, err := cmd.CombinedOutput()
		if err != nil {
			return fmt.Errorf("failed to create database: %w\nOutput: %s", err, string(output))
		}
	}

	return nil
}
