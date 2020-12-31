pipeline {
    agent any
    tools {
        go 'Go 1.15 Compiler'
    }
    environment {
        GO115MODULE = 'on'
        CGO_ENABLED = 0 
        GOPATH = "${JENKINS_HOME}/jobs/${JOB_NAME}/builds/${BUILD_ID}"
    }
    stages {        
        stage('Pre Test') {
            steps {
                echo 'Installing dependencies'
                sh 'go version'
                sh 'go get -u golang.org/x/lint/golint'
                sh 'go mod download'
            }
        }

        stage('Test') {
            steps {
                withEnv(["PATH+GO=${GOPATH}/bin"]){
                    echo 'Running vetting'
                    sh 'go vet .'
                    echo 'Running linting'
                    sh 'golint .'
                    echo 'Running test'
                    sh 'go test'
                }
            }
        }

        stage('Pre integration test') {
            steps {
                echo 'Bringing up docker container for integration test'
                sh 'docker-compose up -d'
            }
        }

        stage('Run integration tests') {
            steps{
                withEnv(["PATH+GO=${GOPATH}/bin"]){
                    echo 'Running vetting'
                    sh 'go vet .'
                    echo 'Running linting'
                    sh 'golint .'
                    echo 'Running test'
                    sh 'go test --tags=integration'
                }
            }
        }

        stage('Stop containers') {
            steps{
                sh 'docker-compose down'
            }
        }

        stage('Build') {
            when {
                branch 'master'
            }
            steps {
                echo 'Compiling and building'
                sh 'go build'
            }
        }
        
    }
    // Add email notifications
    // post {
    //     always {
    //         emailext body: "${currentBuild.currentResult}: Job ${env.JOB_NAME} build ${env.BUILD_NUMBER}\n More info at: ${env.BUILD_URL}",
    //             recipientProviders: [[$class: 'DevelopersRecipientProvider'], [$class: 'RequesterRecipientProvider']],
    //             to: "${params.RECIPIENTS}",
    //             subject: "Jenkins Build ${currentBuild.currentResult}: Job ${env.JOB_NAME}"
            
    //     }
    // }  
}