package com.peanut996.cloudflarewarpspeedtest

import android.os.Bundle
import android.util.Log
import android.widget.Button
import android.widget.TextView
import androidx.appcompat.app.AppCompatActivity
import androidx.lifecycle.lifecycleScope
import kotlinx.coroutines.Dispatchers
import kotlinx.coroutines.delay
import kotlinx.coroutines.launch
import kotlinx.coroutines.withContext

class MainActivity : AppCompatActivity() {
    private val TAG = "MainActivity"
    private lateinit var speedTest: SpeedTest
    private lateinit var resultTextView: TextView
    private lateinit var startButton: Button
    private lateinit var stopButton: Button

    override fun onCreate(savedInstanceState: Bundle?) {
        super.onCreate(savedInstanceState)
        setContentView(R.layout.activity_main)

        speedTest = SpeedTest()
        setupViews()
        setupClickListeners()
        configureSpeedTest()
    }

    private fun setupViews() {
        resultTextView = findViewById(R.id.resultText)
        startButton = findViewById(R.id.startButton)
        stopButton = findViewById(R.id.stopButton)
        
        // Initially disable stop button
        stopButton.isEnabled = false
    }

    private fun setupClickListeners() {
        startButton.setOnClickListener {
            startSpeedTest()
        }

        stopButton.setOnClickListener {
            stopSpeedTest()
        }
    }

    private fun configureSpeedTest() {
        val defaultConfig = """
            {
                "threadCount": 200,
                "pingTimes": 10,
                "maxScanCount": 5000,
                "maxDelay": 300,
                "minDelay": 0,
                "maxLossRate": 1.0,
                "testAllCombos": false,
                "ipv6Mode": false,
                "resultDisplayCount": 10
            }
        """.trimIndent()

        try {
            speedTest.configure(defaultConfig)
        } catch (e: Exception) {
            Log.e(TAG, "Failed to configure speed test: ${e.message}")
            resultTextView.text = "Configuration error: ${e.message}"
        }
    }

    private fun startSpeedTest() {
        lifecycleScope.launch {
            try {
                startButton.isEnabled = false
                stopButton.isEnabled = true
                resultTextView.text = "Starting speed test..."

                speedTest.start()

                // Poll for results while the test is running
                while (speedTest.isRunning()) {
                    val results = speedTest.getResults()
                    if (results != "[]") {
                        withContext(Dispatchers.Main) {
                            resultTextView.text = formatResults(results)
                        }
                    }
                    delay(1000) // Update every second
                }

                // Get final results
                val finalResults = speedTest.getResults()
                if (finalResults != "[]") {
                    withContext(Dispatchers.Main) {
                        resultTextView.text = formatResults(finalResults)
                    }
                }
            } catch (e: Exception) {
                Log.e(TAG, "Error during speed test: ${e.message}")
                withContext(Dispatchers.Main) {
                    resultTextView.text = "Error: ${e.message}"
                }
            } finally {
                withContext(Dispatchers.Main) {
                    startButton.isEnabled = true
                    stopButton.isEnabled = false
                }
            }
        }
    }

    private fun stopSpeedTest() {
        speedTest.stop()
        startButton.isEnabled = true
        stopButton.isEnabled = false
        resultTextView.text = "Speed test stopped"
    }

    private fun formatResults(jsonResults: String): String {
        // TODO: Implement better results formatting
        return "Test Results:\n$jsonResults"
    }
}
