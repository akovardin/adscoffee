package com.adscoffee.sdk

import android.content.Context
import com.yandex.mobile.ads.common.InitializationListener
import com.yandex.mobile.ads.common.YandexAds

class CoffeeAds {
    companion object {
        private var initialized = false
        internal var baseUrl: String = "https://platform.ads.coffee"

        fun initialize(context: Context, onReady: () -> Unit) {
            if (initialized) {
                onReady()
                return
            }

            initialized = true

            YandexAds.initialize(context.applicationContext, object : InitializationListener {
                override fun onInitializationCompleted() {
                    onReady()
                }
            })
        }

        fun setBaseUrl(url: String) {
            baseUrl = url
        }
    }
}
