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
    }
}

