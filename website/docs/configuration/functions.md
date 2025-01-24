# Functions
```
Functions List
modify_file: Modify the file at path with the contents of content_text. Modified file must be full file content including modified content
    path
        Path of the file to be modified

    content_text
        The new content of the file

submit_files: Submit the modified files by GitHub Pull Request
    commit_message_short
        Short Commit message indicating purpose to change the file

    commit_message_detail
        Detail commit message indicating changes to the file

    pull_request_content
        Pull Request Content

get_pull_request: Get a GitHub Pull Request
    pr_number
        Pull Request Number to get

search_files: Search for files containing specific keyword (e.g., "xxx") within a directory path recursively
    keyword
        The keyword to search for.

    path
        The path to search within its directory

open_file: Open the file full content
    path
        The path of the file to open

list_files: List the files within the directory like Unix ls command. Each line contains the file mode, byte size, and name
    path
        The valid path to list within its directory

put_file: Put new file content to path
    path
        Path of the file to be changed to the new content

    content_text
        The new content of the file

```
