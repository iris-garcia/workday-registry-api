+++
title = "Continuous Delivery (GitHub)"
author = ["Iris Garcia"]
lastmod = 2019-12-06T23:33:09+01:00
tags = ["openshift", "cd"]
draft = false
weight = 3
asciinema = true
+++

## Step 1: Create a secret encoded in base64 {#step-1-create-a-secret-encoded-in-base64}

We need to create a secret which need to be known by our deployment in
OpenShift and the Webhook in GitHub:

```bash
$ echo 'supersecret' | base64
c3VwZXJzZWNyZXQ=
```


## Step 2: Deploy in OpenShift using the encoded secret {#step-2-deploy-in-openshift-using-the-encoded-secret}

The following resources in the OpenShift deployment template are the
ones creating the secret and a trigger for GitHub.

```yaml
# ...
- kind: Secret
  apiVersion: v1
  metadata:
    name: gh-secret
    creationTimestamp:
  data:
    WebHookSecretKey: "${GITHUB_SECRET}"
# ...
- kind: BuildConfig
  apiVersion: v1
  metadata:
    name: api
    annotations:
      description: "Defines how to build the application"
  spec:
    source:
      type: Git
      git:
        uri: "${SOURCE_REPOSITORY_URL}"
        ref: "${SOURCE_REPOSITORY_REF}"
      contextDir: "${CONTEXT_DIR}"
    strategy:
      type: Docker
      dockerStrategy: {}
    output:
      to:
        kind: ImageStreamTag
        name: api:latest
    postCommit:
      script: "GIN_MODE=release go test -v ./..."
    resources:
      limits:
        cpu: 100m
        memory: 1Gi
    triggers:
    - type: "GitHub"
      github:
        secretReference:
          name: "gh-secret"
#...
parameters:
- name: SOURCE_REPOSITORY_URL
  description: "The URL of the repository with your application source code"
  value: "https://github.com/iris-garcia/workday.git"
- name: SOURCE_REPOSITORY_REF
  description: "Set this to a branch name, tag or other ref of your repository if you are not using the default branch"
- name: CONTEXT_DIR
  description: "Set this to the relative path to your project if it is not in the root of your repository"
- name: GITHUB_SECRET
  description: "Github webhook secret"
```

Then we simply need to run the deploy passing the encoded secret as a
parameter:

```bash
oc new-app deployment/openshift.yml -p GITHUB_SECRET='c3VwZXJzZWNyZXQ='
```


## Step 3: Create a GitHub webhook {#step-3-create-a-github-webhook}

In this step we will create a new GitHub webhook which will send a
_POST_ request to our OpenShift's app endpoint everytime there is a
new push.

To get the enpoint generated by OpenShift we just need to run the
following command:

```bash
$ oc describe bc api endpoint
Name:		api
Namespace:	workday
Created:	5 days ago
Labels:		app=api
Description:	Defines how to build the application
Annotations:	openshift.io/generated-by=OpenShiftNewApp
Latest Version:	21

Strategy:	Docker
URL:		https://github.com/iris-garcia/workday.git
Output to:	ImageStreamTag api:latest

Build Run Policy:	Serial
Triggered by:		<none>
Webhook GitHub:
        URL:	https://api.us-east-2.starter.openshift-online.com:6443/apis/build.openshift.io/v1/namespaces/workday/buildconfigs/api/webhooks/<secret>/github
Builds History Limit:
        Successful:	5
        Failed:		5
```

{{% notice note %}}
**Save the URL, it will be needed in the creation of the GitHub's webhook, and replace the <secret> with the real secret used.**
{{% /notice %}}

To create the webhook browse to the GitHub's repository and click in
**Settings**.
![](/images/gh_oc_1.png)

In the **Settings** page, click in **Webhooks** then **Add webhook**.
![](/images/gh_oc_2.png)

Then we need to fill the following fields:

-   **Payload URL**: The one we copied in the [Step 3](https://iris-garcia.github.io/workday/howto/github-cd/#step-3-create-a-github-webhook).
-   **Content type**: It has to be `application/json`
-   **Secret**: Leave it empty (the secret is included in the payload url).

Then finally click in **Add webhook**.
![](/images/gh_oc_3.png)