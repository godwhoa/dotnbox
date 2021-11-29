## dotnbox

> Like [Lichess](https://en.wikipedia.org/wiki/Lichess) but for dots and boxes.

### Status
Playable and useable. Needs major UI/UX improvements.

### Running
```bash
# Backend
cd backend/
go run main.go
# Or if you want pretty logs
go run main.go 2>&1 | jq -R 'fromjson? | select(type == "object")'

# Frontend
cd frontend
npm ci
npm start
```


### Screenshot
![screenshot of the game](https://i.imgur.com/tOox3FV.png)