# jclubtakeaways-api

Deploys a heroku app to pair with https://github.com/lukemassa/nickandluke.

Runs an API that takes as input guest names, and returns whether or not they are valid guests, and which google form to send them.

## Deployment

Heroku app is called nickandluke-api in lukefmassa@hotmail.com account

### Update

`git push heroku main`

### Add guest

1. `./run --action download`
1. Edit `staging/guests.csv`
1. `./run --action upload`
  1. This should update `hashed.txt`
1. Commit `hashed.txt` and run `git push heroku main`.
