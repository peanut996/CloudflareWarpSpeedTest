package com.peanut996.cloudflarewarpspeedtest.models

/**
 * Configuration for speed test
 *
 * @property threadCount Number of concurrent threads for testing (1-1000)
 * @property pingTimes Number of ping attempts per endpoint (1-10)
 * @property maxScanCount Maximum number of endpoints to scan (1-10000)
 * @property maxDelay Maximum acceptable delay in milliseconds (50-2000)
 * @property minDelay Minimum acceptable delay in milliseconds (0-1000)
 * @property maxLossRate Maximum acceptable packet loss rate (0.0-1.0)
 * @property testAllCombos Whether to test all IP:Port combinations
 * @property ipv6Mode Whether to use IPv6 addresses
 * @property resultDisplayCount Number of results to display (1-100)
 */
data class SpeedTestConfig(
    val threadCount: Int = 200,
    val pingTimes: Int = 1,
    val maxScanCount: Int = 100,
    val maxDelay: Int = 300,
    val minDelay: Int = 0,
    val maxLossRate: Double = 1.0,
    val testAllCombos: Boolean = false,
    val ipv6Mode: Boolean = false,
    val resultDisplayCount: Int = 10
)

/**
 * Builder class for SpeedTestConfig with validation
 */
class SpeedTestConfigBuilder {
    private var threadCount: Int = 200
    private var pingTimes: Int = 1
    private var maxScanCount: Int = 100
    private var maxDelay: Int = 300
    private var minDelay: Int = 0
    private var maxLossRate: Double = 1.0
    private var testAllCombos: Boolean = false
    private var ipv6Mode: Boolean = false
    private var resultDisplayCount: Int = 10

    /**
     * Set number of concurrent threads (1-1000)
     */
    fun setThreadCount(value: Int) = apply {
        require(value in 1..1000) { "Thread count must be between 1 and 1000" }
        threadCount = value
    }

    /**
     * Set number of ping attempts per endpoint (1-10)
     */
    fun setPingTimes(value: Int) = apply {
        require(value in 1..10) { "Ping times must be between 1 and 10" }
        pingTimes = value
    }

    /**
     * Set maximum number of endpoints to scan (1-10000)
     */
    fun setMaxScanCount(value: Int) = apply {
        require(value in 1..10000) { "Max scan count must be between 1 and 10000" }
        maxScanCount = value
    }

    /**
     * Set maximum acceptable delay in milliseconds (50-2000)
     */
    fun setMaxDelay(value: Int) = apply {
        require(value in 50..2000) { "Max delay must be between 50 and 2000 ms" }
        maxDelay = value
    }

    /**
     * Set minimum acceptable delay in milliseconds (0-1000)
     */
    fun setMinDelay(value: Int) = apply {
        require(value in 0..1000) { "Min delay must be between 0 and 1000 ms" }
        minDelay = value
    }

    /**
     * Set maximum acceptable packet loss rate (0.0-1.0)
     */
    fun setMaxLossRate(value: Double) = apply {
        require(value in 0.0..1.0) { "Max loss rate must be between 0.0 and 1.0" }
        maxLossRate = value
    }

    /**
     * Set whether to test all IP:Port combinations
     */
    fun setTestAllCombos(value: Boolean) = apply {
        testAllCombos = value
    }

    /**
     * Set whether to use IPv6 addresses
     */
    fun setIpv6Mode(value: Boolean) = apply {
        ipv6Mode = value
    }

    /**
     * Set number of results to display (1-100)
     */
    fun setResultDisplayCount(value: Int) = apply {
        require(value in 1..100) { "Result display count must be between 1 and 100" }
        resultDisplayCount = value
    }

    /**
     * Build SpeedTestConfig with validation
     */
    fun build(): SpeedTestConfig {
        require(minDelay < maxDelay) { "Min delay must be less than max delay" }
        
        return SpeedTestConfig(
            threadCount = threadCount,
            pingTimes = pingTimes,
            maxScanCount = maxScanCount,
            maxDelay = maxDelay,
            minDelay = minDelay,
            maxLossRate = maxLossRate,
            testAllCombos = testAllCombos,
            ipv6Mode = ipv6Mode,
            resultDisplayCount = resultDisplayCount
        )
    }

    companion object {
        /**
         * Create default config
         */
        fun createDefault() = SpeedTestConfig()

        /**
         * Create config optimized for speed (less accurate)
         */
        fun createFastConfig() = SpeedTestConfigBuilder()
            .setThreadCount(400)
            .setPingTimes(1)
            .setMaxScanCount(50)
            .setMaxDelay(500)
            .setMinDelay(0)
            .setMaxLossRate(1.0)
            .setResultDisplayCount(5)
            .build()

        /**
         * Create config optimized for accuracy (slower)
         */
        fun createAccurateConfig() = SpeedTestConfigBuilder()
            .setThreadCount(100)
            .setPingTimes(3)
            .setMaxScanCount(200)
            .setMaxDelay(1000)
            .setMinDelay(0)
            .setMaxLossRate(0.5)
            .setResultDisplayCount(20)
            .build()
    }
}
