package com.adscoffee.example

import android.app.Activity
import android.content.Intent
import android.net.Uri
import android.os.Bundle
import android.widget.ImageView
import android.widget.TextView
import android.view.View
import com.adscoffee.sdk.CoffeeAds
import com.adscoffee.sdk.inapp.*

class MainActivity : Activity() {

    private var currentAd: InAppAd? = null
    private lateinit var imageView: ImageView
    private lateinit var descriptionView: TextView
    private lateinit var infoView: TextView
    private lateinit var adView: InAppAdView

    override fun onCreate(savedInstanceState: Bundle?) {
        super.onCreate(savedInstanceState)

        CoffeeAds.initialize(this) { }
        CoffeeAds.setBaseUrl("http://10.0.2.2:8081")

        imageView = ImageView(this)
        descriptionView = TextView(this)
        infoView = TextView(this)

        adView = InAppAdView(this)
        adView.addView(imageView)
        adView.addView(descriptionView)
        adView.addView(infoView)
        setContentView(adView)

        loadAd()
    }

    private fun loadAd() {
        val loader = InAppAdLoader(this)
        val request = InAppAdRequest.Builder(1).build()

        loader.loadAd(request, object : InAppAdLoadListener {
            override fun onAdLoaded(ad: InAppAd) {
                currentAd = ad

                val binder = InAppAdViewBinder.Builder(adView)
                    .setImageView(imageView)
                    .setDescriptionView(descriptionView)
                    .setInformationView(infoView)
                    .build()

                val result = ad.bindInAppAd(binder)

                when (result) {
                    is AdBindingResult.Success -> {
                        ad.setInAppAdEventListener(object : InAppAdEventListener {
                            override fun onAdClicked() {
                                openTarget(ad.target)
                            }
                            override fun onImpression(data: ImpressionData?) { }
                        })
                    }
                    is AdBindingResult.Failure -> {
                        descriptionView.text = "Bind failed: ${result.exception.message}"
                    }
                }
            }

            override fun onAdFailedToLoad(error: AdRequestError) {
                descriptionView.text = "Error: ${error.message}"
            }
        })
    }

    private fun openTarget(url: String) {
        try {
            startActivity(Intent(Intent.ACTION_VIEW, Uri.parse(url)))
        } catch (_: Exception) { }
    }
}