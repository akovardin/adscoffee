package com.adscoffee.sdk.inapp

sealed class AdBindingResult {
    object Success : AdBindingResult()
    class Failure(val exception: Exception) : AdBindingResult()
}