package com.adscoffee.sdk.internal

import java.net.HttpURLConnection
import java.net.URL

class TrackerService {

    fun fire(urls: List<String>) {
        urls.forEach { urlStr ->
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
}