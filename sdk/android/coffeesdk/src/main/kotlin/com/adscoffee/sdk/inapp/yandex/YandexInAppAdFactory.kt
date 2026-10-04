package com.adscoffee.sdk.inapp.yandex

import android.content.Context
import com.adscoffee.sdk.inapp.AdRequestError
import com.adscoffee.sdk.inapp.AdResponse
import com.adscoffee.sdk.inapp.InAppAd
import com.adscoffee.sdk.inapp.InAppAdFactory

internal class YandexInAppAdFactory : InAppAdFactory {

    override val network: String = YandexInAppAd.NETWORK

    override fun create(
        context: Context,
        response: AdResponse,
        onLoaded: (InAppAd) -> Unit,
        onError: (AdRequestError) -> Unit
    ) {
        var data = YandexAdData.parse(response.data)
        if (data.block.isBlank()) {
            onError(AdRequestError("Yandex response does not contain block id"))
            return
        }

        // The block payload usually has no format: fall back to the
        // format declared by the platform for this placement.
        if (data.format.isNullOrBlank() && response.format.isNotBlank()) {
            data = data.copy(format = response.format)
        }

        val ad = YandexInAppAd(data, response.impressions, response.clicks)
        ad.price = if (data.price > 0.0) data.price else response.price

        onLoaded(ad)
    }
}
