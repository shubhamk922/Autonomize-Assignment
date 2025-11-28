AI Agents Assignment 

docker build command 

    docker build -t team-monitoring -f build/docker/Dockerfile .

docker run command 

   docker run -d -p 8080:8080 \
  -e AWS_ACCESS_KEY_ID=$AWS_ACCESS_KEY_ID \
  -e AWS_SECRET_ACCESS_KEY=$AWS_SECRET_ACCESS_KEY \
  -e AWS_REGION=$AWS_REGION \
  -e OPENAI_API_KEY=$OPENAI_API_KEY \
  team-monitoring


docker logs 
    docker exec -it 3557e7f777ca sh

docker compose -f build/docker/docker-compose.yml up -d --build

docker compose -f build/docker/docker-compose.yml down
export PATH=$PATH:/usr/local/go/bin
