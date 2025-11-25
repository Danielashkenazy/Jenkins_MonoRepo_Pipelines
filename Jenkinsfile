// ========================================
// Reusable Groovy Functions 
// ========================================
def runLintForServices() {
    sh "bash shared/ci/lint.sh '${env.SERVICES_CHANGED}'"
}

def runTestsForServices() {
    sh "bash shared/ci/test.sh '${env.SERVICES_CHANGED}'"
}

def runSecurityScanForServices() {
    sh "bash shared/ci/scan.sh '${env.SERVICES_CHANGED}'"
}

def publishReportsForService(service) {
    // JUnit report (if exists)
    if (fileExists("${service}/junit.xml")) {
        junit "${service}/junit.xml"
    }
    
    // Coverage report
    if (fileExists("${service}/coverage/index.html")) {
        publishHTML(target: [
            reportDir: "${service}/coverage",
            reportFiles: 'index.html',
            reportName: "${service} Coverage Report"
        ])
    } else if (fileExists("${service}/htmlcov/index.html")) {
        publishHTML(target: [
            reportDir: "${service}/htmlcov",
            reportFiles: 'index.html',
            reportName: "${service} Coverage Report"
        ])
    } else if (fileExists("${service}/coverage.html")) {
        publishHTML(target: [
            reportDir: service,
            reportFiles: 'coverage.html',
            reportName: "${service} Coverage Report"
        ])
    }
}

// ========================================
// Pipeline Definition
// ========================================
pipeline {
    agent { label 'linux' }

    stages {

        stage('Detect Changed Services') {
            steps {
                script {
                    def changes = sh(script: "bash shared/ci/detect_changes.sh", returnStdout: true).trim()
                    echo "Changed Services: ${changes}"

                    if (changes == "") {
                        echo "No services changed. Skipping..."
                        env.SERVICES_CHANGED = ""
                        currentBuild.result = 'SUCCESS'
                        return
                    }

                    env.SERVICES_CHANGED = changes
                }
            }
        }

        stage('Secrets Detection') {
            when { expression { env.SERVICES_CHANGED?.trim() } }
            steps {
                script {
                    echo "Running secrets detection with TruffleHog..."
                    sh """
                    docker run --rm -v "\$(pwd):/scan" trufflesecurity/trufflehog:latest filesystem /scan --fail --no-update \
                    --exclude-paths '**/.venv' '**/node_modules' '**/package-lock.json' '**/*.pyc' '**/*dist-info/RECORD'
                    """
                }
            }
        }

        stage('Run Services in Parallel') {
            when { expression { env.SERVICES_CHANGED?.trim() } }
            
            parallel {

                stage('User Service Pipeline') {
                    when { expression { env.SERVICES_CHANGED.contains("user-service") } }
                    
                    stages {
                        stage('Lint User Service') {
                            steps {
                                script {
                                    sh "bash shared/ci/lint.sh 'user-service'"
                                }
                            }
                        }

                        stage('Test User Service') {
                            steps {
                                script {
                                    sh "bash shared/ci/test.sh 'user-service'"
                                }
                            }
                            post {
                                always {
                                    script {
                                        publishReportsForService('user-service')
                                    }
                                }
                            }
                        }

                        stage('Security Scan User Service') {
                            steps {
                                retry(2) {
                                    script {
                                        sh "bash shared/ci/scan.sh 'user-service'"
                                    }
                                }
                            }
                        }
                    }
                }

                stage('Transaction Service Pipeline') {
                    when { expression { env.SERVICES_CHANGED.contains("transaction-service") } }
                    
                    stages {
                        stage('Lint Transaction Service') {
                            steps {
                                script {
                                    sh "bash shared/ci/lint.sh 'transaction-service'"
                                }
                            }
                        }

                        stage('Test Transaction Service') {
                            steps {
                                script {
                                    sh "bash shared/ci/test.sh 'transaction-service'"
                                }
                            }
                            post {
                                always {
                                    script {
                                        publishReportsForService('transaction-service')
                                    }
                                }
                            }
                        }

                        stage('Security Scan Transaction Service') {
                            steps {
                                retry(2) {
                                    script {
                                        sh "bash shared/ci/scan.sh 'transaction-service'"
                                    }
                                }
                            }
                        }
                    }
                }

                stage('Notification Service Pipeline') {
                    when { expression { env.SERVICES_CHANGED.contains("notification-service") } }
                    
                    stages {
                        stage('Lint Notification Service') {
                            steps {
                                script {
                                    sh "bash shared/ci/lint.sh 'notification-service'"
                                }
                            }
                        }

                        stage('Test Notification Service') {
                            steps {
                                script {
                                    sh "bash shared/ci/test.sh 'notification-service'"
                                }
                            }
                            post {
                                always {
                                    script {
                                        publishReportsForService('notification-service')
                                    }
                                }
                            }
                        }

                        stage('Security Scan Notification Service') {
                            steps {
                                retry(2) {
                                    script {
                                        sh "bash shared/ci/scan.sh 'notification-service'"
                                    }
                                }
                            }
                        }
                    }
                }
            }
        }

        stage('Ready for Deployment') {
            when { expression { env.SERVICES_CHANGED?.trim() } }
            steps {
                script {
                    def userInput = false
                    try {
                        timeout(time: 10, unit: 'MINUTES') {
                            userInput = input(
                                id: 'Proceed1',
                                message: 'Deploy to Docker Hub?',
                                parameters: [
                                    [$class: 'BooleanParameterDefinition',
                                     defaultValue: true,
                                     description: '',
                                     name: 'Please confirm you agree with this']
                                ]
                            )
                        }
                    } catch(err) {
                        currentBuild.result = 'ABORTED'
                        error("Deployment cancelled")
                    }
                }
            }
        }

        stage('Build Docker Images') {
            when { expression { env.SERVICES_CHANGED } }
            steps {
                script {
                    def shortSha = sh(script: "git rev-parse --short HEAD", returnStdout: true).trim()
                    def services = env.SERVICES_CHANGED.split(" ")

                    withCredentials([usernamePassword(credentialsId: 'dockerhub', usernameVariable: 'DOCKER_USER', passwordVariable: 'DOCKER_PASS')]) {
                        sh """
                        echo "$DOCKER_PASS" | docker login -u "$DOCKER_USER" --password-stdin
                        """
                    }

                    services.each { svc ->
                        sh """
                        cd ${svc}
                        docker build -t danielashkenazy1/${svc}:ci-${shortSha} .
                        docker push danielashkenazy1/${svc}:ci-${shortSha}
                        """
                    }
                }
            }
        }
    }

    post {
        always {
            script {
                def status = currentBuild.result ?: "SUCCESS"
                def emoji = (status == "SUCCESS") ? "✅" : "❌"

                withCredentials([string(credentialsId: 'slack_webhook', variable: 'SLACK_URL')]) {
                    sh """
                    curl -X POST -H 'Content-type: application/json' \
                    --data '{"text": "${emoji} *Pipeline Status:* ${status}\n*Branch:* ${env.GIT_BRANCH}\n*Commit:* ${env.GIT_COMMIT}"}' \
                    $SLACK_URL
                    """
                }
            }
        }
    }
}