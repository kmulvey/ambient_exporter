#!/bin/sh
systemctl daemon-reload
systemctl enable ambient_exporter.service
systemctl start ambient_exporter.service

# WiFi watchdog for brcmfmac silent hang (Raspberry Pi BCM4345)
cp /usr/share/ambient_exporter/wifi-watchdog.sh /usr/local/bin/wifi-watchdog.sh
chmod +x /usr/local/bin/wifi-watchdog.sh
cp /usr/share/ambient_exporter/wifi-watchdog.timer /etc/systemd/system/wifi-watchdog.timer
cp /usr/share/ambient_exporter/wifi-watchdog.service /etc/systemd/system/wifi-watchdog.service
systemctl daemon-reload
systemctl enable --now wifi-watchdog.timer
