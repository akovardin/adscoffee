package com.adscoffee.sdk.inapp.coffee.formats

import android.app.Activity
import android.app.Dialog
import android.content.Context
import android.graphics.Bitmap
import android.graphics.BitmapFactory
import android.graphics.Color
import android.graphics.drawable.ColorDrawable
import android.os.Handler
import android.os.Looper
import android.util.TypedValue
import android.view.Gravity
import android.view.ViewGroup
import android.widget.FrameLayout
import android.widget.ImageView
import android.widget.LinearLayout
import android.widget.TextView
import com.adscoffee.sdk.inapp.InAppAdViewBinder
import com.adscoffee.sdk.inapp.coffee.CoffeeAdCallbacks
import com.adscoffee.sdk.inapp.coffee.CoffeeAdData
import com.adscoffee.sdk.inapp.coffee.CoffeeAdFormat
import com.adscoffee.sdk.inapp.coffee.CoffeeAdUnit
import java.net.URL

internal class CoffeeInterstitialFormat : CoffeeAdFormat {

    override val name: String = CoffeeAdData.FORMAT_INTERSTITIAL

    override fun matches(data: CoffeeAdData): Boolean = CoffeeAdData.isInterstitial(data.format)

    override fun create(context: Context, data: CoffeeAdData): CoffeeAdUnit {
        return CoffeeInterstitialUnit(data)
    }
}

private class CoffeeInterstitialUnit(
    private val data: CoffeeAdData
) : CoffeeAdUnit {

    private var dialog: Dialog? = null

    override fun render(binder: InAppAdViewBinder, callbacks: CoffeeAdCallbacks) {
        val context = binder.adView.context
        if (context !is Activity) {
            callbacks.onError("Interstitial requires an Activity context")
            return
        }

        val view = CoffeeInterstitialView(context, data, object : CoffeeInterstitialView.Listener {
            override fun onClose() {
                dialog?.dismiss()
                dialog = null
                callbacks.onClosed()
            }

            override fun onClick() = callbacks.onClick()

            override fun onImageError(message: String) = callbacks.onError(message)
        })

        val interstitial = Dialog(context, android.R.style.Theme_Black_NoTitleBar_Fullscreen).apply {
            setContentView(view)
            setCancelable(false)
            window?.setBackgroundDrawable(ColorDrawable(Color.WHITE))
            window?.setLayout(
                ViewGroup.LayoutParams.MATCH_PARENT,
                ViewGroup.LayoutParams.MATCH_PARENT
            )
        }

        interstitial.show()
        dialog = interstitial

        callbacks.onImpression()
    }

    override fun destroy() {
        dialog?.dismiss()
        dialog = null
    }
}

internal class CoffeeInterstitialView(
    context: Context,
    private val data: CoffeeAdData,
    private val listener: Listener
) : FrameLayout(context) {

    interface Listener {
        fun onClose()
        fun onClick()
        fun onImageError(message: String)
    }

    private val mainHandler = Handler(Looper.getMainLooper())

    private val imageView = ImageView(context)
    private val informationView = TextView(context)
    private val descriptionView = TextView(context)

    init {
        setBackgroundColor(Color.WHITE)
        setOnClickListener { listener.onClick() }

        imageView.scaleType = ImageView.ScaleType.CENTER_CROP

        informationView.setTextColor(TEXT_COLOR)
        informationView.setTextSize(TypedValue.COMPLEX_UNIT_SP, 13f)
        informationView.setPadding(0, 0, 0, dp(4))

        descriptionView.setTextColor(TEXT_COLOR)
        descriptionView.setTextSize(TypedValue.COMPLEX_UNIT_SP, 18f)

        informationView.text = data.information
        descriptionView.text = data.description

        val content = LinearLayout(context).apply {
            orientation = LinearLayout.VERTICAL
            setPadding(dp(16), dp(16), dp(16), dp(16))
            addView(informationView)
            addView(descriptionView)
        }

        val root = LinearLayout(context).apply {
            orientation = LinearLayout.VERTICAL
            addView(imageView, LinearLayout.LayoutParams(LayoutParams.MATCH_PARENT, 0, 1f))
            addView(content, LinearLayout.LayoutParams(LayoutParams.MATCH_PARENT, LayoutParams.WRAP_CONTENT))
        }

        addView(root, LayoutParams(LayoutParams.MATCH_PARENT, LayoutParams.MATCH_PARENT))

        val closeView = TextView(context).apply {
            text = "✕"
            setTextColor(TEXT_COLOR)
            setTextSize(TypedValue.COMPLEX_UNIT_SP, 24f)
            setPadding(dp(16), dp(8), dp(16), dp(8))
            setOnClickListener { listener.onClose() }
        }

        val closeParams = LayoutParams(LayoutParams.WRAP_CONTENT, LayoutParams.WRAP_CONTENT)
        closeParams.gravity = Gravity.TOP or Gravity.END
        addView(closeView, closeParams)

        loadImage(data.image)
    }

    private fun loadImage(url: String) {
        if (url.isBlank()) return

        Thread {
            try {
                val bitmap = BitmapFactory.decodeStream(URL(url).openStream())
                mainHandler.post { imageView.setImageBitmap(bitmap) }
            } catch (e: Exception) {
                mainHandler.post { listener.onImageError(e.message ?: "Failed to load interstitial image") }
            }
        }.start()
    }

    private fun dp(value: Int): Int = (value * resources.displayMetrics.density).toInt()

    companion object {
        private val TEXT_COLOR = Color.parseColor("#212121")
    }
}
