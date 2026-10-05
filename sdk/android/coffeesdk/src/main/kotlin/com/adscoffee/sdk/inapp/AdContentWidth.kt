package com.adscoffee.sdk.inapp

import android.view.View
import android.view.ViewTreeObserver
import kotlin.math.roundToInt

/**
 * Ширина, доступная для баннера внутри [container], с учётом паддингов
 * (safe area). Возвращается в dp один раз — после первой раскладки.
 */
internal fun observeAdContentWidth(container: View, onWidth: (Int) -> Unit) {
    container.viewTreeObserver.addOnGlobalLayoutListener(
        object : ViewTreeObserver.OnGlobalLayoutListener {
            override fun onGlobalLayout() {
                if (container.viewTreeObserver.isAlive) {
                    container.viewTreeObserver.removeOnGlobalLayoutListener(this)
                }

                var widthPixels = container.width - container.paddingLeft - container.paddingRight
                if (widthPixels <= 0) {
                    widthPixels = container.resources.displayMetrics.widthPixels
                }

                val widthInDp =
                    (widthPixels / container.resources.displayMetrics.density).roundToInt()

                onWidth(widthInDp)
            }
        }
    )
}
