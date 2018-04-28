pipeline {
  agent any
  stages {
    stage('Build') {
      steps {
        sh '''ls -lsa
docker build -t gempir/gempbotgo .
docker push gempir/gempbotgo'''
      }
    }
    stage('Deploy') {
      steps {
        sh 'cp ./prod.yml /home/gempir/gempbotgo'
        sh '''cd /home/gempir/gempbotgo
docker-compose -f prod.yml pull
docker-compose -f prod.yml up -d'''
      }
    }
  }
}