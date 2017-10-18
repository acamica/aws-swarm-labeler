pipeline {
    agent any
    tools {
        nodejs 'Node 6.x'
    }
    environment {
        REGISTRY      = credentials('acamica-registry-password')
    }
    stages {
        stage('Build Docker') {
            steps {
                sh "echo //registry.npmjs.org/:_authToken=${env.ACAMICA_TOKEN} > .npmrc" //TODO: replace with folder credentials
                sh './docker-build.sh ${JOB_NAME} jenkins-${BUILD_NUMBER}'
                sh 'rm .npmrc'
            }
        }

        stage('Push to Docker Registry'){
            when { branch 'release' }
            steps {
                sh 'docker login -u ${REGISTRY_USR} -p ${REGISTRY_PSW} registry.acamica.com' //TODO: replace with folder credentials
                sh './docker-push.sh ${JOB_NAME} jenkins-${BUILD_NUMBER}'
            }
        }
    }

    post {
        success {
            slackSend channel: env.SLACK_SUCC_CHANNEL, color: 'good', message: "Build finished successfully : ${env.JOB_NAME} - ${env.BUILD_NUMBER} (<${env.JOB_URL}|Open>)"
        }
        failure {
            slackSend channel: env.SLACK_ERR_CHANNEL, color: '#FF0000', message: "Build failed: ${env.JOB_NAME} - ${env.BUILD_NUMBER} (<${env.JOB_URL}|Open>)"
        }
    }
}

