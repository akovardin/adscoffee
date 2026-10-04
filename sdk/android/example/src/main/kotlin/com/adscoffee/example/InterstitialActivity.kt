package com.adscoffee.example

import android.app.Activity
import android.os.Bundle
import com.adscoffee.sdk.CoffeeAds
import com.adscoffee.sdk.inapp.AdRequestError
import com.adscoffee.sdk.inapp.InAppAd
import com.adscoffee.sdk.inapp.InAppAdLoadListener
import com.adscoffee.sdk.inapp.InAppAdLoader
import com.adscoffee.sdk.inapp.InAppAdRequest
import com.adscoffee.sdk.inapp.InAppInterstitial

class InterstitialActivity : Activity() {

    private var currentAd: InAppAd? = null
    private var interstitial: InAppInterstitial? = null

    override fun onCreate(savedInstanceState: Bundle?) {
        super.onCreate(savedInstanceState)

        CoffeeAds.setBaseUrl(MainActivity.BASE_URL)
        CoffeeAds.initialize(this) {
            loadInterstitial()
        }
    }

    private fun loadInterstitial() {
        val loader = InAppAdLoader(this)
        val request = InAppAdRequest.Builder(MainActivity.INTERSTITIAL_PLACEMENT).build()

        loader.loadAd(request, object : InAppAdLoadListener {
            override fun onAdLoaded(ad: InAppAd) {
                currentAd = ad

                interstitial = InAppInterstitial.show(
                    context = this@InterstitialActivity,
                    ad = ad
                ) {
                    interstitial = null
                    finish()
                }
            }

            override fun onAdFailedToLoad(error: AdRequestError) {
                finish()
            }
        })
    }

    override fun onDestroy() {
        interstitial?.dismiss()
        interstitial = null

        currentAd?.destroy()
        currentAd = null

        super.onDestroy()
    }
}
