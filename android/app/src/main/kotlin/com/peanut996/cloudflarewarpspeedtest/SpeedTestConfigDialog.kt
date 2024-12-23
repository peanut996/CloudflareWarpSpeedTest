package com.peanut996.cloudflarewarpspeedtest

import android.app.Dialog
import android.os.Bundle
import android.view.LayoutInflater
import android.view.View
import android.view.ViewGroup
import android.widget.CheckBox
import android.widget.RadioButton
import android.widget.Toast
import androidx.fragment.app.DialogFragment
import com.google.android.material.dialog.MaterialAlertDialogBuilder
import com.google.android.material.textfield.TextInputEditText
import com.peanut996.cloudflarewarpspeedtest.models.SpeedTestConfig
import com.peanut996.cloudflarewarpspeedtest.models.SpeedTestConfigBuilder

class SpeedTestConfigDialog : DialogFragment() {
    private var currentConfig = SpeedTestConfigBuilder.createDefault()
    private var onConfigUpdated: ((SpeedTestConfig) -> Unit)? = null

    private lateinit var threadCountInput: TextInputEditText
    private lateinit var pingTimesInput: TextInputEditText
    private lateinit var maxScanCountInput: TextInputEditText
    private lateinit var maxDelayInput: TextInputEditText
    private lateinit var minDelayInput: TextInputEditText
    private lateinit var maxLossRateInput: TextInputEditText
    private lateinit var resultDisplayCountInput: TextInputEditText
    private lateinit var testAllCombosCheckbox: CheckBox
    private lateinit var ipv6ModeCheckbox: CheckBox
    private lateinit var defaultConfigRadio: RadioButton
    private lateinit var fastConfigRadio: RadioButton
    private lateinit var accurateConfigRadio: RadioButton

    override fun onCreateDialog(savedInstanceState: Bundle?): Dialog {
        val view = LayoutInflater.from(requireContext())
            .inflate(R.layout.dialog_speed_test_config, null)

        initializeViews(view)
        setupPresetConfigs()
        loadCurrentConfig()

        return MaterialAlertDialogBuilder(requireContext())
            .setTitle("Speed Test Configuration")
            .setView(view)
            .setPositiveButton("Apply") { _, _ -> applyConfig() }
            .setNegativeButton("Cancel", null)
            .create()
    }

    private fun initializeViews(view: View) {
        threadCountInput = view.findViewById(R.id.threadCountInput)
        pingTimesInput = view.findViewById(R.id.pingTimesInput)
        maxScanCountInput = view.findViewById(R.id.maxScanCountInput)
        maxDelayInput = view.findViewById(R.id.maxDelayInput)
        minDelayInput = view.findViewById(R.id.minDelayInput)
        maxLossRateInput = view.findViewById(R.id.maxLossRateInput)
        resultDisplayCountInput = view.findViewById(R.id.resultDisplayCountInput)
        testAllCombosCheckbox = view.findViewById(R.id.testAllCombosCheckbox)
        ipv6ModeCheckbox = view.findViewById(R.id.ipv6ModeCheckbox)
        defaultConfigRadio = view.findViewById(R.id.defaultConfigRadio)
        fastConfigRadio = view.findViewById(R.id.fastConfigRadio)
        accurateConfigRadio = view.findViewById(R.id.accurateConfigRadio)
    }

    private fun setupPresetConfigs() {
        defaultConfigRadio.setOnCheckedChangeListener { _, isChecked ->
            if (isChecked) updateInputs(SpeedTestConfigBuilder.createDefault())
        }
        fastConfigRadio.setOnCheckedChangeListener { _, isChecked ->
            if (isChecked) updateInputs(SpeedTestConfigBuilder.createFastConfig())
        }
        accurateConfigRadio.setOnCheckedChangeListener { _, isChecked ->
            if (isChecked) updateInputs(SpeedTestConfigBuilder.createAccurateConfig())
        }
    }

    private fun loadCurrentConfig() {
        updateInputs(currentConfig)
    }

    private fun updateInputs(config: SpeedTestConfig) {
        threadCountInput.setText(config.threadCount.toString())
        pingTimesInput.setText(config.pingTimes.toString())
        maxScanCountInput.setText(config.maxScanCount.toString())
        maxDelayInput.setText(config.maxDelay.toString())
        minDelayInput.setText(config.minDelay.toString())
        maxLossRateInput.setText(config.maxLossRate.toString())
        resultDisplayCountInput.setText(config.resultDisplayCount.toString())
        testAllCombosCheckbox.isChecked = config.testAllCombos
        ipv6ModeCheckbox.isChecked = config.ipv6Mode
    }

    private fun applyConfig() {
        try {
            val builder = SpeedTestConfigBuilder()
                .setThreadCount(threadCountInput.text.toString().toInt())
                .setPingTimes(pingTimesInput.text.toString().toInt())
                .setMaxScanCount(maxScanCountInput.text.toString().toInt())
                .setMaxDelay(maxDelayInput.text.toString().toInt())
                .setMinDelay(minDelayInput.text.toString().toInt())
                .setMaxLossRate(maxLossRateInput.text.toString().toDouble())
                .setResultDisplayCount(resultDisplayCountInput.text.toString().toInt())
                .setTestAllCombos(testAllCombosCheckbox.isChecked)
                .setIpv6Mode(ipv6ModeCheckbox.isChecked)

            val newConfig = builder.build()
            currentConfig = newConfig
            onConfigUpdated?.invoke(newConfig)
        } catch (e: Exception) {
            Toast.makeText(context, "Invalid input: ${e.message}", Toast.LENGTH_LONG).show()
        }
    }

    fun setConfig(config: SpeedTestConfig) {
        currentConfig = config
        if (this::threadCountInput.isInitialized) {
            loadCurrentConfig()
        }
    }

    fun setOnConfigUpdatedListener(listener: (SpeedTestConfig) -> Unit) {
        onConfigUpdated = listener
    }

    companion object {
        fun newInstance(config: SpeedTestConfig): SpeedTestConfigDialog {
            return SpeedTestConfigDialog().apply {
                currentConfig = config
            }
        }
    }
}
