# secret-reset

The `secret-reset` repository contains code that communicates with an external authentication service to obtain a valid token. This token is utilized by other components, such as [`cert-external-issuer`](https://github.com/dana-team/cert-external-issuer) for authentication purposes. 

The token is stored in a Kubernetes `Secret`. If the `Secret` exists then it updates its content, and if it doesn't exist then it creates it.

The code is packaged as a Helm Chart to run as a Kubernetes `CronJob`, designed to run on a schedule to ensure that the token in the `Secret` remains valid.

## Quickstart

### Environment Variables

Before running `secret-reset`, ensure the following environment variables are set:

| Key                  | Description                                     |
|----------------------|-------------------------------------------------|
| `AUTH_USERNAME`      | The authentication username.                    |
| `AUTH_CLIENT_SECRET` | The authentication client secret.               |
| `AUTH_URL`           | The URL of the authentication service.          |
| `SECRET_NAME`        | The name of the secret to create or update.     |
| `SECRET_NAMESPACE`   | The namespace where the secret will be created or updated.|

### Deploy

Ensure that the `secret-reset` `CronJob` is deployed in the same namespace as the `Secret` that needs to be created or updated.

Use the provided Helm Chart in this repository to deploy the `secret-reset` on your Kubernetes cluster:

```bash
$ helm upgrade --install <RELEASE_NAME> \
  --namespace <SECRET_NAMESPACE> \
  --create-namespace \
  oci://ghcr.io/dana-team/helm-charts/secret-reset \
  --version <release> \
  --set config.env.AUTH_USERNAME=<USERNAME> \
  --set config.env.AUTH_CLIENT_SECRET=<CLIENT_SECRET> \
  --set config.env.AUTH_URL=<URL> \
  --set config.env.SECRET_NAME=<SECRET_NAME> \
  --set config.env.SECRET_NAMESPACE=<SECRET_NAMESPACE>
```

Alternatively, you can deploy the latest non-released version using the Chart located at `charts/secret-reset` in this repository. Execute the following `Makefile` target:

```bash
$ make deploy AUTH_USERNAME=<USERNAME> AUTH_CLIENT_SECRET=<CLIENT_SECRET> AUTH_URL=<URL> SECRET_NAME=<SECRET_NAME>
```

### Undeploy

To undeploy the `secret-reset` `CronJob`, use:

```bash
$ helm uninstall <RELEASE_NAME> --namespace <SECRET_NAMESPACE>
```

Or alternatively, use:

```bash
$ make undeploy
```