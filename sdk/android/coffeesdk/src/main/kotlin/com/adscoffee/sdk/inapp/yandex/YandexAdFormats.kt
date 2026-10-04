package com.adscoffee.sdk.inapp.yandex

import com.adscoffee.sdk.inapp.yandex.formats.YandexBannerFormat

internal object YandexAdFormats {

    private val formats: List<YandexAdFormat> = listOf(
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
