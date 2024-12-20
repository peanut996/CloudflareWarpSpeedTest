package com.peanut996.cloudflarewarpspeedtest

import android.os.Bundle
import android.util.Log
import android.view.View
import android.widget.Button
import android.widget.ProgressBar
import android.widget.TextView
import androidx.appcompat.app.AppCompatActivity
import androidx.lifecycle.lifecycleScope
import com.peanut996.cloudflarewarpspeedtest.models.SpeedTestConfig
import com.peanut996.cloudflarewarpspeedtest.models.SpeedTestResult
import kotlinx.coroutines.Dispatchers
import kotlinx.coroutines.delay
import kotlinx.coroutines.launch
import kotlinx.coroutines.withContext

class MainActivity : AppCompatActivity() {
    private val TAG = "MainActivity"
    private lateinit var speedTest: SpeedTest
    private lateinit var resultTextView: TextView
    private lateinit var progressTextView: TextView
    private lateinit var progressBar: ProgressBar
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
        progressTextView = findViewById(R.id.progressText)
        progressBar = findViewById(R.id.progressBar)
        startButton = findViewById(R.id.startButton)
        stopButton = findViewById(R.id.stopButton)
        
        // Initially disable stop button and hide progress
        stopButton.isEnabled = false
        progressBar.visibility = View.GONE
        progressTextView.visibility = View.GONE
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
        val config = SpeedTestConfig()
        try {
            speedTest.configure(config)
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
                progressBar.visibility = View.VISIBLE
                progressTextView.visibility = View.VISIBLE
                resultTextView.text = "Starting speed test..."

                speedTest.start()

                // Poll for results while the test is running
                while (speedTest.isRunning()) {
                    val result = speedTest.getResults()
                    if (result != null) {
                        withContext(Dispatchers.Main) {
                            updateProgress(result)
                        }
                    }
                    delay(500) // Update every 500ms
                }

                // Get final results
                val finalResult = speedTest.getResults()
                if (finalResult != null) {
                    withContext(Dispatchers.Main) {
                        updateProgress(finalResult)
                    }
                }
            } catch (e: Exception) {
                Log.e(TAG, "Error during speed test: ${e.message}")
                withContext(Dispatchers.Main) {
                    resultTextView.text = "Error: ${e.message}"
                    progressBar.visibility = View.GONE
                    progressTextView.visibility = View.GONE
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
        progressBar.visibility = View.GONE
        progressTextView.visibility = View.GONE
        resultTextView.text = "Speed test stopped"
    }

    private fun updateProgress(result: Any) {
        when (result) {
            is List<*> -> {
                @Suppress("UNCHECKED_CAST")
                val results = result as List<SpeedTestResult>
                resultTextView.text = formatResults(results)
                progressBar.visibility = View.GONE
                progressTextView.visibility = View.GONE
            }
            is String -> {
                val progressMatch = Regex("Progress: (\\d+)/(\\d+) \\((\\d+)%\\)").find(result)
                if (progressMatch != null) {
                    val (current, total, percentage) = progressMatch.destructured
                    progressBar.progress = percentage.toInt()
                    progressTextView.text = "Testing endpoints: $current/$total ($percentage%)"
                } else if (result.startsWith("Found working endpoint")) {
                    // Append the new endpoint to existing text while keeping progress
                    val currentText = resultTextView.text.toString()
                    resultTextView.text = if (currentText.startsWith("Starting")) {
                        result + "\n"
                    } else {
                        currentText + result + "\n"
                    }
                } else {
                    resultTextView.text = result
                }
            }
        }
    }

    private fun formatResults(results: List<SpeedTestResult>): String {
        if (results.isEmpty()) {
            return "No results yet"
        }

        return buildString {
            appendLine("Test Results:")
            results.forEachIndexed { index, result ->
                appendLine("${index + 1}. IP: ${result.ip}:${result.port}")
                appendLine("   Delay: ${result.delay}ms")
                appendLine("   Loss Rate: ${String.format("%.2f", result.lossRate * 100)}%")
            }
        }
    }
}
