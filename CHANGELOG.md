# Changelog

## [2.0.0](https://github.com/meysam81/jaeger-postgresql/compare/v1.7.0...v2.0.0) (2024-06-28)


### ⚠ BREAKING CHANGES

* renamed prom metric, and added disk size prom metric
* removed configuration file and moved to conf via env vars

### Features

* added a docker container to be used as an init container ([83b7b8e](https://github.com/meysam81/jaeger-postgresql/commit/83b7b8e2d5b0bf47e61b4d351df7e69c1141ebc9))
* added a timeout to initial db conn ([6aab74c](https://github.com/meysam81/jaeger-postgresql/commit/6aab74c3e19f824bfd02c7d8afd14edb915ff07d))
* added extra configuration methods ([#40](https://github.com/meysam81/jaeger-postgresql/issues/40)) ([947140b](https://github.com/meysam81/jaeger-postgresql/commit/947140b5888e8719dfd3d3e4c3af7833b15435ad))
* added helm chart ([8634167](https://github.com/meysam81/jaeger-postgresql/commit/86341674d85b199dea028a6a2ed176ecf1ba8a76))
* added indexes to make querying performant ([#17](https://github.com/meysam81/jaeger-postgresql/issues/17)) ([f92ef04](https://github.com/meysam81/jaeger-postgresql/commit/f92ef04bedb020f147ff72662082dc9e73c705af))
* added new span cleaner ([3405a96](https://github.com/meysam81/jaeger-postgresql/commit/3405a967ac61b79f0fa7fdef9301990ff1b817b3))
* added prometheus scrap annotations to service ([24d9201](https://github.com/meysam81/jaeger-postgresql/commit/24d9201677471e33565b6ee4ecee46f25757ad8f))
* added server host-port logging to init ([814a657](https://github.com/meysam81/jaeger-postgresql/commit/814a65772be8406b17229e85321cfec9889bfb7a))
* added span gauge to prometheus ([#36](https://github.com/meysam81/jaeger-postgresql/issues/36)) ([f891462](https://github.com/meysam81/jaeger-postgresql/commit/f891462b44ca2b9f284c1149c14cecda7a9c2fc9))
* enabling maximum start time ([df64c75](https://github.com/meysam81/jaeger-postgresql/commit/df64c75ec02ae1398768bb902ea38f07ab6b2e35))
* removed configuration file and moved to conf via env vars ([e12b256](https://github.com/meysam81/jaeger-postgresql/commit/e12b25613ec96332ac67ec24f3985930e651e590))
* renamed prom metric, and added disk size prom metric ([8c431e4](https://github.com/meysam81/jaeger-postgresql/commit/8c431e4554666c15363675fc7c3ffd41c905327f))
* rewrote to use pgx and support modern versions of jaeger ([74ef9ad](https://github.com/meysam81/jaeger-postgresql/commit/74ef9ad1684c0ae7128ba87f21b2da3532719a71))
* rewrote with grpc and instrumentation ([5069d0a](https://github.com/meysam81/jaeger-postgresql/commit/5069d0a5e1ef951da0de518f6fdd0bca92afe40f))


### Bug Fixes

* finished supporting integration tests and resulting compliance ([40b491a](https://github.com/meysam81/jaeger-postgresql/commit/40b491ad4daeb37987b53e2ffe3c426ead43bef7))
* fixed issue with dockerfile ([d7a9c2a](https://github.com/meysam81/jaeger-postgresql/commit/d7a9c2a2da84a186bd2a1875f95b695969c67be6))
* fixed issue with filtering not functioning for start_time ([#23](https://github.com/meysam81/jaeger-postgresql/issues/23)) ([46cdd8d](https://github.com/meysam81/jaeger-postgresql/commit/46cdd8d50a960be9c9dd0058131e91232db3eb43))
* fixing pgxpool connection timeout ([22a89b4](https://github.com/meysam81/jaeger-postgresql/commit/22a89b408cda42001587b685769b471f1062575a))
* moving to using simple protocol for short term pgbouncer fix ([fdaaed8](https://github.com/meysam81/jaeger-postgresql/commit/fdaaed880e648f4742c5a388e46eca5322405ade))
* tweaking ci ([f483a04](https://github.com/meysam81/jaeger-postgresql/commit/f483a04806745f0286e4b944daa9014d1915ba21))
* undoing simple protocol to fix jsonb ([b496cb3](https://github.com/meysam81/jaeger-postgresql/commit/b496cb30eb3f993389d02223de5d4c2d6c858b93))

## [1.7.0](https://github.com/robbert229/jaeger-postgresql/compare/v1.6.0...v1.7.0) (2024-04-04)


### Features

* added extra configuration methods ([#40](https://github.com/robbert229/jaeger-postgresql/issues/40)) ([947140b](https://github.com/robbert229/jaeger-postgresql/commit/947140b5888e8719dfd3d3e4c3af7833b15435ad))

## [1.6.0](https://github.com/robbert229/jaeger-postgresql/compare/v1.5.1...v1.6.0) (2024-03-28)


### Features

* added span gauge to prometheus ([#36](https://github.com/robbert229/jaeger-postgresql/issues/36)) ([f891462](https://github.com/robbert229/jaeger-postgresql/commit/f891462b44ca2b9f284c1149c14cecda7a9c2fc9))

## [1.5.1](https://github.com/robbert229/jaeger-postgresql/compare/v1.5.0...v1.5.1) (2024-03-18)


### Bug Fixes

* fixed issue with filtering not functioning for start_time ([#23](https://github.com/robbert229/jaeger-postgresql/issues/23)) ([46cdd8d](https://github.com/robbert229/jaeger-postgresql/commit/46cdd8d50a960be9c9dd0058131e91232db3eb43))

## [1.5.0](https://github.com/robbert229/jaeger-postgresql/compare/v1.4.0...v1.5.0) (2024-03-17)


### Features

* added indexes to make querying performant ([#17](https://github.com/robbert229/jaeger-postgresql/issues/17)) ([f92ef04](https://github.com/robbert229/jaeger-postgresql/commit/f92ef04bedb020f147ff72662082dc9e73c705af))
