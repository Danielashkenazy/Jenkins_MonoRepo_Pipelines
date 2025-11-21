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
                set -e
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
                set -e
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
                set -e
                cd notification-service
                go mod tidy
                go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
                golangci-lint run
                """
            }
        }

        // ========== UNIT TESTING STAGES ==========

        stage('Test User Service') {
            when { expression { env.SERVICES_CHANGED.contains("user-service") } }
            steps {
                sh """
                set -e
                cd user-service
                npm install
                npm test -- --coverage --reporters=default --reporters=jest-junit
                """
            }
            post {
                always {
                    junit 'user-service/junit.xml'
                    publishHTML(target: [
                        allowMissing: false,
                        alwaysLinkToLastBuild: true,
                        keepAll: true,
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
                set -e
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
                        allowMissing: false,
                        alwaysLinkToLastBuild: true,
                        keepAll: true,
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
                set -e
                cd notification-service
                go mod tidy
                go test -v -coverprofile=coverage.out -covermode=atomic ./...
                go tool cover -html=coverage.out -o coverage.html
                """
            }
            post {
                always {
                    publishHTML(target: [
                        allowMissing: false,
                        alwaysLinkToLastBuild: true,
                        keepAll: true,
                        reportDir: 'notification-service',
                        reportFiles: 'coverage.html',
                        reportName: 'Notification Service Coverage Report'
                    ])
                }
            }
        }
        // Security check stages //
        stage('Security Scan User Service') {
            when { expression { env.SERVICES_CHANGED.contains("user-service") } }
            steps {
                sh """
                set -e
                cd user-service
                npm audit --audit-level=moderate
                """
            }
        }
        stage('Security Scan Transaction Service') {
            when { expression { env.SERVICES_CHANGED.contains("transaction-service") } }
            steps {
                sh """
                set -e
                cd transaction-service
                if [ ! -d ".venv" ]; then
                    python3 -m venv .venv
                fi
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
                set -e
                cd notification-service
                gosec ./...
                """ 
            }
        }
    }
}