package com.peanut996.cloudflarewarpspeedtest

import android.util.Log
import android.widget.TextView
import com.peanut996.cloudflarewarpspeedtest.models.SpeedTestConfig
import com.peanut996.cloudflarewarpspeedtest.models.SpeedTestResult
import kotlinx.coroutines.CoroutineScope
import kotlinx.coroutines.Dispatchers
import kotlinx.coroutines.Job
import kotlinx.coroutines.SupervisorJob
import kotlinx.coroutines.async
import kotlinx.coroutines.awaitAll
import kotlinx.coroutines.channels.Channel
import kotlinx.coroutines.launch
import kotlinx.coroutines.withContext
import java.net.DatagramPacket
import java.net.DatagramSocket
import java.net.InetAddress
import java.util.concurrent.atomic.AtomicBoolean
import kotlin.math.log

class SpeedTest {
    private val TAG = "SpeedTest"
    private val isRunning = AtomicBoolean(false)
    private var job: Job? = null
    private val coroutineScope = CoroutineScope(Dispatchers.Default + SupervisorJob())
    private val resultQueue = Channel<Any>(Channel.UNLIMITED)
    private var config = SpeedTestConfig()

    companion object {
        private val warpHandshakePacket: ByteArray by lazy {
            val hexString = "013cbdafb4135cac96a29484d7a0175ab152dd3e59be35049beadf758b8d48af14ca65f25a168934746fe8bc8867b1c17113d71c0fac5c141ef9f35783ffa5357c9871f4a006662b83ad71245a862495376a5fe3b4f2e1f06974d748416670e5f9b086297f652e6dfbf742fbfc63c3d8aeb175a3e9b7582fbc67c77577e4c0b32b05f92900000000000000000000000000000000"
            hexString.chunked(2).map { it.toInt(16).toByte() }.toByteArray()
        }

        private val ports = listOf(500, 854, 859, 864, 878, 880, 890, 891, 894, 903, 908, 928, 934, 939, 942, 943, 945, 946, 955, 968, 987, 988, 1002, 1010, 1014, 1018, 1070, 1074, 1180, 1387, 1701, 2408, 4500, 5050, 5242, 6515, 7103, 7152, 7156, 7281, 7559, 8319, 8742, 8854, 8886)
    }

    fun configure(config: SpeedTestConfig) {
        this.config = config
    }

    fun isRunning(): Boolean = isRunning.get()

    suspend fun getResults(): Any? = resultQueue.tryReceive().getOrNull()

    fun stop() {
        job?.cancel()
        isRunning.set(false)
    }

    private fun loadWarpIPRanges(): List<String> {
        val ipRanges = mutableListOf<String>()
        
        // Add IPv4 CIDR ranges
        val ipv4Ranges = if (!config.ipv6Mode) {
            listOf(
                "162.159.192.0/24",
                "162.159.193.0/24",
                "162.159.195.0/24",
                "162.159.204.0/24",
                "188.114.96.0/24",
                "188.114.97.0/24",
                "188.114.98.0/24",
                "188.114.99.0/24"
            )
        } else emptyList()

        // Process IPv4 ranges
        for (cidr in ipv4Ranges) {
            try {
                val (baseIP, mask) = cidr.split("/")
                val baseIPParts = baseIP.split(".").map { it.toInt() }
                val maskBits = mask.toInt()
                
                if (maskBits < 0 || maskBits > 32) continue
                
                val hostBits = 32 - maskBits
                val numHosts = 1 shl hostBits
                
                // Only generate IPs within the subnet
                for (i in 0 until numHosts) {
                    val lastOctet = baseIPParts[3] + i
                    if (lastOctet > 255) break
                    
                    val ip = "${baseIPParts[0]}.${baseIPParts[1]}.${baseIPParts[2]}.$lastOctet"
                    ipRanges.add(ip)
                }
            } catch (e: Exception) {
                Log.e(TAG, "Error processing CIDR $cidr: ${e.message}")
                continue
            }
        }
        
        // Add IPv6 CIDR ranges if enabled
        val ipv6Ranges = if (config.ipv6Mode) {
            listOf(
                "2606:4700:d0::/48",
                "2606:4700:d1::/48"
            )
        } else emptyList()

        // TODO: Add IPv6 support
        return ipRanges
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
                val results = mutableListOf<SpeedTestResult>()
                var testedCount = 0
                val allEndpoints = ipRanges.flatMap { ip -> 
                    ports.map { port -> Pair(ip, port) }
                }
                // Randomly select maxScanCount endpoints
                val endpoints = allEndpoints.shuffled().take(config.maxScanCount)
                val totalCount = endpoints.size
                Log.d(TAG, "total count: $totalCount")

                withContext(Dispatchers.IO) {
                    val jobs = endpoints.map { (ip, port) ->
                        async {
                            val result = testEndpoint(ip, port) { 
                                testedCount++
                                // Calculate percentage and update progress
                                val percentage = (testedCount.toFloat() / totalCount.toFloat() * 100).toInt()
                                launch { 
                                    // Send intermediate result if all tests are complete
                                    if (testedCount == totalCount) {
                                        if (results.isEmpty()) {
                                            resultQueue.send("No available endpoints found after testing ${endpoints.size} combinations")
                                        }
                                    }
                                }
                            }
                            result
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

                // Send final results or "No available endpoints" message
                if (filteredResults.isEmpty()) {
                    resultQueue.send("No available endpoints found after testing ${endpoints.size} combinations")
                } else {
                    resultQueue.send(filteredResults)
                }
            } catch (e: Exception) {
                Log.e(TAG, "Error during speed test: ${e.message}")
                resultQueue.send("Error during speed test: ${e.message}")
            } finally {
                isRunning.set(false)
            }
        }
    }

    private suspend fun testEndpoint(ipAddr: String, port: Int, onTested: () -> Unit): SpeedTestResult? = withContext(Dispatchers.IO) {
        Log.d(TAG, "Testing endpoint: $ipAddr:$port")
        var socket: DatagramSocket? = null
        
        try {
            socket = DatagramSocket()
            socket.soTimeout = 2000 // 2 second timeout
            
            val address = InetAddress.getByName(ipAddr)
            var successfulPings = 0
            var totalDelay = 0L
            var lastError: String? = null
            
            repeat(config.pingTimes) {
                val startTime = System.currentTimeMillis()
                val packet = DatagramPacket(warpHandshakePacket, warpHandshakePacket.size, address, port)
                
                try {
                    socket.send(packet)
                    
                    val receiveData = ByteArray(92) // WireGuard handshake response size
                    val receivePacket = DatagramPacket(receiveData, receiveData.size)
                    socket.receive(receivePacket)
                    val delay = System.currentTimeMillis() - startTime
                    if (delay <= config.maxDelay) { // Only count responses within maxDelay
                        totalDelay += delay
                        successfulPings++
                        // Send immediate success notification
                        if (successfulPings == 1) { // Only send on first success to avoid spam
                            resultQueue.send("Found working endpoint - IP: $ipAddr, Port: $port, Latency: ${delay}ms")
                        }
                    }
                } catch (e: Exception) {
                    lastError = e.message
                    Log.d(TAG, "Ping attempt ${it + 1}/${config.pingTimes} failed for $ipAddr:$port: ${e.message}")
                }
            }
            
            if (successfulPings > 0) {
                val avgDelay = totalDelay / successfulPings
                val lossRate = 1.0 - (successfulPings.toDouble() / config.pingTimes)
                
                if (avgDelay <= config.maxDelay && lossRate <= config.maxLossRate) {
                    Log.d(TAG, "Ping successful for $ipAddr:$port")
                    return@withContext SpeedTestResult(
                        ip = ipAddr,
                        port = port,
                        delay = avgDelay.toInt(),
                        lossRate = lossRate
                    )
                } else {
                    Log.d(TAG, "Endpoint $ipAddr:$port excluded: delay=$avgDelay, lossRate=$lossRate")
                }
            } else {
                Log.d(TAG, "All pings failed for $ipAddr:$port. Last error: $lastError")
            }
        } catch (e: Exception) {
            Log.e(TAG, "Error testing endpoint $ipAddr:$port: ${e.message}")
        } finally {
            socket?.close()
            onTested() // Ensure we count the attempt
        }
        null
    }
}
