# Work Distribution System

#### Simple work distribution system written in Golang and MongoDB

### Assumptions :
* Agent is a manual process/person (UI) who will manually accept and start working on the task in his queue.
* Lets consider that the Agent will manually mark the task as work in progress using some API interface.
* Another assumption is that the system doesnt account for time GAP in Agent completing and beginning a new task.

## Running the application
`docker-compose up -d`

You will see three container will come up.

```λ docker-compose up -d
 Creating network "task-scheduler_dev" with driver "bridge"
 Creating mongo_dev            ... done
 Creating mongo_seed           ... done
 Creating mongo_seed           ... done
```

The Go server is ready to accept request when the container log prints

`Connected to MongoDB!`

## Restarting the application

`docker-compose restart`

## Shutting down the application

`docker-compose down`

```λ docker-compose down
   Stopping task-scheduler_app_1 ... done
   Stopping mongo_dev            ... done
   Removing mongo_seed           ... done
   Removing task-scheduler_app_1 ... done
   Removing mongo_dev            ... done
   Removing network task-scheduler_dev
```

## Testing the application.

playground.http file in the repo can be used to test the below API endpoint in VSCode or Intellij Editor.

Create a Task
`POST http://localhost:3000/task`

Begin a task
`PUT http://localhost:3000/task/begin/5df9a8f0bd02f9890959e10a`

Complete a task
`PUT http://localhost:3000/task/complete/5df9a8f0bd02f9890959e10a`

List of Agents
`GET http://localhost:3000/agents`



