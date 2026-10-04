package com.adscoffee.example

import android.app.Activity
import android.content.Intent
import android.net.Uri
import android.os.Bundle
import com.adscoffee.sdk.CoffeeAds
import com.adscoffee.sdk.inapp.*

class MainActivity : Activity() {

    private var currentAd: InAppAd? = null
    private lateinit var adView: InAppAdView

    override fun onCreate(savedInstanceState: Bundle?) {
        super.onCreate(savedInstanceState)

        adView = InAppAdView(this)
        setContentView(adView)

        CoffeeAds.setBaseUrl("https://platform.ads.coffee")
        CoffeeAds.initialize(this) {
            loadAd()
        }
    }

    private fun loadAd() {
        val loader = InAppAdLoader(this)
        val request = InAppAdRequest.Builder(2).build()

        loader.loadAd(request, object : InAppAdLoadListener {
            override fun onAdLoaded(ad: InAppAd) {
                currentAd = ad

                ad.setInAppAdEventListener(object : InAppAdEventListener {
                    override fun onAdClicked() {
                        ad.target?.let { openTarget(it) }
                    }

                    override fun onImpression(data: ImpressionData?) { }

                    override fun onAdError(message: String) {
                        title = "Ad error: $message"
                    }
                })

                adView.removeAllViews()

                val binder = InAppAdViewBinder.Builder(adView).build()

                when (val result = ad.bindInAppAd(binder)) {
                    is AdBindingResult.Success -> Unit
                    is AdBindingResult.Failure -> {
                        title = "Bind failed: ${result.exception.message}"
                    }
                }
            }

            override fun onAdFailedToLoad(error: AdRequestError) {
                title = "Error: ${error.message}"
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
