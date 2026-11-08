// Ported from .github/workflows/ci.yml (Go 1.24: download, vet, build, test).
pipeline {
    agent { label 'paperflow-docker' }
    options {
        skipDefaultCheckout(true)
        disableConcurrentBuilds()
        timeout(time: 30, unit: 'MINUTES')
        buildDiscarder(logRotator(numToKeepStr: '20', daysToKeepStr: '30'))
        timestamps()
    }
    stages {
        stage('Checkout and prepare') {
            options { timeout(time: 10, unit: 'MINUTES') }
            steps {
                deleteDir()
                checkout scm
                script {
                    env.CI_RUN_ID = sh(returnStdout: true, script: 'printf "ci-%s" "$(printf %s "$BUILD_TAG" | sha256sum | cut -c1-24)"').trim()
                }
                sh 'bash ci/run.sh init'
            }
        }
        stage('Download dependencies') {
            options { timeout(time: 10, unit: 'MINUTES') }
            steps { sh 'bash ci/run.sh deps' }
        }
        stage('Vet') {
            options { timeout(time: 10, unit: 'MINUTES') }
            steps { sh 'bash ci/run.sh vet' }
        }
        stage('Build') {
            options { timeout(time: 10, unit: 'MINUTES') }
            steps { sh 'bash ci/run.sh build' }
        }
        stage('Test') {
            options { timeout(time: 15, unit: 'MINUTES') }
            steps { sh 'bash ci/run.sh test' }
        }
    }
    post {
        always {
            script {
                try {
                    if (env.CI_RUN_ID && fileExists('ci/run.sh')) {
                        timeout(time: 3, unit: 'MINUTES') { sh 'bash ci/run.sh cleanup' }
                    }
                } finally {
                    deleteDir()
                }
            }
        }
    }
}
