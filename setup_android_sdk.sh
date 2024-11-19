#!/bin/bash

# Set up environment variables
export ANDROID_HOME=$HOME/Library/Android/sdk
export PATH=$PATH:$ANDROID_HOME/tools
export PATH=$PATH:$ANDROID_HOME/platform-tools

# Create SDK directory if it doesn't exist
mkdir -p $ANDROID_HOME

# Install required SDK components
sdkmanager --sdk_root=$ANDROID_HOME "platform-tools" "platforms;android-34" "build-tools;34.0.0"

# Accept licenses
yes | sdkmanager --sdk_root=$ANDROID_HOME --licenses

echo "Android SDK setup completed!"
