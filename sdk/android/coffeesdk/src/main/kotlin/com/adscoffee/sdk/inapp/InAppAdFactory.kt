package com.adscoffee.sdk.inapp

import android.content.Context
import com.adscoffee.sdk.inapp.coffee.CoffeeInAppAdFactory
import com.adscoffee.sdk.inapp.yandex.YandexInAppAdFactory

internal interface InAppAdFactory {
    val network: String

    fun create(
        context: Context,
        response: AdResponse,
        onLoaded: (InAppAd) -> Unit,
        onError: (AdRequestError) -> Unit
    )
}

internal class InAppAdFactoryRegistry(factories: List<InAppAdFactory>) {

    private val byNetwork: Map<String, InAppAdFactory> = factories.associateBy { it.network }
    private val fallback: InAppAdFactory = factories.first()

    fun resolve(network: String): InAppAdFactory = byNetwork[network] ?: fallback

    companion object {
        fun default(): InAppAdFactoryRegistry = InAppAdFactoryRegistry(
            listOf(
                CoffeeInAppAdFactory(),
                YandexInAppAdFactory()
            )
        )
    }
}
