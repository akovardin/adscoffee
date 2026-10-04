package com.adscoffee.sdk.inapp.yandex

import com.google.gson.Gson
import com.google.gson.JsonObject

internal data class YandexAdData(
    val block: String,
    val format: String?,
    val price: Double
) {
    fun isInterstitial(): Boolean {
        // Формат определяется только полем format (задаётся платформой),
        // а не префиксом рекламного блока.
        return format.equals(FORMAT_INTERSTITIAL, ignoreCase = true)
    }

    companion object {
        const val FORMAT_INTERSTITIAL = "interstitial"

        private val gson = Gson()

        fun parse(raw: String): YandexAdData {
            if (raw.isBlank()) return YandexAdData("", null, 0.0)

            return try {
                val obj = gson.fromJson(raw, JsonObject::class.java)
                YandexAdData(
                    block = obj.get("block")?.asString.orEmpty(),
                    format = obj.get("format")?.asString,
                    price = obj.get("price")?.asDouble ?: 0.0
                )
            } catch (_: Exception) {
                YandexAdData("", null, 0.0)
            }
        }
    }
}
