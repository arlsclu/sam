pipeline {
    agent  any
    tools ('go' , 'docker')
    environment {
        tag = 'sam:v0.0.1' 
    }
    stage('testing')
        steps {
            sh 'hello world ! '
            go test ./...
        }
    stage('build')
        steps {
            docker build -t ${env.tag} .
            docker login ...
            docker push --tag  ..
        }
    stage('deploy')
        steps {
            docker  run  -d ${env.tag}
        }

    always {
        post {
            docker logout
        }
    }



}