package com.adscoffee.sdk.inapp.yandex

import com.adscoffee.sdk.inapp.InAppAd
import com.adscoffee.sdk.inapp.InAppAdViewBinder

internal class YandexInAppAd(
    private val data: YandexAdData,
    impressions: List<String>,
    clicks: List<String>
) : InAppAd(NETWORK, impressions, clicks) {

    private var unit: YandexAdUnit? = null

    override fun render(binder: InAppAdViewBinder) {
        val format = YandexAdFormats.resolve(data)
            ?: throw IllegalStateException("Unsupported yandex ad format for block ${data.block}")

        val created = format.create(binder.adView.context, data)
        unit = created

        created.render(binder, object : YandexAdCallbacks {
            override fun onImpression() = fireImpressionTrackers()

            override fun onClick() = fireClickTrackers()

            override fun onError(message: String) = notifyError(message)
        })
    }

    override fun destroy() {
        unit?.destroy()
        unit = null
    }

    companion object {
        const val NETWORK = "yandex"
    }
}
