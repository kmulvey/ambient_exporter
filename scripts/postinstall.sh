#!/bin/sh
systemctl daemon-reload
systemctl enable ambient_exporter.service
systemctl start ambient_exporter.service
