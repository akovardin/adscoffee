package com.adscoffee.sdk.inapp

interface InAppAdLoadListener {
    fun onAdLoaded(ad: InAppAd)
    fun onAdFailedToLoad(error: AdRequestError)
}