#!/usr/bin/env pwsh
# Простой скрипт для настройки кодировки

# Устанавливаем UTF-8 кодировку
$OutputEncoding = [System.Text.Encoding]::UTF8
[Console]::OutputEncoding = [System.Text.Encoding]::UTF8

Write-Host "UTF-8 encoding set successfully"