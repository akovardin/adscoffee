package com.adscoffee.sdk.internal

import java.math.BigDecimal
import java.net.HttpURLConnection
import java.net.URL

class TrackerService {

    fun fire(urls: List<String>, price: Double = 0.0) {
        val priceValue = formatPrice(price)

        urls.forEach { raw ->
            val urlStr = raw.replace(PRICE_MACRO, priceValue)

            Thread {
                try {
                    val url = URL(urlStr)
                    val conn = url.openConnection() as HttpURLConnection
                    conn.requestMethod = "GET"
                    conn.connectTimeout = 3000
                    conn.readTimeout = 3000
                    conn.responseCode
                    conn.disconnect()
                } catch (_: Exception) { }
            }.start()
        }
    }

    private fun formatPrice(price: Double): String {
        return BigDecimal.valueOf(price).stripTrailingZeros().toPlainString()
    }

    companion object {
        private const val PRICE_MACRO = "{price}"
    }
}
