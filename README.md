# mass-pr
Tool to create PRs to multiple repos (do not use for evil)
It is built for adding files to our student's repo to add new chapters of courses using a PR.

This tools is created to be used next to GitHub Classroom or [repo-create](https://github.com/meyskens/repo-create).

## How to use
```console
go run ./cmd/mass-pr/ add-directory -o itfactory-tm --files ./1ITF_2020_2021_Webdesign_Essentials/2_CSS/ --to 2_CSS --prefix '2020-2021-*-webdesign-essentials-meyskens'
```
Note: some parts are still quite hard coded... maybe don't use this yet