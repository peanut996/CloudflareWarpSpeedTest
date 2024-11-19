package com.peanut996.cloudflarewarpspeedtest

import android.os.Bundle
import android.widget.Toast
import androidx.appcompat.app.AppCompatActivity
import androidx.lifecycle.lifecycleScope
import com.peanut996.cloudflarewarpspeedtest.databinding.ActivityMainBinding
import kotlinx.coroutines.Dispatchers
import kotlinx.coroutines.launch
import kotlinx.coroutines.withContext
import mobile.SpeedTest
import mobile.TestConfig
import org.json.JSONArray
import org.json.JSONObject

class MainActivity : AppCompatActivity() {
    private lateinit var binding: ActivityMainBinding
    private var speedTest: SpeedTest? = null

    override fun onCreate(savedInstanceState: Bundle?) {
        super.onCreate(savedInstanceState)
        binding = ActivityMainBinding.inflate(layoutInflater)
        setContentView(binding.root)

        setupSpeedTest()
        setupUI()
    }

    private fun setupSpeedTest() {
        speedTest = SpeedTest().apply {
            // Initialize any default configuration here if needed
        }
    }

    private fun setupUI() {
        binding.startButton.setOnClickListener {
            if (speedTest?.isRunning() == true) {
                speedTest?.stop()
                binding.startButton.text = "Start Speed Test"
                return@setOnClickListener
            }

            startSpeedTest()
        }
    }

    private fun startSpeedTest() {
        val config = TestConfig().apply {
            threadCount = (binding.threadCountInput.text.toString().toIntOrNull() ?: 200).toLong()
            maxDelay = binding.maxDelayInput.text.toString().toLongOrNull() ?: 300L
            minDelay = 0L
            maxLossRate = 0.0
            iPv6Mode = binding.ipv6Switch.isChecked
            pingTimes = 10
            maxScanCount = 5000
            resultDisplayCount = 10
        }

        lifecycleScope.launch {
            try {
                binding.startButton.text = "Stop Test"
                binding.resultText.text = "Starting speed test...\n"

                withContext(Dispatchers.IO) {
                    speedTest?.configure(JSONObject().apply {
                        put("thread_count", config.threadCount)
                        put("max_delay", config.maxDelay)
                        put("min_delay", config.minDelay)
                        put("max_loss_rate", config.maxLossRate)
                        put("ipv6_mode", config.iPv6Mode)
                        put("ping_times", config.pingTimes)
                        put("max_scan_count", config.maxScanCount)
                        put("result_display_count", config.resultDisplayCount)
                    }.toString())

                    speedTest?.start()

                    while (speedTest?.isRunning() == true) {
                        val results = speedTest?.getResults()
                        withContext(Dispatchers.Main) {
                            updateResults(results)
                        }
                        kotlinx.coroutines.delay(1000)
                    }
                }
            } catch (e: Exception) {
                withContext(Dispatchers.Main) {
                    Toast.makeText(this@MainActivity, "Error: ${e.message}", Toast.LENGTH_LONG).show()
                    binding.startButton.text = "Start Speed Test"
                }
            }
        }
    }

    private fun updateResults(resultsJson: String?) {
        if (resultsJson == null) return

        try {
            val results = JSONArray(resultsJson)
            val formattedResults = StringBuilder()
            
            for (i in 0 until results.length()) {
                val result = results.getJSONObject(i)
                formattedResults.append(
                    String.format(
                        "IP: %-30s Loss: %3d%%  Latency: %.2fms\n",
                        result.getString("ip"),
                        (result.getDouble("loss") * 100).toInt(),
                        result.getDouble("latency")
                    )
                )
            }

            binding.resultText.text = formattedResults.toString()
        } catch (e: Exception) {
            if (resultsJson.contains("status")) {
                // Still running
                return
            }
            binding.resultText.append("Error parsing results: ${e.message}\n")
        }
    }

    override fun onDestroy() {
        super.onDestroy()
        speedTest?.stop()
    }
}
