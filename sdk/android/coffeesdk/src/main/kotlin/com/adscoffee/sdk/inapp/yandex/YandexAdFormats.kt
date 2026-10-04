package com.adscoffee.sdk.inapp.yandex

import com.adscoffee.sdk.inapp.yandex.formats.YandexBannerFormat
import com.adscoffee.sdk.inapp.yandex.formats.YandexInterstitialFormat

internal object YandexAdFormats {

    private val formats: List<YandexAdFormat> = listOf(
        YandexInterstitialFormat(),
        YandexBannerFormat()
    )

    fun resolve(data: YandexAdData): YandexAdFormat? {
        data.format?.takeIf { it.isNotBlank() }?.let { name ->
            return formats.firstOrNull { it.name.equals(name, ignoreCase = true) }
        }

        return formats.firstOrNull { it.matches(data) }
            ?: formats.firstOrNull()
    }
}
