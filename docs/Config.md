
# Config File Format

## Helm

```yaml
secrets:
  foo:
    bar: baz # This will set the `foo.bar` entry in the values file to `baz`
config:
  HELLO: world # This will replace all instances of `$HELLO` in the rendered template with `world`
```

## Kustomize

```yaml
secrets:
  foo:
    bar: baz # This will write `baz` to the `files/bar` file in the `foo` overlay
config:
  HELLO: world # This will replace all instances of `$HELLO` in the rendered template with `world`
```
