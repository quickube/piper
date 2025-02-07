## How To Contribute

We appreciate contributions from the community to make Piper even better. To contribute, follow the steps below:

1. Fork the Piper repository to your GitHub account.
2. Clone the forked repository to your local machine:
```bash
git clone https://github.com/your-username/Piper.git
```
3. Create a new branch to work on your feature or bug fix:
```bash
git checkout -b my-feature
```
4. Make your changes, following the coding guidelines outlined in this document.
5. Commit your changes with clear and descriptive commit messages and sign it:
```bash
git commit -s -m "fix: Add new feature"
```
* Please make sure you commit as described in [conventional commit](https://www.conventionalcommits.org/en/v1.0.0/)
6. Push your changes to your forked repository:
```bash
git push origin my-feature
```
7. Open a [pull request](#pull-requests) against the main branch of the original Piper repository.

## Pull Requests

We welcome and appreciate contributions from the community. If you have developed a new feature, improvement, or bug fix for Piper, follow these steps to submit a pull request:

1. Make sure you have forked the Piper repository and created a new branch for your changes. Checkout [How To Contribute](#How-to-contribute).
2. Commit your changes and push them to your forked repository.
3. Go to the Piper repository on GitHub.
4. Click on the "New Pull Request" button.
5. Select your branch and provide a [descriptive title](#pull-request-naming) and a detailed description of your changes.
6. If your pull request relates to an open issue, reference the issue in the description using the GitHub issue syntax (e.g., Fixes #123).
7. Submit the pull request, and our team will review your changes. We appreciate your patience during the review process and may provide feedback or request further modifications.

### Pull Request Naming

The name should follow Conventional Commits naming.

## Coding Guidelines

To maintain a consistent codebase and ensure readability, we follow a set of coding guidelines in Piper. Please adhere to the following guidelines when making changes:

* Follow the [Effective Go](https://go.dev/doc/effective_go) guide for Go code.
* Follow the [Folder convention](https://github.com/golang-standards/project-layout) guide for Go code.
* Write clear and concise comments to explain the code's functionality.
* Use meaningful variable and function names.
* Make sure your code is properly formatted and free of syntax errors.
* Run tests locally.
* Check that the feature is documented.
* Add new packages only if necessary and if an already existing one can't be used.
* Add tests for new features or modification.

## Helm Chart Development

To make sure that the documentation is updated use [helm-docs](https://github.com/norwoodj/helm-docs) comment convention. The pipeline will execute `helm-docs` command and update the version of the chart.

Also, please make sure to run those commands locally to debug the chart before merging:

```bash
make helm
```

## Local deployment

To make it easy to develop locally, please run the following

Prerequisites :
1. install helm
```bash
brew install helm
```
2. install kubectl
```bash
brew install kubectl
```
3. install docker

4. install ngrok
```bash
brew install ngrok
```
5. install kind
```bash
brew install kind
```

Deployment:
1. Make sure Docker is running.
2. Create a tunnel with ngrok using `make ngrok`, and save the `Forwarding` address.
3. Create a `values.dev.yaml` file that contains a subset of the chart's `values.yaml` file. Check the [example of values file](https://github.com/quickube/piper/tree/main/examples/template.values.dev.yaml), rename it to `values.dev.yaml`, and put it in the root directory.
4. Use `make deploy`. It will do the following:
   * Deploy a local registry as a container
   * Deploy a kind cluster as a container with configuration
   * Deploy nginx reverse proxy in the kind cluster
   * Deploy Piper with the local helm chart
5. Validate using `curl localhost/piper/healthz`.

## Debugging

For debugging, the best practice is to use Rookout. To enable this function, pass a Rookout token in the chart `rookout.token` or as an existing secret `rookout.existingSecret`.