package com.adscoffee.sdk.inapp

import android.content.Context
import android.os.Handler
import android.os.Looper
import com.adscoffee.sdk.CoffeeAds
import com.adscoffee.sdk.internal.ApiClient

class InAppAdLoader(private val context: Context) {

    private val mainHandler = Handler(Looper.getMainLooper())
    private val factories = InAppAdFactoryRegistry.default()

    fun loadAd(request: InAppAdRequest, listener: InAppAdLoadListener) {
        Thread {
            try {
                val client = ApiClient(CoffeeAds.baseUrl)
                val banners = client.fetchBanners(request.placementId)
                val response = banners.firstOrNull()

                if (response == null) {
                    postFailure(
                        listener,
                        AdRequestError("No banners returned for placement ${request.placementId}")
                    )
                    return@Thread
                }

                val factory = factories.resolve(response.network)
                factory.create(
                    context = context,
                    response = response,
                    onLoaded = { ad ->
                        ad.placementId = request.placementId
                        mainHandler.post { listener.onAdLoaded(ad) }
                    },
                    onError = { error -> postFailure(listener, error) }
                )
            } catch (e: Exception) {
                postFailure(listener, AdRequestError("Failed to load ad: ${e.message}", e))
            }
        }.start()
    }

    private fun postFailure(listener: InAppAdLoadListener, error: AdRequestError) {
        mainHandler.post { listener.onAdFailedToLoad(error) }
    }
}
