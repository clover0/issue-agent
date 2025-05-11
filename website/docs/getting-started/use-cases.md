# Use Cases

This section describe common use cases for Issue Agent.


## Simple but difficult to automate tasks

- Tasks that do not require immediate developer attention but are simple and expected to be accomplished asynchronously.
- Migrate deprecated statements that tools cannot handle.
- Update forgotten documentation.
- Delete files that are no longer needed.


### Let's take a closer look

Additions or modifications of wording following feature additions or changes.

**GitHub Issue #1**

```markdown
Change "wording" to "new wording" in the `dir1/` directory.
```

Issue Agent thinks and executes the functions:

- Repository and code analysis
    - list_files
    - open_file
    - ...
- Decide the changes
    - modify_file
    - ...
- Submit the changes
    - submit_files

Issue Agent creates a Pull Request with the following changes:

```diff
--- dir1/example.txt
+++ dir1/example.txt
@@ -1,5 +1,5 @@

-Sometimes xxx is written multiple times: xxx.
+Sometimes yyy is written multiple times: yyy.
 Here is the end of the example.
```


**GitHub Issue #2**

```markdown
Fix all typos present in the comments under the `dir1/` directory.
```

Issue Agent will create a Pull Request:

```diff
--- dir1/document.txt
+++ dir1/document.txt
@@ -1,5 +1,5 @@
-Please make sure to recieve the package on time.
+Please make sure to receive the package on time.

 This document is an example with a typo.

-We often see common typos like "recieve."
+We often see common typos like "receive."
```


## Horizontal deployment of tasks that are challenging to automate

For changes that require wide-ranging adjustments,
create a Pull Request for some parts handled by a human developer.
Then, a developer apply similar adjustments to other areas.

GitHub Issue #1

```markdown
Make similar changes to `path/to/dir`, as shown in the Pull Request below.
https://github.com/clover0/example-repository/pull/80
```

Issue Agent get pull request written in the issue:

- get_pull_request
    - from clover0/example-repository
    - number 80

Issue Agent thinks and executes the functions:

- Repository and code analysis
    - list_files
    - open_file
    - ...
- Decide the changes
    - modify_file
    - ...
- Submit the changes
    - submit_files

Issue Agent will create a Pull Request:

- Like the Pull Request #80, the Issue Agent will create a Pull Request with similar changes in the `path/to/dir`
  directory.

...

Repeat Issue #1 and the creation of pull requests by Issue Agent for the range which we want to apply.


## Code review with guidelines from another repository

For code reviews, it's common to have review guidelines or checklists in a separate repository. Issue Agent can load these guidelines and use them to perform a review.


!!! warning "Another repository"

    When referring to "another repository", it is limited to repositories owned by your organization or yourself, 
    in order to avoid retrieving content from untrusted repositories.
    If you want to refer to a public repository, you need to copy the file to your own repository.

    ### Usable Patterns
    - your-org/repo1
    - your-org/repo2

    When repo2 refers to repo1 (e.g., using guidelines or making references), this is **allowed**.

    ### Not Supported Patterns
    - public-user/repoA
    - your-org/repo2

    When repo2 refers to repoA (e.g., trying to access content from a public user's repository), this is **not permitted**.


There is the following GitHub Issue.

```markdown
Review the pull request #123 using the review guidelines 
from the repository `organization/review-guidelines` at path `guidelines/code-review-checklist.md`.
```

Issue Agent thinks and executes the functions:

Get review guidelines from another repository.

- get_repository_content
    - from `organization/review-guidelines`
    - path `guidelines/code-review-checklist.md`

Get pull request details.

- get_pull_request
    - number 123

Repository and code analysis.

- list_files
- open_file


Perform review based on guidelines.

- create_pull_request_review_comment
- ...

Issue Agent will create review comments on the pull request:

```
[In file src/main.js, line 45]
According to the code review guidelines, error handling should be implemented for all API calls.
Consider adding try/catch block here to handle potential network errors.

[In file src/utils/helpers.js, line 120]
The review guidelines recommend using descriptive variable names.
Consider renaming `x` to something more descriptive of its purpose.
```

This use case demonstrates how Issue Agent can leverage existing documentation 
and guidelines from different repositories to perform more standardized and thorough reviews.
