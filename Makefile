run_dev:
	docker-compose -f docker-compose.yml --env-file ./docker.env.dev up -d

logs_dev:
	docker attach taskmanager

kill_dev:
	docker-compose -f dev.docker-compose.yml --env-file ./docker.env.dev down
	docker-compose -f dev.docker-compose.yml --env-file ./docker.env.dev rm
	docker rmi server

lint:
	cd cmd/taskmanager
	go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest && golangci-lint run

run_tests:
	go test ./cmd/taskmanager -coverprofile=coverage.out && go tool cover -html=coverage.out

# create volumes
create_volumes:
	docker volume create ${CONTAINER_VOLUME_DEPLOYMENTS}
	docker volume create ${CONTAINER_VOLUME_CONFIGS}
	docker volume create ${SERVICE_MONGODB_VOLUME_DBDATA}
	docker volume create ${SERVICE_MONGODB_VOLUME_INITDB}

kill:
	docker stop ${CONTAINER_NAME} && docker rm ${CONTAINER_NAME}
	docker stop ${SERVICE_MONGODB} && docker rm ${SERVICE_MONGODB}
	docker stop mongo-express && docker rm mongo-express
	docker image rm ${IMAGE_NAME}

create_network:
	docker network create -d bridge ${NETWORK_NAME}

remove_network:
	docker network rm ${NETWORK_NAME}

run:
	##this process is just to copy the ansible config in the volume
	# create configs volume - this contains ansible config
	#[[ -z $(docker volume ls | grep ${CONTAINER_VOLUME_CONFIGS}) ]] && docker volume create ${CONTAINER_VOLUME_CONFIGS}
	# create a temp container with volume that maps data folder
	docker container create --name temp-ansible-copy -v ${CONTAINER_VOLUME_CONFIGS}:/data alpine:latest
	# copy the ansible config in the data folder of the container
	docker cp ${DEV_ANSIBLE_CONFIG_R25} temp-ansible-copy:/data/ansible.cfg
	# remove the container
	docker rm temp-ansible-copy

	# create the mongo db data volume
	#[[ -z $(docker volume ls | grep ${SERVICE_MONGODB_VOLUME_DBDATA}) ]] && docker volume create ${SERVICE_MONGODB_VOLUME_DBDATA}

	# create init db config - creating user
	#[[ -z $(docker volume ls | grep ${SERVICE_MONGODB_VOLUME_INITDB}) ]] && docker volume create ${SERVICE_MONGODB_VOLUME_INITDB}
	# again create a container and later copy the file in the mongo config file inside the container and then kill the container - this prcess will add the data in the volume
	docker container create --name temp-mongo-config-copy -v ${SERVICE_MONGODB_VOLUME_INITDB}:/data alpine:latest
	docker cp ${DEV_MONGO_INIT} temp-mongo-config-copy:/data/dev-mongo-init.js
	docker rm temp-mongo-config-copy

	# run mongo container
	docker run -d \
            --name=${SERVICE_MONGODB} \
            --restart=always \
            --network=${NETWORK_NAME} \
            -v ${SERVICE_MONGODB_VOLUME_DBDATA}:/data/db \
            -v ${SERVICE_MONGODB_VOLUME_INITDB}:/docker-entrypoint-initdb.d \
            -e MONGO_INITDB_ROOT_USERNAME=${DEV_MONGO_USERNAME} \
            -e MONGO_INITDB_ROOT_PASSWORD=${DEV_MONGO_USERPASS} \
            -e MONGO_INITDB_DATABASE=${DEV_MONGO_INITDB_DATABASE} \
            ${SERVICE_MONGODB_IMAGE}
    # run mongo-express for DB UI
	docker run -d \
          --name=mongo-express \
          --restart=always \
          --network=${NETWORK_NAME} \
          -p 8201:8081 \
          -e ME_CONFIG_OPTIONS_EDITORTHEME="ambiance" \
          -e ME_CONFIG_MONGODB_ENABLE_ADMIN="false" \
          -e ME_CONFIG_MONGODB_URL=mongodb://${DEV_MONGO_USERNAME}:${DEV_MONGO_USERPASS}@${SERVICE_MONGODB}:27017/${DEV_MONGO_INITDB_DATABASE} \
          -e ME_CONFIG_MONGODB_AUTH_USERNAME=${DEV_MONGO_USERNAME} \
          -e ME_CONFIG_MONGODB_AUTH_PASSWORD=${DEV_MONGO_USERPASS} \
          -e ME_CONFIG_MONGODB_AUTH_DATABASE=${DEV_MONGO_INITDB_DATABASE} \
          mongo-express
    # run server container
	docker run -d \
            --name=${CONTAINER_NAME} \
            --restart=always \
            --network=${NETWORK_NAME} \
            -v ${CONTAINER_VOLUME_DEPLOYMENTS}:/root/deployments \
            -v ${CONTAINER_VOLUME_CONFIGS}:/appconfigs \
            -p ${PORT}:80 \
            -e GIN_MODE=release \
            -e TASK_MANAGEMENT_SERVICE_DATABASE_CONNECTION_URI=mongodb://${DEV_MONGO_USERNAME}:${DEV_MONGO_USERPASS}@${SERVICE_MONGODB}:27017/${DEV_MONGO_INITDB_DATABASE} \
            -e ANSIBLE_CONFIG=/appconfigs/ansible.cfg \
            -e TASK_MANAGEMENT_SERVICE_DATABASE_DATABASE=${DEV_MONGO_INITDB_DATABASE} \
            ${IMAGE_NAME}

