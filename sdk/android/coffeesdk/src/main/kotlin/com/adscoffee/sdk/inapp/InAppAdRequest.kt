package com.adscoffee.sdk.inapp

class InAppAdRequest private constructor(
    val placementId: Int
) {
    class Builder(private val placementId: Int) {
        fun build(): InAppAdRequest {
            return InAppAdRequest(placementId)
        }
    }
}