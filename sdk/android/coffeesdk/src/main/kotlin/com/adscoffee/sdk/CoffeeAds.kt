package com.adscoffee.sdk

class CoffeeAds {
    companion object {
        private var initialized = false
        internal var baseUrl: String = "http://localhost:8081"

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