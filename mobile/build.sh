#!/bin/bash

# Install gomobile if not already installed
go install golang.org/x/mobile/cmd/gomobile@latest
go install golang.org/x/mobile/cmd/gobind@latest

# Initialize gomobile
gomobile init

# Build for Android
gomobile bind -target=android -o=android/app/libs/speedtest.aar ./mobile

# Build for iOS (optional)
# gomobile bind -target=ios -o build/CloudflareWarpSpeedTest.framework ./mobile
