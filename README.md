# `Charite Data Donation`

## Use Case

Users of the CovApp have the option to contribute data for the population level assessment of the pandemic.
They can provide their geolocation in form of their postal code along with their risk score.
These data points can help medical experts to assess where in the country hotspots might be emerging.
The Charité data donation services handles the forwarding of the risk level and the postal code to a Charité database.
Users can optionally consent to this data donation.
If they do, they need to provide their postal code.
The risk level is calculated based on the answers they provide.
Both data points are written into a Charité database.

## Building, Running, Testing

```bash
make build
make run
make test
```

## Ops checks

Liveness and Readiness checks are meant for Kubernetes.

- Liveness is a check that responds with HTTP code 200 if the application has started
- Readiness is a check that responds with HTTP code 200 if the application is ready to serve requests (e.g., connected to the database)

## Run in local Kubernetes

1. **Add localhost alias**

    Add `data-donation.local` to `/etc/hosts' file. For example in the kubernetes.docker section:

    ```txt
    # To allow the same kube context to work on the host and the container:
    127.0.0.1 kubernetes.docker.internal galaxy.local
    # End of section
    ```

1. **Build the image**

    ```bash
    export GITHUB_USER_TOKEN=<your-GH-API-token>
    make docker-build
    ```

1. **Deploy to local Kubernetes**

    - make sure you have the right `kubectl` context selected (by default docker-desktop)

        ```bash
        kubectl config use-context docker-desktop
        ```

    - render templates and deploy manifests

        ```bash
        make kube-deploy
        ```

    - check that the pod is running

        ```bash
        kubectl get pod
        ```

## Swagger API definition

The API specification can be found in `/swagger/api.yml`. To preview the specification:

1. add a [swagger viewer](https://marketplace.visualstudio.com/items?itemName=Arjun.swagger-viewer) to VSCode
1. open the `yml` file
1. open the preview with `SHIFT + OPTION + P`
