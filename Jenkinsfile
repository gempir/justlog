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
    stage('Copy compose.yml') {
      steps {
        sh 'cp ./prod.yml /home/gempir/gempbotgo'
      }
    }
    stage('Change Dir') {
      steps {
        sh ' cd /home/gempir/gempbotgo'
      }
    }
    stage('Pull Image') {
      steps {
        sh 'docker-compose -f prod.yml pull'
      }
    }
    stage('Restart App') {
      steps {
        sh 'docker-compose -f prod.yml up -d'
      }
    }
  }
}