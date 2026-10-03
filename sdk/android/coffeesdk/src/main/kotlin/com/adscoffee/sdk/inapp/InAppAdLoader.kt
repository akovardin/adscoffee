package com.adscoffee.sdk.inapp

import android.os.Handler
import android.os.Looper
import com.adscoffee.sdk.CoffeeAds
import com.adscoffee.sdk.internal.ApiClient

class InAppAdLoader(context: Any) {

    private val mainHandler = Handler(Looper.getMainLooper())

    fun loadAd(request: InAppAdRequest, listener: InAppAdLoadListener) {
        Thread {
            try {
                val client = ApiClient(CoffeeAds.baseUrl)
                val banners = client.fetchBanners(request.placementId)
                if (banners.isNotEmpty()) {
                    val ad = banners[0]
                    ad.setPlacementId(request.placementId)
                    mainHandler.post { listener.onAdLoaded(ad) }
                } else {
                    mainHandler.post {
                        listener.onAdFailedToLoad(
                            AdRequestError("No banners returned for placement ${request.placementId}")
                        )
                    }
                }
            } catch (e: Exception) {
                mainHandler.post {
                    listener.onAdFailedToLoad(
                        AdRequestError("Failed to load ad: ${e.message}", e)
                    )
                }
            }
        }.start()
    }
}