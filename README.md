DesPlops -- (D)eploym(e)nt tool for kubernete(s) a(P)p(l)ications using s(op)s to configure secret(s)

# About

# Requirements

- [SOPS](https://github.com/getsops/sops/releases)
- [Helm](https://helm.sh/docs/intro/install/)
- [Kustomize](https://github.com/kubernetes-sigs/kustomize?tab=readme-ov-file#kubectl-integration) -- Requires `kubectl` client version > 1.14
- [kapp](https://carvel.dev/kapp/docs/v0.63.x/install/)

# Getting Started

To install the latest version of `DesPlops`, run the following:

```
curl -L https://github.com/jfcarter2358/DesPlops/releases/download/latest/desplops.linux.amd64 > /usr/bin/desplops && chmod +x /usr/bin/desplops
```

To use `DesPlops`, first setup your `sops.yaml` file to pull from the shared KMS keys like so:

```yaml
creation_rules:
    - kms: '<KMS key arn 1>,<KMS key arn 2>,...'
```

Then create your SOPS config file

```sh
sops config.yaml
```

Once you are in the editor for your config file, setup the structure like the following:

## Helm Deployments

```yaml
secrets:
  foo:
    bar: baz # This will set the `foo.bar` entry in the values file to `baz`
config:
  HELLO: world # This will replace all instances of `$HELLO` in the rendered template with `world`
```

## Kustomize Deployments

```yaml
secrets:
  foo:
    bar: baz # This will write `baz` to the `files/bar` file in the `foo` overlay
config:
  HELLO: world # This will replace all instances of `$HELLO` in the rendered template with `world`
```

After setting up your configuration, either create or use an existing a local helm chart/kustomize template and deploy

## Helm Deployments

Repo structure:

```
root
|--- sops.yaml
|--- config.yaml
|--- charts
     |--- foo
          |--- templates
               |--- ...
          |--- Chart.yaml
          |--- values.yaml
```

To deploy, use

```sh
desplops -mode helm -template ./charts/foo -yes
```

Alternatively, to do a dry-run use

```sh
desplops -mode helm -template ./charts/foo
```

Finally, if you want automated rollbacks on your deployment use

```sh
desplops -m helm -template ./charts/foo -yes -rollback
```

## Kustomize Deployments

Repo structure:

```
root
|--- sops.yaml
|--- config.yaml
|--- kustomize
     |--- base
          |--- ...
     |--- overlays
          |--- foo
               |--- files
                    |--- ...
               |--- kustomize.yaml
               |--- ...
```

To deploy, use

```sh
desplops -mode kustomize -template ./kustomize/overlays/foo -yes
```

Alternatively, to do a dry-run use

```sh
desplops -mode kustomize -template ./kustomize/overlays/foo
```

Finally, if you want automated rollbacks on your deployment use

```sh
desplops -m kustomize -template ./kustomize/overlays/foo -yes -rollback
```

# Usage

```
Usage of desplops:
  -backup string
        Path to output backup manifest on deploy. Defaults to './backup-manifest.yaml' (default "./backup-manifest.yaml")
  -config string
        Path to sops config file. Defaults to 'config.yaml' (default "config.yaml")
  -loglevel string
        Log level to use. Valid values are 'NONE', 'FATAL', 'SUCCESS', 'ERROR', 'WARN', 'INFO', 'DEBUG', and 'TRACE'. Defaults to 'WARN' (default "WARN")
  -mode string
        Path to write values file for Helm deployment. Valid values are 'kustomize' and 'helm'. Required
  -output string
        Path to output rendered manifest. Defaults to './rendered-manifest.yaml' (default "./rendered-manifest.yaml")
  -rollback
        Should rollback on deploy failure. Defaults to false
  -template string
        Path to either the Kustomize overlay (e.g. ${PWD}/kustomize/overlays/foo) or the Helm chart directory containing the Chart.yaml file (e.g. ${PWD}/helm/foo). Required
  -values string
        Path to write values file for Helm deployment. Defaults to './values.yaml' (default "./values.yaml")
  -yes
        Should a non-dry run deploy be performed. Defaults to false
```

# Contact

This software is written and maintained by John Carter, if you have any questions or issues please don't hesitate to reach out at jfcarter2358@gmail.com or submit issues at [https://github.com/jfcarter2358/DesPlops/issues](https://github.com/jfcarter2358/DesPlops/issues)
