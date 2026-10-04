import org.gradle.api.JavaVersion

plugins {
    kotlin("android") version "2.0.21"
    id("com.android.application") version "8.7.3"
}

group = "com.adscoffee.example"
version = "1.0.0"

repositories {
    google()
    mavenCentral()
}

dependencies {
    implementation(project(":coffeesdk"))
}

android {
    namespace = "com.adscoffee.example"
    compileSdk = 35
    defaultConfig {
        applicationId = "com.adscoffee.example"
        minSdk = 21
        targetSdk = 35
    }
    lint {
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