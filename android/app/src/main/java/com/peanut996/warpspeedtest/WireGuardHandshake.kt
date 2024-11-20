package com.peanut996.warpspeedtest

import android.util.Base64
import java.nio.ByteBuffer
import java.security.MessageDigest
import java.security.SecureRandom
import javax.crypto.Mac
import javax.crypto.spec.SecretKeySpec
import kotlin.experimental.and

class WireGuardHandshake {
    companion object {
        private const val NOISE_PUBLIC_KEY_SIZE = 32
        private const val NOISE_TIMESTAMP_SIZE = 12
        private const val BLAKE2S_SIZE = 16
        private const val WARP_PUBLIC_KEY = "bmXOC+F1FxEMF9dyiK2H5/1SUtzH0JuVo51h2wPfgyo="
        private val RESERVED_BYTES = byteArrayOf(60, -67, -81) // [60, 189, 175] in signed bytes

        fun createHandshakeInitiation(privateKey: String = "", publicKey: String = ""): ByteArray {
            val random = SecureRandom()
            val buffer = ByteBuffer.allocate(148) // Size of handshake initiation message

            // Message type (1 for handshake initiation)
            buffer.put(1.toByte())
            
            // Reserved bytes
            buffer.put(RESERVED_BYTES)
            
            // Sender index (random 32-bit number)
            buffer.putInt(random.nextInt())
            
            // Ephemeral key (random 32 bytes)
            val ephemeralKey = ByteArray(NOISE_PUBLIC_KEY_SIZE)
            random.nextBytes(ephemeralKey)
            buffer.put(ephemeralKey)
            
            // Static key (encrypted)
            val staticKey = ByteArray(NOISE_PUBLIC_KEY_SIZE)
            random.nextBytes(staticKey)
            buffer.put(staticKey)
            buffer.put(ByteArray(16)) // Poly1305 tag
            
            // Timestamp (current time in TAI64N format)
            val timestamp = createTAI64NTimestamp()
            buffer.put(timestamp)
            buffer.put(ByteArray(16)) // Poly1305 tag
            
            // MAC1 (Blake2s hash)
            buffer.put(ByteArray(BLAKE2S_SIZE))
            
            // MAC2 (Blake2s hash)
            buffer.put(ByteArray(BLAKE2S_SIZE))

            return buffer.array()
        }

        private fun createTAI64NTimestamp(): ByteArray {
            val buffer = ByteBuffer.allocate(NOISE_TIMESTAMP_SIZE)
            val currentTime = System.currentTimeMillis() / 1000L
            val tai64Time = currentTime + 4611686018427387914L // TAI epoch offset
            
            buffer.putLong(tai64Time)
            buffer.putInt((System.nanoTime() % 1_000_000_000).toInt())
            
            return buffer.array()
        }

        private fun decodeBase64(input: String): ByteArray {
            return Base64.decode(input, Base64.DEFAULT)
        }

        private fun blake2sHash(data: ByteArray): ByteArray {
            // Note: This is a simplified version. In production, use a proper Blake2s implementation
            val digest = MessageDigest.getInstance("SHA-256")
            return digest.digest(data).copyOfRange(0, BLAKE2S_SIZE)
        }

        private fun poly1305Mac(key: ByteArray, data: ByteArray): ByteArray {
            val mac = Mac.getInstance("HmacSHA256")
            val secretKey = SecretKeySpec(key, "HmacSHA256")
            mac.init(secretKey)
            return mac.doFinal(data).copyOfRange(0, 16)
        }
    }
}
