pipeline {
  agent any
  stages {
    stage('Build Image') {
      steps {
        sh 'docker build .'
      }
    }
    stage('Tag Image') {
      steps {
        sh 'docker tag gempbotgo_gempbotgo gempir/gempbotgo'
      }
    }
    stage('Push Image') {
      steps {
        sh 'docker push gempir/gempbotgo'
      }
    }
    stage('Deploy') {
      steps {
        sh '''cp ./prod.yml /home/gempir/gempbotgo && cd /home/gempir/gempbotgo &&
docker-compose -f prod.yml pull
&& docker-compose -f prod.yml up -d'''
      }
    }
  }
}