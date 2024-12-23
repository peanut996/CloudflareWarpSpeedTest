package com.peanut996.cloudflarewarpspeedtest

import android.content.ClipData
import android.content.ClipboardManager
import android.content.Context
import android.view.LayoutInflater
import android.view.View
import android.view.ViewGroup
import android.widget.TextView
import android.widget.Toast
import androidx.recyclerview.widget.RecyclerView
import com.peanut996.cloudflarewarpspeedtest.models.SpeedTestResult

class SpeedTestResultAdapter : RecyclerView.Adapter<SpeedTestResultAdapter.ViewHolder>() {
    private val results = mutableListOf<SpeedTestResult>()

    class ViewHolder(view: View) : RecyclerView.ViewHolder(view) {
        val ipPortText: TextView = view.findViewById(R.id.ipPortText)
        val detailsText: TextView = view.findViewById(R.id.detailsText)
    }

    override fun onCreateViewHolder(parent: ViewGroup, viewType: Int): ViewHolder {
        val view = LayoutInflater.from(parent.context)
            .inflate(R.layout.item_speed_test_result, parent, false)
        return ViewHolder(view)
    }

    override fun onBindViewHolder(holder: ViewHolder, position: Int) {
        val result = results[position]
        val ipPort = "${result.ip}:${result.port}"
        holder.ipPortText.text = ipPort
        holder.detailsText.text = "Delay: ${result.delay}ms, Loss Rate: ${String.format("%.1f", result.lossRate * 100)}%"
        
        holder.itemView.setOnClickListener {
            val context = holder.itemView.context
            val clipboard = context.getSystemService(Context.CLIPBOARD_SERVICE) as ClipboardManager
            val clip = ClipData.newPlainText("IP:Port", ipPort)
            clipboard.setPrimaryClip(clip)
            Toast.makeText(context, "Copied: $ipPort", Toast.LENGTH_SHORT).show()
        }
    }

    override fun getItemCount() = results.size

    fun updateResults(newResults: List<SpeedTestResult>) {
        results.clear()
        results.addAll(newResults)
        notifyDataSetChanged()
    }

    fun clearResults() {
        results.clear()
        notifyDataSetChanged()
    }
}
