[![License](https://img.shields.io/github/license/jaimeph/helm-lint-kcl.svg)](https://github.com/jaimeph/helm-lint-kcl/blob/main/LICENSE)
[![Current Release](https://img.shields.io/github/release/jaimeph/helm-lint-kcl.svg?logo=github)](https://github.com/jaimeph/helm-lint-kcl/releases/latest)
[![GitHub Repo stars](https://img.shields.io/github/stars/jaimeph/helm-lint-kcl?style=flat&logo=github)](https://github.com/jaimeph/helm-lint-kcl/stargazers)
[![GitHub issues](https://img.shields.io/github/issues/jaimeph/helm-lint-kcl.svg)](https://github.com/jaimeph/helm-lint-kcl/issues)
[![GitHub pull requests](https://img.shields.io/github/issues-pr/jaimeph/helm-lint-kcl.svg)](https://github.com/jaimeph/helm-lint-kcl/pulls)
[![codecov](https://codecov.io/gh/jaimeph/helm-lint-kcl/branch/main/graph/badge.svg?token=4qAukyB2yX)](https://codecov.io/gh/jaimeph/helm-lint-kcl)
[![Artifact Hub](https://img.shields.io/endpoint?url=https://artifacthub.io/badge/repository/lint-kcl)](https://artifacthub.io/packages/search?repo=lint-kcl)
# Helm Lint KCL Plugin

A Helm plugin that uses KCL schemas to validate values. Instead of using JSON Schema in `values.schema.json`, you can write schema KCL validation in `schemas.k`.

## Installation

```bash
helm plugin install https://github.com/jaimeph/helm-lint-kcl
```

## Usage

### Validation

Validate your chart values using the validate command:
```bash
helm lint-kcl [NAME] [CHART] [flags]
```

Flags:
```bash
Flags:
      --debug            enable debug
  -h, --help             help for helm
      --set strings      set values on the command line (can specify multiple or separate values with commas: key1=val1,key2=val2)
      --show-values      show values
  -f, --values strings   specify values in a YAML file or a URL (can specify multiple)
  -v, --version string   chart version
```

## Contributing

1. Fork the repository
2. Create your feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add some amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## License

Distributed under the MIT License. See `LICENSE` for more information.
