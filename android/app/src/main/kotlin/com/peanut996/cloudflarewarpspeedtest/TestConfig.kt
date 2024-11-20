package com.peanut996.cloudflarewarpspeedtest

data class TestConfig(
    var threadCount: Int = 200,
    var pingTimes: Int = 10,
    var maxScanCount: Int = 1000,
    var maxDelay: Long = 3,
    var minDelay: Long = 0,
    var maxLossRate: Double = 0.0,
    var testAllCombos: Boolean = false,
    var ipv6Mode: Boolean = false,
    var resultDisplayCount: Int = 10,
    var customIPFile: String = "",
    var customIPText: String = "",
    var privateKey: String = "",
    var publicKey: String = "bmXOC+F1FxEMF9dyiK2H5/1SUtzH0JuVo51h2wPfgyo="
)
