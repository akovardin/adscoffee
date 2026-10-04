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
        val data = YandexAdData.parse(response.data)
        if (data.block.isBlank()) {
            onError(AdRequestError("Yandex response does not contain block id"))
            return
        }

        val ad = YandexInAppAd(data, response.impressions, response.clicks)
        ad.price = if (data.price > 0.0) data.price else response.price

        onLoaded(ad)
    }
}
