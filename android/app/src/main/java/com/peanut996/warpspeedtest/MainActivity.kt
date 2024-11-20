package com.peanut996.warpspeedtest

import android.os.Bundle
import android.util.Log
import android.widget.Button
import android.widget.TextView
import androidx.appcompat.app.AppCompatActivity
import androidx.lifecycle.lifecycleScope
import kotlinx.coroutines.Dispatchers
import kotlinx.coroutines.launch
import kotlinx.coroutines.withContext

class MainActivity : AppCompatActivity() {
    private lateinit var speedTest: SpeedTest
    private lateinit var resultText: TextView
    private lateinit var startButton: Button
    private lateinit var stopButton: Button

    override fun onCreate(savedInstanceState: Bundle?) {
        super.onCreate(savedInstanceState)
        setContentView(R.layout.activity_main)

        speedTest = SpeedTest()
        setupViews()
        setupListeners()
        
        // Configure the speed test with default settings
        val defaultConfig = """{
            "threadCount": 200,
            "pingTimes": 10,
            "maxScanCount": 5000,
            "maxDelay": 300,
            "minDelay": 0,
            "maxLossRate": 1.0,
            "testAllCombos": false,
            "ipv6Mode": false,
            "resultDisplayCount": 10
        }"""
        
        try {
            speedTest.configure(defaultConfig)
        } catch (e: Exception) {
            Log.e("MainActivity", "Failed to configure speed test: ${e.message}")
        }
    }

    private fun setupViews() {
        resultText = findViewById(R.id.resultText)
        startButton = findViewById(R.id.startButton)
        stopButton = findViewById(R.id.stopButton)
    }

    private fun setupListeners() {
        startButton.setOnClickListener {
            startButton.isEnabled = false
            stopButton.isEnabled = true
            resultText.text = "Testing..."
            
            lifecycleScope.launch {
                speedTest.start()
                updateResults()
            }
        }

        stopButton.setOnClickListener {
            speedTest.stop()
            startButton.isEnabled = true
            stopButton.isEnabled = false
            resultText.text = "Test stopped"
        }
    }

    private suspend fun updateResults() {
        while (speedTest.isRunning()) {
            withContext(Dispatchers.Main) {
                val results = speedTest.getResults()
                if (results != "[]") {
                    resultText.text = results
                }
            }
            kotlinx.coroutines.delay(1000)
        }
        
        withContext(Dispatchers.Main) {
            startButton.isEnabled = true
            stopButton.isEnabled = false
        }
    }
}
