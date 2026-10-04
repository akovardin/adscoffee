package com.adscoffee.sdk.inapp.coffee

import android.content.Context
import kotlin.math.max
import kotlin.math.roundToInt

internal class CoffeeBannerSize private constructor(
    val widthInDp: Int,
    val heightInDp: Int
) {
    fun getWidthInPixels(context: Context): Int {
        return (widthInDp * context.resources.displayMetrics.density).roundToInt()
    }

    fun getHeightInPixels(context: Context): Int {
        return (heightInDp * context.resources.displayMetrics.density).roundToInt()
    }

    companion object {
        private const val MIN_HEIGHT_IN_DP = 50
        private const val MAX_SCREEN_HEIGHT_RATIO = 0.15f

        fun stickySize(
            context: Context,
            widthInDp: Int,
            creativeWidth: Int,
            creativeHeight: Int
        ): CoffeeBannerSize {
            val metrics = context.resources.displayMetrics

            val ratio = if (creativeWidth > 0) {
                creativeHeight.toFloat() / creativeWidth.toFloat()
            } else {
                1f
            }

            val maxHeightInDp = (metrics.heightPixels / metrics.density * MAX_SCREEN_HEIGHT_RATIO).roundToInt()
            val minHeightInDp = MIN_HEIGHT_IN_DP.coerceAtMost(maxHeightInDp)
            val heightInDp = (widthInDp * ratio).roundToInt().coerceIn(minHeightInDp, max(minHeightInDp, maxHeightInDp))

            return CoffeeBannerSize(widthInDp = widthInDp, heightInDp = heightInDp)
        }
    }
}
