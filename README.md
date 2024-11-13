# Tesla Fetch
Teslafetch is a command line tool written in go 1.22.5. Teslafetch can be used to view the status of your vehicle from the 
command line in a convient and aesthetically pleasing way. you can also send commands to your vehicle straight from your terminal. 
For a list of currently implimented commands see below.

The repo is split into two parts. One is the web server you can run locally to interact with the Tesla API. Two is the client app 
that will display and draw vehicle data to your terminal.

Teslafetch does not support any Model X or S released before 2018.

## Commands
### Get Data

![Kapture 2024-10-20 at 16 01 41](https://github.com/user-attachments/assets/5f7ef76a-180c-4e66-b64d-700bc429fa1a)

running `tfetch` with no command line argument first checks the status of your car then polls the `wake` api endpoint 
from Tesla to wake your car. It calls the `vehicle_data` endpoint and displays your cars current status.

### Commands
#### Lock
`tfetch lock`
will lock your car and will return an error if already locked
#### Unlock
`tfetch unlock`
will unlock your car and return an error if its already unlocked
#### Climate
`tfetch climate`
will toggle climate on if off and off if on
