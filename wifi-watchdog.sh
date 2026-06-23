#!/bin/bash
# Watchdog for brcmfmac WiFi firmware silent hang (BCM4345 on Raspberry Pi).
# The driver reports associated to NetworkManager but stops passing packets.
# This script bounces wlan0 if no data has been received in 15 minutes.
# Installed via: systemd timer (wifi-watchdog.timer)
LAST=$(journalctl -u ambient_exporter --since "15 minutes ago" --no-pager | grep "received a report" | wc -l)
if [ "$LAST" -eq 0 ]; then
    logger "wifi-watchdog: no data in 15 min, bouncing wlan0"
    nmcli device disconnect wlan0
    sleep 5
    nmcli device connect wlan0
    logger "wifi-watchdog: wlan0 reconnected"
fi
