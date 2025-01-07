## Global Variables

Piper will automatically add Workflow scope parameters that can be referenced from any template.
The parameters are taken from webhook metadata and will be populated according to the GitProvider and the event that triggered the workflow.

1. `{{ workflow.parameters.event }}` The event that triggered the workflow.

2. `{{ workflow.parameters.action }}` The action that triggered the workflow.

3. `{{ workflow.parameters.dest_branch }}` The destination branch for the pull request.

4. `{{ workflow.parameters.commit }}` The commit that triggered the workflow.

5. `{{ workflow.parameters.repo }}` The repository name that triggered the workflow.

6. `{{ workflow.parameters.user }}` The username that triggered the workflow.

7. `{{ workflow.parameters.user_email }}` The user's email that triggered the workflow.

8. `{{ workflow.parameters.pull_request_url }}` The URL of the pull request that triggered the workflow.

9. `{{ workflow.parameters.pull_request_title }}` The title of the pull request that triggered the workflow.

10. `{{ workflow.parameters.pull_request_labels }}` Comma-separated labels of the pull request that triggered the workflow.