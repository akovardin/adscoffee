package com.adscoffee.example

import android.app.Activity
import android.content.Intent
import android.net.Uri
import android.os.Bundle
import com.adscoffee.sdk.CoffeeAds
import com.adscoffee.sdk.inapp.*

class InterstitialActivity : Activity() {

    private var currentAd: InAppAd? = null
    private lateinit var adView: InAppAdView

    override fun onCreate(savedInstanceState: Bundle?) {
        super.onCreate(savedInstanceState)

        adView = InAppAdView(this)
        setContentView(adView)

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

                ad.setInAppAdEventListener(object : InAppAdEventListener {
                    override fun onAdClicked() {
                        ad.target?.let { openTarget(it) }
                    }

                    override fun onImpression(data: ImpressionData?) { }

                    override fun onAdError(message: String) {
                        finish()
                    }

                    override fun onAdClosed() {
                        finish()
                    }
                })

                val binder = InAppAdViewBinder.Builder(adView).build()
                ad.bindInAppAd(binder)
            }

            override fun onAdFailedToLoad(error: AdRequestError) {
                finish()
            }
        })
    }

    override fun onDestroy() {
        currentAd?.destroy()
        currentAd = null
        super.onDestroy()
    }

    private fun openTarget(url: String) {
        try {
            startActivity(Intent(Intent.ACTION_VIEW, Uri.parse(url)))
        } catch (_: Exception) { }
    }
}
