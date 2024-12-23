package com.peanut996.cloudflarewarpspeedtest

import android.os.Bundle
import android.util.Log
import android.view.View
import android.widget.Button
import android.widget.TextView
import androidx.appcompat.app.AppCompatActivity
import androidx.lifecycle.lifecycleScope
import androidx.recyclerview.widget.LinearLayoutManager
import androidx.recyclerview.widget.RecyclerView
import com.peanut996.cloudflarewarpspeedtest.models.SpeedTestConfigBuilder
import com.peanut996.cloudflarewarpspeedtest.models.SpeedTestResult
import kotlinx.coroutines.Dispatchers
import kotlinx.coroutines.delay
import kotlinx.coroutines.launch
import kotlinx.coroutines.withContext

class MainActivity : AppCompatActivity() {
    private val TAG = "MainActivity"
    private lateinit var speedTest: SpeedTest
    private lateinit var progressTextView: TextView
    private lateinit var startButton: Button
    private lateinit var stopButton: Button
    private lateinit var ipPortTextView: TextView
    private lateinit var settingsButton: View
    private lateinit var resultRecyclerView: RecyclerView
    private lateinit var resultHeader: View
    private lateinit var resultAdapter: SpeedTestResultAdapter
    private var currentConfig = SpeedTestConfigBuilder.createDefault()

    override fun onCreate(savedInstanceState: Bundle?) {
        super.onCreate(savedInstanceState)
        setContentView(R.layout.activity_main)

        speedTest = SpeedTest()
        setupViews()
        setupRecyclerView()
        setupClickListeners()
        configureSpeedTest()
    }

    private fun setupViews() {
        progressTextView = findViewById(R.id.progressText)
        startButton = findViewById(R.id.startButton)
        stopButton = findViewById(R.id.stopButton)
        ipPortTextView = findViewById(R.id.ipPortText)
        settingsButton = findViewById(R.id.settingsButton)
        resultRecyclerView = findViewById(R.id.resultRecyclerView)
        resultHeader = findViewById(R.id.resultHeader)
        
        // Initially disable stop button and hide progress
        stopButton.isEnabled = false
        progressTextView.visibility = View.GONE
        resultHeader.visibility = View.GONE
        resultRecyclerView.visibility = View.GONE
    }

    private fun setupRecyclerView() {
        resultAdapter = SpeedTestResultAdapter()
        resultRecyclerView.apply {
            layoutManager = LinearLayoutManager(context)
            adapter = resultAdapter
            setHasFixedSize(true)
        }
    }

    private fun setupClickListeners() {
        startButton.setOnClickListener {
            startSpeedTest()
        }

        stopButton.setOnClickListener {
            stopSpeedTest()
        }

        settingsButton.setOnClickListener {
            showConfigDialog()
        }
    }

    private fun showConfigDialog() {
        SpeedTestConfigDialog.newInstance(currentConfig).apply {
            setOnConfigUpdatedListener { config ->
                currentConfig = config
                configureSpeedTest()
            }
        }.show(supportFragmentManager, "config_dialog")
    }

    private fun configureSpeedTest() {
        try {
            speedTest.configure(currentConfig)
        } catch (e: Exception) {
            Log.e(TAG, "Failed to configure speed test: ${e.message}")
            progressTextView.text = "Configuration error: ${e.message}"
        }
    }

    private fun startSpeedTest() {
        lifecycleScope.launch {
            try {
                startButton.isEnabled = false
                stopButton.isEnabled = true
                progressTextView.visibility = View.VISIBLE
                resultAdapter.clearResults()
                resultHeader.visibility = View.GONE
                resultRecyclerView.visibility = View.GONE
                ipPortTextView.text = ""

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

                // Process any remaining results in the channel
                var finalResult: Any?
                do {
                    finalResult = speedTest.getResults()
                    if (finalResult != null) {
                        withContext(Dispatchers.Main) {
                            updateProgress(finalResult)
                        }
                    }
                } while (finalResult != null)

            } catch (e: Exception) {
                Log.e(TAG, "Error during speed test: ${e.message}")
                withContext(Dispatchers.Main) {
                    progressTextView.text = "Error: ${e.message}"
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
        stopButton.isEnabled = false
        progressTextView.visibility = View.GONE
        resultHeader.visibility = View.GONE
        resultRecyclerView.visibility = View.GONE
        progressTextView.text = "Speed test stopped"
    }

    private fun updateProgress(result: Any) {
        when (result) {
            is String -> {
                progressTextView.text = result
            }
            is List<*> -> {
                @Suppress("UNCHECKED_CAST")
                val speedTestResults = result as List<SpeedTestResult>
                resultAdapter.updateResults(speedTestResults)
                if (speedTestResults.isNotEmpty()) {
                    progressTextView.visibility = View.GONE
                    resultHeader.visibility = View.VISIBLE
                    resultRecyclerView.visibility = View.VISIBLE
                }
            }
        }
    }
}
