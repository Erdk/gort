pipeline {
    agent any

    tools {
        go '1.19.1'
    }

    stages {
        stage('SCM') {
            steps {
                git branch: 'main', url: 'https://github.com/Erdk/gort'
            }
        }
        stage('Build') {
            steps {
                sh "go build ."
            }
        }
        stage('Tests') {
            environment {
                CGO_ENABLED = '0'
            }
            steps {
                sh "go test ./..."
            }
        }
    }
}
