# Welcome to Issue Agent

## Introduction

Issue Agent is a lightweight tool powered by a Large Language Model (LLM).

When given an issue, the Agent autonomously attempts to solve the issue and submit the results as a Pull Request on GitHub.


## Why Issue Agent?

### Ready to use immediately

Issue Agent is a command line tool. 

It can be used as a GitHub Action or installed on your machine.


### Very limited scope

Issue Agent has a very limited scope because it is designed to be secure and practical in use.

What we mean by secure here are the following:

* Do not execute the response returned from the LLM as-is. For example, avoid directly executing shell commands. Only predefined functions are allowed to run._
* Control the credentials given to the Agent. Control the use of credentials.


### Handles simple but difficult to automate tasks

Unlike AI tools like Copilot, which collaborate with developers to create deliverables,
this Agent handles tasks autonomously from start to finish.

The agent comprehends the initial instructions and works toward the goal based on those instructions.
Once the agent starts working, there is no interaction between the agent and the person who gave the instructions. 

Therefore, tasks are completed quickly.
Tasks that require human evaluation of deliverables from various perspectives are not ideal for this tool.


### What the Issue Agent can and cannot do

* ✅ Pull requests are created only in repositories configured in the configuration file or specified via CLI flags
* ✅ To read an issue in one GitHub repository and submit a PR to that repository
* 🚫 Interactive development work between an Agent and the human who directs the Agent
* 🚫 Commits or pull requests to an unconfigured repository
* 🚫 Free operation using shell commands such as find, curl, etc.
