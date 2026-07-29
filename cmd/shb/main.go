// SPDX-License-Identifier: AGPL-3.0-or-later
// Copyright (C) 2026 Siyovush Hamidov and The Hadaf Contributors

package main

import (
	"shb/internal/application"

	_ "github.com/joho/godotenv/autoload"
)

func main() {
	a := application.NewApplication()
	a.Start()
}
