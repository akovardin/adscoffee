package com.adscoffee.sdk.inapp.coffee

import com.adscoffee.sdk.inapp.InAppAd
import com.adscoffee.sdk.inapp.InAppAdViewBinder

internal class CoffeeInAppAd(
    private val data: CoffeeAdData,
    impressions: List<String>,
    clicks: List<String>
) : InAppAd(NETWORK, impressions, clicks) {

    private var unit: CoffeeAdUnit? = null

    override val target: String = data.target

    override fun render(binder: InAppAdViewBinder) {
        val format = CoffeeAdFormats.resolve(data)
            ?: throw IllegalStateException("Unsupported coffee ad format")

        val created = format.create(binder.adView.context, data)
        unit = created

        created.render(binder, object : CoffeeAdCallbacks {
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
        const val NETWORK = "coffee"
    }
}
