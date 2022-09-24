pipeline {
    agent any

    tools {
        go '1.19.1'
    }

    stages {
        stage('Build') {
            steps {
                git branch: 'main', url: 'https://github.com/Erdk/gort'
                sh "go build ."
            }
        }
    }
}
