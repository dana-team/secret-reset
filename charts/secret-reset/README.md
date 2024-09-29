# secret-reset

![Version: 0.0.0](https://img.shields.io/badge/Version-0.0.0-informational?style=flat-square) ![Type: application](https://img.shields.io/badge/Type-application-informational?style=flat-square) ![AppVersion: latest](https://img.shields.io/badge/AppVersion-latest-informational?style=flat-square)

A Helm chart for secret-reset.

## Values

| Key | Type | Default | Description |
|-----|------|---------|-------------|
| affinity | object | `{}` | Node affinity rules for scheduling pods. Allows you to specify advanced node selection constraints. |
| config.cron.activeDeadlineSeconds | int | `120` | Specifies the duration in seconds relative to the startTime that the job may be continuously active before the system tries to terminate |
| config.cron.concurrencyPolicy | string | `"Replace"` | Specifies how to treat concurrent executions of a Job. |
| config.cron.failedJobsHistoryLimit | int | `3` | The number of failed finished jobs to retain. Value must be non-negative integer. |
| config.cron.restartPolicy | string | `"OnFailure"` | Restart policy for all containers within the pod. |
| config.cron.schedule | string | `"0 * * * *"` | The schedule of the CronJob. |
| config.cron.successfulJobsHistoryLimit | int | `3` | The number of successful finished jobs to retain. Value must be non-negative integer. |
| config.cron.suspend | bool | `false` | Whether the Crob is active or not. |
| config.env.AUTH_CLIENT_SECRET | string | `"secretive"` | The auth client secret. |
| config.env.AUTH_URL | string | `"https://127.0.0.1:8080"` | The auth service URL. |
| config.env.AUTH_USERNAME | string | `"dana"` | The auth username. |
| config.env.INSECURE_SKIP_TLS_VERIFY | bool | `false` | Optionally disable tls verification |
| config.env.SECRET_NAME | string | `"top-secret"` | The name of the secret to create or update. |
| config.env.SECRET_NAMESPACE | string | `"top-secret-ns"` | The namespace of the secret to create or update. |
| fullnameOverride | string | `""` |  |
| image.pullPolicy | string | `"IfNotPresent"` | The pull policy for the image. |
| image.repository | string | `"ghcr.io/dana-team/reset-secret"` | The repository of the container image. |
| image.tag | string | `""` | The tag of the container image. |
| nameOverride | string | `""` |  |
| nodeSelector | object | `{}` | Node selector for scheduling pods. Allows you to specify node labels for pod assignment. |
| resources | object | `{"limits":{"cpu":"500m","memory":"128Mi"},"requests":{"cpu":"10m","memory":"64Mi"}}` | Resource requests and limits for the container. |
| securityContext | object | `{}` | Pod-level security context for the entire pod. |
| tolerations | list | `[]` | Node tolerations for scheduling pods. Allows the pods to be scheduled on nodes with matching taints. |

