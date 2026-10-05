package com.adscoffee.sdk.inapp

import android.os.Handler
import android.os.Looper
import com.adscoffee.sdk.internal.TrackerService

abstract class InAppAd internal constructor(
    val network: String,
    val impressions: List<String>,
    val clicks: List<String>
) {
    private val trackerService = TrackerService()
    private var listener: InAppAdEventListener? = null

    protected val mainHandler = Handler(Looper.getMainLooper())

    internal var placementId: Int = 0
    internal var revenue: Double = 0.0

    open val target: String? = null

    /** "banner" или "interstitial"; определяет как рендерить креатив. */
    open val format: String? = null

    fun bindInAppAd(binder: InAppAdViewBinder): AdBindingResult {
        return try {
            render(binder)
            AdBindingResult.Success
        } catch (e: Exception) {
            AdBindingResult.Failure(e)
        }
    }

    protected abstract fun render(binder: InAppAdViewBinder)

    fun setInAppAdEventListener(listener: InAppAdEventListener) {
        this.listener = listener
    }

    fun fireImpressionTrackers(revenue: Double? = null) {
        val value = revenue ?: this.revenue
        listener?.onImpression(ImpressionData(placementId, network, value))
        trackerService.fire(impressions, value)
    }

    fun fireClickTrackers() {
        listener?.onAdClicked()
        trackerService.fire(clicks, revenue)
    }

    protected fun notifyError(message: String) {
        listener?.onAdError(message)
    }

    protected fun notifyClosed() {
        listener?.onAdClosed()
    }

    open fun destroy() {}
}
