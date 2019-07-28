# KITA is a Kubernetes operator for building Instructor-led Tutorials and Awesome demos

## General Idea

Enable declarative setups for demos, hand-on sessions etc and combine that setup in an all-browser-based environment using [Coder](https://coder.com/).

## Kita Space CR

```yaml
apiVersion: kita.danistrebel.io/v1alpha1
kind: KitaSpace
metadata:
  name: awesome-space
spec:
  owner:
    name: john
    email: john.doe@example.com
  repos:
    - https://github.com/ramitsurana/awesome-kubernetes.git
    - https://github.com/operator-framework/awesome-operators
```

## Result

![Result Screenshot](/documentation/editor-screenshot.png?raw=true "Editor Screenshot")

## Optional Sendgrid Integration for Space Login Token

To send an email to the Kita Space Owner containing the login token create a Sendgrid API Key and Set the following Env Parameters

```bash
export SENDGRID_API_KEY='SG.xxxxxREPLACE_THIS_TOKENxxxxx'
export SENDGRID_EMAIL_SENDER='sender@example.com'
```
