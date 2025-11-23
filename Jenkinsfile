pipeline {
    agent { label 'linux' }

    stages {

        stage('Detect Changed Services') {
            steps {
                script {
                    def changes = sh(script: "bash shared/ci/detect_changes.sh", returnStdout: true).trim()
                    echo "Changed Services: ${changes}"

                    if (!changes) {
                        env.SERVICES_CHANGED = ""
                        currentBuild.result = 'SUCCESS'
                        return
                    }

                    env.SERVICES_CHANGED = changes
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
                                sh """
                                cd user-service
                                npm install
                                npx eslint .
                                """
                            }
                        }

                        stage('Test User Service') {
                            steps {
                                sh """
                                cd user-service
                                npm install
                                npm test -- --coverage --reporters=default --reporters=jest-junit
                                """
                            }
                            post {
                                always {
                                    junit 'user-service/junit.xml'
                                    publishHTML(target: [
                                        reportDir: 'user-service/coverage',
                                        reportFiles: 'index.html',
                                        reportName: 'User Service Coverage Report'
                                    ])
                                }
                            }
                        }

                        stage('Security Scan User Service') {
                            steps {
                                sh """
                                cd user-service
                                npm audit --audit-level=moderate
                                """
                            }
                        }
                    }
                }

                stage('Notification Service Pipeline') {
                    when { expression { env.SERVICES_CHANGED.contains("notification-service") } }

                    stages {

                        stage('Lint Notification Service') {
                            steps {
                                sh """
                                cd notification-service
                                go mod tidy
                                go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
                                golangci-lint run
                                """
                            }
                        }

                        stage('Test Notification Service') {
                            steps {
                                sh """
                                cd notification-service
                                go mod tidy
                                go test -v -coverprofile=coverage.out -covermode=atomic ./...
                                go tool cover -html=coverage.out -o coverage.html
                                """
                            }
                            post {
                                always {
                                    publishHTML(target: [
                                        reportDir: 'notification-service',
                                        reportFiles: 'coverage.html',
                                        reportName: 'Notification Service Coverage Report'
                                    ])
                                }
                            }
                        }

                        stage('Security Scan Notification Service') {
                            steps {
                                sh """
                                go install github.com/securego/gosec/v2/cmd/gosec@latest
                                export PATH=\$(go env GOPATH)/bin:\$PATH
                                cd notification-service
                                gosec -severity medium -confidence medium -fmt json -out gosec-report.json ./...
                                """
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
