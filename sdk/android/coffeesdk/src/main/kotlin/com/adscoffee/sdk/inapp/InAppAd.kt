package com.adscoffee.sdk.inapp

import android.graphics.BitmapFactory
import android.os.Handler
import android.os.Looper
import android.view.View
import android.widget.ImageView
import android.widget.TextView
import com.adscoffee.sdk.CoffeeAds
import com.adscoffee.sdk.internal.ApiClient
import com.adscoffee.sdk.internal.TrackerService
import java.net.URL

class InAppAd internal constructor(
    val description: String,
    val information: String,
    val image: String,
    val target: String,
    val impressions: List<String>,
    val clicks: List<String>,
    val network: String
) {
    private var listener: InAppAdEventListener? = null
    private val trackerService = TrackerService()
    private val mainHandler = Handler(Looper.getMainLooper())
    private var placementId: Int = 0

    internal fun setPlacementId(id: Int) {
        placementId = id
    }

    fun bindInAppAd(binder: InAppAdViewBinder): AdBindingResult {
        return try {
            binder.descriptionView?.let { it.text = description }
            binder.informationView?.let { it.text = information }

            binder.imageView?.let { imageView ->
                loadImageAsync(imageView)
            }

            binder.adView.setOnClickListener(View.OnClickListener {
                fireClickTrackers()
            })

            fireImpressionTrackers()
            AdBindingResult.Success
        } catch (e: Exception) {
            AdBindingResult.Failure(e)
        }
    }

    fun setInAppAdEventListener(listener: InAppAdEventListener) {
        this.listener = listener
    }

    fun fireImpressionTrackers() {
        listener?.onImpression(ImpressionData(placementId))
        trackerService.fire(impressions)
    }

    fun fireClickTrackers() {
        listener?.onAdClicked()
        trackerService.fire(clicks)
    }

    private fun loadImageAsync(imageView: ImageView) {
        Thread {
            try {
                val bitmap = BitmapFactory.decodeStream(URL(image).openStream())
                mainHandler.post { imageView.setImageBitmap(bitmap) }
            } catch (_: Exception) { }
        }.start()
    }
}