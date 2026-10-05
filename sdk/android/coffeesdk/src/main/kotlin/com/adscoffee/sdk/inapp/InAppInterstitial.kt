package com.adscoffee.sdk.inapp

import android.app.Activity
import android.app.Dialog
import android.content.Context
import android.content.ContextWrapper
import android.content.Intent
import android.graphics.Color
import android.net.Uri
import android.util.TypedValue
import android.view.Gravity
import android.view.ViewGroup.LayoutParams.MATCH_PARENT
import android.view.ViewGroup.LayoutParams.WRAP_CONTENT
import android.widget.FrameLayout
import android.widget.TextView

/**
 * Полноэкранный показ in-app рекламы в качестве интерстишела.
 *
 * Если креатив имеет формат `interstitial`, отрисовку выполняет сам SDK
 * (собственный диалог для coffee, активность для Yandex). Баннерный креатив
 * показывается в полноэкранном контейнере, чтобы его можно было использовать
 * как интерстишел.
 */
class InAppInterstitial private constructor(
    private val dialog: Dialog?
) {
    /** Закрывает полноэкранный контейнер (если он есть). */
    fun dismiss() {
        if (dialog?.isShowing == true) {
            dialog.dismiss()
        }
    }

    companion object {
        private const val FORMAT_INTERSTITIAL = "interstitial"

        /**
         * Показывает [ad] во весь экран.
         *
         * @param context контекст (Activity или обёртка над ним).
         * @param ad загруженная реклама.
         * @param onClosed вызовется один раз, когда реклама закрыта или
         *   не удалось её показать.
         */
        fun show(context: Context, ad: InAppAd, onClosed: () -> Unit): InAppInterstitial? {
            val activity = unwrapActivity(context)
            if (activity == null) {
                onClosed()
                return null
            }

            var done = false
            var dialog: Dialog? = null

            fun close() {
                if (done) return
                done = true
                if (dialog?.isShowing == true) {
                    dialog?.dismiss()
                }
                onClosed()
            }

            ad.setInAppAdEventListener(object : InAppAdEventListener {
                override fun onImpression(data: ImpressionData?) {}
                override fun onAdClicked() {
                    ad.target?.let { openTarget(activity, it) }
                }
                override fun onAdClosed() = close()
                override fun onAdError(message: String) = close()
            })

            val adView = InAppAdView(activity)
            val binder = InAppAdViewBinder.Builder(adView).build()

            // Формат interstitial — SDK рисует сам (диалог/активность).
            if (ad.format.equals(FORMAT_INTERSTITIAL, ignoreCase = true)) {
                return if (ad.bindInAppAd(binder) is AdBindingResult.Success) {
                    ad.fireImpressionTrackers()
                    InAppInterstitial(null)
                } else {
                    close()
                    null
                }
            }

            // Баннерный креатив показываем полноэкранно.
            val container = FrameLayout(activity).apply { setBackgroundColor(Color.WHITE) }
            container.addView(adView, FrameLayout.LayoutParams(MATCH_PARENT, MATCH_PARENT))

            val closeButton = TextView(activity).apply {
                text = "✕"
                setTextColor(Color.DKGRAY)
                setTextSize(TypedValue.COMPLEX_UNIT_SP, 24f)
                setPadding(40, 40, 40, 40)
                setOnClickListener { close() }
            }
            container.addView(
                closeButton,
                FrameLayout.LayoutParams(WRAP_CONTENT, WRAP_CONTENT).apply {
                    gravity = Gravity.TOP or Gravity.END
                }
            )

            dialog = Dialog(activity, android.R.style.Theme_Black_NoTitleBar_Fullscreen).apply {
                setContentView(container)
                setCancelable(false)
                setOnDismissListener {
                    if (!done) {
                        done = true
                        onClosed()
                    }
                }
                show()
            }

            return if (ad.bindInAppAd(binder) is AdBindingResult.Success) {
                ad.fireImpressionTrackers()
                InAppInterstitial(dialog)
            } else {
                close()
                null
            }
        }

        private fun openTarget(context: Context, url: String) {
            try {
                context.startActivity(Intent(Intent.ACTION_VIEW, Uri.parse(url)))
            } catch (_: Exception) {
            }
        }

        private fun unwrapActivity(context: Context): Activity? {
            var c: Context? = context
            while (c is ContextWrapper) {
                if (c is Activity) return c
                c = c.baseContext
            }
            return null
        }
    }
}
