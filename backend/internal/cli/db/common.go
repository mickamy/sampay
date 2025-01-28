package db

import (
	"fmt"
	"os"
	"os/exec"
	"path"
	"strings"

	"mickamy.com/sampay/config"
)

func runPSQL(fileName string, variables map[string]string) error {
	db := config.Database()
	filePath := path.Join(config.Common().PackageRoot, "db", fileName)
	opts := []string{
		"-U", db.AdminUser.Escape(),
		"-h", db.Host,
		"-d", "postgres",
		"-f", filePath,
	}

	cmd := exec.Command("psql", opts...)
	cmd.Env = append(os.Environ(), "PGPASSWORD="+db.AdminPassword.Escape())

	var pgOpts []string
	for k, v := range variables {
		pgOpts = append(pgOpts, fmt.Sprintf("-c %s=%s", k, v))
	}

	cmd.Env = append(cmd.Env, "PGOPTIONS="+strings.Join(pgOpts, " "))

	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("failed to execute %s: %w\nOutput: %s", fileName, err, string(output))
	}

	return nil
}
