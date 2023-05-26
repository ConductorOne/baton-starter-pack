# baton-starter-pack
A simple template for quickly building your own Baton connector.

## How to use
1. Click on the "use this template" button, and you'll have a new repository for your connector.
2. On your new repository, ensure that workflow actions have read and write access. This is available under the Actions settings for the repository. This setting can be reverted once the repository has been provisioned.
3. Update `cookiecutter.json` with the appropriate configuration
    ```json
    {
       "repo_owner": "conductorone",
       "repo_name": "baton-example",
       "name": "baton-example"
    }
    ```
10. Commit this change, and a Github action will process the update and initialize the repo for you.
