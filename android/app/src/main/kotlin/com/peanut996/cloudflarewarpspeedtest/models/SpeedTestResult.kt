package com.peanut996.cloudflarewarpspeedtest.models

data class SpeedTestResult(
    val ip: String,
    val port: Int,
    val delay: Int,
    val lossRate: Double
)
