# Core Concept

## Design Principles

**Secure by restricting**:

Ensure security by restricting functionality and preventing the execution of  
shell commands within the Agent environment. 
Instead of trying to have AI handle complex tasks well through detailed instructions, 
the focus should be on limiting what it is allowed to do.


**Simple, small things**: 

We do not engage in interactive development with the Agents.  
Instead, we provide them with initial instructions to complete tasks and deliver results.  
The goal is to structure development tasks with sufficient granularity so they can be completed by instructing the Agent once,  
without requiring us to open an editor, make changes, or create a Pull Request.

**No interactive, quick work**:

Our goal is not to spend time interacting with Agents, but to delegate tasks that are essential yet have been challenging to automate,  
allowing us to optimize how we use our time.
