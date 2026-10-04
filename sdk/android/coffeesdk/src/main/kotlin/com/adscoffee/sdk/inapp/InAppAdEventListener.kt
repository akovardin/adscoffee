package com.adscoffee.sdk.inapp

interface InAppAdEventListener {
    fun onAdClicked()
    fun onImpression(data: ImpressionData?)
    fun onAdError(message: String) {}
    fun onAdClosed() {}
}
