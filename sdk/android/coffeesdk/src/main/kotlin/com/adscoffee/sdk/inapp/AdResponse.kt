package com.adscoffee.sdk.inapp

internal data class AdResponse(
    val description: String,
    val information: String,
    val image: String,
    val target: String,
    val impressions: List<String>,
    val clicks: List<String>,
    val data: String,
    val network: String,
    val format: String,
    val revenue: Double
)
