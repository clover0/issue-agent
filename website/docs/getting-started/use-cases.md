# Use Cases

This section describe common use cases for Issue Agent.


## Simple but difficult to automate tasks

- Tasks that do not require immediate developer attention but are simple and expected to be accomplished asynchronously.
- Migrate deprecated statements that cannot be handled by tools
- Update forgotten documentation
- Delete files that are no longer needed


### Let's take a closer look.

Addition or modification of wording accompanying feature additions or changes.

GitHub Issue #1

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

GitHub Issue #2

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


## Horizontal deployment of tasks that are difficult to automate

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
