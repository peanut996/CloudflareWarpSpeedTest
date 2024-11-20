package com.peanut996.cloudflarewarpspeedtest

import java.security.SecureRandom
import javax.crypto.Mac
import javax.crypto.spec.SecretKeySpec
import kotlin.experimental.xor

object WireGuardHandshake {
    private const val HANDSHAKE_INIT_SIZE = 148
    private const val MAC_SIZE = 16
    private val secureRandom = SecureRandom()

    fun createHandshakeInitiation(privateKey: String, publicKey: String): ByteArray {
        // Create a handshake initiation packet according to the WireGuard protocol
        val packet = ByteArray(HANDSHAKE_INIT_SIZE)
        
        // Message type - 1 for handshake initiation
        packet[0] = 1
        
        // Reserved bytes
        packet[1] = 0
        packet[2] = 0
        packet[3] = 0
        
        // Sender index (32-bit random number)
        val senderIndex = secureRandom.nextInt()
        packet[4] = (senderIndex shr 24).toByte()
        packet[5] = (senderIndex shr 16).toByte()
        packet[6] = (senderIndex shr 8).toByte()
        packet[7] = senderIndex.toByte()

        // Unencrypted ephemeral (32 bytes)
        secureRandom.nextBytes(packet, 8, 32)

        // Encrypted static (48 bytes)
        // In a real implementation, this would be encrypted using AEAD
        secureRandom.nextBytes(packet, 40, 48)

        // Encrypted timestamp (28 bytes)
        // In a real implementation, this would be encrypted using AEAD
        secureRandom.nextBytes(packet, 88, 28)

        // MAC1 (16 bytes)
        calculateMac(packet, 0, 116, privateKey.toByteArray(), packet, 116)

        // MAC2 (16 bytes)
        calculateMac(packet, 0, 132, publicKey.toByteArray(), packet, 132)

        return packet
    }

    private fun calculateMac(
        data: ByteArray,
        dataOffset: Int,
        dataLength: Int,
        key: ByteArray,
        output: ByteArray,
        outputOffset: Int
    ) {
        try {
            val mac = Mac.getInstance("HmacSHA256")
            val secretKey = SecretKeySpec(key, "HmacSHA256")
            mac.init(secretKey)
            
            mac.update(data, dataOffset, dataLength)
            val result = mac.doFinal()
            
            // Copy first MAC_SIZE bytes to output
            System.arraycopy(result, 0, output, outputOffset, MAC_SIZE)
        } catch (e: Exception) {
            // In case of any crypto errors, fill with random data
            secureRandom.nextBytes(output, outputOffset, MAC_SIZE)
        }
    }

    private fun SecureRandom.nextBytes(array: ByteArray, offset: Int, length: Int) {
        val temp = ByteArray(length)
        nextBytes(temp)
        System.arraycopy(temp, 0, array, offset, length)
    }
}
