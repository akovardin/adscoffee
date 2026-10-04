import org.gradle.api.JavaVersion

plugins {
    kotlin("android") version "2.0.21"
    id("com.android.library") version "8.7.3"
}

group = "com.adscoffee"
version = "1.0.0"

repositories {
    google()
    mavenCentral()
}

dependencies {
    implementation("com.google.code.gson:gson:2.11.0")
    implementation("com.yandex.android:mobileads:7.14.1")
    testImplementation("org.jetbrains.kotlin:kotlin-test:2.0.21")
}

android {
    namespace = "com.adscoffee.sdk"
    compileSdk = 35
    defaultConfig {
        minSdk = 21
    }
    testOptions {
        targetSdk = 35
    }
    compileOptions {
        targetCompatibility = JavaVersion.VERSION_21
        sourceCompatibility = JavaVersion.VERSION_21
    }
}

tasks.withType<org.jetbrains.kotlin.gradle.tasks.KotlinCompile>().configureEach {
    kotlinOptions {
        jvmTarget = "21"
    }
}