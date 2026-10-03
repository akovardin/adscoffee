package com.adscoffee.sdk.inapp

class AdRequestError(message: String, cause: Exception? = null) :
    RuntimeException(message, cause)