# Architecture

## Overview
![overview.png](overview.png)

1. We run the runner on the host machine or GitHub Actions
2. The runner starts the agent container using Docker
3. The agent container executes the agent binary
4. The agent binary communicates with LLM
5. The agent only executes the restricted functions
6. Finally, the agent creates a pull request to GitHub


## Runner

The runner is a program that simply runs the agent container for users.


## Use Docker Container

Containers are used to isolate the file system, processes, and other resources from the host machine.

We commonly use Docker for this purpose.


## Restricted Functions

Instead of executing a shell on the agent, function calling is used to process responses from the LLM.

To prevent unintended information leakage to external sources or the execution of unsafe commands, 
shell execution is strictly avoided.
