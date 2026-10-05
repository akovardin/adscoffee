package com.adscoffee.sdk.inapp.coffee

import android.content.Context
import com.adscoffee.sdk.inapp.AdRequestError
import com.adscoffee.sdk.inapp.AdResponse
import com.adscoffee.sdk.inapp.InAppAd
import com.adscoffee.sdk.inapp.InAppAdFactory

internal class CoffeeInAppAdFactory : InAppAdFactory {

    override val network: String = CoffeeInAppAd.NETWORK

    override fun create(
        context: Context,
        response: AdResponse,
        onLoaded: (InAppAd) -> Unit,
        onError: (AdRequestError) -> Unit
    ) {
        val ad = CoffeeInAppAd(
            data = CoffeeAdData(
                description = response.description,
                information = response.information,
                image = response.image,
                target = response.target,
                format = response.format
            ),
            impressions = response.impressions,
            clicks = response.clicks
        )
        ad.revenue = response.revenue

        onLoaded(ad)
    }
}
