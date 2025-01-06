# Introduction

<p>
  <img src="https://www.rookout.com/wp-content/uploads/2022/10/ArgoPipeline_1.0_Hero.png.webp?raw=true"  alt="ArgoPipeline"/>
</p>

Welcome to Piper!

Piper is an open-source project aimed at providing multibranch pipeline functionality to Argo Workflows. This allows users to create distinct Workflows based on Git branches. We support GitHub and Bitbucket.

## General Explanation

<p>
  <img src="https://raw.githubusercontent.com/quickube/piper/main/docs/img/flow.svg"  alt="flow"/>
</p>

Piper handles the hard work of configuring multibranch pipelines for us! At initialization, it will load all configuration and create a webhook in the repository or organization scope. Then, for each branch that has a `.workflows` folder, Piper will create a Workflow CRD out of the files in this folder. Finally, when Piper detects changes in the repository via the webhook, it triggers the workflows that match the branch and event.

![type:video](./img/piper-demo-1080.mp4)