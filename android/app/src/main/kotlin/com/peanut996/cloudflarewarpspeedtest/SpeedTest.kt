package com.peanut996.cloudflarewarpspeedtest

import android.util.Log
import com.google.gson.Gson
import kotlinx.coroutines.*
import java.net.DatagramPacket
import java.net.DatagramSocket
import java.net.InetAddress
import java.util.concurrent.ConcurrentLinkedQueue
import java.util.concurrent.atomic.AtomicBoolean

class SpeedTest {
    private val TAG = "SpeedTest"
    private val gson = Gson()
    private var job: Job? = null
    private val isRunning = AtomicBoolean(false)
    private val resultQueue = ConcurrentLinkedQueue<String>()
    private val coroutineScope = CoroutineScope(Dispatchers.Default + Job())

    data class TestConfig(
        val threadCount: Int = 200,
        val pingTimes: Int = 10,
        val maxScanCount: Int = 5000,
        val maxDelay: Int = 300,
        val minDelay: Int = 0,
        val maxLossRate: Double = 1.0,
        val testAllCombos: Boolean = false,
        val ipv6Mode: Boolean = false,
        val resultDisplayCount: Int = 10,
        val customIpFile: String = "",
        val customIpText: String = "",
        val privateKey: String = "",
        val publicKey: String = ""
    )

    data class TestResult(
        val ip: String,
        val port: Int,
        val delay: Long,
        val lossRate: Double
    )

    private var config: TestConfig = TestConfig()

    fun configure(configJson: String) {
        try {
            config = gson.fromJson(configJson, TestConfig::class.java)
        } catch (e: Exception) {
            Log.e(TAG, "Failed to parse config: ${e.message}")
            throw IllegalArgumentException("Invalid configuration JSON")
        }
    }

    fun start() {
        if (isRunning.get()) {
            Log.w(TAG, "Speed test is already running")
            return
        }

        isRunning.set(true)
        job = coroutineScope.launch {
            try {
                val ipRanges = loadWarpIPRanges()
                val results = mutableListOf<TestResult>()

                withContext(Dispatchers.IO) {
                    val jobs = ipRanges.map { ipAddr ->
                        async {
                            testEndpoint(ipAddr)
                        }
                    }
                    jobs.awaitAll().filterNotNull().forEach { results.add(it) }
                }

                // Sort results by delay and filter based on config
                val filteredResults = results
                    .filter { it.delay >= config.minDelay && it.delay <= config.maxDelay }
                    .filter { it.lossRate <= config.maxLossRate }
                    .sortedBy { it.delay }
                    .take(config.resultDisplayCount)

                resultQueue.offer(gson.toJson(filteredResults))
            } catch (e: Exception) {
                Log.e(TAG, "Error during speed test: ${e.message}")
            } finally {
                isRunning.set(false)
            }
        }
    }

    fun stop() {
        job?.cancel()
        isRunning.set(false)
    }

    fun getResults(): String {
        return resultQueue.poll() ?: "[]"
    }

    fun isRunning(): Boolean = isRunning.get()

    private suspend fun testEndpoint(ipAddr: String): TestResult? = withContext(Dispatchers.IO) {
        try {
            val socket = DatagramSocket()
            socket.soTimeout = 1000 // 1 second timeout
            
            val address = InetAddress.getByName(ipAddr)
            val handshakePacket = WireGuardHandshake.createHandshakeInitiation(
                privateKey = config.privateKey,
                publicKey = config.publicKey
            )
            
            var successfulPings = 0
            var totalDelay = 0L
            
            repeat(config.pingTimes) {
                val startTime = System.currentTimeMillis()
                val packet = DatagramPacket(handshakePacket, handshakePacket.size, address, 2408)
                
                try {
                    socket.send(packet)
                    
                    val receiveData = ByteArray(92) // WireGuard handshake response size
                    val receivePacket = DatagramPacket(receiveData, receiveData.size)
                    socket.receive(receivePacket)
                    
                    val delay = System.currentTimeMillis() - startTime
                    totalDelay += delay
                    successfulPings++
                } catch (e: Exception) {
                    Log.d(TAG, "Ping failed for $ipAddr: ${e.message}")
                }
            }
            
            socket.close()
            
            if (successfulPings > 0) {
                val avgDelay = totalDelay / successfulPings
                val lossRate = 1.0 - (successfulPings.toDouble() / config.pingTimes)
                
                return@withContext TestResult(
                    ip = ipAddr,
                    port = 2408,
                    delay = avgDelay,
                    lossRate = lossRate
                )
            }
            
            null
        } catch (e: Exception) {
            Log.e(TAG, "Error testing endpoint $ipAddr: ${e.message}")
            null
        }
    }

    private fun loadWarpIPRanges(): List<String> {
        // TODO: Implement IP range loading from resources or config
        // This is a placeholder - actual implementation needs to match the Go version's IP range loading
        return listOf(
            "162.159.192.1",
            "162.159.193.1",
            "162.159.195.1",
            "162.159.204.1"
        )
    }
}
