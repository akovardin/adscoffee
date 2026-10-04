package com.adscoffee.example

import android.app.Activity
import android.content.Intent
import android.net.Uri
import android.os.Bundle
import android.view.Gravity
import android.widget.Button
import android.widget.LinearLayout
import com.adscoffee.sdk.CoffeeAds
import com.adscoffee.sdk.inapp.*

class MainActivity : Activity() {

    private var currentAd: InAppAd? = null
    private lateinit var adView: InAppAdView

    override fun onCreate(savedInstanceState: Bundle?) {
        super.onCreate(savedInstanceState)

        adView = InAppAdView(this)

        val interstitialButton = Button(this).apply {
            text = "Показать интерстишел"
            setOnClickListener {
                startActivity(Intent(this@MainActivity, InterstitialActivity::class.java))
            }
        }

        val root = LinearLayout(this).apply {
            orientation = LinearLayout.VERTICAL
            gravity = Gravity.CENTER_HORIZONTAL
            addView(adView)
            addView(interstitialButton)
        }

        setContentView(root)

        CoffeeAds.setBaseUrl(BASE_URL)
        CoffeeAds.initialize(this) {
            loadBanner()
        }
    }

    private fun loadBanner() {
        val loader = InAppAdLoader(this)
        val request = InAppAdRequest.Builder(BANNER_PLACEMENT).build()

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

    companion object {
        const val BASE_URL = "https://platform.ads.coffee"
        const val BANNER_PLACEMENT = 4
        const val INTERSTITIAL_PLACEMENT = 5
    }
}
