package main

import "time"

const CookieExpirationTime = 15 * time.Minute
const TimeFormat string = "2006-01-02 15:04:05"
const refreshTokenExpirationDays = 30
const refreshTokenExpiration = refreshTokenExpirationDays * 24 * time.Hour
const MinPasswordLength = 8

const defaultMaxCartSize = 1000

const (
	minInt32 = -2147483648
	maxInt32 = 2147483647
)
