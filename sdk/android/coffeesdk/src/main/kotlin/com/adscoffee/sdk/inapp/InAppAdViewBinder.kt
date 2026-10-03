package com.adscoffee.sdk.inapp

import android.widget.ImageView
import android.widget.TextView

class InAppAdViewBinder private constructor(
    val adView: InAppAdView,
    val imageView: ImageView?,
    val descriptionView: TextView?,
    val informationView: TextView?
) {
    class Builder(private val adView: InAppAdView) {
        private var imageView: ImageView? = null
        private var descriptionView: TextView? = null
        private var informationView: TextView? = null

        fun setImageView(view: ImageView): Builder {
            imageView = view
            return this
        }

        fun setDescriptionView(view: TextView): Builder {
            descriptionView = view
            return this
        }

        fun setInformationView(view: TextView): Builder {
            informationView = view
            return this
        }

        fun build(): InAppAdViewBinder {
            return InAppAdViewBinder(adView, imageView, descriptionView, informationView)
        }
    }
}