package com.adscoffee.sdk.inapp

class ImpressionData(
    val placementId: Int,
    val network: String? = null,
    val price: Double = 0.0,
    val timestamp: Long = System.currentTimeMillis()
)
