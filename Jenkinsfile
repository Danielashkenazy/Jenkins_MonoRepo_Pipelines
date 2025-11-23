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

        stage('Lint User Service') {
            when { expression { env.SERVICES_CHANGED.contains("user-service") } }
            steps {
                sh """
                cd user-service
                npm install
                npx eslint .
                """
            }
        }

        stage('Lint Transaction Service') {
            when { expression { env.SERVICES_CHANGED.contains("transaction-service") } }
            steps {
                sh """
                cd transaction-service
                python3 -m venv .venv
                . .venv/bin/activate
                pip install -r requirements.txt
                pip install flake8
                flake8 . --exclude=.venv,__pycache__,.git
                """
            }
        }

        stage('Lint Notification Service') {
            when { expression { env.SERVICES_CHANGED.contains("notification-service") } }
            steps {
                sh """
                cd notification-service
                go mod tidy
                go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
                golangci-lint run
                """
            }
        }

        stage('Test User Service') {
            when { expression { env.SERVICES_CHANGED.contains("user-service") } }
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

        stage('Test Transaction Service') {
            when { expression { env.SERVICES_CHANGED.contains("transaction-service") } }
            steps {
                sh """
                cd transaction-service
                python3 -m venv .venv || true
                . .venv/bin/activate
                pip install pytest pytest-cov
                pytest --cov=. --cov-report=html --cov-report=xml --junitxml=junit.xml
                """
            }
            post {
                always {
                    junit 'transaction-service/junit.xml'
                    publishHTML(target: [
                        reportDir: 'transaction-service/htmlcov',
                        reportFiles: 'index.html',
                        reportName: 'Transaction Service Coverage Report'
                    ])
                }
            }
        }

        stage('Test Notification Service') {
            when { expression { env.SERVICES_CHANGED.contains("notification-service") } }
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

        stage('Security Scan User Service') {
            when { expression { env.SERVICES_CHANGED.contains("user-service") } }
            steps {
                sh """
                cd user-service
                npm audit --audit-level=moderate
                """
            }
        }

        stage('Security Scan Transaction Service') {
            when { expression { env.SERVICES_CHANGED.contains("transaction-service") } }
            steps {
                sh """
                cd transaction-service
                if [ ! -d ".venv" ]; then python3 -m venv .venv; fi
                . .venv/bin/activate
                pip install -r requirements.txt
                pip install bandit
                bandit -r app -x .venv,tests,__pycache__,**/site-packages/** -ll
                """
            }
        }

        stage('Security Scan Notification Service') {
            when { expression { env.SERVICES_CHANGED.contains("notification-service") } }
            steps {
                sh """
                go install github.com/securego/gosec/v2/cmd/gosec@latest
                export PATH=\$(go env GOPATH)/bin:\$PATH
                cd notification-service
                gosec -severity medium -confidence medium -fmt json -out gosec-report.json ./...
                """
            }
        }
        stage('Ready for Deployment') {
            when { expression { env.SERVICES_CHANGED?.trim() } }
            agent none
            steps {
                script {
                    input message: "Approve deployment?", ok: "Proceed"
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
