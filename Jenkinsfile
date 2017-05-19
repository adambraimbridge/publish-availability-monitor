@Library('k8s-pipeline-lib@first-version') _

import com.ft.up.BuildConfig

BuildConfig config = new BuildConfig()
config.appDockerImageId = "coco/publish-availability-monitor"
config.useInternalDockerReg = false

entryPointForReleaseAndDev(config)
