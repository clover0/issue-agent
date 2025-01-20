# Core Concept

## Design Principles

* **Secure by Limiting**: Ensure security by limiting functionality and avoiding the execution of 
various shell commands within the Agent environment. We want AI to do well, or rather control well.

* **Simple, small things**: We do not do any interactive development work with the Agents,
but only give initial instructions to them to work on the tasks and deliver the results.
The goal is to do development work with a granularity that can be resolved by instructing the Agent once,
rather than us having to open an editor to make changes and create a Pull Request.

* **No interactive, quick work**:
  It is not that we want to spend our time talking with Agents, but to ask Agents to perform tasks that are necessary but have been difficult to automate so far,
so that we can make the best use of our time.
