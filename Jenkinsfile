String registryEndpoint = 'registry.comp.ystv.co.uk'

def vaultConfig = [vaultUrl: 'https://vault.comp.ystv.co.uk',
                  vaultCredentialId: 'vault-ansible',
                  engineVersion: 2]

def branch = env.BRANCH_NAME.replaceAll("/", "_")
def image
String imageName = "ystv/computing-site:${branch}-${env.BUILD_ID}"

pipeline {
  agent {
    label 'docker'
  }

  environment {
    DOCKER_BUILDKIT = '1'
  }

  stages {
    stage('Build image') {
      steps {
        script {
          GIT_COMMIT_HASH = sh (script: "git log -n 1 --pretty=format:'%H'", returnStdout: true)
          def secrets = [
            [path: "ci/ystv-internal-certs", engineVersion: 2, secretValues: [
              [envVar: 'COMP_SITE_CERT_PEM', vaultKey: 'cert'],
              [envVar: 'COMP_SITE_KEY_PEM', vaultKey: 'key']
            ]]
          ]
          withVault([configuration: vaultConfig, vaultSecrets: secrets]) {
            docker.withRegistry('https://' + registryEndpoint, 'docker-registry') {
              image = docker.build(imageName, "--build-arg COMP_SITE_VERSION_ARG=${env.BRANCH_NAME}-${env.BUILD_ID} --build-arg COMP_SITE_COMMIT_ARG=${GIT_COMMIT_HASH} --build-arg COMP_SITE_CERT_PEM='${COMP_SITE_CERT_PEM}' --build-arg COMP_SITE_KEY_PEM='${COMP_SITE_KEY_PEM}' .")
            }
          }
        }
      }
    }

    stage('Push image to registry') {
      steps {
        script {
          docker.withRegistry('https://' + registryEndpoint, 'docker-registry') {
            image.push()
            if (env.BRANCH_IS_PRIMARY) {
              image.push('latest')
            }
          }
        }
      }
    }

    stage('Deploy') {
      stages {
        stage('Development') {
          when {
            expression { env.BRANCH_IS_PRIMARY }
          }
          steps {
            build(job: 'Deploy Nomad Job', parameters: [
              string(name: 'JOB_FILE', value: 'computing-site-dev.nomad'),
              text(name: 'TAG_REPLACEMENTS', value: "${registryEndpoint}/${imageName}")
            ])
          }
        }

        stage('Production') {
          when {
            // Checking if it is semantic version release.
            expression { return env.TAG_NAME ==~ /v(0|[1-9]\d*)\.(0|[1-9]\d*)\.(0|[1-9]\d*)/ }
          }
          steps {
            build(job: 'Deploy Nomad Job', parameters: [
              string(name: 'JOB_FILE', value: 'computing-site-prod.nomad'),
              text(name: 'TAG_REPLACEMENTS', value: "${registryEndpoint}/${imageName}")
            ])
          }
        }
      }
    }
  }
}
