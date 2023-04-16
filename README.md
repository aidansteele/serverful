# serverful

## Putting the servers back in serverless

Scenario: it's opposite day. Your organisation is entirely serverless, but you
want a server. You are too forgetful for EC2 or Fargate: you never remember to
shut them down. What about a Lambda-powered server that automatically shuts off
after 1-900 seconds? Enter `serverful`.

## Usage

Your app has a Dockerfile. You add these lines to your Dockerfile:

```dockerfile
COPY --from=ghcr.io/aidansteele/serverful:main /opt/extensions/serverful /opt/extensions/serverful
# don't complain about this next line. it's hardly the worst idea here
ENV NGROK_AUTHTOKEN=your-token-from-ngrok.com 
```

Next, deploy your app to Lambda, perhaps using this SAM template and 
`sam build && sam deploy`:

```yaml
Transform: AWS::Serverless-2016-10-31

Resources:
  BadIdea:
    Type: AWS::Serverless::Function
    Metadata:
      DockerContext: ./example
      Dockerfile: Dockerfile
    Properties:
      PackageType: Image
      MemorySize: 1769
      Timeout: 40
      Architectures: [arm64]
      FunctionUrlConfig:
        AuthType: NONE

Outputs:
  Url:
    Value: !GetAtt BadIdeaUrl.FunctionUrl
```

Now you can navigate to your Lambda [function URL][furl]. It will start an ngrok
tunnel to your serverful app and issue a 302 Redirect to the tunneled domain. 
Enjoy the next 40 seconds of serverful silliness.

[furl]: https://docs.aws.amazon.com/lambda/latest/dg/lambda-urls.html
