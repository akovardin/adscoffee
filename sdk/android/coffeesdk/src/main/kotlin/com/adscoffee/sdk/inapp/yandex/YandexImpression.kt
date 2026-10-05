package com.adscoffee.sdk.inapp.yandex

import com.google.gson.JsonObject
import com.google.gson.JsonParser
import com.yandex.mobile.ads.common.ImpressionData

/**
 * Yandex отдаёт данные о показе как JSON в [ImpressionData.getRawData]:
 *
 * ```json
 * {"currency":"RUB","revenue":"0.0031","revenueUSD":"0.0000368",
 *  "precision":"estimated","adType":"banner", ...}
 * ```
 *
 * Возвращает доход за показ (в валюте сети, обычно RUB). Если поля нет —
 * 0.0. Значение подставляется в трекер по макросу `{revenue}`.
 */
internal fun ImpressionData?.revenue(): Double {
    val raw = this?.rawData ?: return 0.0
    return try {
        val obj = JsonParser.parseString(raw).asJsonObject
        obj.doubleOrNull("revenue") ?: obj.doubleOrNull("revenueUSD") ?: 0.0
    } catch (_: Exception) {
        0.0
    }
}

private fun JsonObject.doubleOrNull(key: String): Double? {
    val element = get(key) ?: return null
    if (!element.isJsonPrimitive) return null
    return try {
        element.asDouble
    } catch (_: Exception) {
        null
    }
}
