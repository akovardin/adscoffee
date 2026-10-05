package com.adscoffee.sdk.inapp.coffee.formats

import android.content.Context
import android.graphics.Bitmap
import android.graphics.BitmapFactory
import android.graphics.Color
import android.os.Handler
import android.os.Looper
import android.util.TypedValue
import android.view.Gravity
import android.widget.ImageView
import android.widget.LinearLayout
import android.widget.TextView
import com.adscoffee.sdk.inapp.InAppAdViewBinder
import com.adscoffee.sdk.inapp.coffee.CoffeeAdCallbacks
import com.adscoffee.sdk.inapp.coffee.CoffeeAdData
import com.adscoffee.sdk.inapp.coffee.CoffeeAdFormat
import com.adscoffee.sdk.inapp.coffee.CoffeeAdUnit
import com.adscoffee.sdk.inapp.coffee.CoffeeBannerSize
import com.adscoffee.sdk.inapp.observeAdContentWidth
import java.net.URL

internal class CoffeeBannerFormat : CoffeeAdFormat {

    override val name: String = "banner"

    override fun matches(data: CoffeeAdData): Boolean = !CoffeeAdData.isInterstitial(data.format)

    override fun create(context: Context, data: CoffeeAdData): CoffeeAdUnit {
        return CoffeeBannerUnit(data)
    }
}

private class CoffeeBannerUnit(
    private val data: CoffeeAdData
) : CoffeeAdUnit {

    private var bannerView: CoffeeBannerView? = null

    override fun render(binder: InAppAdViewBinder, callbacks: CoffeeAdCallbacks) {
        val container = binder.adView
        val view = CoffeeBannerView(container.context)
        view.setOnClickListener { callbacks.onClick() }
        bannerView = view

        container.addView(view)

        var widthInDp = 0
        var bitmap: Bitmap? = null

        fun applySize() {
            val loaded = bitmap ?: return
            if (widthInDp <= 0) return

            val size = CoffeeBannerSize.stickySize(view.context, widthInDp, loaded.width, loaded.height)
            view.applySize(size, loaded)
        }

        observeAdContentWidth(container) { dp ->
            widthInDp = dp
            applySize()
        }

        view.setContent(data, callbacks) { loaded ->
            bitmap = loaded
            applySize()
        }

        callbacks.onImpression()
    }

    override fun destroy() {
        bannerView = null
    }
}

internal class CoffeeBannerView(context: Context) : LinearLayout(context) {

    private val mainHandler = Handler(Looper.getMainLooper())

    private val contentView = LinearLayout(context)
    private val imageView = ImageView(context)
    private val descriptionView = TextView(context)
    private val informationView = TextView(context)

    init {
        setBackgroundColor(Color.WHITE)
        setPadding(dp(12), dp(12), dp(12), dp(12))

        informationView.setTextColor(TEXT_COLOR)
        informationView.setTextSize(TypedValue.COMPLEX_UNIT_SP, 12f)
        informationView.setPadding(0, 0, 0, dp(4))

        descriptionView.setTextColor(TEXT_COLOR)
        descriptionView.setTextSize(TypedValue.COMPLEX_UNIT_SP, 16f)

        contentView.orientation = VERTICAL
        contentView.addView(informationView)
        contentView.addView(descriptionView)
    }

    fun setContent(data: CoffeeAdData, callbacks: CoffeeAdCallbacks, onImageLoaded: (Bitmap) -> Unit) {
        descriptionView.text = data.description
        informationView.text = data.information

        loadImageAsync(data.image, callbacks, onImageLoaded)
    }

    fun applySize(size: CoffeeBannerSize, bitmap: Bitmap) {
        layoutParams = LayoutParams(size.getWidthInPixels(context), size.getHeightInPixels(context))

        imageView.setImageBitmap(bitmap)
        arrange(size.getWidthInPixels(context) > size.getHeightInPixels(context))

        requestLayout()
    }

    private fun arrange(landscape: Boolean) {
        removeAllViews()

        imageView.scaleType = ImageView.ScaleType.FIT_CENTER
        contentView.gravity = Gravity.CENTER_VERTICAL

        if (landscape) {
            orientation = HORIZONTAL

            val imageParams = LayoutParams(0, LayoutParams.MATCH_PARENT, 2f)
            imageView.layoutParams = imageParams
            addView(imageView)

            val contentParams = LayoutParams(0, LayoutParams.MATCH_PARENT, 3f)
            addView(contentView, contentParams)
        } else {
            orientation = VERTICAL

            imageView.layoutParams = LayoutParams(LayoutParams.MATCH_PARENT, 0, 1f)
            addView(imageView)

            contentView.layoutParams =
                LayoutParams(LayoutParams.MATCH_PARENT, LayoutParams.WRAP_CONTENT)
            addView(contentView)
        }
    }

    private fun loadImageAsync(
        url: String,
        callbacks: CoffeeAdCallbacks,
        onImageLoaded: (Bitmap) -> Unit
    ) {
        Thread {
            try {
                val bitmap = BitmapFactory.decodeStream(URL(url).openStream())
                mainHandler.post { onImageLoaded(bitmap) }
            } catch (e: Exception) {
                mainHandler.post { callbacks.onError(e.message ?: "Failed to load banner image") }
            }
        }.start()
    }

    private fun dp(value: Int): Int = (value * resources.displayMetrics.density).toInt()

    companion object {
        private val TEXT_COLOR = Color.parseColor("#212121")
    }
}
