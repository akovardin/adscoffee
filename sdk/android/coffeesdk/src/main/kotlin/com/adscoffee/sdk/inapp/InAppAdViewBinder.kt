package com.adscoffee.sdk.inapp

class InAppAdViewBinder private constructor(
    val adView: InAppAdView
) {
    class Builder(private val adView: InAppAdView) {
        fun build(): InAppAdViewBinder {
            return InAppAdViewBinder(adView)
        }
    }
}
