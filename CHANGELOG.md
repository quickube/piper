# Changelog

## [1.7.0](https://github.com/quickube/piper/compare/v1.6.0...v1.7.0) (2024-12-24)


### Features

* Gitlab provider ([#27](https://github.com/quickube/piper/issues/27)) ([8974887](https://github.com/quickube/piper/commit/8974887966ef429ab7d5ecb92db97bd6060e4d46))

## [1.6.0](https://github.com/quickube/piper/compare/v1.5.1...v1.6.0) (2024-04-06)


### Features

* bitbucket pullrequest:approved event support ([#19](https://github.com/quickube/piper/issues/19)) ([57d39cc](https://github.com/quickube/piper/commit/57d39ccce6d8baadc4c23a56f253ef55ccc5e61e))

## [1.5.1](https://github.com/quickube/piper/compare/v1.5.0...v1.5.1) (2024-04-05)


### Bug Fixes

* change rookout to quickube ([f45a0c9](https://github.com/quickube/piper/commit/f45a0c9d0ce2384e0423e8be66a33e92c1359f18))
* change rookout to quickube ([a77454b](https://github.com/quickube/piper/commit/a77454b28d2b71fbd614bf30fc29d6d241a9f475))
* change rookout to quickube ([53265e1](https://github.com/quickube/piper/commit/53265e112ecaffb1db16b82c8b749fe835d402c8))
* change rookout to quickube ([4cee863](https://github.com/quickube/piper/commit/4cee863c11c0ea3d80ca21bd46963bbd9b63eb58))
* change rookout to quickube ([fd3bc16](https://github.com/quickube/piper/commit/fd3bc160f219ac3598ed7c13552e1add5e9589f2))
* change rookout to quickube ([3ab4b62](https://github.com/quickube/piper/commit/3ab4b628ee62cfd60083e19d545189d7e74bc54a))
* change rookout to quickube ([d990336](https://github.com/quickube/piper/commit/d990336a4d4703ee41b6424dc32535e62ebba8cb))
* trigger ([9d62140](https://github.com/quickube/piper/commit/9d6214031c0da96bb3a8132c58e09785fa427335))
* trigger ([564bdc8](https://github.com/quickube/piper/commit/564bdc85b5e78687a8465fb6b6e661bb379622fb))
* trigger ([53ddeb0](https://github.com/quickube/piper/commit/53ddeb0abc280791c15fa341cd5fa1a7735b9cb1))

## [1.5.0](https://github.com/quickube/piper/compare/v1.4.0...v1.5.0) (2023-09-06)


### Features

* bitbucket cloud support - RK-19712  ([#120](https://github.com/quickube/piper/issues/120)) ([df662d0](https://github.com/quickube/piper/commit/df662d0fcb7ac533c0c31158261d5d265fa2a2c0))
* org belonging enforcment - RK-19484  ([#110](https://github.com/quickube/piper/issues/110)) ([ff9e957](https://github.com/quickube/piper/commit/ff9e957d92e6c5083dd9041ed3ddd3feb979eca6))


### Bug Fixes

* trim github notifiction massge ([#113](https://github.com/quickube/piper/issues/113)) ([75464e3](https://github.com/quickube/piper/commit/75464e3de9ffa90a417bb9e3d16527e24c9d4693))
* watcher sudden closing ([#114](https://github.com/quickube/piper/issues/114)) ([4cc3f0f](https://github.com/quickube/piper/commit/4cc3f0f032c0b42d8a392b3b636a54847112859d))

## [1.4.0](https://github.com/quickube/piper/compare/v1.3.0...v1.4.0) (2023-07-20)


### Features

* add tag handler - RK-19248  ([#92](https://github.com/quickube/piper/issues/92)) ([925dbc5](https://github.com/quickube/piper/commit/925dbc5ebb8ab59f1d5a821235dade36f6f47a33))
* github release event support - RK-19329  ([#101](https://github.com/quickube/piper/issues/101)) ([0a797e8](https://github.com/quickube/piper/commit/0a797e8ffbe417a416f5115c73f5172d74f7123c))
* webhook creation cleanup and health checks- RK-19271  ([#98](https://github.com/quickube/piper/issues/98)) ([cb99cc1](https://github.com/quickube/piper/commit/cb99cc108a260ca7583a907d6abc2d6ef2b30ca8))


### Bug Fixes

* add labels as comma seprated list in global parameters - RK-19217  ([#89](https://github.com/quickube/piper/issues/89)) ([88e1126](https://github.com/quickube/piper/commit/88e11267b8e6878b99ce9261b29967100c219657))
* event watcher failure will close the app - RK-19341  ([#104](https://github.com/quickube/piper/issues/104)) ([e115401](https://github.com/quickube/piper/commit/e1154014c5b85b82db576c07c0faa0c1e3126814))
* pull request url to html url ([#97](https://github.com/quickube/piper/issues/97)) ([4c5f6fd](https://github.com/quickube/piper/commit/4c5f6fd23dc37fe23a781da59c3723b79ac2dc88))
* retry and resumbit notifiction fix ([#96](https://github.com/quickube/piper/issues/96)) ([4d6c815](https://github.com/quickube/piper/commit/4d6c8152cc5456d6c55b3f1843029259cd23c7fb))
* seprate webhook payload from specific git provider impl ([#100](https://github.com/quickube/piper/issues/100)) ([98aa885](https://github.com/quickube/piper/commit/98aa885a16b59dbbb727066d8233605faeb16719))
* server managing improvements - RK-19210  ([#88](https://github.com/quickube/piper/issues/88)) ([af763d4](https://github.com/quickube/piper/commit/af763d4e1411e29f96c0d22a90099e87ea5e526c))

## [1.3.0](https://github.com/quickube/piper/compare/v1.2.0...v1.3.0) (2023-07-08)


### Features

* add gracefull shutdown - RK-19184  ([#84](https://github.com/quickube/piper/issues/84)) ([dcd514c](https://github.com/quickube/piper/commit/dcd514c1a80ccf49dbd253f075dc8a225c5c8c35))
* workflow status listener - RK-19190  ([#85](https://github.com/quickube/piper/issues/85)) ([a8a470e](https://github.com/quickube/piper/commit/a8a470e31b6a5b06ef6d75d8397ace96219b5ee6))


### Bug Fixes

* global parameters added and reorganized - RK-19161  ([#83](https://github.com/quickube/piper/issues/83)) ([61c8fa0](https://github.com/quickube/piper/commit/61c8fa053adb99b1423ceb4eef3b0dd465b1c5cf))

## [1.2.0](https://github.com/quickube/piper/compare/v1.1.3...v1.2.0) (2023-07-02)


### Features

* change commit linter - RK-19094 ([#81](https://github.com/quickube/piper/issues/81)) ([4586284](https://github.com/quickube/piper/commit/45862843efa124be00263619492ba97b0c21e4e5))


### Bug Fixes

* add snyk scan - RK-19072 ([#61](https://github.com/quickube/piper/issues/61)) ([cdf6b90](https://github.com/quickube/piper/commit/cdf6b909471bc66d70d85c329ce96d302e9a12fc))
* e2e init - RK-19093 ([#60](https://github.com/quickube/piper/issues/60)) ([7821785](https://github.com/quickube/piper/commit/7821785d05c7dafce52e2946b63f9f1406996796))
* event selection more specific - RK-19082 ([#76](https://github.com/quickube/piper/issues/76)) ([db58bb3](https://github.com/quickube/piper/commit/db58bb35123df023a7a0be813231968a48e5d4a6))

## [1.1.3](https://github.com/quickube/piper/compare/v1.1.2...v1.1.3) (2023-06-27)


### Bug Fixes

* generate name conversion and labels added - RK-19070 ([#74](https://github.com/quickube/piper/issues/74)) ([ca3fad0](https://github.com/quickube/piper/commit/ca3fad063c0bcf8ea5a7f03a673180eac56f0c8f))

## [1.1.2](https://github.com/quickube/piper/compare/v1.1.1...v1.1.2) (2023-06-27)


### Bug Fixes

* config selection improved logic ([#54](https://github.com/quickube/piper/issues/54)) ([59f25dc](https://github.com/quickube/piper/commit/59f25dc23e934fd8f935b59a734f96c5b4691ff4))
* git tests added - RK-19030 ([#50](https://github.com/quickube/piper/issues/50)) ([7d011c7](https://github.com/quickube/piper/commit/7d011c7b982f686b2c27bb74c8d1eb8d6997af26))
* json yaml conflic umarshling template ref to DAG task - RK-19034 ([#56](https://github.com/quickube/piper/issues/56)) ([68ab378](https://github.com/quickube/piper/commit/68ab3787fc30380dca82cd67138e54e4421889d5))
* makefile improvence - RK-19036 ([#59](https://github.com/quickube/piper/issues/59)) ([017d653](https://github.com/quickube/piper/commit/017d65312473a99cf1e5bd880d21765009de4ac5))
* remove argo configuration requierments ([#55](https://github.com/quickube/piper/issues/55)) ([3a8ed9c](https://github.com/quickube/piper/commit/3a8ed9c2abf2a2f711da6642f1617ef89d1d3119))
* stop healthz logs spam ([#57](https://github.com/quickube/piper/issues/57)) ([3f44edc](https://github.com/quickube/piper/commit/3f44edc27e12cb52adde8137f5790f0d96277fd7))

## [1.1.1](https://github.com/quickube/piper/compare/v1.1.0...v1.1.1) (2023-06-23)


### Bug Fixes

* blank config map support ([#52](https://github.com/quickube/piper/issues/52)) ([fc4a4b2](https://github.com/quickube/piper/commit/fc4a4b28617afdd618125d1fbcd796b4967b2650))
* config injection new config validation ([#49](https://github.com/quickube/piper/issues/49)) ([c84ea9a](https://github.com/quickube/piper/commit/c84ea9ac8542810b156bf2ae62fafa83a6fc1de0))
* refactor names ([#51](https://github.com/quickube/piper/issues/51)) ([36e066f](https://github.com/quickube/piper/commit/36e066f8f738547af825b7833d97f9f4c4750f08))

## [1.1.0](https://github.com/quickube/piper/compare/v1.0.14...v1.1.0) (2023-06-19)


### Features

* add app version to helm chart ([#47](https://github.com/quickube/piper/issues/47)) ([12e592d](https://github.com/quickube/piper/commit/12e592d2222ed74fc7fee6d23ba20271491102ac))

## [1.0.14](https://github.com/quickube/piper/compare/v1.0.13...v1.0.14) (2023-06-19)


### Bug Fixes

* change chart deployment flow ([#45](https://github.com/quickube/piper/issues/45)) ([bcb329c](https://github.com/quickube/piper/commit/bcb329c6a2b77f6d40e37ee0827b0f40a260e749))

## [1.0.13](https://github.com/quickube/piper/compare/v1.0.12...v1.0.13) (2023-06-19)


### Bug Fixes

* change chart deployment flow ([#43](https://github.com/quickube/piper/issues/43)) ([4585e0f](https://github.com/quickube/piper/commit/4585e0f81fb79ee0a4dddde19e09e108b64bdfa3))

## [1.0.12](https://github.com/quickube/piper/compare/v1.0.11...v1.0.12) (2023-06-19)


### Bug Fixes

* move chart path ([#41](https://github.com/quickube/piper/issues/41)) ([cbdd58a](https://github.com/quickube/piper/commit/cbdd58a026d747d2ed64e32d76da1f8fa7cbf399))

## [1.0.11](https://github.com/quickube/piper/compare/v1.0.10...v1.0.11) (2023-06-19)


### Bug Fixes

* move chart path ([#39](https://github.com/quickube/piper/issues/39)) ([50af217](https://github.com/quickube/piper/commit/50af2178642f36ffff10087a1be69021d0497d2c))

## [1.0.10](https://github.com/quickube/piper/compare/v1.0.9...v1.0.10) (2023-06-19)


### Bug Fixes

* chart release flow ([#38](https://github.com/quickube/piper/issues/38)) ([89bb175](https://github.com/quickube/piper/commit/89bb175707ae7be73374340510148b0b4fe19e6b))
* move chart path ([#35](https://github.com/quickube/piper/issues/35)) ([f9d3a42](https://github.com/quickube/piper/commit/f9d3a4253f2a32b5582ee0048386ccddd4d64937))
* move chart path ([#37](https://github.com/quickube/piper/issues/37)) ([c4da10b](https://github.com/quickube/piper/commit/c4da10b3039fe0ac4dba310e3545b8ab2d6553f6))

## [1.0.9](https://github.com/quickube/piper/compare/v1.0.8...v1.0.9) (2023-06-19)


### Bug Fixes

* minor chart fix ([#33](https://github.com/quickube/piper/issues/33)) ([cf1fdb9](https://github.com/quickube/piper/commit/cf1fdb9dfb1f2b2de3778255ce60c2f56fcf800f))

## [1.0.8](https://github.com/quickube/piper/compare/v1.0.7...v1.0.8) (2023-06-19)


### Bug Fixes

* minor chart fix ([#31](https://github.com/quickube/piper/issues/31)) ([f538110](https://github.com/quickube/piper/commit/f538110309b3c2d0ab0c9c1a7a8ce64d8ae32ec2))

## [1.0.7](https://github.com/quickube/piper/compare/v1.0.6...v1.0.7) (2023-06-19)


### Bug Fixes

* minor chart fix ([#29](https://github.com/quickube/piper/issues/29)) ([def910b](https://github.com/quickube/piper/commit/def910b5efbb5c199e314fae1d0545cdb441ece9))

## [1.0.6](https://github.com/quickube/piper/compare/v1.0.5...v1.0.6) (2023-06-19)


### Bug Fixes

* minor chart fix ([#27](https://github.com/quickube/piper/issues/27)) ([e87d7e5](https://github.com/quickube/piper/commit/e87d7e5dcf7afc1118bd2d13757adfdb0a8525e1))

## [1.0.5](https://github.com/quickube/piper/compare/v1.0.4...v1.0.5) (2023-06-19)


### Bug Fixes

* minor chart fix ([#25](https://github.com/quickube/piper/issues/25)) ([d8c9804](https://github.com/quickube/piper/commit/d8c9804416acddfdb384850e7c29fb52b7de82e2))

## [1.0.4](https://github.com/quickube/piper/compare/v1.0.3...v1.0.4) (2023-06-19)


### Bug Fixes

* minor chart fix ([#23](https://github.com/quickube/piper/issues/23)) ([3569238](https://github.com/quickube/piper/commit/35692382a4f31d74646cbe2110071624f108ea25))

## [1.0.3](https://github.com/quickube/piper/compare/v1.0.2...v1.0.3) (2023-06-19)


### Bug Fixes

* minor chart fix ([#20](https://github.com/quickube/piper/issues/20)) ([596a242](https://github.com/quickube/piper/commit/596a242a4beb871617f67dd5d6a06f04039e4f46))
* minor chart fix ([#22](https://github.com/quickube/piper/issues/22)) ([037dda7](https://github.com/quickube/piper/commit/037dda7e8e845b3b3e5d6576cf2bb0b63e74c4c0))

## [1.0.2](https://github.com/quickube/piper/compare/v1.0.1...v1.0.2) (2023-06-19)


### Bug Fixes

* minor chart fix ([#18](https://github.com/quickube/piper/issues/18)) ([4e34925](https://github.com/quickube/piper/commit/4e34925591e795972c5d4bae9315666e5abadbc9))

## [1.0.1](https://github.com/quickube/piper/compare/v1.0.0...v1.0.1) (2023-06-19)


### Bug Fixes

* minor chart fix ([#16](https://github.com/quickube/piper/issues/16)) ([14cfa19](https://github.com/quickube/piper/commit/14cfa193d3f9151f05ce0e77a0e2f1416cd3ccf7))

## 1.0.0 (2023-06-19)


### Features

* UT | Add test to AddFilesToTemplates ([#5](https://github.com/quickube/piper/issues/5)) ([c316e0e](https://github.com/quickube/piper/commit/c316e0e301494e3edf48a614ac84fcca3f77a688))


### Bug Fixes

* add UT for utils ([#9](https://github.com/quickube/piper/issues/9)) ([75819de](https://github.com/quickube/piper/commit/75819dec4d5e4bab4da1dd440dde8dbe6e865be5))
* create ci ([#11](https://github.com/quickube/piper/issues/11)) ([9e193dd](https://github.com/quickube/piper/commit/9e193ddc5450de973b7c1d1b7c069c78ac371ca7))
* edge cases with nil pointer dereferences ([#8](https://github.com/quickube/piper/issues/8)) ([b36a8ef](https://github.com/quickube/piper/commit/b36a8ef40877f33c2ffdede0694188bfdca572b1))
* Improve UT to IsOrgWebhookEnabled ([#7](https://github.com/quickube/piper/issues/7)) ([eb5de70](https://github.com/quickube/piper/commit/eb5de701e895dfa51a77ada42184325479c93198))
* linting issues ([206bc9e](https://github.com/quickube/piper/commit/206bc9eaa1f403d68992d65ff6f870e1930ca844))
* linting issues ([#14](https://github.com/quickube/piper/issues/14)) ([3c9724a](https://github.com/quickube/piper/commit/3c9724a3b33ed83ecaa5ecac5973ba67dd2be6b3))
* quickube default value fix ([#15](https://github.com/quickube/piper/issues/15)) ([a6b1158](https://github.com/quickube/piper/commit/a6b1158dc491bc7c6d47b4784ad1afcc80d64389))
* typo in package name ([#12](https://github.com/quickube/piper/issues/12)) ([114d12a](https://github.com/quickube/piper/commit/114d12a79853b1db675dbb912b9881ce7f3c4795))
