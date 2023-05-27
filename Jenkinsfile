String registryEndpoint = 'registry.comp.ystv.co.uk'

def image
String imageName = "ystv/computing-site:${env.BRANCH_NAME}-${env.BUILD_ID}"

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
        sh "git log --format="%H" -n 1"
        script {
          docker.withRegistry('https://' + registryEndpoint, 'docker-registry') {
            image = docker.build(imageName, "--build-arg COMP_SITE_VERSION_ARG=${env.BRANCH_NAME}-${env.BUILD_ID} --build-arg COMP_SITE_COMMIT_ARG=$(git log --format="%H" -n 1) .")
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



// pipeline {
//     agent any
//
//     environment {
//         REGISTRY_ENDPOINT = credentials('docker-comp-registry-endpoint')
//     }
//
//     stages {
//         stage('Update Components') {
//             steps {
//                 withCredentials([[$class: 'UsernamePasswordMultiBinding', credentialsId: 'comp-docker', usernameVariable: 'USERNAME', passwordVariable: 'PASSWORD']]) {
//                     sh 'docker login --username $USERNAME --password $PASSWORD $REGISTRY_ENDPOINT'
//                     sh 'docker pull golang:1.19-alpine'
//                 }
//             }
//         }
//         stage('Build') {
//             steps {
//                 sh 'docker build -t $REGISTRY_ENDPOINT/ystv/computing:$BUILD_ID .'
//             }
//         }
//         stage('Registry Upload') {
//             steps {
//                 sh 'docker push $REGISTRY_ENDPOINT/ystv/computing:$BUILD_ID' // Uploaded to registry
//             }
//         }
//         stage('Deploy') {
//             stages {
//                 stage('Staging') {
//                     when {
//                         branch 'master'
//                         not {
//                             expression { return env.TAG_NAME ==~ /v(0|[1-9]\d*)\.(0|[1-9]\d*)\.(0|[1-9]\d*)/ }
//                         }
//                     }
//                     environment {
//                         APP_ENV = credentials('comp-env')
//                         TARGET_SERVER = credentials('comp-server-address')
//                         TARGET_PATH = credentials('comp-server-path')
//                     }
//                     steps {
//                         sshagent(credentials : ['comp-server-key']) {
//                             script {
//                                 sh 'rsync -av $APP_ENV deploy@$TARGET_SERVER:$TARGET_PATH/computing/.env'
//                                 sh '''ssh -tt deploy@$TARGET_SERVER << EOF
//                                     docker pull $REGISTRY_ENDPOINT/ystv/computing:$BUILD_ID
//                                     docker rm -f ystv-computing
//                                     docker run -d -p 7075:7075 --env-file $TARGET_PATH/computing/.env --name ystv-computing --restart=always $REGISTRY_ENDPOINT/ystv/computing:$BUILD_ID
//                                     docker image prune -a -f --filter "label=site=computing"
//                                     exit 0
//                                 EOF'''
//                             }
//                         }
//                     }
//                 }
//                 /*stage('Production') {
//                     when {
//                         expression { return env.TAG_NAME ==~ /v(0|[1-9]\d*)\.(0|[1-9]\d*)\.(0|[1-9]\d*)/ } // Checking if it is main semantic version release
//                     }
//                     environment {
//                         APP_ENV = credentials('computing-env')
//                         TARGET_SERVER = credentials('prod-server-address')
//                         TARGET_PATH = credentials('prod-server-path')
//                     }
//                     steps {
//                         sshagent(credentials : ['prod-server-key']) {
//                             script {
//                                 sh 'rsync -av $APP_ENV deploy@$TARGET_SERVER:$TARGET_PATH/computing/.env'
//                                 sh '''ssh -tt deploy@$TARGET_SERVER << EOF
//                                     docker pull $REGISTRY_ENDPOINT/ystv/computing:$BUILD_ID
//                                     docker rm -f ystv-computing
//                                     docker run -d -p 7075:7075 --env-file $TARGET_PATH/computing/.env --name ystv-computing --restart=always $REGISTRY_ENDPOINT/ystv/computing:$BUILD_ID
//                                     docker image prune -a -f --filter "label=site=computing"
//                                     exit 0
//                                 EOF'''
//                             }
//                         }
//                     }
//                 }*/
//             }
//         }
//     }
//     post {
//         success {
//             echo 'Very cash-money'
//         }
//         failure {
//             echo 'That is not ideal, cheeky bugger'
//         }
//         always {
//             sh "docker image prune -f --filter label=site=computing --filter label=stage=builder" // Removing the local builder image
//             sh 'docker image prune -a -f --filter "label=site=computing"' // remove old image
//         }
//     }
// }
