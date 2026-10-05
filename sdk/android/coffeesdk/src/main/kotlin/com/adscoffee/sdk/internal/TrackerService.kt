package com.adscoffee.sdk.internal

import java.math.BigDecimal
import java.net.HttpURLConnection
import java.net.URL

class TrackerService {

    fun fire(urls: List<String>, revenue: Double = 0.0) {
        val revenueValue = formatRevenue(revenue)

        urls.forEach { raw ->
            val urlStr = applyRevenue(raw, revenueValue, revenue)

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

    private fun applyRevenue(raw: String, revenueValue: String, revenue: Double): String {
        if (raw.contains(REVENUE_MACRO)) {
            return raw.replace(REVENUE_MACRO, revenueValue)
        }

        if (revenue <= 0.0) {
            return raw
        }

        // Старые версии сервера не подставляют {revenue} в URL. Трекер читает
        // доход из query-параметра revenue, поэтому добавляем его сами.
        val separator = if (raw.contains('?')) '&' else '?'
        return "$raw$separator$REVENUE_QUERY=$revenueValue"
    }

    private fun formatRevenue(revenue: Double): String {
        return BigDecimal.valueOf(revenue).stripTrailingZeros().toPlainString()
    }

    companion object {
        private const val REVENUE_MACRO = "{revenue}"
        private const val REVENUE_QUERY = "revenue"
    }
}
