package main

import "time"

const CookieExpirationTime = 24 * time.Hour
const TimeFormat string = "2006-01-02 15:04:05"
const refreshTokenExpirationDays = 60
const refreshTokenExpiration = refreshTokenExpirationDays * 24 * time.Hour
