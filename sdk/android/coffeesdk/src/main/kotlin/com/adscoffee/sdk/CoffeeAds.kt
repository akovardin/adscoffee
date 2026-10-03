package com.adscoffee.sdk

class CoffeeAds {
    companion object {
        private var initialized = false
        internal var baseUrl: String = "https://platform.ads.coffee"

        fun initialize(context: Any, onReady: () -> Unit) {
            if (initialized) return
            initialized = true
            onReady()
        }

        fun setBaseUrl(url: String) {
            baseUrl = url
        }
    }
}