#!/bin/bash

# Create necessary directories
mkdir -p app/src/main/res/values
mkdir -p app/src/main/kotlin/com/peanut996/cloudflarewarpspeedtest
mkdir -p app/libs

# Copy the AAR file
cp ../mobile/build/cloudflare_warp_speedtest.aar app/libs/

# Create strings.xml
cat > android/app/src/main/res/values/strings.xml << EOL
<?xml version="1.0" encoding="utf-8"?>
<resources>
    <string name="app_name">Cloudflare WARP Speed Test</string>
</resources>
EOL

# Create styles.xml
cat > android/app/src/main/res/values/themes.xml << EOL
<?xml version="1.0" encoding="utf-8"?>
<resources>
    <style name="Theme.CloudflareWarpSpeedTest" parent="Theme.MaterialComponents.DayNight.DarkActionBar">
        <item name="colorPrimary">@color/design_default_color_primary</item>
        <item name="colorPrimaryDark">@color/design_default_color_primary_dark</item>
        <item name="colorAccent">@color/design_default_color_secondary</item>
    </style>
</resources>
EOL

echo "Android project setup completed!"
echo "To build the project:"
echo "1. Open the 'android' folder in Android Studio"
echo "2. Let the project sync and download dependencies"
echo "3. Click 'Run' to build and run the app on your device or emulator"
