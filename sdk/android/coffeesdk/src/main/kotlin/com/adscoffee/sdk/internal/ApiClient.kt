package com.adscoffee.sdk.internal

import com.adscoffee.sdk.inapp.InAppAd
import com.google.gson.Gson
import com.google.gson.JsonArray
import com.google.gson.JsonObject
import java.net.HttpURLConnection
import java.net.URL

class ApiClient(private val baseUrl: String) {

    private val gson = Gson()

    fun fetchBanners(placementId: Int): List<InAppAd> {
        val url = URL(normalizedUrl(placementId))
        val conn = url.openConnection() as HttpURLConnection
        conn.requestMethod = "GET"
        conn.connectTimeout = 10000
        conn.readTimeout = 10000

        try {
            val code = conn.responseCode
            if (code != HttpURLConnection.HTTP_OK) {
                throw RuntimeException("Server returned HTTP $code for placement $placementId")
            }
            val body = conn.inputStream.bufferedReader().readText()
            return parseBanners(body)
        } finally {
            conn.disconnect()
        }
    }

    private fun normalizedUrl(placementId: Int): String {
        val base = if (baseUrl.endsWith("/")) baseUrl else "$baseUrl/"
        return "$base/inapp/$placementId"
    }

    private fun parseBanners(json: String): List<InAppAd> {
        val arr = gson.fromJson(json, JsonArray::class.java)
        return arr.map { parseBanner(it.asJsonObject) }
    }

    private fun parseBanner(obj: JsonObject): InAppAd {
        return InAppAd(
            description = getString(obj, "description"),
            information = getString(obj, "information"),
            image = getString(obj, "image"),
            target = getString(obj, "target"),
            impressions = getStringList(obj, "impressions"),
            clicks = getStringList(obj, "clicks"),
            network = getString(obj, "network")
        )
    }

    private fun getString(obj: JsonObject, key: String): String {
        return obj.get(key)?.asString ?: ""
    }

    private fun getStringList(obj: JsonObject, key: String): List<String> {
        val element = obj.get(key)
        if (element == null || !element.isJsonArray) return emptyList()
        return element.asJsonArray.map { it.asString ?: "" }
    }
}