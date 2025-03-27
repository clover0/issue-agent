# Welcome to Issue Agent

## Introduction

Issue Agent is a lightweight tool powered by a Large Language Model (LLM).

When given an issue, the Agent autonomously attempts to solve the issue and submit the results as a Pull Request on GitHub.


## Why Issue Agent?

### Ready to use immediately

Issue Agent is a command line tool. 

It can be used as a GitHub Actions or installed on your machine.


### Very limited scope

Issue Agent has a very limited scope because it is designed to be secure and practical in use.

What we mean by secure here are the following:

* Issue Agent does not execute the response returned from the LLM. For example, avoid directly executing shell commands. Only predefined functions are allowed to use.
* Control the credentials given to the Agent. The agent can not retrieve and send confidential information in prompt from environment variables


### Handle simple tasks but difficult to automate tasks

Unlike AI tools like Copilot, which collaborate with developers to create deliverables,
this Agent handles tasks autonomously from start to finish.

The agent comprehends the initial instructions and works toward the goal based on those instructions.
Once the agent starts working, there is no interaction between the agent and the person who gave the instructions. 


## What the Issue Agent can and cannot do

* ✅ Pull requests are created only in repository belonging to GitHub issue.
* ✅ To read an issue in one GitHub repository and submit a PR to that repository.
* 🚫 Interactive development work between an Agent and the human who directs the Agent
* 🚫 The Agent can not execute arbitrary code. The Agent can only execute predefined functions.
