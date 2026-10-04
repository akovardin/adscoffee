package com.adscoffee.sdk.inapp.yandex.formats

import android.content.Context
import com.adscoffee.sdk.inapp.InAppAdViewBinder
import com.adscoffee.sdk.inapp.yandex.YandexAdCallbacks
import com.adscoffee.sdk.inapp.yandex.YandexAdData
import com.adscoffee.sdk.inapp.yandex.YandexAdFormat
import com.adscoffee.sdk.inapp.yandex.YandexAdUnit
import com.yandex.mobile.ads.banner.BannerAdEventListener
import com.yandex.mobile.ads.banner.BannerAdSize
import com.yandex.mobile.ads.banner.BannerAdView
import com.yandex.mobile.ads.common.AdRequest
import com.yandex.mobile.ads.common.AdRequestError
import com.yandex.mobile.ads.common.ImpressionData

internal class YandexBannerFormat : YandexAdFormat {

    override val name: String = "banner"

    override fun matches(data: YandexAdData): Boolean = !data.isInterstitial()

    override fun create(context: Context, data: YandexAdData): YandexAdUnit {
        return YandexBannerUnit(context, data)
    }
}

private class YandexBannerUnit(
    private val context: Context,
    private val data: YandexAdData
) : YandexAdUnit {

    private var bannerAdView: BannerAdView? = null

    override fun render(binder: InAppAdViewBinder, callbacks: YandexAdCallbacks) {
        val view = BannerAdView(context)
        view.setAdSize(BannerAdSize.sticky(context, widthInDp()))
        view.setBannerAdEventListener(object : BannerAdEventListener {
            override fun onAdLoaded() {}

            override fun onAdFailedToLoad(error: AdRequestError) {
                callbacks.onError(error.description)
            }

            override fun onAdClicked() = callbacks.onClick()

            override fun onImpression(impressionData: ImpressionData?) = callbacks.onImpression()
        })

        bannerAdView = view
        binder.adView.addView(view)
        view.loadAd(AdRequest.Builder(data.block).build())
    }

    override fun destroy() {
        bannerAdView?.destroy()
        bannerAdView = null
    }

    private fun widthInDp(): Int {
        val metrics = context.resources.displayMetrics
        return (metrics.widthPixels / metrics.density).toInt()
    }
}
