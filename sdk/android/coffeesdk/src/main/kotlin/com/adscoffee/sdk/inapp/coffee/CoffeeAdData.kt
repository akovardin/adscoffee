package com.adscoffee.sdk.inapp.coffee

internal data class CoffeeAdData(
    val description: String,
    val information: String,
    val image: String,
    val target: String,
    val format: String
) {
    companion object {
        const val FORMAT_BANNER = "banner"
        const val FORMAT_INTERSTITIAL = "interstitial"

        fun isInterstitial(format: String): Boolean {
            return format.equals(FORMAT_INTERSTITIAL, ignoreCase = true)
        }
    }
}
