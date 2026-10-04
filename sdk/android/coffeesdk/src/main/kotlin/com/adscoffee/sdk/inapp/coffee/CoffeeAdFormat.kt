package com.adscoffee.sdk.inapp.coffee

import android.content.Context
import com.adscoffee.sdk.inapp.InAppAdViewBinder

internal interface CoffeeAdFormat {
    val name: String

    fun matches(data: CoffeeAdData): Boolean

    fun create(context: Context, data: CoffeeAdData): CoffeeAdUnit
}

internal interface CoffeeAdUnit {
    fun render(binder: InAppAdViewBinder, callbacks: CoffeeAdCallbacks)

    fun destroy()
}

internal interface CoffeeAdCallbacks {
    fun onImpression()

    fun onClick()

    fun onError(message: String)

    fun onClosed()
}
