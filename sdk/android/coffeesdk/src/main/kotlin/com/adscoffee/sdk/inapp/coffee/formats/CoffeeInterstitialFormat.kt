package com.adscoffee.sdk.inapp.coffee.formats

import android.app.Activity
import android.app.Dialog
import android.content.Context
import android.graphics.BitmapFactory
import android.graphics.Color
import android.graphics.drawable.ColorDrawable
import android.os.Build
import android.os.Handler
import android.os.Looper
import android.util.TypedValue
import android.view.View
import android.view.ViewGroup
import android.view.Window
import android.view.WindowInsets
import android.view.WindowManager
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
import kotlin.math.max

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

        val interstitial = Dialog(context, android.R.style.Theme_Translucent_NoTitleBar).apply {
            setContentView(view)
            setCancelable(false)
            window?.setBackgroundDrawable(ColorDrawable(Color.TRANSPARENT))
            window?.setLayout(
                ViewGroup.LayoutParams.MATCH_PARENT,
                ViewGroup.LayoutParams.MATCH_PARENT
            )
            window?.setDimAmount(0f)
        }

        interstitial.window?.configureFullscreenTranslucent()

        interstitial.show()
        dialog = interstitial

        callbacks.onImpression()
    }

    override fun destroy() {
        dialog?.dismiss()
        dialog = null
    }
}

@Suppress("DEPRECATION")
private fun Window.configureFullscreenTranslucent() {
    if (Build.VERSION.SDK_INT < Build.VERSION_CODES.VANILLA_ICE_CREAM) {
        addFlags(WindowManager.LayoutParams.FLAG_DRAWS_SYSTEM_BAR_BACKGROUNDS)
        statusBarColor = Color.TRANSPARENT
        navigationBarColor = Color.TRANSPARENT
    }

    if (Build.VERSION.SDK_INT in Build.VERSION_CODES.Q until Build.VERSION_CODES.VANILLA_ICE_CREAM) {
        isStatusBarContrastEnforced = false
        isNavigationBarContrastEnforced = false
    }

    if (Build.VERSION.SDK_INT >= Build.VERSION_CODES.R) {
        if (Build.VERSION.SDK_INT < Build.VERSION_CODES.VANILLA_ICE_CREAM) {
            setDecorFitsSystemWindows(false)
        }
        attributes.layoutInDisplayCutoutMode =
            WindowManager.LayoutParams.LAYOUT_IN_DISPLAY_CUTOUT_MODE_ALWAYS
    } else {
        decorView.systemUiVisibility =
            View.SYSTEM_UI_FLAG_LAYOUT_STABLE or
                View.SYSTEM_UI_FLAG_LAYOUT_FULLSCREEN or
                View.SYSTEM_UI_FLAG_LAYOUT_HIDE_NAVIGATION

        if (Build.VERSION.SDK_INT >= Build.VERSION_CODES.P) {
            attributes.layoutInDisplayCutoutMode =
                WindowManager.LayoutParams.LAYOUT_IN_DISPLAY_CUTOUT_MODE_SHORT_EDGES
        }
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

    private val content = AdContentView(context)
    private val imageView = ImageView(context)
    private val textContainer = LinearLayout(context)
    private val informationView = TextView(context)
    private val descriptionView = TextView(context)

    init {
        setBackgroundColor(OVERLAY_COLOR)
        setOnClickListener { listener.onClick() }

        imageView.scaleType = ImageView.ScaleType.CENTER_CROP

        informationView.setTextColor(TEXT_COLOR)
        informationView.setTextSize(TypedValue.COMPLEX_UNIT_SP, 13f)
        informationView.setPadding(0, 0, 0, dp(4))

        descriptionView.setTextColor(TEXT_COLOR)
        descriptionView.setTextSize(TypedValue.COMPLEX_UNIT_SP, 18f)

        informationView.text = data.information
        descriptionView.text = data.description

        textContainer.orientation = LinearLayout.VERTICAL
        textContainer.setBackgroundColor(Color.WHITE)
        textContainer.setPadding(dp(16), dp(12), dp(16), dp(16))
        textContainer.addView(informationView)
        textContainer.addView(descriptionView)

        content.setBackgroundColor(Color.WHITE)
        content.imageView = imageView
        content.textView = textContainer
        content.addView(imageView)
        content.addView(textContainer)

        addView(
            content,
            LayoutParams(LayoutParams.MATCH_PARENT, LayoutParams.MATCH_PARENT)
        )

        val closeView = TextView(context).apply {
            text = "✕"
            setTextColor(TEXT_COLOR)
            setTextSize(TypedValue.COMPLEX_UNIT_SP, 24f)
            setPadding(dp(16), dp(8), dp(16), dp(8))
            setOnClickListener { listener.onClose() }
        }
        content.closeView = closeView
        content.addView(closeView)

        setOnApplyWindowInsetsListener { _, insets -> applySafeArea(insets) }
        requestApplyInsets()

        loadImage(data.image)
    }

    @Suppress("DEPRECATION")
    private fun applySafeArea(insets: WindowInsets): WindowInsets {
        val left: Int
        val top: Int
        val right: Int
        val bottom: Int

        if (Build.VERSION.SDK_INT >= Build.VERSION_CODES.R) {
            val safeArea = insets.getInsets(
                WindowInsets.Type.systemBars() or WindowInsets.Type.displayCutout()
            )
            left = safeArea.left
            top = safeArea.top
            right = safeArea.right
            bottom = safeArea.bottom
        } else {
            var l = insets.systemWindowInsetLeft
            var t = insets.systemWindowInsetTop
            var r = insets.systemWindowInsetRight
            var b = insets.systemWindowInsetBottom

            if (Build.VERSION.SDK_INT >= Build.VERSION_CODES.P) {
                insets.displayCutout?.let { cutout ->
                    l = max(l, cutout.safeInsetLeft)
                    t = max(t, cutout.safeInsetTop)
                    r = max(r, cutout.safeInsetRight)
                    b = max(b, cutout.safeInsetBottom)
                }
            }

            left = l
            top = t
            right = r
            bottom = b
        }

        setPadding(left, top, right, bottom)
        return insets
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
        private val OVERLAY_COLOR = Color.parseColor("#66000000")
    }
}

/**
 * Контент интерстишела: картинка на всю безопасную область, текст прижат к низу
 * поверх картинки. Высота текста не ограничивается родителем, поэтому при
 * большом количестве текста он растёт вверх и никогда не обрезается снизу.
 */
private class AdContentView(context: Context) : FrameLayout(context) {

    var imageView: ImageView? = null
    var textView: View? = null
    var closeView: View? = null

    override fun onMeasure(widthMeasureSpec: Int, heightMeasureSpec: Int) {
        val width = MeasureSpec.getSize(widthMeasureSpec)
        val height = MeasureSpec.getSize(heightMeasureSpec)

        imageView?.measure(
            MeasureSpec.makeMeasureSpec(width, MeasureSpec.EXACTLY),
            MeasureSpec.makeMeasureSpec(height, MeasureSpec.EXACTLY)
        )

        textView?.measure(
            MeasureSpec.makeMeasureSpec(width, MeasureSpec.EXACTLY),
            MeasureSpec.makeMeasureSpec(0, MeasureSpec.UNSPECIFIED)
        )

        closeView?.measure(
            MeasureSpec.makeMeasureSpec(width, MeasureSpec.AT_MOST),
            MeasureSpec.makeMeasureSpec(height, MeasureSpec.AT_MOST)
        )

        setMeasuredDimension(
            resolveSize(width, widthMeasureSpec),
            resolveSize(height, heightMeasureSpec)
        )
    }

    override fun onLayout(changed: Boolean, left: Int, top: Int, right: Int, bottom: Int) {
        val width = right - left
        val height = bottom - top

        imageView?.layout(0, 0, width, height)

        textView?.let { view ->
            val textHeight = view.measuredHeight
            val textTop = height - textHeight
            view.layout(0, textTop, width, textTop + textHeight)
        }

        closeView?.let { view ->
            view.layout(
                width - view.measuredWidth,
                0,
                width,
                view.measuredHeight
            )
        }
    }
}
