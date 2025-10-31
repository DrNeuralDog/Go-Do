package assets

import (
    "embed"
    "log"

    "fyne.io/fyne/v2"
)

// Icons are stored under doc/Icons. We embed the whole folder to support filenames with spaces.
//go:embed Icons/*
var iconsFS embed.FS

// Exported themed header icons as fyne resources
var (
    HeaderIconDark  fyne.Resource
    HeaderIconLight fyne.Resource
)

func init() {
    // Read Dark icon
    if data, err := iconsFS.ReadFile("Icons/Icon_Work_Version for Dark Header.png"); err == nil {
        HeaderIconDark = fyne.NewStaticResource("header-dark.png", data)
    } else {
        log.Printf("assets: failed to load dark header icon: %v", err)
    }

    // Read Light icon
    if data, err := iconsFS.ReadFile("Icons/Icon_Work_Version for Light Header.png"); err == nil {
        HeaderIconLight = fyne.NewStaticResource("header-light.png", data)
    } else {
        log.Printf("assets: failed to load light header icon: %v", err)
    }
}


