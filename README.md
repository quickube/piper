# Piper
![alt text](https://www.rookout.com/wp-content/uploads/2022/10/ArgoPipeline_1.0_Hero.png.webp?raw=true)

Welcome to Piper! Piper is an open-source project aimed at providing multibranch pipeline functionality to Argo Workflows, allowing users to create distinct Workflows based on Git branches.

## Table of Contents

- [Piper](#piper)
  - [Table of Contents](#table-of-contents)
  - [Getting Started](#getting-started)
  - [Reporting Issues](#reporting-issues)
  - [How to Contribute](#how-to-contribute)
  - [License](#license)

## Getting Started

Piper configures a webhook in the git provider and listens to the webhooks sent. It will create a Workflow CRD out of branches that contain a `.workflows` folder.
This folder should contain declarations of the templates and the main DAG that will be running.
Finally, it will submit the Workflow as a K8s resource in the cluster.
To access more detailed explanations, please navigate to the [Documentation site](https://piper.quickube.com).

https://github.com/quickube/piper/assets/106976988/09b3a5d8-3428-4bdc-9146-3034d81164bf

## Reporting Issues

If you encounter any issues or bugs while using Piper, please help us improve by reporting them. Follow these steps to report an issue:

1. Go to the [Piper Issues](https://github.com/quickube/Piper/issues) page on GitHub.
2. Click on the "New Issue" button.
3. Provide a descriptive title and detailed description of the issue, including any relevant error messages or steps to reproduce the problem.
4. Add appropriate labels to categorize the issue (e.g., bug, enhancement, question).
5. Submit the issue, and our team will review and address it as soon as possible.

## How to Contribute

If you're interested in contributing to this project, please feel free to submit a pull request. We welcome all contributions and feedback.
Please check out our [contribution guidelines for this project](docs/CONTRIBUTING.md).

## License

This project is licensed under the Apache License. Please see the [LICENSE](LICENSE) file for details.