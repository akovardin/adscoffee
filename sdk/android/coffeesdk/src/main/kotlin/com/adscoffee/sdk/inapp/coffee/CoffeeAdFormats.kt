package com.adscoffee.sdk.inapp.coffee

import com.adscoffee.sdk.inapp.coffee.formats.CoffeeBannerFormat
import com.adscoffee.sdk.inapp.coffee.formats.CoffeeInterstitialFormat

internal object CoffeeAdFormats {

    private val formats: List<CoffeeAdFormat> = listOf(
        CoffeeInterstitialFormat(),
        CoffeeBannerFormat()
    )

    fun resolve(data: CoffeeAdData): CoffeeAdFormat? {
        return formats.firstOrNull { it.matches(data) }
            ?: formats.firstOrNull()
    }
}
