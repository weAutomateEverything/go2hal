#! /bin/bash
# Push only if it's not a pull request
if [ -z "$TRAVIS_PULL_REQUEST" ] || [ "$TRAVIS_PULL_REQUEST" == "false" ]; then
  # Push only if we're testing the master branch
  if [ "$TRAVIS_BRANCH" == "master" ]; then

    # This is needed to login on AWS and push the image on ECR
    # Change it accordingly to your docker repo
    pip install --user awscli
    export PATH=$PATH:$HOME/.local/bin
    eval $(aws ecr get-login --region $AWS_DEFAULT_REGION)

    # Build and push
    docker build -t $IMAGE_NAME:1.$TRAVIS_JOB_NUMBER .
    echo "Pushing $IMAGE_NAME:2.$TRAVIS_JOB_NUMBER"
    docker tag $IMAGE_NAME:2.$TRAVIS_JOB_NUMBER "$IMAGE_NAME:1.$TRAVIS_JOB_NUMBER"
    docker login -u $DOCKER_USERNAME -p $DOCKER_PASSWORD
    docker push "$IMAGE_NAME:1.$TRAVIS_JOB_NUMBER"
    echo "Pushed $IMAGE_NAME:1.$TRAVIS_JOB_NUMBER"
  else
    echo "Skipping deploy because branch is not 'master'"
  fi
else
  echo "Skipping deploy because it's a pull request"
fi