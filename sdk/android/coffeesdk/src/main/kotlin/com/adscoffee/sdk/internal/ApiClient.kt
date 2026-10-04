package com.adscoffee.sdk.internal

import com.adscoffee.sdk.inapp.AdResponse
import com.google.gson.Gson
import com.google.gson.JsonArray
import com.google.gson.JsonObject
import java.net.HttpURLConnection
import java.net.URL

internal class ApiClient(private val baseUrl: String) {

    private val gson = Gson()

    fun fetchBanners(placementId: Int): List<AdResponse> {
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
            return parseResponses(body)
        } finally {
            conn.disconnect()
        }
    }

    private fun normalizedUrl(placementId: Int): String {
        val base = if (baseUrl.endsWith("/")) baseUrl else "$baseUrl/"
        return "$base/inapp/$placementId"
    }

    private fun parseResponses(json: String): List<AdResponse> {
        val arr = gson.fromJson(json, JsonArray::class.java)
        return arr.map { parseResponse(it.asJsonObject) }
    }

    private fun parseResponse(obj: JsonObject): AdResponse {
        return AdResponse(
            description = getString(obj, "description"),
            information = getString(obj, "information"),
            image = getString(obj, "image"),
            target = getString(obj, "target"),
            impressions = getStringList(obj, "impressions"),
            clicks = getStringList(obj, "clicks"),
            data = getString(obj, "data"),
            network = getString(obj, "network"),
            format = getString(obj, "format").ifBlank { FORMAT_BANNER },
            price = getDouble(obj, "price")
        )
    }

    private fun getString(obj: JsonObject, key: String): String {
        return obj.get(key)?.asString ?: ""
    }

    private fun getDouble(obj: JsonObject, key: String): Double {
        val element = obj.get(key) ?: return 0.0
        if (element.isJsonNull) return 0.0
        return try {
            element.asDouble
        } catch (_: Exception) {
            0.0
        }
    }

    private fun getStringList(obj: JsonObject, key: String): List<String> {
        val element = obj.get(key)
        if (element == null || !element.isJsonArray) return emptyList()
        return element.asJsonArray.map { it.asString ?: "" }
    }

    companion object {
        internal const val FORMAT_BANNER = "banner"
    }
}
