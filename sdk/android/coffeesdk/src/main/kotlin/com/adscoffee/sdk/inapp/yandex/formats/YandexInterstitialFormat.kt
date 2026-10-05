package com.adscoffee.sdk.inapp.yandex.formats

import android.app.Activity
import android.content.Context
import com.adscoffee.sdk.inapp.InAppAdViewBinder
import com.adscoffee.sdk.inapp.yandex.YandexAdCallbacks
import com.adscoffee.sdk.inapp.yandex.YandexAdData
import com.adscoffee.sdk.inapp.yandex.YandexAdFormat
import com.adscoffee.sdk.inapp.yandex.YandexAdUnit
import com.adscoffee.sdk.inapp.yandex.revenue
import com.yandex.mobile.ads.common.AdError
import com.yandex.mobile.ads.common.AdRequest
import com.yandex.mobile.ads.common.AdRequestError
import com.yandex.mobile.ads.common.ImpressionData
import com.yandex.mobile.ads.interstitial.InterstitialAd
import com.yandex.mobile.ads.interstitial.InterstitialAdEventListener
import com.yandex.mobile.ads.interstitial.InterstitialAdLoadListener
import com.yandex.mobile.ads.interstitial.InterstitialAdLoader

internal class YandexInterstitialFormat : YandexAdFormat {

    override val name: String = YandexAdData.FORMAT_INTERSTITIAL

    override fun matches(data: YandexAdData): Boolean = data.isInterstitial()

    override fun create(context: Context, data: YandexAdData): YandexAdUnit {
        return YandexInterstitialUnit(data)
    }
}

private class YandexInterstitialUnit(
    private val data: YandexAdData
) : YandexAdUnit {

    private var adLoader: InterstitialAdLoader? = null
    private var interstitialAd: InterstitialAd? = null

    override fun render(binder: InAppAdViewBinder, callbacks: YandexAdCallbacks) {
        val context = binder.adView.context
        val activity = context as? Activity
        if (activity == null) {
            callbacks.onError("Interstitial requires an Activity context")
            return
        }

        val loader = InterstitialAdLoader(context)
        adLoader = loader

        android.util.Log.i("CoffeeAdsSDK", "yandex interstitial loadAd block=${data.block}")
        loader.loadAd(
            AdRequest.Builder(data.block).build(),
            object : InterstitialAdLoadListener {
                override fun onAdLoaded(ad: InterstitialAd) {
                    android.util.Log.i("CoffeeAdsSDK", "yandex interstitial loaded, showing")
                    interstitialAd = ad
                    ad.setAdEventListener(object : InterstitialAdEventListener {
                        override fun onAdShown() {}

                        override fun onAdFailedToShow(error: AdError) {
                            callbacks.onError(error.description)
                        }

                        override fun onAdDismissed() {
                            callbacks.onClosed()
                        }

                        override fun onAdClicked() {
                            callbacks.onClick()
                        }

                        override fun onAdImpression(data: ImpressionData?) {
                            callbacks.onImpression(data.revenue())
                        }
                    })

                    ad.show(activity)
                }

                override fun onAdFailedToLoad(error: AdRequestError) {
                    android.util.Log.e("CoffeeAdsSDK", "yandex interstitial failed: ${error.description}")
                    callbacks.onError(error.description)
                }
            }
        )
    }

    override fun destroy() {
        adLoader?.cancelLoading()
        adLoader = null
        interstitialAd = null
    }
}
