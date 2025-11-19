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
    }
}
