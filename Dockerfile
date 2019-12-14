FROM mongo:latest
COPY mongo/agents.json /agents.json
CMD mongoimport --host mongo --db tasker --collection agents --drop --file /agents.json --jsonArray
