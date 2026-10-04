package com.adscoffee.sdk.inapp.yandex

import android.content.Context
import com.adscoffee.sdk.inapp.InAppAdViewBinder

internal interface YandexAdFormat {
    val name: String

    fun matches(data: YandexAdData): Boolean

    fun create(context: Context, data: YandexAdData): YandexAdUnit
}

internal interface YandexAdUnit {
    fun render(binder: InAppAdViewBinder, callbacks: YandexAdCallbacks)

    fun destroy()
}

internal interface YandexAdCallbacks {
    fun onImpression()

    fun onClick()

    fun onError(message: String)

    fun onClosed()
}
