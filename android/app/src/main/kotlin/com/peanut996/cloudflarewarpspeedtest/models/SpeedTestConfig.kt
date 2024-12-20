package com.peanut996.cloudflarewarpspeedtest.models

data class SpeedTestConfig(
    val threadCount: Int = 200,
    val pingTimes: Int = 1,
    val maxScanCount: Int = 10,
    val maxDelay: Int = 300,
    val minDelay: Int = 0,
    val maxLossRate: Double = 1.0,
    val testAllCombos: Boolean = false,
    val ipv6Mode: Boolean = false,
    val resultDisplayCount: Int = 10
)
